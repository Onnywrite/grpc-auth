package transport

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService interface {
	Register(context.Context, *gen.Credentials) (*emptypb.Empty, error)
	Login(context.Context, *gen.Credentials) (*gen.SsoResponse, error)
	Logout(context.Context, *gen.RefreshToken) (*emptypb.Empty, error)
	Relogin(context.Context, *gen.RefreshToken) (*gen.SsoResponse, error)
	Check(context.Context, *gen.SuperAccessToken) (*emptypb.Empty, error)
	SetProfile(context.Context, *gen.ProfileChangeRequest) (*emptypb.Empty, error)
	GetProfile(context.Context, *gen.SuperAccessToken) (*gen.Profile, error)
	GetApps(context.Context, *gen.SuperAccessToken) (*gen.Apps, error)
	SetPassword(context.Context, *gen.PasswordChangeRequest) (*emptypb.Empty, error)
	Delete(context.Context, *gen.DangerousRequest) (*emptypb.Empty, error)
	Recover(context.Context, *gen.Credentials) (*emptypb.Empty, error)
}

type AppService interface {
	Login(context.Context, *gen.AppRequest) (*gen.AppResponse, error)
	Logout(context.Context, *gen.RefreshToken) (*emptypb.Empty, error)
	Relogin(context.Context, *gen.RefreshToken) (*gen.AppResponse, error)
	Check(context.Context, *gen.AccessToken) (*emptypb.Empty, error)
	SetProfile(context.Context, *gen.ProfileChangeRequest) (*emptypb.Empty, error)
	GetProfile(context.Context, *gen.AccessToken) (*gen.Profile, error)
	GetSessions(context.Context, *gen.AccessToken) (*gen.Sessions, error)
	Delete(context.Context, *gen.AccessToken) (*emptypb.Empty, error)
	Recover(context.Context, *gen.AppRequest) (*emptypb.Empty, error)
}
