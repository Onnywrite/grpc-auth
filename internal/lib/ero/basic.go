package ero

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

const (
	CurrentService = "SSO"
)

type Basic struct {
	Service string
	Code    string
	Errors  []string
	mu      *sync.Mutex
}

func New(code string, errs ...string) *Basic {
	return &Basic{
		Service: CurrentService,
		Code:    code,
		Errors:  errs,
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

func (e *Basic) lock() {
	e.mu.Lock()
}

func (e *Basic) unlock() {
	e.mu.Unlock()
}

func (e *Basic) Add(err ...string) *Basic {
	e.lock()
	e.addWithoutLock(err...)
	e.unlock()
	return e
}
func (e *Basic) addWithoutLock(err ...string) {
	e.Errors = append(e.Errors, err...)
}

func (e *Basic) SetCode(code string) *Basic {
	e.lock()
	e.Code = code
	e.unlock()
	return e
}

func Unmarshal(body []byte, errTypes map[string]interface{}) error {
	var err Basic
	json.Unmarshal(body, &err)

	var tp interface{}
	var ok bool
	if tp, ok = errTypes[err.Code]; !ok {
		return NewInternal(fmt.Sprintf("could not find type for '%s' error", err.Code))
	}

	initErr := reflect.New(reflect.TypeOf(tp))
	json.Unmarshal(body, initErr.Interface())

	return initErr.Elem().Interface().(error)
}
