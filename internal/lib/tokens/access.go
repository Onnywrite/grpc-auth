package tokens

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/golang-jwt/jwt"
)

func Access(token *models.AccessToken, secret string) (string, ero.Error) {
	return New(jwt.MapClaims{
		"id":         token.Id,
		"service_id": token.ServiceId,
		"roles":      token.Roles,
		"exp":        token.Exp,
	}, secret)
}

func ParseAccess(tkn, secret string) (*models.AccessToken, ero.Error) {
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

		for k, v := range claims {
			fmt.Println(k, v, reflect.TypeOf(v).Name())
		}

		id, ok := claims["id"].(float64)
		if !ok {
			return nil, ero.NewClient(ero.CodeUnauthorized, ErrInvalidData)
		}
		serviceId, ok := claims["service_id"].(float64)
		if !ok {
			return nil, ero.NewClient(ero.CodeUnauthorized, ErrInvalidData)
		}
		rolesInterface, ok := claims["roles"].([]interface{})
		if !ok {
			return nil, ero.NewClient(ero.CodeUnauthorized, ErrInvalidData)
		}

		roles := make([]string, 0, len(rolesInterface))
		for _, role := range rolesInterface {
			if roleStr, ok := role.(string); !ok {
				return nil, ero.NewClient(ero.CodeUnauthorized, ErrInvalidData)
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
			return token, ero.NewClient(ero.CodeUnauthorized, ErrTokenExpired)
		}
		return token, nil
	}

	return nil, ero.NewInternal(ero.CodeInternal, "could not parse token")
}
