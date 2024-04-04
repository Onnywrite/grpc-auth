package models

import "time"

type Signup struct {
	UserId    int64 `db:"user_fk"`
	ServiceId int64 `db:"service_fk"`
}

type SavedSignup struct {
	Id        int64      `db:"signup_id"`
	UserId    int64      `db:"user_fk"`
	ServiceId int64      `db:"service_fk"`
	CreatedAt time.Time  `db:"at"`
	BannedAt  *time.Time `db:"banned_at"`
}
