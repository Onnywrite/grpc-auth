package ero

type Internal struct {
	*Basic
	CallStack []string
}

func NewInternal(code int, function string, errors ...string) *Internal {
	return &Internal{
		Basic:     New(code, errors...),
		CallStack: []string{function},
	}
}

func InternalFrom(function string, err Error) *Internal {
	switch terr := err.(type) {
	case *Internal:
		terr.CallStack = append(terr.CallStack, function)
		return terr
	case *Server:
		terr.CallStack = append(terr.CallStack, function)
		return &Internal{
			Basic:     terr.Basic,
			CallStack: terr.CallStack,
		}
	case *Client:
		return &Internal{
			Basic:     terr.Basic,
			CallStack: []string{function},
		}
	case *Basic:
		return &Internal{
			Basic:     terr,
			CallStack: []string{function},
		}
	default:
		panic("invalid error type")
	}
}

func InternalFromWithCode(code int, function string, err Error) *Internal {
	switch terr := err.(type) {
	case *Internal:
		terr.CallStack = append(terr.CallStack, function)
		terr.Code = code
		return terr
	case *Server:
		terr.CallStack = append(terr.CallStack, function)
		terr.Code = code
		return &Internal{
			Basic:     terr.Basic,
			CallStack: terr.CallStack,
		}
	case *Client:
		terr.Code = code
		return &Internal{
			Basic:     terr.Basic,
			CallStack: []string{function},
		}
	case *Basic:
		terr.Code = code
		return &Internal{
			Basic:     terr,
			CallStack: []string{function},
		}
	default:
		panic("invalid error type")
	}
}
