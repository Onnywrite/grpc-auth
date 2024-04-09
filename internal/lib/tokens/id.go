package tokens

import (
	"fmt"
	"os"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/golang-jwt/jwt"
)

func Id(token *models.IdToken) (string, error) {
	return New(jwt.MapClaims{
		"id":           token.Id,
		"session_uuid": token.SessionUUID,
		"exp":          token.Exp,
	})
}

func ParseId(tkn string) (*models.IdToken, error) {
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

		id, ok := claims["id"].(float64)
		if !ok {
			return nil, fmt.Errorf("could not convert 'id' to int64")
		}
		sessionUUID, ok := claims["session_uuid"].(string)
		if !ok {
			return nil, ErrInvalidData
		}

		token := &models.IdToken{
			Id:          int64(id),
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
