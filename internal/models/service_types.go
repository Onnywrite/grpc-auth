package models

import "time"

type Credentials struct {
	User
	Info SessionInfo
}

type AppCredentials struct {
	SuperAccessToken string
	ServiceId        int64
	Info             *SessionInfo
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	Profile      *Profile
}

type App struct {
	Service      *SavedService
	RegisteredAt *time.Time
	Sessions     []SessionInfo
}
