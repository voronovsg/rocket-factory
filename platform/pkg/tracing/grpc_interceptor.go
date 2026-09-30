package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Константы для заголовков трассировки
const (
	TraceIDHeader = "x-trace-id"
)

// UnaryServerInterceptor создает gRPC unary interceptor для трассировки входящих запросов
func UnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	// Возвращаем функцию-interceptor
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		ctx = propagator.Extract(ctx, metadataCarrier(md))

		ctx, span := tracer.Start(
			ctx,
			info.FullMethod, // Используем полный метод gRPC как имя спана
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// Добавляем trace ID в исходящие метаданные для возврата клиенту
		ctx = AddTraceIDToResponse(ctx)
		resp, err := handler(ctx, req)
		if err != nil {
			span.RecordError(err)
		}

		return resp, err
	}
}

// UnaryClientInterceptor создает gRPC unary interceptor для трассировки исходящих запросов
func UnaryClientInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	// Возвращаем функцию-interceptor
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Определяем имя спана в зависимости от наличия контекста трассировки
		spanName := formatSpanName(ctx, method)

		// Создаем новый спан с подготовленным именем
		ctx, span := tracer.Start(
			ctx,
			spanName,
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()

		// Создаем переносчик метаданных для пропагации трассировки
		carrier := metadataCarrier(extractOutgoingMetadata(ctx))

		// Пропагатор добавляет информацию о текущем трейсе в метаданные запроса
		propagator.Inject(ctx, carrier)

		// Обновляем метаданные в контексте
		ctx = metadata.NewOutgoingContext(ctx, metadata.MD(carrier))

		// Вызываем следующий обработчик с обогащенным контекстом
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			trace.SpanFromContext(ctx).RecordError(err)
		}

		return err
	}
}

// formatSpanName формирует имя спана в зависимости от наличия контекста трассировки
func formatSpanName(ctx context.Context, method string) string {
	if !trace.SpanContextFromContext(ctx).IsValid() {
		return "client." + method
	}

	return method
}

// extractOutgoingMetadata извлекает исходящие метаданные из контекста и создает их копию
func extractOutgoingMetadata(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.New(nil)
	}

	return md.Copy()
}

// GetTraceIDFromContext извлекает trace ID из контекста.
func GetTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}

// AddTraceIDToResponse добавляет trace ID в исходящие метаданные gRPC ответа.
func AddTraceIDToResponse(ctx context.Context) context.Context {
	traceID := GetTraceIDFromContext(ctx)
	if traceID == "" {
		return ctx
	}

	md := extractOutgoingMetadata(ctx)

	md.Set(TraceIDHeader, traceID)
	return metadata.NewOutgoingContext(ctx, md)
}
