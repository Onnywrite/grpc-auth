package storage

import "errors"

var (
	ErrUserExists       = errors.New("user already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrSignupExists     = errors.New("signup already exists")
	ErrNoSuchPrimaryKey = errors.New("no such primary key")
)
