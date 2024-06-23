package grpc

import (
	"github.com/Onnywrite/grpc-auth/gen"
	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToSessionInfo(info *gen.SessionInfo) models.SessionInfo {
	return models.SessionInfo{
		Browser: info.Browser,
		OS:      info.Os,
		Ip:      info.Ip,
	}
}

func ToCredentials(creds *gen.Credentials) models.Credentials {
	return models.Credentials{
		User: models.User{
			Nickname: creds.Nickname,
			Email:    creds.Email,
			Phone:    creds.Phone,
			Password: creds.Password,
		},
		Info: ToSessionInfo(creds.Info),
	}
}

func ToGrpcProfile(profile *models.Profile) *gen.Profile {
	return &gen.Profile{
		Id:       profile.Id,
		Nickname: profile.Nickname,
		Email:    profile.Email,
		Phone:    profile.Phone,
	}
}

func ToGrpcSsoResponse(resp *models.LoginResponse) *gen.SsoResponse {
	return &gen.SsoResponse{
		Access: &gen.SuperAccessToken{
			Token: resp.AccessToken,
		},
		Refresh: &gen.RefreshToken{
			Token: resp.RefreshToken,
		},
		Profile: ToGrpcProfile(resp.Profile),
	}
}

func ToErrorStatus(err ero.Error) error {
	return status.Error(codes.Code(ero.ToGrpcCode(err.GetCode())), err.Error())
}
