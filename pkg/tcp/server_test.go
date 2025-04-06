package tcp

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		factory func() (*Server, error)
		wantErr bool
	}{
		{
			name: "empty listener",
			factory: func() (*Server, error) {
				return New("", nil)
			},
			wantErr: true,
		},
		{
			name: "empty handler",
			factory: func() (*Server, error) {
				return New("localhost:9900", nil)
			},
			wantErr: true,
		},
		{
			name: "success",
			factory: func() (*Server, error) {
				return New("localhost:9900",
					ConnectionHandlerFunc(func(_ context.Context, _ io.Reader) error {
						return nil
					}))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := tt.factory()
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, server)
			} else {
				require.NoError(t, err)
				require.NotNil(t, server)
			}
		})
	}
}

func TestNew_ErrorHost(t *testing.T) {
	noopConnectionHandlerFunc := ConnectionHandlerFunc(func(_ context.Context, _ io.Reader) error {
		return nil
	})

	srv, err := New("localhost32:1900", noopConnectionHandlerFunc)
	require.NoError(t, err)
	require.NotNil(t, srv)
	require.Error(t, srv.Run(context.Background()))
}

func TestServer_RunGraceFullShutdown(t *testing.T) {
	noopConnectionHandlerFunc := ConnectionHandlerFunc(func(_ context.Context, _ io.Reader) error {
		return nil
	})

	srv, err := New("localhost:9900", noopConnectionHandlerFunc)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = srv.Run(ctx)
	require.ErrorIs(t, err, context.Canceled)
}
