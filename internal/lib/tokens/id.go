package tokens

import (
	"fmt"
	"os"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/golang-jwt/jwt"
)

func Id(token models.IdToken) (string, error) {
	return NewWithSecret(jwt.MapClaims{
		"id":  token.Id,
		"iss": token.Iss,
		"sub": token.Sub,
		"exp": token.Exp,
	}, os.Getenv(IdEnv))
}

func ParseId(tkn string) (*models.IdToken, error) {
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return []byte(os.Getenv(IdEnv)), nil
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
		iss, ok := claims["iss"].(string)
		if !ok {
			return nil, fmt.Errorf("could not convert 'iss' to string")
		}
		sub, ok := claims["sub"].(string)
		if !ok {
			return nil, fmt.Errorf("could not convert 'sub' to string")
		}

		token := &models.IdToken{
			Id:  int64(id),
			Iss: iss,
			Sub: sub,
			Exp: int64(exp),
		}

		if float64(time.Now().Unix()) > exp {
			return token, ErrTokenExpired
		}
		return token, nil
	}

	return nil, err
}
