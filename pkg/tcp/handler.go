package tcp

import (
	"context"
	"io"
)

type ConnectionHandler interface {
	Accept(ctx context.Context, r io.Reader) error
}

type ConnectionHandlerFunc func(ctx context.Context, r io.Reader) error

func (h ConnectionHandlerFunc) Accept(ctx context.Context, r io.Reader) error {
	return h(ctx, r)
}
