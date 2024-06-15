package ero

type Client struct {
	*Basic
}

func NewClient(errors ...string) *Client {
	return &Client{
		Basic: New(errors...),
	}
}

func ClientFrom(basic *Basic) *Client {
	return &Client{
		Basic: basic,
	}
}
