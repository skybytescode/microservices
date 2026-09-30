package ports

import (
	"context"

	"github.com/skybytescode/microservices/shipping-service/internal/application/core/domain"
)

type DBPort interface {
	Save(ctx context.Context, shipping *domain.Shipping) error
}
