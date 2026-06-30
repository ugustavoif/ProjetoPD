package api

import (
	"github.com/ruandg/microservices/shipping/internal/application/core/domain"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (a Application) CalculateShipping(shipping domain.Shipping) (domain.Shipping, error) {
	// The shipping struct is already created with the delivery time calculated
	// in domain.NewShipping, so we just return it.
	return shipping, nil
}
