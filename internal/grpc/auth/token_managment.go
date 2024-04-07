package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (authServer) Relogin(context.Context, *gen.RefreshToken) (*gen.Tokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Relogin not implemented")
}

func (authServer) Check(context.Context, *gen.AccessToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
