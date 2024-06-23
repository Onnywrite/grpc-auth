package ero

import (
	"encoding/json"
	"slices"
	"sync"
)

const (
	CurrentService = "SSO"
)

type Error interface {
	error
	Has(errText string) bool
	GetCode() int
}

type Basic struct {
	Service string
	Errors  []string
	Code    int
	mu      *sync.Mutex
}

func New(code int, errs ...string) *Basic {
	return &Basic{
		Service: CurrentService,
		Errors:  errs,
		Code:    code,
		mu:      &sync.Mutex{},
	}
}

func (e Basic) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (e Basic) Unwrap() error {
	return e
}

func (e Basic) GetCode() int {
	return e.Code
}

func (e *Basic) Add(err ...string) *Basic {
	e.lock()
	e.addWithoutLock(err...)
	e.unlock()
	return e
}

func (e *Basic) Has(errText string) bool {
	return slices.Contains(e.Errors, errText)
}

// checks if not nil err has specific error.
// If err is nil, returns false
func Has(err Error, errText string) bool {
	if err != nil {
		return err.Has(errText)
	}
	return false
}

func (e *Basic) lock() {
	e.mu.Lock()
}

func (e *Basic) unlock() {
	e.mu.Unlock()
}

func (e *Basic) addWithoutLock(err ...string) {
	e.Errors = append(e.Errors, err...)
}

// TODO: ero.HasAll, ero.HasAny
