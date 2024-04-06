package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (authServer) LogIn(context.Context, *gen.AuthRequest) (*gen.Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method is not implemented")
}

func (authServer) LogOut(context.Context, *gen.Token) (*gen.NullResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method is not implemented")
}
