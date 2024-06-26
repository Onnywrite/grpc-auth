package grpcauth

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//var (
//invalidCredentials     = "invalid credentials"
//invalidLoginOrPassword = "invalid login or password"
//)

func (a authServer) Register(c context.Context, r *gen.InRequest) (*gen.IdTokens, error) {
	user := &models.User{
		Login:    r.Credentials.Login,
		Email:    r.Credentials.Email,
		Phone:    r.Credentials.Phone,
		Password: r.Credentials.Password,
	}

	var info models.SessionInfo
	if r.Info != nil {
		info = models.SessionInfo{
			Browser: r.Info.Browser,
			Ip:      r.Info.Ip,
			OS:      r.Info.Os,
		}
	}

	resp, err := a.service.Register(c, user, info)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}

func (authServer) Recover(context.Context, *gen.InRequest) (*gen.IdTokens, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Recover not implemented")
}

func (a authServer) Login(c context.Context, r *gen.InRequest) (*gen.IdTokens, error) {
	user := &models.User{
		Login:    r.Credentials.Login,
		Email:    r.Credentials.Email,
		Phone:    r.Credentials.Phone,
		Password: r.Credentials.Password,
	}

	var info models.SessionInfo
	if r.Info != nil {
		info = models.SessionInfo{
			Browser: r.Info.Browser,
			Ip:      r.Info.Ip,
			OS:      r.Info.Os,
		}
	}

	resp, err := a.service.Login(c, user, info)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}
