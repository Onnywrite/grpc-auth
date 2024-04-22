package se

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

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

func From(ves validator.ValidationErrors) Errors {
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

func Error(code Code, errList ...string) Errors {
	return ErrorSlice(code, errList)
}

func ErrorSlice(code Code, errors []string) Errors {
	return Errors{
		Code:   code,
		Errors: errors,
	}
}

type ValidateFn func(*validator.Validate, any) error

func ValidateWith(ctx context.Context, validate ValidateFn, structs ...any) *Errors {
	wg := sync.WaitGroup{}
	wg.Add(2)
	errors := make(chan Errors)

	go func() {
		defer wg.Done()
		defer close(errors)
		v := validator.New()
		for i := range structs {
			select {
			case <-ctx.Done():
				return
			default:
				if err := validate(v, structs[i]); err != nil {
					errs := From(err.(validator.ValidationErrors))
					errors <- errs
				}
			}
		}
	}()

	out := &Errors{Code: CodeValidation}

	go func() {
		defer wg.Done()
		for err := range errors {
			select {
			case <-ctx.Done():
				return
			default:
				out.Errors = append(out.Errors, err.Errors...)
			}
		}
	}()

	wg.Wait()

	if len(out.Errors) > 0 {
		return out
	}
	return nil
}

func Validate(ctx context.Context, structs ...any) *Errors {
	return ValidateWith(ctx, func(v *validator.Validate, s any) error {
		return v.StructCtx(ctx, s)
	}, structs...)
}
