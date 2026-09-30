package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	DefaultCompressor           = "gzip"
	DefaultRetryEnabled         = true
	DefaultRetryInitialInterval = 500 * time.Millisecond
	DefaultRetryMaxInterval     = 5 * time.Second
	DefaultRetryMaxElapsedTime  = 30 * time.Second
	DefaultTimeout              = 5 * time.Second
)

// serviceName - имя сервиса для трассировки
var serviceName string

type Config interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}

// InitTracer инициализирует глобальный трейсер OpenTelemetry
func InitTracer(ctx context.Context, cfg Config) error {
	serviceName = cfg.ServiceName()

	// Создаем экспортер для отправки трейсов в OpenTelemetry Collector через gRPC
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.CollectorEndpoint()),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(DefaultTimeout),
		otlptracegrpc.WithCompressor(DefaultCompressor),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         DefaultRetryEnabled,
			InitialInterval: DefaultRetryInitialInterval,
			MaxInterval:     DefaultRetryMaxInterval,
			MaxElapsedTime:  DefaultRetryMaxElapsedTime,
		}),
	)
	if err != nil {
		return err
	}

	// Создаем ресурс с метаданными сервиса
	attributeResource, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName()),
			semconv.ServiceVersion(cfg.ServiceVersion()),
			attribute.String("environment", cfg.Environment()),
		),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return err
	}

	// Создаем провайдер трейсов с настроенным экспортером и ресурсом
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(attributeResource),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
	)

	// Устанавливаем глобальный провайдер трейсов
	otel.SetTracerProvider(tracerProvider)

	// Настраиваем пропагацию контекста для передачи между сервисами
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return nil
}

// ShutdownTracer закрывает глобальный трейсер OpenTelemetry
func ShutdownTracer(ctx context.Context) error {
	provider := otel.GetTracerProvider()
	if provider == nil {
		return nil
	}

	// Приводим к конкретному типу для вызова Shutdown
	tracerProvider, ok := provider.(*sdktrace.TracerProvider)
	if !ok {
		return nil
	}

	err := tracerProvider.Shutdown(ctx)
	if err != nil {
		// Ошибки при закрытии не критичны, но могут привести к потере последних трейсов
		return err
	}

	return nil
}

// StartSpan создает новый спан и возвращает его вместе с новым контекстом
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return otel.Tracer(serviceName).Start(ctx, name, opts...)
}

// SpanFromContext возвращает текущий активный спан из контекста
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// TraceIDFromContext извлекает trace ID из контекста
func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}
