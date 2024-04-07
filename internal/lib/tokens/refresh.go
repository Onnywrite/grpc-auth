package tokens

import (
	"os"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/golang-jwt/jwt"
)

func Refresh(token *models.RefreshToken) (string, error) {
	return New(jwt.MapClaims{
		"session_uuid": token.SessionUUID,
		"exp":          token.Exp,
	})
}

func ParseRefresh(tkn string) (*models.RefreshToken, error) {
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return []byte(os.Getenv(Env)), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := claims["exp"].(float64)

		sessionUUID, ok := claims["session_uuid"].(string)
		if !ok {
			return nil, ErrInvalidData
		}

		token := &models.RefreshToken{
			SessionUUID: sessionUUID,
			Exp:         int64(exp),
		}

		if float64(time.Now().Unix()) > exp {
			return token, ErrTokenExpired
		}
		return token, nil
	}

	return nil, err
}
