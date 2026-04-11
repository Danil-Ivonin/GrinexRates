package observability

import (
	"context"
	"errors"

	prometheusexporter "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// NewMeterProvider creates a Prometheus-backed OTel MeterProvider
func NewMeterProvider(ctx context.Context) (*sdkmetric.MeterProvider, error) {
	exp, err := prometheusexporter.New()
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
	)
	if err != nil && !errors.Is(err, resource.ErrPartialResource) && !errors.Is(err, resource.ErrSchemaURLConflict) {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exp),
		sdkmetric.WithResource(res),
	)
	return mp, nil
}
