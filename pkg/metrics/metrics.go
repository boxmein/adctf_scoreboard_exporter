package metrics

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func Setup() (func(), error) {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatalf("failed to initialize metric provider: %v", err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	otel.SetMeterProvider(provider)

	return func() {
		ctx := context.Background()
		err = provider.Shutdown(ctx)
		if err != nil {
			log.Printf("failed to clean up metrics: %v", err)
		}
	}, nil
}
