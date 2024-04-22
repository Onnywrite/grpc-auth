package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"github.com/Onnywrite/grpc-auth/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authServer struct {
	gen.UnimplementedAuthServer
	service transport.AuthService
}

func Register(server *grpc.Server, service transport.AuthService) {
	gen.RegisterAuthServer(server, &authServer{service: service})
}

func (authServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
