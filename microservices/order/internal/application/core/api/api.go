package api

import (
	"fmt"
	"log"

	"github.com/ruandg/microservices/order/internal/application/core/domain"
	"github.com/ruandg/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db       ports.DBPort
	payment  ports.PaymentPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
	return &Application{
		db:       db,
		payment:  payment,
		shipping: shipping,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	// Validate total items (max 50)
	var totalItems int32
	for _, item := range order.OrderItems {
		totalItems += item.Quantity
	}
	if totalItems > 50 {
		return domain.Order{}, status.Error(codes.InvalidArgument, fmt.Sprintf("Pedido com mais de 50 itens não é permitido. Total de itens: %d", totalItems))
	}

	var productCodes []string
	for _, item := range order.OrderItems {
		productCodes = append(productCodes, item.ProductCode)
	}

	if err := a.db.CheckProductsExist(productCodes); err != nil {
		return domain.Order{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err := a.db.Save(&order)
	if err != nil {
		a.db.UpdateStatus(fmt.Sprintf("%d", order.ID), "Canceled")
		return domain.Order{}, err
	}
	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		if status.Code(paymentErr) == codes.DeadlineExceeded {
			log.Printf("Tempo limite (deadline) excedido ao chamar o serviço de pagamento: %v", paymentErr)
		}
		a.db.UpdateStatus(fmt.Sprintf("%d", order.ID), "Canceled")
		return domain.Order{}, paymentErr
	}
	a.db.UpdateStatus(fmt.Sprintf("%d", order.ID), "Paid")

	shippingDays, shippingErr := a.shipping.CalculateShipping(&order)
	if shippingErr != nil {
		if status.Code(shippingErr) == codes.DeadlineExceeded {
			log.Printf("Tempo limite (deadline) excedido ao chamar o serviço de envio: %v", shippingErr)
		}
	} else {
		log.Printf("====================================================")
		log.Printf("Pedido %d processado com sucesso!", order.ID)
		log.Printf("Status do Pagamento: Aprovado")
		log.Printf("Previsão de Entrega: %d dias", shippingDays)
		log.Printf("Total de itens no pedido: %d", totalItems)
		log.Printf("====================================================")
	}

	return order, nil
}