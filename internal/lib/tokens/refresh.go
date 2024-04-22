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
		// TODO: rotation
		"rotation": 0,
		"exp":      token.Exp,
	})
}

func ParseRefresh(tkn string) (*models.RefreshToken, error) {
	parser := jwt.Parser{SkipClaimsValidation: true}
	token, err := parser.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return []byte(os.Getenv(Env)), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, ErrInvalidData
		}

		sessionUUID, ok := claims["session_uuid"].(string)
		if !ok {
			return nil, ErrInvalidData
		}
		rotation, ok := claims["rotation"].(float64)
		if !ok {
			return nil, ErrInvalidData
		}

		token := &models.RefreshToken{
			SessionUUID: sessionUUID,
			Rotation:    int32(rotation),
			Exp:         int64(exp),
		}

		if float64(time.Now().Unix()) > exp {
			return token, ErrTokenExpired
		}
		return token, nil
	}

	return nil, err
}
