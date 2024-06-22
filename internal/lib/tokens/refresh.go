package tokens

import (
	"os"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/golang-jwt/jwt"
)

func Refresh(token *models.RefreshToken) (string, ero.Error) {
	return New(jwt.MapClaims{
		"session_uuid": token.SessionUUID,
		// TODO: rotation
		"rotation": 0,
		"exp":      token.Exp,
	})
}

func ParseRefresh(tkn string) (*models.RefreshToken, ero.Error) {
	parser := jwt.Parser{SkipClaimsValidation: true}
	token, err := parser.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ero.NewServer(ErrUnexpectedSigningMethod)
		}

		return []byte(os.Getenv(Env)), nil
	})

	if err != nil {
		return nil, ero.NewInternal(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, ero.NewClient(ErrInvalidData)
		}

		sessionUUID, ok := claims["session_uuid"].(string)
		if !ok {
			return nil, ero.NewClient(ErrInvalidData)
		}
		rotation, ok := claims["rotation"].(float64)
		if !ok {
			return nil, ero.NewClient(ErrInvalidData)
		}

		token := &models.RefreshToken{
			SessionUUID: sessionUUID,
			Rotation:    int32(rotation),
			Exp:         int64(exp),
		}

		if float64(time.Now().Unix()) > exp {
			return token, ero.NewClient(ErrTokenExpired)
		}
		return token, nil
	}

	return nil, ero.NewInternal("tokens.ParseRefresh", "could not parse token")
}
