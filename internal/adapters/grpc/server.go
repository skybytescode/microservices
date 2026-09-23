package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/skybytescode/microservices-proto/golang/order"
	"github.com/skybytescode/order-service/config"
	"github.com/skybytescode/order-service/internal/application/core/domain"
	"github.com/skybytescode/order-service/internal/ports"
)

type Adapter struct {
	api  ports.APIPort
	port int
	order.UnimplementedOrderServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a Adapter) Create(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	var items []domain.OrderItem
	for _, it := range req.Items {
		items = append(items, domain.OrderItem{ProductCode: it.Name})
	}
	newOrder := domain.NewOrder(req.UserId, items)
	result, err := a.api.PlaceOrder(newOrder)
	if err != nil {
		return nil, err
	}
	return &order.CreateOrderResponse{OrderId: result.ID}, nil
}

func (a Adapter) Run() {
	listen, _ := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	grpcServer := grpc.NewServer()
	order.RegisterOrderServer(grpcServer, a)
	if config.GetEnv() == "development" {
		reflection.Register(grpcServer) // lets grpcurl introspect the API
	}
	grpcServer.Serve(listen)
}
