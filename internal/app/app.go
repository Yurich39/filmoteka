// Package app configures and runs application.
package app

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/exp/slog"

	"people-finder/pkg/logger"

	"people-finder/config"
	"people-finder/internal/controller/api"
	"people-finder/internal/usecase"
	"people-finder/internal/usecase/repo"
	"people-finder/internal/usecase/webapi"
	"people-finder/pkg/httpserver"
	"people-finder/pkg/postgres"

	"github.com/go-chi/chi/v5"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Env)

	// Repository
	db, err := postgres.New(cfg.StorageConfig)
	if err != nil {
		l.Debug("failed to init storage", l.Err(err))
		os.Exit(1)
	}

	// Use case 'person'
	UseCase := usecase.New(
		repo.New(db),
		webapi.New(),
	)

	// HTTP Server
	r := chi.NewRouter()
	api.NewRouter(r, l, UseCase)

	l.Info("starting server", slog.String("address", cfg.Address))

	httpServer := httpserver.New(r, cfg)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Debug("Failed to start server", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Debug("Server shutdown", err)
	}
}
