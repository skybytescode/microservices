# microservices

Go microservices that talk to each other over gRPC.

| Folder | What it is |
|---|---|
| [`order-service/`](order-service) | Order service: gRPC server with a MySQL (GORM) store |
| [`payment-service/`](payment-service) | Payment service: gRPC server with a MySQL (GORM) store |
| [`e2e/`](e2e) | End-to-end test: runs the whole stack with Docker Compose |
| [`mysql/`](mysql) | MySQL manifest for Kubernetes |
| [`jaeger/`](jaeger) | Jaeger manifest for Kubernetes (trace collector and UI) |
| [`kind/`](kind) | Config for the local kind cluster |

`order-service`, `payment-service` and `e2e` are each their own Go module,
tied together by `go.work`.

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

Runs on a local [kind](https://kind.sigs.k8s.io/) cluster named `microservices`,
deployed with [Skaffold](https://skaffold.dev/). Requests reach the services
through ingress-nginx on `localhost:443`.

### 1. Create the cluster (once)

```sh
kind create cluster --config kind/kind-config.yaml
```

[`kind/kind-config.yaml`](kind/kind-config.yaml) maps `127.0.0.1:80` and
`127.0.0.1:443` into the cluster. Port mappings can only be set when a cluster
is created, so changing them means deleting and recreating it
(`kind delete cluster --name microservices`).

### 2. Install the ingress controller (once)

```sh
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.15.1/deploy/static/provider/kind/deploy.yaml
kubectl -n ingress-nginx wait --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller --timeout=180s
```

### 3. Deploy

```sh
skaffold dev
```

This builds the `order` and `payment` images, loads them into the cluster and
deploys Jaeger, MySQL, Order and Payment. It rebuilds when you save a file and
removes everything on Ctrl+C. Use `skaffold run` to deploy once and
`skaffold delete` to remove it. [`skaffold.yaml`](skaffold.yaml) always deploys
to the `kind-microservices` context, whatever your current context is.

### 4. Call the Order service

```sh
cd ../microservices-proto   # a clone of skybytescode/microservices-proto
grpcurl -insecure -import-path order -proto order.proto \
  -d '{"user_id":1,"order_items":[{"product_code":"CAM","unit_price":10,"quantity":2}]}' \
  localhost:443 Order/Create
```

- The Ingress sends paths starting with `/Order` to the `order` service and
  `/Payment` to the `payment` service.
- gRPC through ingress-nginx needs TLS. The controller answers with its built-in
  self-signed certificate, hence `-insecure`.
- The pods run with `ENV=prod`, which turns off gRPC reflection, so grpcurl
  needs the `.proto` files.

### 5. Open the Jaeger UI

```sh
kubectl -n jaeger port-forward svc/jaeger-otel 16686:16686
```

Open http://localhost:16686, pick service **order** and click **Find Traces**.
An `Order/Create` call shows up as one trace across Order and Payment. Jaeger
keeps traces in memory, so they are lost when its pod restarts. So is the MySQL
data: MySQL has no persistent volume.

### Troubleshooting

| Symptom | Likely cause |
|---|---|
| `kind create` fails with "address already in use" | Something else uses port 80 or 443: `sudo ss -ltnp \| grep -E ':80 \|:443 '` |
| grpcurl: `connection refused` | The ingress controller is not ready yet; rerun the `wait` from step 2 |
| grpcurl: `404` or `Unimplemented` | The Ingress is missing: `kubectl get ingress` |
| Pods stuck in `Init` | MySQL is not ready yet; the init container waits for it |

## Tests

```sh
cd order-service && go test ./...   # unit tests + MySQL integration test (needs Docker)
cd e2e && go test ./...             # builds both images and runs the full stack (needs Docker)
```

## Tracing

Both services export OpenTelemetry traces over OTLP when
`OTEL_EXPORTER_OTLP_ENDPOINT` is set (for example `http://localhost:4318` for a
local Jaeger). Logs are JSON and include `trace_id` and `span_id`.
