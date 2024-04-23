package validation_test

import (
	"context"
	"testing"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/lib/validation"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	ptr := func(s string) *string {
		return &s
	}

	tests := []struct {
		name      string
		structs   []interface{}
		timeout   time.Duration
		expErr    bool
		expFields []ero.FieldError
		expCode   string
	}{
		{
			name: "OK",
			structs: []interface{}{
				models.User{
					Login:    ptr("aboba"),
					Email:    ptr("some@mail.test"),
					Phone:    ptr("+79998887766"),
					Password: "12345678",
				},
			},
			timeout:   time.Second,
			expErr:    false,
			expFields: []ero.FieldError{},
			expCode:   "",
		},
		{
			name: "1 Field",
			structs: []interface{}{
				models.User{
					Login:    ptr("aboba"),
					Phone:    ptr("one"),
					Password: "12345678",
				},
			},
			timeout: time.Second,
			expErr:  true,
			expFields: []ero.FieldError{
				{
					Field:      "Phone",
					Constraint: "e164",
					Value:      "one",
				},
			},
			expCode: ero.CodeValidation,
		},
		{
			name: "2 Fields",
			structs: []interface{}{
				models.User{
					Login:    ptr("aboba"),
					Phone:    ptr("one"),
					Password: "two",
				},
			},
			timeout: time.Second,
			expErr:  true,
			expFields: []ero.FieldError{
				{
					Field:      "Phone",
					Constraint: "e164",
					Value:      "one",
				},
				{
					Field:      "Password",
					Constraint: "gte",
					Value:      "two",
				},
			},
			expCode: ero.CodeValidation,
		},
		{
			name: "Timeout",
			structs: []interface{}{
				models.User{
					Login:    ptr("aboba"),
					Password: "12345678",
				},
			},
			timeout:   0,
			expErr:    true,
			expFields: []ero.FieldError{},
			expCode:   ero.CodeValidationTimeout,
		},
		{
			name: "3 Structs",
			structs: []interface{}{
				models.User{
					Login:    ptr("aboba"),
					Phone:    ptr("+79998887766"),
					Password: "12345678",
				},
				models.Session{
					ServiceId: -1,
					UserId:    0,
					Info: models.SessionInfo{
						Browser: ptr("aaa"),
						Ip:      ptr("0.0.0.0"),
						OS:      ptr("bbb"),
					},
				},
				models.Session{
					ServiceId: 0,
					UserId:    0,
					Info: models.SessionInfo{
						Browser: ptr("aaa "),
						Ip:      ptr("ip address"),
						OS:      ptr("! bbb"),
					},
				},
			},
			timeout: time.Second,
			expErr:  true,
			expFields: []ero.FieldError{
				{
					Field:      "ServiceId",
					Constraint: "gte",
					Value:      int64(-1),
				},
				{
					Field:      "Browser",
					Constraint: "alphanum",
					Value:      "aaa ",
				},
				{
					Field:      "Ip",
					Constraint: "ip",
					Value:      "ip address",
				},
				{
					Field:      "OS",
					Constraint: "alphanum",
					Value:      "! bbb",
				},
			},
			expCode: ero.CodeValidation,
		},
	}

	t.Parallel()

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			ctx, c := context.WithTimeout(context.Background(), tc.timeout)
			defer c()

			err := validation.Validate(ctx, tc.structs...)

			assert.Equal(tt, tc.expErr, err != nil)
			if ve, ok := err.(ero.ValidationError); ok {
				assert.ElementsMatch(tt, tc.expFields, ve.Fields)
				assert.Equal(tt, tc.expCode, ve.Code)
			}
		})
	}
}

func TestValidateWith(t *testing.T) {
	ptr := func(s string) *string {
		return &s
	}

	tests := []struct {
		name      string
		structs   []interface{}
		fn        validation.ValidateFn
		expErr    bool
		expFields []ero.FieldError
		expCode   string
	}{
		{
			name: "OK",
			structs: []interface{}{
				models.User{
					Login:    ptr("aboba"),
					Email:    ptr("some@mail.test"),
					Phone:    ptr("+79998887766"),
					Password: "12345678",
				},
			},
			fn: func(v *validator.Validate, a any) error {
				return v.StructExcept(a, "Ip")
			},
			expErr:    false,
			expFields: []ero.FieldError{},
			expCode:   "",
		},
		{
			name: "1 Except",
			structs: []interface{}{
				models.User{
					Login:    ptr("aboba"),
					Phone:    ptr("+79998887766"),
					Password: "12345678",
				},
				models.Session{
					ServiceId: -1,
					UserId:    0,
					Info: models.SessionInfo{
						Browser: ptr("aaa"),
						Ip:      ptr("0.0.0.0"),
						OS:      ptr("bbb"),
					},
				},
				models.Session{
					ServiceId: 0,
					UserId:    0,
					Info: models.SessionInfo{
						Browser: ptr("aaa "),
						Ip:      ptr("ip address"),
						OS:      ptr("! bbb"),
					},
				},
			},
			fn: func(v *validator.Validate, a any) error {
				return v.StructExcept(a, "Info.Ip")
			},
			expErr: true,
			expFields: []ero.FieldError{
				{
					Field:      "ServiceId",
					Constraint: "gte",
					Value:      int64(-1),
				},
				{
					Field:      "Browser",
					Constraint: "alphanum",
					Value:      "aaa ",
				},
				{
					Field:      "OS",
					Constraint: "alphanum",
					Value:      "! bbb",
				},
			},
			expCode: ero.CodeValidation,
		},
		{
			name: "2 Excepts",
			structs: []interface{}{
				models.User{
					Login:    ptr("aboba"),
					Phone:    ptr("+79998887766"),
					Password: "12345678",
				},
				models.Session{
					ServiceId: -1,
					UserId:    0,
					Info: models.SessionInfo{
						Browser: ptr("aaa"),
						Ip:      ptr("0.0.0.0"),
						OS:      ptr("bbb"),
					},
				},
				models.Session{
					ServiceId: 0,
					UserId:    0,
					Info: models.SessionInfo{
						Browser: ptr("aaa "),
						Ip:      ptr("ip address"),
						OS:      ptr("! bbb"),
					},
				},
			},
			fn: func(v *validator.Validate, a any) error {
				return v.StructExcept(a, "Info.Ip", "ServiceId")
			},
			expErr: true,
			expFields: []ero.FieldError{
				{
					Field:      "Browser",
					Constraint: "alphanum",
					Value:      "aaa ",
				},
				{
					Field:      "OS",
					Constraint: "alphanum",
					Value:      "! bbb",
				},
			},
			expCode: ero.CodeValidation,
		},
	}

	t.Parallel()

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			err := validation.ValidateWith(context.Background(), tc.fn, tc.structs...)

			assert.Equal(tt, tc.expErr, err != nil)
			if ve, ok := err.(ero.ValidationError); ok {
				assert.ElementsMatch(tt, tc.expFields, ve.Fields)
				assert.Equal(tt, tc.expCode, ve.Code)
			}
		})
	}
}
