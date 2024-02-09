package cmd

import (
	"log/slog"
	"os"
)

func InitLogger() {
	var lvl slog.Level
	if _, ok := os.LookupEnv("DEBUG"); ok {
		lvl = slog.LevelDebug
	}

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
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
