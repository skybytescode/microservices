# microservices

Go microservices that talk to each other over gRPC.

| Folder | What it is |
|---|---|
| [`microservices-proto/`](microservices-proto) | `.proto` contracts and the Go code generated from them |
| [`order-service/`](order-service) | Order service: gRPC server with a MySQL (GORM) store |

Each folder is its own Go module. The services use the proto module from this
repo through a `replace` directive in their `go.mod`, so changes to a `.proto`
file are picked up without publishing a new version.

## Regenerate the Go code after editing a .proto file

```sh
cd microservices-proto && make proto
```

## Run the order service

```sh
cd order-service
DATA_SOURCE_URL='user:pass@tcp(localhost:3306)/order?parseTime=true' \
APPLICATION_PORT=9000 go run ./cmd
```
