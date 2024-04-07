package grpcauth

import (
	"context"
	"errors"

	"github.com/Onnywrite/grpc-auth/gen"
	"github.com/Onnywrite/grpc-auth/internal/lib/ve"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/Onnywrite/grpc-auth/internal/services/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	invalidCredentials     = "invalid credentials"
	invalidLoginOrPassword = "invalid login or password"
)

func (a authServer) Register(c context.Context, u *gen.UserCredentials) (*emptypb.Empty, error) {
	err := a.service.Register(c, u.Login, u.Email, u.Phone, u.Password)
	if ve, ok := err.(ve.ValidationErrorsList); ok {
		return nil, status.Error(codes.InvalidArgument, ve.JSON())
	}
	if errors.Is(err, auth.ErrUserAlreadyRegistered) {
		switch {
		case u.Login != nil:
			return nil, status.Error(codes.AlreadyExists, "login is occupied")
		case u.Email != nil:
			return nil, status.Error(codes.AlreadyExists, "email is occupied")
		case u.Phone != nil:
			return nil, status.Error(codes.AlreadyExists, "phone is occupied")
		}
	}
	if errors.Is(err, auth.ErrUserDeleted) {
		return nil, status.Error(codes.FailedPrecondition, "account has been deleted and can be recovered")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "SSO service internal error")
	}

	return nil, nil
}

func (authServer) Recover(context.Context, *gen.UserCredentials) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Recover not implemented")
}

func (a authServer) Signup(c context.Context, r *gen.SignupRequest) (*emptypb.Empty, error) {
	if r.Credentials == nil {
		return nil, status.Error(codes.InvalidArgument, "credentials are null")
	}

	var err error
	switch {
	case r.Credentials.Login != nil:
		err = a.service.Signup(c, models.UserIdentifier{
			Key:      "login",
			Value:    *r.Credentials.Login,
			Password: r.Credentials.Password,
		}, r.ServiceId)
	case r.Credentials.Email != nil:
		err = a.service.Signup(c, models.UserIdentifier{
			Key:      "email",
			Value:    *r.Credentials.Login,
			Password: r.Credentials.Password,
		}, r.ServiceId)
	case r.Credentials.Phone != nil:
		err = a.service.Signup(c, models.UserIdentifier{
			Key:      "phone",
			Value:    *r.Credentials.Login,
			Password: r.Credentials.Password,
		}, r.ServiceId)
	}

	if errors.Is(err, auth.ErrInvalidCredentials) {
		return nil, status.Error(codes.NotFound, "user or service not found")
	}
	if errors.Is(err, auth.ErrAlreadySignedUp) {
		return nil, status.Error(codes.AlreadyExists, "user has already signed up")
	}
	if errors.Is(err, auth.ErrSignupDeleted) {
		return nil, status.Error(codes.FailedPrecondition, "account in this service has been deleted and can be recovered")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "SSO service internal error")
	}
	return nil, nil
}

func (authServer) RecoverSignup(context.Context, *gen.SignupRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecoverSignup not implemented")
}

func (a authServer) Login(c context.Context, r *gen.LoginRequest) (*gen.Tokens, error) {
	if r.Signup == nil || r.Signup.Credentials == nil {
		return nil, status.Error(codes.InvalidArgument, "signup or credentials are null")
	}
	
	if ve, ok := err.(ve.ValidationErrorsList); ok {
		return nil, status.Error(codes.InvalidArgument, ve.JSON())
	}

	if errors.Is(err, auth.ErrInvalidCredentials) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if errors.Is(err, auth.ErrSignupNotExists) {
		return nil, status.Error(codes.NotFound, "user has never signed up to the service")
	}
	if errors.Is(err, auth.ErrAlreadyLoggedIn) {
		return nil, status.Error(codes.AlreadyExists, "already logged in")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "SSO service internal error")
	}
	return &gen.Tokens{
		Refresh: tokens.Refresh,
		Access:  tokens.Access,
	}, nil
}

func (authServer) Logout(context.Context, *gen.RefreshToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
