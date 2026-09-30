package api

import (
	"context"

	"github.com/skybytescode/microservices/shipping-service/internal/application/core/domain"
	"github.com/skybytescode/microservices/shipping-service/internal/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

func (a Application) Ship(ctx context.Context, shipping domain.Shipping) (domain.Shipping, error) {
	shipping.DeliveryDays = shipping.EstimateDeliveryDays()
	err := a.db.Save(ctx, &shipping)
	if err != nil {
		return domain.Shipping{}, err
	}
	return shipping, nil
}
