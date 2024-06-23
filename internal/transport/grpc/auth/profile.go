package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"github.com/Onnywrite/grpc-auth/internal/transport/grpc"
)

func (a *authServer) Register(ctx context.Context, creds *gen.Credentials) (*gen.SsoResponse, error) {
	credentials := grpc.ToCredentials(creds)
	resp, err := a.service.Register(ctx, &credentials)
	if err != nil {
		return nil, grpc.ToErrorStatus(err)
	}
	return grpc.ToGrpcSsoResponse(resp), nil
}

func (a *authServer) Login(ctx context.Context, creds *gen.Credentials) (*gen.SsoResponse, error) {
	credentials := grpc.ToCredentials(creds)
	resp, err := a.service.Login(ctx, &credentials)
	if err != nil {
		return nil, grpc.ToErrorStatus(err)
	}
	return grpc.ToGrpcSsoResponse(resp), nil
}
