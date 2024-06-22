package validation_test

import (
	"context"
	"testing"
	"time"

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
		name    string
		structs []interface{}
		timeout time.Duration
		expErr  bool
	}{
		{
			name: "OK",
			structs: []interface{}{
				models.User{
					Nickname: "aboba",
					Email:    ptr("some@mail.test"),
					Phone:    ptr("+79998887766"),
					Password: "12345678",
				},
			},
			timeout: time.Second,
			expErr:  false,
		},
		{
			name: "1 Field",
			structs: []interface{}{
				models.User{
					Nickname: "aboba",
					Phone:    ptr("one"),
					Password: "12345678",
				},
			},
			timeout: time.Second,
			expErr:  true,
		},
		{
			name: "2 Fields",
			structs: []interface{}{
				models.User{
					Nickname: "aboba",
					Phone:    ptr("one"),
					Password: "two",
				},
			},
			timeout: time.Second,
			expErr:  true,
		},
		{
			name: "Timeout",
			structs: []interface{}{
				models.User{
					Nickname: "aboba",
					Password: "12345678",
				},
			},
			timeout: 0,
			expErr:  true,
		},
		{
			name: "3 Structs",
			structs: []interface{}{
				models.User{
					Nickname: "aboba",
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
		},
	}

	t.Parallel()

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			ctx, c := context.WithTimeout(context.Background(), tc.timeout)
			defer c()

			err := validation.Validate(ctx, tc.structs...)

			assert.Equal(tt, tc.expErr, err != nil)
		})
	}
}

func TestValidateWith(t *testing.T) {
	ptr := func(s string) *string {
		return &s
	}

	tests := []struct {
		name    string
		structs []interface{}
		fn      validation.ValidateFn
		expErr  bool
	}{
		{
			name: "OK",
			structs: []interface{}{
				models.User{
					Nickname: "aboba",
					Email:    ptr("some@mail.test"),
					Phone:    ptr("+79998887766"),
					Password: "12345678",
				},
			},
			fn: func(v *validator.Validate, a any) error {
				return v.StructExcept(a, "Ip")
			},
			expErr: false,
		},
		{
			name: "1 Except",
			structs: []interface{}{
				models.User{
					Nickname: "aboba",
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
		},
		{
			name: "2 Excepts",
			structs: []interface{}{
				models.User{
					Nickname: "aboba",
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
		},
	}

	t.Parallel()

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			err := validation.ValidateWith(context.Background(), tc.fn, tc.structs...)

			assert.Equal(tt, tc.expErr, err != nil)
		})
	}
}
