package transfer

import (
	"context"

	"github.com/Onnywrite/grpc-auth/gen"
	"github.com/Onnywrite/grpc-auth/internal/models"
)

type AuthService interface {
	Register(ctx context.Context, user *models.User, info models.SessionInfo) (*gen.IdTokens, error)
	// Recover(ctx context.Context, user *models.User, info models.SessionInfo) (*gen.IdTokens, error)
	Login(ctx context.Context, user *models.User, info models.SessionInfo) (*gen.IdTokens, error)
	// Logout(ctx context.Context, idToken string) error

	Signup(ctx context.Context, idToken string, serviceId int64, info models.SessionInfo) (*gen.AppTokens, error)
	// RecoverSignup(ctx context.Context, idToken string, serviceId int64, info models.SessionInfo) (*gen.AppTokens, error)
	Signin(ctx context.Context, idToken string, serviceId int64, info models.SessionInfo) (*gen.AppTokens, error)
	Resignin(ctx context.Context, refresh string) error
	Check(ctx context.Context, access string) error
	// Unsign(ctx context.Context, access string) error
}
