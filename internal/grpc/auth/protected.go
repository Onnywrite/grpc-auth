package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (authServer) GetSessions(context.Context, *gen.GetSessionsRequest) (*gen.Sessions, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSessions not implemented")
}

func (authServer) GetProfile(context.Context, *gen.AccessToken) (*gen.UserProfile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfile not implemented")
}

func (authServer) EditProfile(context.Context, *gen.EditProfileRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditProfile not implemented")
}