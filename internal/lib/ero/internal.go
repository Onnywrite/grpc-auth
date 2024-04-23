package ero

import (
	"time"

	"github.com/google/uuid"
)

type Internal struct {
	*Basic
	Time          *time.Time
	UniqueErrorId *uuid.UUID
}

func NewInternal(errs ...string) *Internal {
	return &Internal{
		Basic: New(CodeInternal, errs...),
	}
}

func (i *Internal) WithTime() *Internal {
	t := time.Now()
	i.lock()
	i.Time = &t
	i.unlock()
	return i
}

func (i *Internal) WithUUID(id uuid.UUID) *Internal {
	i.lock()
	i.UniqueErrorId = &id
	i.unlock()
	return i
}
