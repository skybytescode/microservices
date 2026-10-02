package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/skybytescode/microservices/order-service/config"
	"github.com/skybytescode/microservices/order-service/internal/adapters/db"
	"github.com/skybytescode/microservices/order-service/internal/adapters/grpc"
	"github.com/skybytescode/microservices/order-service/internal/adapters/metrics"
	"github.com/skybytescode/microservices/order-service/internal/adapters/payment"
	"github.com/skybytescode/microservices/order-service/internal/adapters/shipping"
	"github.com/skybytescode/microservices/order-service/internal/application/core/api"
)

const (
	service     = "order"
	environment = "dev"
	id          = 1
)

// tracerProvider exports spans over OTLP to the endpoint in
// OTEL_EXPORTER_OTLP_ENDPOINT (Jaeger accepts OTLP directly). Without that
// variable spans are still created, so trace IDs appear in logs, but not exported.
func tracerProvider(ctx context.Context) (*tracesdk.TracerProvider, error) {
	opts := []tracesdk.TracerProviderOption{
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		exp, err := otlptracehttp.New(ctx)
		if err != nil {
			return nil, err
		}
		opts = append(opts, tracesdk.WithBatcher(exp))
	}
	return tracesdk.NewTracerProvider(opts...), nil
}

func init() {
	log.SetFormatter(customLogger{
		formatter: log.JSONFormatter{FieldMap: log.FieldMap{
			"msg": "message",
		}},
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

type customLogger struct {
	formatter log.JSONFormatter
}

func (l customLogger) Format(entry *log.Entry) ([]byte, error) {
	span := trace.SpanFromContext(entry.Context)
	entry.Data["trace_id"] = span.SpanContext().TraceID().String()
	entry.Data["span_id"] = span.SpanContext().SpanID().String()
	return l.formatter.Format(entry)
}

func main() {
	tp, err := tracerProvider(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))

	if port := config.GetMetricsPort(); port != 0 {
		metrics.Serve(port)
	}

	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	paymentAdapter, err := payment.NewAdapter(config.GetPaymentServiceUrl())
	if err != nil {
		log.Fatalf("Failed to initialize payment stub. Error: %v", err)
	}

	shippingAdapter, err := shipping.NewAdapter(config.GetShippingServiceUrl())
	if err != nil {
		log.Fatalf("Failed to initialize shipping stub. Error: %v", err)
	}

	application := api.NewApplication(dbAdapter, paymentAdapter, shippingAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
