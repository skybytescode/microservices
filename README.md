# microservices

Go microservices that talk to each other over gRPC.

| Folder | What it is |
|---|---|
| [`order-service/`](order-service) | Order service: gRPC server with a MySQL (GORM) store |
| [`payment-service/`](payment-service) | Payment service (not implemented yet) |

Each folder is its own Go module.

The `.proto` contracts and the Go code generated from them live in
[skybytescode/microservices-proto](https://github.com/skybytescode/microservices-proto).
Each service's stubs are a separate module there, released with tags like
`golang/order/v1.0.0`. To use a new version in a service:

```sh
cd order-service
go get github.com/skybytescode/microservices-proto/golang/order@v1.0.1
```

## Run the order service

```sh
cd order-service
DATA_SOURCE_URL='user:pass@tcp(localhost:3306)/order?parseTime=true' \
APPLICATION_PORT=9000 go run ./cmd
```
