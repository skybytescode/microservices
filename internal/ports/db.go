package ports

import "github.com/skybytescode/order-service/internal/application/core/domain"

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(*domain.Order) error
}
