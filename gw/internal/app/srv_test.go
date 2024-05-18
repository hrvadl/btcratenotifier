package app

import (
	"log"
	"net"
	"net/http"
	"reflect"
	"testing"
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
				Handler:  http.NewServeMux(),
				Addr:     net.JoinHostPort("0.0.0.0", "80"),
				ErrorLog: &log.Logger{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := newServer(tt.args.h, tt.args.addr, tt.args.log)
			if !reflect.DeepEqual(got.ErrorLog, tt.want.ErrorLog) {
				t.Fatal("logger doesn't match")
			}

			if !reflect.DeepEqual(got.Addr, tt.want.Addr) {
				t.Fatal("addr doesn't match")
			}

			if !reflect.DeepEqual(got.Handler, tt.want.Handler) {
				t.Fatal("handler doesn't match")
			}
		})
	}
}
