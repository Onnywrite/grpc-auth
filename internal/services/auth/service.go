package auth

import (
	"log/slog"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/services/auth/wrap"
)

type AuthService struct {
	log                                            *slog.Logger
	db                                             *wrap.Wrapper
	tokenTTL, refreshTokenTTL, superAccessTokenTTL time.Duration
}

func New(logger *slog.Logger, db wrap.Storage, tokenTTL, refreshTokenTTL, superAccessTokenTTL time.Duration) *AuthService {
	return &AuthService{
		log:                 logger,
		db:                  wrap.New(logger, db),
		tokenTTL:            tokenTTL,
		superAccessTokenTTL: superAccessTokenTTL,
		refreshTokenTTL:     refreshTokenTTL,
	}
}
