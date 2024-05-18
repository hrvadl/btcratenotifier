package app

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/hrvadl/btcratenotifier/gw/internal/cfg"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		cfg cfg.Config
		log *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want *App
	}{
		{
			name: "Should create app with correct fields",
			args: args{
				cfg: cfg.Config{},
				log: slog.Default(),
			},
			want: &App{
				cfg: cfg.Config{},
				log: slog.Default(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := New(tt.args.cfg, tt.args.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
