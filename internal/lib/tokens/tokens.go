package tokens

import (
	"os"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/golang-jwt/jwt"
)

var (
	ErrUnexpectedSigningMethod = "unexpected signing method"
	ErrTokenExpired            = "token has expired"
	ErrInvalidData             = "invalid data provided"
	ErrSigning                 = "cannot sign token"
)

const (
	Env = "TOKEN_SECRET"
)

func New(claims jwt.MapClaims) (string, ero.Error) {
	return NewWithSecret(claims, os.Getenv(Env))
}

func NewWithSecret(claims jwt.MapClaims, secret string) (string, ero.Error) {
	refreshTkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tkn, err := refreshTkn.SignedString([]byte(secret))
	if err != nil {
		return "", ero.NewInternal(err.Error())
	}

	return tkn, nil
}
