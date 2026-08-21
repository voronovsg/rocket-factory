package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	authV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/auth/v1"
	commonV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/common/v1"
)

type stubAuthServer struct {
	authV1.UnimplementedAuthServiceServer
}

func (s *stubAuthServer) Whoami(_ context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	return &authV1.WhoamiResponse{
		Session: &commonV1.Session{Uuid: req.GetSessionUuid()},
		User:    &commonV1.User{Uuid: "00000000-0000-0000-0000-000000000001"},
	}, nil
}

func main() {
	host := os.Getenv("GRPC_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50053"
	}

	listener, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		panic(fmt.Errorf("failed to listen: %w", err))
	}

	server := grpc.NewServer()
	authV1.RegisterAuthServiceServer(server, &stubAuthServer{})

	fmt.Printf("stub IAM listening on %s\n", listener.Addr())

	if err = server.Serve(listener); err != nil {
		panic(fmt.Errorf("failed to serve: %w", err))
	}
}
