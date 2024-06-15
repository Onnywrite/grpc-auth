package ero

type Server struct {
	*Basic
	CallStack []string
}

func NewServer(function string, errors ...string) *Server {
	return &Server{
		Basic:     New(errors...),
		CallStack: []string{function},
	}
}

func ServerFrom(function string, err Error) *Server {
	switch terr := err.(type) {
	case *Internal:
		terr.CallStack = append(terr.CallStack, function)
		return &Server{
			Basic:     terr.Basic,
			CallStack: terr.CallStack,
		}
	case *Server:
		terr.CallStack = append(terr.CallStack, function)
		return terr
	case *Client:
		return &Server{
			Basic:     terr.Basic,
			CallStack: []string{function},
		}
	case *Basic:
		return &Server{
			Basic:     terr,
			CallStack: []string{function},
		}
	default:
		panic("invalid error type")
	}
}
