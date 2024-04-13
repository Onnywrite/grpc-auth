package models

import "time"

type Signup struct {
	UserId    int64 `db:"user_fk" validate:"gte=0"`
	ServiceId int64 `db:"service_fk" validate:"gte=0"`
}

type SavedSignup struct {
	Id        int64      `db:"signup_id"`
	UserId    int64      `db:"user_fk"`
	ServiceId int64      `db:"service_fk"`
	CreatedAt time.Time  `db:"at"`
	BannedAt  *time.Time `db:"banned_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (su *SavedSignup) IsDeleted() bool {
	return su.DeletedAt != nil
}

func (su *SavedSignup) IsBanned() bool {
	return su.BannedAt != nil
}
