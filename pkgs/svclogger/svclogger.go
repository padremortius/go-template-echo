package svclogger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// Logger -.
type Log struct {
	Logger *slog.Logger
	lvl    *slog.LevelVar
	Level  string `env-required:"true" yaml:"level" json:"level" env:"LOG_LEVEL"`
}

// New -.
func New(level string) *Log {
	renameFields := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Key = "timestamp"
			return a
		}
		if a.Key == slog.MessageKey {
			a.Key = "message"
			return a
		}
		return a
	}
	lVar := new(slog.LevelVar)
	lVar.Set(GetLevelByString(level))
	locLog := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       lVar,
		ReplaceAttr: renameFields,
	}))

	return &Log{
		Logger: locLog,
		lvl:    lVar,
	}
}

func GetLevelByString(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "error":
		return slog.LevelError
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

func (l *Log) ChangeLogLevel(ctx context.Context, lvl string) {
	l.Logger.InfoContext(ctx, fmt.Sprint("Change log level to ", lvl))
	l.lvl.Set(GetLevelByString(lvl))
}
