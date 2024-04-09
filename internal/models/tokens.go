package models

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
	SessionUUID string
	Exp         int64
}

type IdToken struct {
	Id          int64
	SessionUUID string
	Exp         int64
}
