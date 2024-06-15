package validation

import (
	"context"
	"fmt"
	"sync"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/go-playground/validator/v10"
)

const (
	ErrContextDone = "context done"
)

type ValidateFn func(*validator.Validate, any) error

func ValidateWith(ctx context.Context, validate ValidateFn, structs ...any) ero.Error {
	const op = "validation.ValidateWith"

	wg := sync.WaitGroup{}
	wg.Add(len(structs))
	errors := ero.New()

	v := validator.New()
	v.RegisterValidation("nickname", validateNickname)

	for i := range structs {
		go func(ii int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				errors.Add(ErrContextDone)
				// TODO: add logging
				return
			default:
				if err := validate(v, structs[ii]); err != nil {
					errs := err.(validator.ValidationErrors)
					for _, e := range errs {
						errors.Add(fmt.Sprintf("field %s has invalid value %s", e.StructField(), fmt.Sprint(e.Value())))
					}
				}
			}
		}(i)
	}

	wg.Wait()
	switch {
	case errors.Has(ErrContextDone):
		return ero.ServerFrom(op, errors)
	case len(errors.Errors) > 0:
		return ero.ClientFrom(errors)
	default:
		return nil
	}
}

func Validate(ctx context.Context, structs ...any) ero.Error {
	return ValidateWith(ctx, func(v *validator.Validate, s any) error {
		return v.StructCtx(ctx, s)
	}, structs...)
}
