package storage_common

import "errors"

var (
	ErrUniqueConstraint = errors.New("object already exists")
	ErrEmptyResult      = errors.New("got empty result")
	ErrFKConstraint     = errors.New("no row in referenced table")
)
