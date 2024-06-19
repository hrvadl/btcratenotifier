//go:build !integration

package logger

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapLevels(t *testing.T) {
	t.Parallel()
	type args struct {
		lvl string
	}
	tests := []struct {
		name string
		args args
		want slog.Level
	}{
		{
			name: "Should map debug level correctly",
			args: args{
				lvl: "DEBUG",
			},
			want: slog.LevelDebug,
		},
		{
			name: "Should map info level correctly",
			args: args{
				lvl: "INFO",
			},
			want: slog.LevelInfo,
		},
		{
			name: "Should map error level correctly",
			args: args{
				lvl: "ERROR",
			},
			want: slog.LevelError,
		},
		{
			name: "Should map unknown level to debug",
			args: args{
				lvl: "OGHEisg",
			},
			want: slog.LevelDebug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := MapLevels(tt.args.lvl)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		lvl string
		w   io.Writer
	}
	tests := []struct {
		name string
		args args
		want *slog.Logger
	}{
		{
			name: "Shold create slogger with correct fields",
			args: args{
				lvl: "INFO",
				w:   &bytes.Buffer{},
			},
			want: slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})),
		},
		{
			name: "Shold create slogger with correct fields",
			args: args{
				lvl: "ERROR",
				w:   os.Stderr,
			},
			want: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: slog.LevelError,
			})),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := New(tt.args.w, tt.args.lvl)
			require.Equal(t, tt.want, got)
		})
	}
}
