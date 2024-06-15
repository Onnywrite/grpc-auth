package models

import "time"

type Session struct {
	ServiceId int64 `validate:"gte=0"`
	UserId    int64 `validate:"gte=0"`
	Info      SessionInfo
}

type SavedSession struct {
	UUID      string    `db:"session_uuid"`
	ServiceId int64     `db:"service_fk"`
	UserId    int64     `db:"user_fk"`
	Browser   *string   `db:"browser"`
	IP        *string   `db:"ip"`
	OS        *string   `db:"os"`
	CreatedAt time.Time `db:"at"`
}

type SessionInfo struct {
	Browser *string `validate:"omitempty,alphanum"`
	Ip      *string `validate:"omitempty,ip"`
	OS      *string `validate:"omitempty,alphanum"`
}
