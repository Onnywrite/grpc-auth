package validation

import (
	"context"
	"sync"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/go-playground/validator/v10"
)

type ValidateFn func(*validator.Validate, any) error

func ValidateWith(ctx context.Context, validate ValidateFn, structs ...any) error {
	wg := sync.WaitGroup{}
	wg.Add(len(structs))
	validationErrors := ero.NewValidation()

	v := validator.New()
	for i := range structs {
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				validationErrors.SetCode(ero.CodeValidationTimeout)
				return
			default:
				if err := validate(v, structs[i]); err != nil {
					errs := err.(validator.ValidationErrors)
					for _, e := range errs {
						validationErrors.AddField(e.StructField(), e.Tag(), e.Value())
					}
				}
			}
		}()
	}

	wg.Wait()

	if len(validationErrors.Errors) > 0 || validationErrors.Code != ero.CodeValidation {
		return validationErrors
	}
	return nil
}

func Validate(ctx context.Context, structs ...any) error {
	return ValidateWith(ctx, func(v *validator.Validate, s any) error {
		return v.StructCtx(ctx, s)
	}, structs...)
}
