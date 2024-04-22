package grpcapp

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"github.com/Onnywrite/grpc-auth/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type appServer struct {
	gen.UnimplementedAppServer
	service transport.AppService
}

func Register(server *grpc.Server, service transport.AppService) {
	gen.RegisterAppServer(server, &appServer{service: service})
}

func (appServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
