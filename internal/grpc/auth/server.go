package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/grpc"
)

type AuthService interface {
	SignUp(ctx context.Context, login, email, password string) (token string, err error)
	SignUpLogin(ctx context.Context, login, password string) (token string, err error)
	SignUpEmail(ctx context.Context, email, password string) (token string, err error)
	LogIn(ctx context.Context, loginOrEmail, password string) (token string, err error)
	Logout(ctx context.Context, token string) error
}

type authServerImpl struct {
	gen.UnimplementedAuthServer
	service AuthService
}

func Register(server *grpc.Server, service AuthService) {
	gen.RegisterAuthServer(server, &authServerImpl{service: service})
}
