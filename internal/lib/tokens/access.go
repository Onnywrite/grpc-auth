package tokens

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/golang-jwt/jwt"
)

func Access(token *models.AccessToken) (string, ero.Error) {
	return New(jwt.MapClaims{
		"id":         token.Id,
		"service_id": token.ServiceId,
		"roles":      token.Roles,
		"exp":        token.Exp,
	})
}

func ParseAccess(tkn string) (*models.AccessToken, ero.Error) {
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

		for k, v := range claims {
			fmt.Println(k, v, reflect.TypeOf(v).Name())
		}

		id, ok := claims["id"].(float64)
		if !ok {
			return nil, ero.NewClient(ErrInvalidData)
		}
		serviceId, ok := claims["service_id"].(float64)
		if !ok {
			return nil, ero.NewClient(ErrInvalidData)
		}
		rolesInterface, ok := claims["roles"].([]interface{})
		if !ok {
			return nil, ero.NewClient(ErrInvalidData)
		}

		roles := make([]string, 0, len(rolesInterface))
		for _, role := range rolesInterface {
			if roleStr, ok := role.(string); !ok {
				return nil, ero.NewClient(ErrInvalidData)
			} else {
				roles = append(roles, roleStr)
			}
		}

		token := &models.AccessToken{
			Id:        int64(id),
			ServiceId: int64(serviceId),
			Roles:     roles,
			Exp:       int64(exp),
		}

		if time.Now().Unix() > int64(exp) {
			return token, ero.NewClient(ErrTokenExpired)
		}
		return token, nil
	}

	return nil, ero.NewInternal("tokens.ParseRefresh", "could not parse token")
}
