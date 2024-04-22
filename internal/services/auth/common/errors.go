package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrUserAlreadyRegistered = errors.New("user already exists")
	ErrUserDeleted           = errors.New("user has unregistred")

	ErrAlreadySignedUp = errors.New("you've already signed up and can sign in")
	ErrSignedOut       = errors.New("you've signed out")
	ErrSignupNotExists = errors.New("signup does not exist")
	ErrSignupBanned    = errors.New("you've been banned")

	ErrAlreadyLoggedIn = errors.New("you've already logged in")

	ErrServiceNotExists = errors.New("service does not exist")

	ErrSessionAlreadyOpened = errors.New("session has already been opened")
	ErrSessionNotExists     = errors.New("sesson does not exist")
	ErrSessionTerminated    = errors.New("session has already been terminated")

	ErrInternal     = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrTokenExpired = errors.New("token has expired")
)
