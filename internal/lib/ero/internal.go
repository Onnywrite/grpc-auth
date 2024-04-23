package ero

type Internal struct {
	Basic
}

func NewInternal(errs ...string) Internal {
	return Internal{
		Basic: New(CodeInternal, errs...),
	}
}
