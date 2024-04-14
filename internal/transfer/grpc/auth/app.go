package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a authServer) Signup(c context.Context, r *gen.AppRequest) (*gen.AppTokens, error) {
	info := models.SessionInfo{
		Browser: r.Info.Browser,
		Ip:      r.Info.Ip,
		OS:      r.Info.Os,
	}

	resp, err := a.service.Signup(c, r.IdToken, r.ServiceId, info)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error: %s", err.Error())
	}

	return resp, nil
}

func (authServer) RecoverSignup(context.Context, *gen.AppRequest) (*gen.AppTokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverSignup not implemented")
}

func (a authServer) Signin(c context.Context, r *gen.AppRequest) (*gen.AppTokens, error) {
	info := models.SessionInfo{
		Browser: r.Info.Browser,
		Ip:      r.Info.Ip,
		OS:      r.Info.Os,
	}

	resp, err := a.service.Signin(c, r.IdToken, r.ServiceId, info)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error: %s", err.Error())
	}

	return resp, nil
}

func (a authServer) Resignin(c context.Context, r *gen.RefreshToken) (*gen.AppTokens, error) {
	resp, err := a.service.Resignin(c, r.Token)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error: %s", err.Error())
	}

	return resp, nil
}

func (a authServer) Signout(c context.Context, r *gen.RefreshToken) (*emptypb.Empty, error) {
	err := a.service.Signout(c, r.Token)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error: %s", err.Error())
	}

	return nil, nil
}

func (a authServer) Check(c context.Context, r *gen.AccessToken) (*emptypb.Empty, error) {
	err := a.service.Check(c, r.Token)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error: %s", err.Error())
	}

	return nil, nil
}

func (authServer) ClearTerminatedSessions(context.Context, *gen.AccessToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverSignup not implemented")
}

func (authServer) Unsign(context.Context, *gen.AccessToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverSignup not implemented")
}
