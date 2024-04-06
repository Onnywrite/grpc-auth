package tokens

import (
	"fmt"
	"os"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/golang-jwt/jwt"
)

func Access(token *models.AccessToken) (string, error) {
	return New(jwt.MapClaims{
		"id":         token.Id,
		"login":      token.Login,
		"service_id": token.ServiceId,
		"roles":      token.Roles,
		"exp":        token.Exp,
	})
}

func ParseAccess(tkn string) (*models.AccessToken, error) {
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

		if float64(time.Now().Unix()) > exp {
			return nil, ErrTokenExpired
		}

		id, ok := claims["id"].(float64)
		if !ok {
			return nil, fmt.Errorf("could not convert 'id' to int64")
		}
		login, ok := claims["login"].(string)
		if !ok {
			return nil, fmt.Errorf("could not convert 'login' to string")
		}
		serviceId, ok := claims["service_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("could not convert 'service_id' to int64")
		}
		roles, ok := claims["roles"].([]string)
		if !ok {
			return nil, fmt.Errorf("could not convert 'roles' to []string")
		}

		token := &models.AccessToken{
			Id:        int64(id),
			Login:     login,
			ServiceId: int64(serviceId),
			Roles:     roles,
			Exp:       int64(exp),
		}

		return token, nil
	}

	return nil, err
}
