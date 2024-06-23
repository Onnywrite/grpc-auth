package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/lib/tokens"
	"github.com/Onnywrite/grpc-auth/internal/lib/validation"
	"github.com/Onnywrite/grpc-auth/internal/models"
)

func (*AuthService) validateCredentials(ctx context.Context, log *slog.Logger, creds *models.Credentials) ero.Error {
	erro := validation.Validate(ctx, creds.User, creds.Info)
	if erro != nil {
		log.Error("validation failed", slog.String("error", erro.Error()))
		return erro
	}
	log.Info("passed validation")
	return nil
}

func (s *AuthService) openSession(ctx context.Context, log *slog.Logger, saved *models.SavedUser, info models.SessionInfo) (*models.LoginResponse, ero.Error) {
	session, erro := s.db.SaveSession(ctx, &models.Session{
		ServiceId: 0,
		UserId:    saved.Id,
		Info:      info,
	})
	if erro != nil {
		log.Error("failed to open session", slog.String("error", erro.Error()))
		return nil, erro
	}

	refresh, erro := tokens.Refresh(&models.RefreshToken{
		SessionUUID: session.UUID,
		Rotation:    1,
		Exp:         time.Now().Add(s.refreshConfig.TTL).Unix(),
	}, s.refreshConfig.Secret)
	if erro != nil {
		log.Error("failed to generate refresh token", slog.String("error", erro.Error()))
		s.db.DeleteSession(ctx, session.UUID)
		return nil, erro
	}

	profile := &models.Profile{
		Nickname: saved.Nickname,
		Email:    saved.Email,
		Phone:    saved.Phone,
		Roles:    []string{},
	}

	superAccess, erro := tokens.Access(&models.AccessToken{
		Id:        saved.Id,
		ServiceId: 0,
		Roles:     []string{},
		Exp:       time.Now().Add(s.accessConfig.TTL).Unix(),
	}, s.accessConfig.Secret)
	if erro != nil {
		log.Error("failed to create super access token", slog.String("error", erro.Error()))
		s.db.DeleteSession(ctx, session.UUID)
		return nil, erro
	}

	return &models.LoginResponse{
		AccessToken:  superAccess,
		RefreshToken: refresh,
		Profile:      profile,
	}, nil
}
