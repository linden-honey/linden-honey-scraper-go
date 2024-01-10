package main

import (
	"log/slog"
	"os"
)

func initLogger() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(l)
}

func warn(err error) {
	slog.Warn(err.Error())
}

func fatal(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}
