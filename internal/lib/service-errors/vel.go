package se

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Code string

var (
	CodeValidation             = Code("000400")
	CodeInternal               = Code("000500")
	CodeInternalRecoveredPanic = Code("000555")
)

type Errors struct {
	Code   Code
	Errors []string
}

var (
	ErrPanicRecoveredGrpc = Error(CodeInternalRecoveredPanic, "internal error")
)

func From(ves validator.ValidationErrors) error {
	errors := make([]string, 0, len(ves))

	for _, e := range ves {
		errors = append(errors, fmt.Sprintf("%s is invalid", e.StructField()))
	}

	return ErrorSlice(CodeValidation, errors)
}

func (e Errors) JSON() string {
	bytes, _ := json.Marshal(&e)
	return string(bytes)
}

func (e Errors) Error() string {
	return e.JSON()
}

func Error(code Code, errList ...string) error {
	return ErrorSlice(code, errList)
}

func ErrorSlice(code Code, errors []string) error {
	return error(Errors{
		Code:   code,
		Errors: errors,
	})
}
