package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (authServerImpl) SignUp(context.Context, *gen.SignUpRequest) (*gen.Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method is not implemented")
}
