package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService interface {
	Register(ctx context.Context, user *models.User, info models.SessionInfo) (*gen.IdTokens, error)
	// Recover(ctx context.Context, user *models.User, info models.SessionInfo) (*gen.IdTokens, error)
	Login(ctx context.Context, user *models.User, info models.SessionInfo) (*gen.IdTokens, error)
	Logout(ctx context.Context, idToken string) error

	Signup(ctx context.Context, idToken string, serviceId int64) (*gen.AppTokens, error)
	RecoverSignup(ctx context.Context, idToken string, serviceId int64) (*gen.AppTokens, error)
	Signin(ctx context.Context, idToken string, serviceId int64) (*gen.AppTokens, error)
	Resignin(ctx context.Context, refresh string) error
	Check(ctx context.Context, access string) error
	Unsign(ctx context.Context, access string) error
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
