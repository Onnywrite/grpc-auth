package ve

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationErrorsList struct {
	Code   string
	Errors []string
}

func From(ves validator.ValidationErrors) *ValidationErrorsList {
	errors := make([]string, 0, len(ves))

	for _, e := range ves {
		errors = append(errors, fmt.Sprintf("%s is invalid", e.StructField()))
	}

	return &ValidationErrorsList{
		Code:   "000400",
		Errors: errors,
	}
}

func (vel ValidationErrorsList) JSON() string {
	bytes, _ := json.Marshal(&vel)
	return string(bytes)
}

func (vel ValidationErrorsList) Error() string {
	return vel.JSON()
}
