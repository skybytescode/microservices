# microservices

Go microservices that talk to each other over gRPC.

| Folder | What it is |
|---|---|
| [`order-service/`](order-service) | Order service: gRPC server with a MySQL (GORM) store |
| [`payment-service/`](payment-service) | Payment service: gRPC server with a MySQL (GORM) store |
| [`e2e/`](e2e) | End-to-end test: runs the whole stack with Docker Compose |
| [`mysql/`](mysql) | MySQL manifest for Kubernetes |
| [`jaeger/`](jaeger) | Jaeger manifest for Kubernetes (trace collector and UI) |

Each folder is its own Go module.

The `.proto` contracts and the Go code generated from them live in
[skybytescode/microservices-proto](https://github.com/skybytescode/microservices-proto).
Each service's stubs are a separate module there, released with tags like
`golang/order/v1.0.0`. To use a new version in a service:

```sh
cd order-service
go get github.com/skybytescode/microservices-proto/golang/order@v1.0.1
```

## How the services talk

A client calls `Order.Create`. Order saves the order, then calls `Payment.Create`
over gRPC to charge the total price. If the charge fails, Order returns
`InvalidArgument` with the payment error in the status details.

## Run everything on Kubernetes

Needs a local cluster (kind or minikube) and [Skaffold](https://skaffold.dev/):

```sh
skaffold dev
```

This deploys Jaeger, MySQL, Order and Payment. To call Order and to open the
Jaeger UI at http://localhost:16686:

```sh
kubectl port-forward svc/order 8080:8080
kubectl -n jaeger port-forward svc/jaeger-otel 16686:16686
```

## Tests

```sh
cd order-service && go test ./...   # unit tests + MySQL integration test (needs Docker)
cd e2e && go test ./...             # builds both images and runs the full stack (needs Docker)
```

## Tracing

Both services export OpenTelemetry traces over OTLP when
`OTEL_EXPORTER_OTLP_ENDPOINT` is set (for example `http://localhost:4318` for a
local Jaeger). Logs are JSON and include `trace_id` and `span_id`.
