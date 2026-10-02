package api

import (
	"context"
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/skybytescode/microservices/order-service/internal/application/core/domain"
	"github.com/skybytescode/microservices/order-service/internal/ports"
)

type Application struct {
	db       ports.DBPort
	payment  ports.PaymentPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
	return &Application{
		db:       db,
		payment:  payment,
		shipping: shipping,
	}
}

func (a Application) PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	err := a.db.Save(ctx, &order)
	if err != nil {
		return domain.Order{}, err
	}
	if err := a.payment.Charge(ctx, &order); err != nil {
		return domain.Order{}, orderCreationFailed("payment", err)
	}
	if err := a.shipping.Ship(ctx, &order); err != nil {
		return domain.Order{}, orderCreationFailed("shipping", err)
	}
	return order, nil
}

// orderCreationFailed wraps an error from a downstream service as
// InvalidArgument, with the service's message as a field violation.
func orderCreationFailed(field string, err error) error {
	st, _ := status.FromError(err)
	fieldErr := &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: st.Message(),
	}
	badReq := &errdetails.BadRequest{}
	badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)
	orderStatus := status.New(codes.InvalidArgument, "order creation failed")
	statusWithDetails, _ := orderStatus.WithDetails(badReq)
	return statusWithDetails.Err()
}

func (a Application) GetOrder(ctx context.Context, id int64) (domain.Order, error) {
	order, err := a.db.Get(ctx, id)
	if errors.Is(err, domain.ErrOrderNotFound) {
		return domain.Order{}, status.Errorf(codes.NotFound, "order %d not found", id)
	}
	return order, err
}
