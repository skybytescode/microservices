package shipping

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skybytescode/microservices-proto/golang/shipping"
	"github.com/skybytescode/microservices/order-service/internal/adapters/metrics"
	"github.com/skybytescode/microservices/order-service/internal/application/core/domain"
)

type Adapter struct {
	shipping shipping.ShippingClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithChainUnaryInterceptor(metrics.Client.UnaryClientInterceptor(metrics.Exemplar)),
	)
	conn, err := grpc.NewClient(shippingServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	client := shipping.NewShippingClient(conn)
	return &Adapter{shipping: client}, nil
}

func (a *Adapter) Ship(ctx context.Context, order *domain.Order) error {
	var items []*shipping.ShippingItem
	for _, orderItem := range order.OrderItems {
		items = append(items, &shipping.ShippingItem{
			ProductCode: orderItem.ProductCode,
			Quantity:    orderItem.Quantity,
		})
	}
	_, err := a.shipping.Create(ctx, &shipping.CreateShippingRequest{
		UserId:  order.CustomerID,
		OrderId: order.ID,
		Items:   items,
	})
	return err
}
