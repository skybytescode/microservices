package grpc

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/skybytescode/microservices-proto/golang/shipping"
	"github.com/skybytescode/microservices/shipping-service/config"
	"github.com/skybytescode/microservices/shipping-service/internal/adapters/metrics"
	"github.com/skybytescode/microservices/shipping-service/internal/ports"
)

type Adapter struct {
	api    ports.APIPort
	port   int
	server *grpc.Server
	shipping.UnimplementedShippingServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a *Adapter) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error: %v", a.port, err)
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(metrics.Server.UnaryServerInterceptor(metrics.Exemplar)),
	)
	a.server = grpcServer
	shipping.RegisterShippingServer(grpcServer, a)
	metrics.Server.InitializeMetrics(grpcServer)
	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}

	log.Printf("starting shipping service on port %d ...", a.port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d", a.port)
	}
}

func (a *Adapter) Stop() {
	a.server.Stop()
}
