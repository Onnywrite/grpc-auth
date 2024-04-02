package models

type Signup struct {
	UserId    int64 `db:"user_fk"`
	ServiceId int64 `db:"service_fk"`
}

type SavedSignup struct {
	
}