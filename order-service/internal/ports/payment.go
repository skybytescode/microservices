package ports

import (
	"context"

	"github.com/skybytescode/microservices/order-service/internal/application/core/domain"
)

type PaymentPort interface {
	Charge(context.Context, *domain.Order) error
}
