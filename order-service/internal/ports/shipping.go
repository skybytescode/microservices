package ports

import (
	"context"

	"github.com/skybytescode/microservices/order-service/internal/application/core/domain"
)

type ShippingPort interface {
	Ship(context.Context, *domain.Order) error
}
