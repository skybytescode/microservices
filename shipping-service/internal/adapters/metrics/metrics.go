// Package metrics exposes the service's Prometheus metrics over HTTP.
package metrics

import (
	"context"
	"fmt"
	"net/http"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

// Server records requests handled by this service's gRPC server.
var Server = grpcprom.NewServerMetrics(grpcprom.WithServerHandlingTimeHistogram())

func init() {
	prometheus.MustRegister(Server)
}

// Exemplar attaches the trace ID to each observation, so a latency sample in
// Prometheus links to its trace in Jaeger.
var Exemplar = grpcprom.WithExemplarFromContext(func(ctx context.Context) prometheus.Labels {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return prometheus.Labels{"trace_id": span.TraceID().String()}
	}
	return nil
})

// Serve starts an HTTP server on port that answers /metrics. It returns
// immediately; the server runs in the background.
func Serve(port int) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		EnableOpenMetrics: true, // exemplars are only sent in the OpenMetrics format
	}))
	go func() {
		log.Printf("serving metrics on port %d ...", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
			log.Fatalf("failed to serve metrics on port %d, error: %v", port, err)
		}
	}()
}
