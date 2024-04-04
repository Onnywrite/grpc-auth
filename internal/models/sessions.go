package models

import (
	"net/netip"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SignupId int64
	IP       netip.Addr
	Browser  *string
	OS       *string
}

type SavedSession struct {
	UUID      uuid.UUID
	SignupId  int64 `db:"signup_fk"`
	IP        netip.Addr
	Browser   string
	OS        string
	CreatedAt time.Time `db:"at"`
	ClosedAt  time.Time
}