package ports

import (
	"context"

	"github.com/skybytescode/microservices/payment-service/internal/application/core/domain"
)

type APIPort interface {
	Charge(ctx context.Context, payment domain.Payment) (domain.Payment, error)
}
