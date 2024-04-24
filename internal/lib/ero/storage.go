package ero

type StorageError struct {
	*Basic
	Method *string
}

func NewStorage(errs ...string) *StorageError {
	return &StorageError{
		Basic: New(CodeStorage, errs...),
	}
}

func (e *StorageError) WithMethod(method string) *StorageError {
	e.lock()
	e.Method = &method
	e.unlock()
	return e
}
