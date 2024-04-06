package tokens

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenExpired            = errors.New("token has expired")
	ErrInvalidData             = errors.New("invalid data provided")
	ErrSigning                 = errors.New("cannot sign token")
)

const (
	Env = "TOKEN_SECRET"
)

func New(claims jwt.MapClaims) (string, error) {
	refreshTkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tkn, err := refreshTkn.SignedString(os.Getenv(Env))
	if err != nil {
		return "", ErrSigning
	}

	return tkn, nil
}
