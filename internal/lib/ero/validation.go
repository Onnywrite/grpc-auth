package ero

import "fmt"

type FieldError struct {
	Field      string
	Constraint string
	Value      interface{}
}

type ValidationError struct {
	Basic
	Fields []FieldError
}

func NewValidationWith(code string) ValidationError {
	return ValidationError{
		Basic:  New(code),
		Fields: make([]FieldError, 0, 2),
	}
}
func NewValidation() ValidationError {
	return NewValidationWith(CodeValidation)
}

func (ve *ValidationError) AddField(field, constraint string, value interface{}) *ValidationError {
	ve.Basic.lock()
	ve.addWithoutLock(fmt.Sprintf("%s is invalid", field))
	ve.Fields = append(ve.Fields, FieldError{
		Field:      field,
		Constraint: constraint,
		Value:      value,
	})
	ve.Basic.unlock()
	return ve
}
