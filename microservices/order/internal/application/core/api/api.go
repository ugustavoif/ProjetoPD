package api

import (
	"fmt"

	"github.com/ruandg/microservices/order/internal/application/core/domain"
	"github.com/ruandg/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	// Validate total items (max 50)
	var totalItems int32
	for _, item := range order.OrderItems {
		totalItems += item.Quantity
	}
	if totalItems > 50 {
		return domain.Order{}, status.Error(codes.InvalidArgument, fmt.Sprintf("Order with more than 50 items is not allowed. Total items: %d", totalItems))
	}

	err := a.db.Save(&order)
	if err != nil {
		a.db.UpdateStatus(fmt.Sprintf("%d", order.ID), "Canceled")
		return domain.Order{}, err
	}
	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		a.db.UpdateStatus(fmt.Sprintf("%d", order.ID), "Canceled")
		return domain.Order{}, paymentErr
	}
	a.db.UpdateStatus(fmt.Sprintf("%d", order.ID), "Paid")
	return order, nil
}
