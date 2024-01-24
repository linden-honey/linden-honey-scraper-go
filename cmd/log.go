package cmd

import (
	"log/slog"
	"os"
)

func InitLogger() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(l)
}

func Warn(err error) {
	slog.Warn(err.Error())
}

func Fatal(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}
