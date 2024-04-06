package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/grpc"
)

type AuthService interface {
	SignUp(ctx context.Context, optLogin, optEmail, optPhone, appToken, password string) (token string, err error)
	LogIn(ctx context.Context, optLogin, optEmail, optPhone, appToken, password string) (token string, err error)
	LogOut(ctx context.Context, token string) error
}

type authServer struct {
	gen.UnimplementedAuthServer
	service AuthService
}

func Register(server *grpc.Server, service AuthService) {
	gen.RegisterAuthServer(server, &authServer{service: service})
}
