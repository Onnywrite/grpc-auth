package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService interface {
	Register(ctx context.Context, optLogin, optEmail, optPhone *string, password string) error
	Signup(ctx context.Context, key, password string, serviceId int64) error
	Login(ctx context.Context, key, password string, serviceId int64, browser, os, ip string) (*models.Tokens, error)
	Logout(ctx context.Context, refresh string) error
}

type authServer struct {
	gen.UnimplementedAuthServer
	service AuthService
}

func Register(server *grpc.Server, service AuthService) {
	gen.RegisterAuthServer(server, &authServer{service: service})
}

func (authServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
