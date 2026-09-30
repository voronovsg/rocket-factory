package metrics

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/voronovsg/rocket-factory/order/internal/config"
	platformMetrics "github.com/voronovsg/rocket-factory/platform/pkg/metrics"
)

const (
	// serviceName определяет уникальный идентификатор сервиса для метрик
	serviceName = "order_service"

	// Префиксы для имён метрик
	metricOrdersTotalByStatus = serviceName + "_orders_by_status_total"
	metricRevenueTotal        = serviceName + "_revenue_total"
	metricRequestDuration     = serviceName + "_http_request_duration_seconds_bucket"
	metricErrorsTotal         = serviceName + "_errors_total"
)

var (
	meter metric.Meter

	// Счетчик заказов по статусам
	orderCounter metric.Int64Counter

	// Сумма выручки
	totalRevenue metric.Float64Counter

	// Гистограмма времени обработки HTTP-запросов
	requestDuration metric.Float64Histogram

	// Счетчик ошибок
	errorsTotal metric.Int64Counter
)

// InitMetrics инициализирует провайдер OpenTelemetry и метрики
func InitMetrics(ctx context.Context, cfg config.MetricServerConfig) error {
	if meter != nil {
		return nil
	}

	err := platformMetrics.InitProvider(ctx, cfg)
	if err != nil {
		return err
	}

	meter = platformMetrics.GetMeterProvider().Meter(serviceName)

	// Инициализируем счетчик заказов по статусам
	orderCounter, err = meter.Int64Counter(
		metricOrdersTotalByStatus,
		metric.WithDescription("Количество созданных заказов по статусам"),
	)
	if err != nil {
		return err
	}

	// Инициализируем счетчик выручки
	totalRevenue, err = meter.Float64Counter(
		metricRevenueTotal,
		metric.WithDescription("Суммарная выручка от всех заказов"),
	)
	if err != nil {
		return err
	}

	// Инициализируем гистограмму времени обработки HTTP-запросов
	requestDuration, err = meter.Float64Histogram(
		metricRequestDuration,
		metric.WithDescription("Время обработки HTTP запросов"),
	)
	if err != nil {
		return err
	}

	// Инициализируем счетчик ошибок
	errorsTotal, err = meter.Int64Counter(
		metricErrorsTotal,
		metric.WithDescription("Общее количество ошибок"),
	)
	if err != nil {
		return err
	}

	return nil
}

// CountOrderCreated увеличивает счетчик заказов со статусом
func CountOrderCreated(ctx context.Context, status string) {
	if meter == nil {
		return
	}

	orderCounter.Add(
		ctx,
		1,
		metric.WithAttributes(
			attribute.String("status", status),
		),
	)
}

// SumOrderRevenue добавляет значение выручки к общей сумме
func SumOrderRevenue(ctx context.Context, revenue float64) {
	if meter == nil {
		return
	}

	totalRevenue.Add(ctx, revenue)
}

// ObserveHttpDuration записывает метрику длительности HTTP запроса
func ObserveHttpDuration(ctx context.Context, method, path, status string, duration float64) {
	if meter == nil {
		return
	}

	requestDuration.Record(
		ctx,
		duration,
		metric.WithAttributes(
			attribute.String("method", method),
			attribute.String("path", path),
			attribute.String("status", status),
		),
	)
}

// CountError увеличивает счетчик ошибок
func CountError(ctx context.Context, source, errorType string) {
	if meter == nil {
		return
	}

	errorsTotal.Add(
		ctx,
		1,
		metric.WithAttributes(
			attribute.String("source", source),
			attribute.String("type", errorType),
		),
	)
}

func ShutdownMetrics(ctx context.Context) error {
	return platformMetrics.Shutdown(ctx)
}
