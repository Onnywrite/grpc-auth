package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a authServer) Signup(c context.Context, r *gen.AppRequest) (*gen.AppTokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}

func (authServer) RecoverSignup(context.Context, *gen.AppRequest) (*gen.AppTokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverSignup not implemented")
}

func (authServer) Signin(context.Context, *gen.AppRequest) (*gen.AppTokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverSignup not implemented")
}

func (authServer) Resignin(context.Context, *gen.RefreshToken) (*gen.AppTokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverSignup not implemented")
}

func (authServer) Signout(context.Context, *gen.RefreshToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverSignup not implemented")
}

func (authServer) Check(context.Context, *gen.AccessToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverSignup not implemented")
}

func (authServer) ClearTerminatedSessions(context.Context, *gen.AccessToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverSignup not implemented")
}

func (authServer) Unsign(context.Context, *gen.AccessToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverSignup not implemented")
}
