package main

import (
	"log"

	"github.com/skybytescode/order-service/config"
	"github.com/skybytescode/order-service/internal/adapters/db"
	"github.com/skybytescode/order-service/internal/adapters/grpc"
	"github.com/skybytescode/order-service/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}
	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
