package observability

import (
	"context"
	"errors"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// NewTraceProvider creates an OTLP gRPC trace provider configured entirely
// via standard OTel environment variables. Returns (nil, nil) when OTEL_EXPORTER_OTLP_ENDPOINT is unset
func NewTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		return nil, nil
	}

	exp, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithAttributes(
			attribute.String(string(semconv.ServiceNameKey), "grinex-rates"),
		),
	)
	if err != nil && !errors.Is(err, resource.ErrPartialResource) && !errors.Is(err, resource.ErrSchemaURLConflict) {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	return tp, nil
}
