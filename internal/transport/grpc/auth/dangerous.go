package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (authServer) LogoutEverywhere(context.Context, *gen.DangerousRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogoutEverywhere not implemented")
}

func (authServer) SignoutEverywhere(context.Context, *gen.DangerousRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignoutEverywhere not implemented")
}

func (authServer) Unregister(context.Context, *gen.DangerousRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unregister not implemented")
}
