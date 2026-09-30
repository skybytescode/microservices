package grpc

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/skybytescode/microservices-proto/golang/shipping"
	"github.com/skybytescode/microservices/shipping-service/internal/application/core/domain"
)

func (a *Adapter) Create(ctx context.Context, request *shipping.CreateShippingRequest) (*shipping.CreateShippingResponse, error) {
	log.WithContext(ctx).Info("Creating shipping...")
	var items []domain.ShippingItem
	for _, item := range request.Items {
		items = append(items, domain.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}
	newShipping := domain.NewShipping(request.UserId, request.OrderId, items)
	result, err := a.api.Ship(ctx, newShipping)
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to ship. %v ", err)).Err()
	}
	return &shipping.CreateShippingResponse{ShippingId: result.ID, DeliveryDays: result.DeliveryDays}, nil
}
