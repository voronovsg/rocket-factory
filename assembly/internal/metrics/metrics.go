package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/voronovsg/rocket-factory/assembly/internal/config"
	platformMetrics "github.com/voronovsg/rocket-factory/platform/pkg/metrics"
)

const (
	// serviceName определяет уникальный идентификатор сервиса для метрик
	serviceName = "assembly_service"

	// Префиксы для имён метрик
	metricAssemblyDuration    = serviceName + "_duration_seconds_bucket"
	metricAssemblyOrdersTotal = serviceName + "_orders_total"
	metricAssemblyErrors      = serviceName + "_errors_total"
)

var (
	meter metric.Meter

	// Гистограмма времени сборки заказов
	assemblyDuration metric.Float64Histogram

	// Суммарное количество собранных заказов
	assemblyOrdersTotal metric.Int64Counter

	// Счетчик ошибок сборки
	assemblyErrors metric.Int64Counter
)

func InitMetrics(ctx context.Context, cfg config.MetricServerConfig) error {
	if meter != nil {
		return nil
	}

	err := platformMetrics.InitProvider(ctx, cfg)
	if err != nil {
		return err
	}

	meter = platformMetrics.GetMeterProvider().Meter(serviceName)

	// Инициализируем гистограмму времени сборки
	assemblyDuration, err = meter.Float64Histogram(
		metricAssemblyDuration,
		metric.WithDescription("Время сборки заказов в секундах"),
	)
	if err != nil {
		return err
	}

	// Инициализируем счетчик собранных заказов
	assemblyOrdersTotal, err = meter.Int64Counter(
		metricAssemblyOrdersTotal,
		metric.WithDescription("Общее количество собранных заказов"),
	)
	if err != nil {
		return err
	}

	// Инициализируем счетчик ошибок сборки
	assemblyErrors, err = meter.Int64Counter(
		metricAssemblyErrors,
		metric.WithDescription("Количество ошибок при сборке заказов"),
	)
	if err != nil {
		return err
	}

	return nil
}

// ObserveAssemblyDuration записывает время сборки заказа
func ObserveAssemblyDuration(ctx context.Context, duration time.Duration) {
	if meter == nil {
		return
	}

	assemblyDuration.Record(ctx, duration.Seconds())
}

// CountAssemblyOrder увеличивает счетчик собранных заказов на 1
func CountAssemblyOrder(ctx context.Context) {
	if meter == nil {
		return
	}

	assemblyOrdersTotal.Add(ctx, 1)
}

// CountAssemblyError увеличивает счетчик ошибок сборки
func CountAssemblyError(ctx context.Context, errorType string) {
	if meter == nil {
		return
	}

	assemblyErrors.Add(
		ctx,
		1,
		metric.WithAttributes(
			attribute.String("error_type", errorType),
		),
	)
}

func ShutdownMetrics(ctx context.Context) error {
	return platformMetrics.Shutdown(ctx)
}
