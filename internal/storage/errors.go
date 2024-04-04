package storage

import "errors"

var (
	ErrUserExists       = errors.New("user already exists")
	ErrSessionExists    = errors.New("session already exists")
	ErrSignupExists     = errors.New("signup already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrNoSuchPrimaryKey = errors.New("no such primary key")
)
