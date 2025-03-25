package tcp

import (
	"fmt"
)

func New(
	address string,
	handler ConnectionHandler,
) (*Server, error) {
	if address == "" {
		return nil, fmt.Errorf("param `address`empty")
	}
	if handler == nil {
		return nil, fmt.Errorf("param `handler` empty")
	}
	return &Server{
		address: address,
		handler: handler,
	}, nil
}
