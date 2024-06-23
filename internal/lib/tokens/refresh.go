package tokens

import (
	"time"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/golang-jwt/jwt"
)

func Refresh(token *models.RefreshToken, secret string) (string, ero.Error) {
	return New(jwt.MapClaims{
		"session_uuid": token.SessionUUID,
		// TODO: rotation
		"rotation": 0,
		"exp":      token.Exp,
	}, secret)
}

func ParseRefresh(tkn, secret string) (*models.RefreshToken, ero.Error) {
	parser := jwt.Parser{SkipClaimsValidation: true}
	token, err := parser.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ero.NewServer(ero.CodeUnauthorized, ErrUnexpectedSigningMethod)
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, ero.NewInternal(ero.CodeInternal, err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, ero.NewClient(ero.CodeUnauthorized, ErrInvalidData)
		}

		sessionUUID, ok := claims["session_uuid"].(string)
		if !ok {
			return nil, ero.NewClient(ero.CodeUnauthorized, ErrInvalidData)
		}
		rotation, ok := claims["rotation"].(float64)
		if !ok {
			return nil, ero.NewClient(ero.CodeUnauthorized, ErrInvalidData)
		}

		token := &models.RefreshToken{
			SessionUUID: sessionUUID,
			Rotation:    int32(rotation),
			Exp:         int64(exp),
		}

		if float64(time.Now().Unix()) > exp {
			return token, ero.NewClient(ero.CodeUnauthorized, ErrTokenExpired)
		}
		return token, nil
	}

	return nil, ero.NewInternal(ero.CodeInternal, "could not parse token")
}
