package ero

type Server struct {
	*Basic
	CallStack []string
}

func NewServer(code int, function string, errors ...string) *Server {
	return &Server{
		Basic:     New(code, errors...),
		CallStack: []string{function},
	}
}

func ServerFrom(function string, err Error) *Server {
	return ServerFromWithCode(err.GetCode(), function, err)
}

func ServerFromWithCode(code int, function string, err Error) *Server {
	switch terr := err.(type) {
	case *Internal:
		terr.CallStack = append(terr.CallStack, function)
		terr.Code = code
		return &Server{
			Basic:     terr.Basic,
			CallStack: terr.CallStack,
		}
	case *Server:
		terr.CallStack = append(terr.CallStack, function)
		terr.Code = code
		return terr
	case *Client:
		terr.Code = code
		return &Server{
			Basic:     terr.Basic,
			CallStack: []string{function},
		}
	case *Basic:
		terr.Code = code
		return &Server{
			Basic:     terr,
			CallStack: []string{function},
		}
	default:
		panic("invalid error type")
	}
}
