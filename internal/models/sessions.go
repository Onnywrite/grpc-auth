package models

import "time"

type Session struct {
	UserId    int64 `validate:"gte=0"`
	ServiceId int64 `validate:"gte=0"`
	Info      SessionInfo
}

type SavedSession struct {
	UUID         string     `db:"session_uuid"`
	SignupId     int64      `db:"signup_fk"`
	Browser      *string    `db:"browser"`
	IP           *string    `db:"ip"`
	OS           *string    `db:"os"`
	CreatedAt    time.Time  `db:"at"`
	TerminatedAt *time.Time `db:"terminated_at"`
}

func (s *SavedSession) IsTerminated() bool {
	return s.TerminatedAt != nil
}

type SessionInfo struct {
	Browser *string `validate:"alphanumunicode"`
	Ip      *string `validate:"ip"`
	OS      *string `validate:"ascii"`
}
