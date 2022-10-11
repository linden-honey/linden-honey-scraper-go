package main

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func newLogger() (logger log.Logger) {
	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowDebug())
	logger = level.NewInjector(logger, level.InfoValue())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	return logger
}

func warn(logger log.Logger, err error) {
	_ = level.Warn(logger).Log("err", err)
}

func fatal(logger log.Logger, err error) {
	_ = level.Error(logger).Log("err", err)
	os.Exit(1)
}
