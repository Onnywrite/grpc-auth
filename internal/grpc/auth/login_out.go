package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (authServerImpl) LogInWithLogin(context.Context, *gen.LogInLoginRequest) (*gen.Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method is not implemented")
}

func (authServerImpl) LogInWithEmail(context.Context, *gen.LogInEmailRequest) (*gen.Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method is not implemented")
}

func (authServerImpl) LogOut(context.Context, *gen.Token) (*gen.NullResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method is not implemented")
}
