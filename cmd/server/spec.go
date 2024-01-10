package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/linden-honey/linden-honey-sdk-go/swaggerui"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/config"
)

func specHTTPHandler(cfg config.SpecConfig) (http.Handler, error) {
	m := http.NewServeMux()

	specURL := path.Join("/", path.Base(cfg.FilePath))
	m.Handle(specURL, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		http.ServeFile(w, r, cfg.FilePath)
	}))

	ui, err := swaggerui.New(
		swaggerui.WithSpecURL(specURL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize swagger ui: %w", err)
	}

	m.Handle("/", ui)

	return m, nil
}
