package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/linden-honey/linden-honey-sdk-go/health"
	sdkmiddleware "github.com/linden-honey/linden-honey-sdk-go/middleware"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/config"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/domain/song"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/application/middleware"
	httptransport "github.com/linden-honey/linden-honey-scraper-go/pkg/application/transport/http"
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/parser"
)

func main() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	{
		ctx = context.Background()
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
	}

	initLogger()

	slog.Info("initialization of the application")

	slog.Info("initializing configuration")

	var cfg *config.Config
	{
		var err error
		cfg, err = config.New()
		if err != nil {
			fatal(fmt.Errorf("failed to initialize a config: %w", err))
		}
	}

	slog.Info("initializing services")

	var songSvc song.Service
	{
		grobScr, err := newScraper(cfg.Scrapers.Grob, parser.NewGrobParser())
		if err != nil {
			fatal(fmt.Errorf("failed to initialize grob scraper: %w", err))
		}

		songSvc = domain.NewSongService(
			domain.SongServiceWithScraper(
				"grob", grobScr,
			),
		)

		songSvc = sdkmiddleware.Compose(
			middleware.SongLoggingMiddleware(),
		)(songSvc)
	}

	slog.Info("initializing http server")

	var httpServer *http.Server
	{
		r := chi.NewRouter()
		r.Use(chimiddleware.Recoverer)

		if cfg.Health.Enabled {
			r.Handle(cfg.Health.Path, health.NewHTTPHandler(health.NewNopService()))
		}

		specHandler, err := specHTTPHandler(cfg.Spec)
		if err != nil {
			fatal(fmt.Errorf("failed to initialize swagger: %w", err))
		}

		r.Mount("/", specHandler)

		r.Route("/api", func(r chi.Router) {
			r.Mount("/songs", httptransport.NewScraperHandler(songSvc))
		})

		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		httpServer = &http.Server{
			Addr:    addr,
			Handler: r,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		}

		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			if err := httpServer.Shutdown(ctx); err != nil {
				warn(fmt.Errorf("failed to shutdown http server: %w", err))
			}
		}()
	}

	errc := make(chan error, 1)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errc <- fmt.Errorf("failed to listen and serve http server: %w", err)
		}
	}()

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-sigc)
	}()

	slog.Info("application started")
	slog.Info("application stopped", "exit", <-errc)
}
