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

//var (
//invalidCredentials     = "invalid credentials"
//invalidLoginOrPassword = "invalid login or password"
//)

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
	var id models.UserIdentifier
	switch {
	case r.GetLogin() != "":
		id = models.UserIdentifier{Key: "login", Value: r.GetLogin()}
	case r.GetEmail() != "":
		id = models.UserIdentifier{Key: "email", Value: r.GetEmail()}
	case r.GetPhone() != "":
		id = models.UserIdentifier{Key: "phone", Value: r.GetPhone()}
	}
	id.Password = r.Password

	err := a.service.Signup(c, id, r.ServiceId)

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

// TODO: need rollback (delete session, if an error occured while creating tokens)
func (a authServer) Login(c context.Context, r *gen.LoginRequest) (*gen.Tokens, error) {
	if r.Signup == nil {
		return nil, status.Error(codes.InvalidArgument, "signup must not be null")
	}

	var id models.UserIdentifier
	switch {
	case r.Signup.GetLogin() != "":
		id = models.UserIdentifier{Key: "login", Value: r.Signup.GetLogin()}
	case r.Signup.GetEmail() != "":
		id = models.UserIdentifier{Key: "email", Value: r.Signup.GetEmail()}
	case r.Signup.GetPhone() != "":
		id = models.UserIdentifier{Key: "phone", Value: r.Signup.GetPhone()}
	}
	id.Password = r.Signup.Password

	tokens, err := a.service.Login(c, id, models.SessionInfo{
		Browser: r.SessionInfo.Browser,
		Ip:      r.SessionInfo.Ip,
		OS:      r.SessionInfo.Os,
	}, r.Signup.ServiceId)

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
