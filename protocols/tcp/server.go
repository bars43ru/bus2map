package tcp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"

	"github.com/bars43ru/bus2map/pkg/xslog"
)

type Server struct {
	address string
	handler ConnectionHandler
}

func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	slog.InfoContext(ctx, "start listener")

	addr, err := net.ResolveTCPAddr("tcp", s.address)
	if err != nil {
		return fmt.Errorf("resolve listener addr: %w", err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen listener: %w", err)
	}

	go func() {
		<-ctx.Done()
		if err := listener.Close(); err != nil {
			slog.ErrorContext(ctx, "close listener", xslog.Error(err))
		}
	}()
	return s.loopAcceptingConnection(ctx, listener)
}

func (s *Server) loopAcceptingConnection(ctx context.Context, listener net.Listener) error {
	slog.InfoContext(ctx, "loop accepting connection")
	var wg sync.WaitGroup

	for ctx.Err() == nil {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) && ctx.Err() != nil {
				return nil
			}
			slog.ErrorContext(ctx, "accept connection", xslog.Error(err))
			continue
		}

		log := slog.With(
			slog.String("remote-addr", conn.RemoteAddr().String()),
			slog.String("local-addr", conn.LocalAddr().String()),
		)

		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			go func() {
				<-ctx.Done()
				if err := conn.Close(); err != nil {
					log.ErrorContext(ctx, "close connection when context cancel", xslog.Error(err))
				}
			}()

			err := s.connectionHandler(ctx, conn)
			if err != nil {
				log.ErrorContext(ctx, "handler connection", xslog.Error(err))
			}
		}()
	}

	wg.Wait()
	return ctx.Err()
}

func (s *Server) connectionHandler(ctx context.Context, r io.Reader) error {
	err := s.handler.Accept(ctx, r)
	if err != nil {
		return fmt.Errorf("connection handler: %w", err)
	}
	return nil
}
