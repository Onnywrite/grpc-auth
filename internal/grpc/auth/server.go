package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/grpc"
)

type AuthService interface {
	SignUpWithLogin(ctx context.Context, login, appToken string) (token string, err error)
	SignUpWithEmail(ctx context.Context, email, appToken string) (token string, err error)
	SignUpWithPhone(ctx context.Context, phone, appToken string) (token string, err error)
	LogInWithLogin(ctx context.Context, login, password string) (token string, err error)
	LogInWithEmail(ctx context.Context, email, password string) (token string, err error)
	LogInWithPassword(ctx context.Context, phone, password string) (token string, err error)
	Logout(ctx context.Context, token string) error
}

type authServerImpl struct {
	gen.UnimplementedAuthServer
	service AuthService
}

func Register(server *grpc.Server, service AuthService) {
	gen.RegisterAuthServer(server, &authServerImpl{service: service})
}
