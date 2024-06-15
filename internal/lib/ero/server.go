package ero

type Server struct {
	*Basic
	Function string
}

func NewServer(function string, errors ...string) *Server {
	return &Server{
		Basic:    New(errors...),
		Function: function,
	}
}

func ServerFrom(function string, basic *Basic) *Server {
	return &Server{
		Basic:    basic,
		Function: function,
	}
}
