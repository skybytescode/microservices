package ports

import "github.com/skybytescode/microservices/order-service/internal/application/core/domain"

type APIPort interface {
	PlaceOrder(order domain.Order) (domain.Order, error)
}
