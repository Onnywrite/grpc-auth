package models

type AccessToken struct {
	Id        int64
	ServiceId int64
	Roles     []string
	Exp       int64
}

type RefreshToken struct {
	SessionUUID string
	Rotation    int32
	Exp         int64
}
