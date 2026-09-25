package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelLogSdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Key string

const (
	traceIDKey      Key = "trace_id"
	spanIDKey       Key = "span_id"
	userIDKey       Key = "user_id"
	shutdownTimeout     = 2 * time.Second // таймаут для graceful shutdown OTLP provider
)

var (
	globalLogger *logger
	initOnce     sync.Once
	dynamicLevel zap.AtomicLevel
	otelProvider *otelLogSdk.LoggerProvider
)

type logger struct {
	zapLogger *zap.Logger
	//otelLogger *otelzap.Logger
}

func SetNopLogger() {
	globalLogger = &logger{
		zapLogger: zap.NewNop(),
	}
}

// Init инициализирует глобальный логгер с Tee архитектурой
func Init(ctx context.Context, levelStr string, asJSON bool, enableOTLP bool, otlpEndpoint string, serviceName string, serviceEnvironment string) error {
	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLevel(levelStr))
		cores := buildCores(ctx, asJSON, enableOTLP, otlpEndpoint, serviceName, serviceEnvironment)
		zapLogger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(2))
		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})
	if globalLogger == nil {
		return fmt.Errorf("⚠️ logger init failed")
	}

	return nil
}

// buildCores создает слайс cores для zapcore.Tee
func buildCores(ctx context.Context, asJSON, enableOTLP bool, otlpEndpoint, serviceName, serviceEnvironment string) []zapcore.Core {
	cores := []zapcore.Core{
		createStdoutCore(asJSON),
	}
	if enableOTLP {
		if otlpCore := createOTLPCore(ctx, otlpEndpoint, serviceName, serviceEnvironment); otlpCore != nil {
			cores = append(cores, otlpCore)
		}
	}

	return cores
}

// createStdoutCore создает core для записи в stdout/stderr
func createStdoutCore(asJSON bool) zapcore.Core {
	config := buildProductionEncoderConfig()
	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), dynamicLevel)
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// createOTLPCore создает core для отправки в OpenTelemetry
func createOTLPCore(ctx context.Context, otlpEndpoint, serviceName, serviceEnvironment string) zapcore.Core {
	otlp, err := createOTLP(ctx, otlpEndpoint, serviceName, serviceEnvironment)
	if err != nil {
		fmt.Printf("⚠️ OTLP core initialization failed: %v\n", err)
		return nil
	}

	otlpCore, err := zapcore.NewIncreaseLevelCore(
		otelzap.NewCore("app", otelzap.WithLoggerProvider(otlp)),
		dynamicLevel,
	)
	if err != nil {
		fmt.Printf("⚠️ Failed to create increased level core: %v\n", err)
		return nil
	}

	return otlpCore
}

// createOTLP создает OTLP логгер с настроенным экспортером и ресурсами
func createOTLP(ctx context.Context, otlpEndpoint, serviceName, serviceEnvironment string) (*otelLogSdk.LoggerProvider, error) {
	exporter, err := createOTLPExporter(ctx, otlpEndpoint)
	if err != nil {
		return nil, err
	}

	rs, err := createResource(ctx, serviceName, serviceEnvironment)
	if err != nil {
		return nil, err
	}

	otelProvider = otelLogSdk.NewLoggerProvider(
		otelLogSdk.WithResource(rs),
		otelLogSdk.WithProcessor(otelLogSdk.NewBatchProcessor(exporter)),
	)

	return otelProvider, nil
}

// createOTLPExporter создает gRPC экспортер для OTLP коллектора
func createOTLPExporter(ctx context.Context, otlpEndpoint string) (*otlploggrpc.Exporter, error) {
	return otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(otlpEndpoint),
		otlploggrpc.WithInsecure(),
	)
}

// createResource создает метаданные сервиса для телеметрии
func createResource(ctx context.Context, serviceName, serviceEnvironment string) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			attribute.String("deployment.environment", serviceEnvironment),
		),
	)
}

// SetLevel динамически меняет уровень логирования
func SetLevel(levelStr string) {
	if dynamicLevel == (zap.AtomicLevel{}) {
		return
	}

	dynamicLevel.SetLevel(parseLevel(levelStr))
}

func Logger() *logger {
	return globalLogger
}

func With(fields ...zap.Field) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fields...),
	}
}

func WithContext(ctx context.Context) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fieldsFromContext(ctx)...),
	}
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Debug(ctx, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Info(ctx, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Warn(ctx, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Error(ctx, msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Fatal(ctx, msg, fields...)
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Debug(msg, allFields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
	//l.otelLogger.Ctx(ctx).Info(msg, allFields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Warn(msg, allFields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Fatal(msg, allFields...)
}

func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	// Извлекаем trace context из активного span'а
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		fields = append(fields,
			zap.String(string(traceIDKey), spanContext.TraceID().String()),
			zap.String(string(spanIDKey), spanContext.SpanID().String()),
		)
	}

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}

// Close корректно завершает работу логгера.
// Останавливает OTLP provider с таймаутом для отправки оставшихся логов.
func Close() error {
	if otelProvider != nil {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := otelProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("⚠️ failed to shutdown OTLP provider: %w", err)
		}
	}

	return nil
}

// CloseIgnoreErrors закрывает логгер, игнорируя ошибки.
// Используется в defer функциях при завершении приложения, где ошибки закрытия не критичны.
//
//nolint:gosec
func CloseIgnoreErrors() {
	_ = Sync()
	_ = Close()
}

// Sync принудительно сбрасывает все буферизованные логи.
// Вызывает sync для всех cores (stdout + OTLP).
func Sync() error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}

	return nil
}
