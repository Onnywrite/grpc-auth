package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (authServer) Register(context.Context, *gen.UserCredentials) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}

func (authServer) Recover(context.Context, *gen.UserCredentials) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Recover not implemented")
}

func (authServer) Signup(context.Context, *gen.SignupRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Signup not implemented")
}

func (authServer) RecoverSignup(context.Context, *gen.SignupRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverSignup not implemented")
}

func (authServer) Login(context.Context, *gen.LoginRequest) (*gen.Tokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}

func (authServer) Logout(context.Context, *gen.RefreshToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
