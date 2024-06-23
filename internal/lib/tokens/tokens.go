package tokens

import (
	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/golang-jwt/jwt"
)

var (
	ErrUnexpectedSigningMethod = "unexpected signing method"
	ErrTokenExpired            = "token has expired"
	ErrInvalidData             = "invalid data provided"
	ErrSigning                 = "cannot sign token"
)

func New(claims jwt.MapClaims, secret string) (string, ero.Error) {
	refreshTkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tkn, err := refreshTkn.SignedString([]byte(secret))
	if err != nil {
		return "", ero.NewInternal(ero.CodeInternal, ErrSigning, err.Error())
	}

	return tkn, nil
}
