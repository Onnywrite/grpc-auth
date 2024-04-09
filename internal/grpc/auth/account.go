package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

//var (
//invalidCredentials     = "invalid credentials"
//invalidLoginOrPassword = "invalid login or password"
//)

func (a authServer) Register(c context.Context, u *gen.InRequest) (*gen.IdTokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}

func (authServer) Recover(context.Context, *gen.InRequest) (*gen.IdTokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Recover not implemented")
}

func (a authServer) Login(c context.Context, r *gen.InRequest) (*gen.IdTokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}

func (authServer) Logout(context.Context, *gen.IdToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
