package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	nicknameRegex = regexp.MustCompile(`^[a-zA-Z0-9_а-яА-Я\-]+$`)
)

func validateNickname(fl validator.FieldLevel) bool {
	return nicknameRegex.MatchString(fl.Field().String())
}
