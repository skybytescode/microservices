package ports

import (
	"context"

	"github.com/skybytescode/microservices/shipping-service/internal/application/core/domain"
)

type APIPort interface {
	Ship(ctx context.Context, shipping domain.Shipping) (domain.Shipping, error)
}
