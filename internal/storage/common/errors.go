package storage_common

import (
	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
)

var (
	ErrUniqueConstraint = ero.NewStorage("object already exists")
	ErrEmptyResult      = ero.NewStorage("got empty result")
	ErrFKConstraint     = ero.NewStorage("no such row in referenced table")
)
