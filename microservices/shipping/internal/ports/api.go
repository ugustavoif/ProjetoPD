package ports

import "github.com/ruandg/microservices/shipping/internal/application/core/domain"

type APIPort interface {
	CalculateShipping(shipping domain.Shipping) (domain.Shipping, error)
}
