package ero

type Client struct {
	*Basic
}

func NewClient(code int, errors ...string) *Client {
	return &Client{
		Basic: New(code, errors...),
	}
}

func ClientFrom(basic *Basic) *Client {
	return &Client{
		Basic: basic,
	}
}

func ClientFromWithCode(code int, basic *Basic) *Client {
	basic.Code = code
	return &Client{
		Basic: basic,
	}
}
