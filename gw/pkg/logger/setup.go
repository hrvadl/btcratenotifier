package logger

import (
	"io"
	"log/slog"
)

func MapLevels(lvl string) slog.Level {
	switch lvl {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func New(w io.Writer, lvl string) *slog.Logger {
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: MapLevels(lvl),
	}))
}
