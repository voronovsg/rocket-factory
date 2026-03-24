package validate

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validatable interface {
	Validate() error
}

func UnaryValidateInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if v, ok := req.(validatable); ok {
			if err := v.Validate(); err != nil {
				log.Printf("validation error: method=%s err=%v", info.FullMethod, err)
				return nil, status.Error(codes.InvalidArgument, "invalid request")
			}
		}
		return handler(ctx, req)
	}
}
