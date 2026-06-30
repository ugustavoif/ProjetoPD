package ports

import "github.com/ruandg/microservices/order/internal/application/core/domain"

type ShippingPort interface {
	CalculateShipping(order *domain.Order) (int32, error)
}
