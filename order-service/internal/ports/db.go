package ports

import "github.com/skybytescode/microservices/order-service/internal/application/core/domain"

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(*domain.Order) error
}
