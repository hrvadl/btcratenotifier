package app

import (
	"log"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	t.Parallel()
	type args struct {
		h    http.Handler
		addr string
		log  *log.Logger
	}
	tests := []struct {
		name string
		args args
		want *http.Server
	}{
		{
			name: "Should create new server correctly",
			args: args{
				h:    http.NewServeMux(),
				addr: net.JoinHostPort("0.0.0.0", "80"),
				log:  &log.Logger{},
			},
			want: &http.Server{
				ReadHeaderTimeout: readHeaderTimeout,
				WriteTimeout:      writeHeaderTimeout,
				IdleTimeout:       idleTimeout,
				Handler:           http.NewServeMux(),
				Addr:              net.JoinHostPort("0.0.0.0", "80"),
				ErrorLog:          &log.Logger{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := newServer(tt.args.h, tt.args.addr, tt.args.log)
			require.Equal(t, tt.want.ErrorLog, got.ErrorLog)
			require.Equal(t, tt.want.Addr, got.Addr)
			require.Equal(t, tt.want.Handler, got.Handler)
		})
	}
}
