package logger

import (
	"context"
	"fmt"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/voronovsg/rocket-factory/platform/pkg/logger"
)

// UnaryLoggerInterceptor создает серверный унарный интерцептор,
// который логирует информацию о времени выполнения методов gRPC сервера.
func UnaryLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		method := path.Base(info.FullMethod)
		logger.Info(ctx, fmt.Sprintf("🚀 Started gRPC method %s\n", info.FullMethod))
		startTime := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(startTime)
		if err != nil {
			st, _ := status.FromError(err)
			logger.Error(ctx,
				fmt.Sprintf("❌ Finished gRPC method %s with code %s: %v (took: %v)\n", method, st.Code(), err, duration))
		} else {
			logger.Info(ctx,
				fmt.Sprintf("✅ Finished gRPC method %s successfully (took: %v)\n", method, duration))
		}

		return resp, err
	}
}
