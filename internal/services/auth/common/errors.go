package auth

var (
	ErrInvalidCredentials = "invalid credentials"

	ErrUserAlreadyRegistered = "user already exists"
	ErrUserDeleted           = "user has unregistred"

	ErrAlreadySignedUp = "already signed up"
	ErrSignedOut       = "signed out"
	ErrSignupNotExists = "signup does not exist"
	ErrSignupBanned    = "user is banned"

	ErrAlreadyLoggedIn = "already logged in"

	ErrServiceNotExists = "service does not exist"

	ErrSessionAlreadyOpened = "session has already been opened"
	ErrSessionNotExists     = "sesson does not exist"

	ErrUnauthorized = "unauthorized"
	ErrTokenExpired = "token has expired"
)
