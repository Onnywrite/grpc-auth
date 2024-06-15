package ero

type Internal struct {
	*Basic
	Function string
}

func NewInternal(function string, errors ...string) *Internal {
	return &Internal{
		Basic:    New(errors...),
		Function: function,
	}
}

func InternalFrom(function string, basic *Basic) *Internal {
	return &Internal{
		Basic:    basic,
		Function: function,
	}
}
