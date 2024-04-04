package models

import "github.com/google/uuid"

type Tokens struct {
	Refresh, Access string
}

type AccessToken struct {
	Id        int64
	Login     string
	ServiceId int64
	Roles     []string
	Exp       int64
}

type RefreshToken struct {
	SessionUUID uuid.UUID
	Exp         int64
}
