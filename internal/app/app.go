package app

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/exp/slog"

	"github.com/go-chi/chi/v5"

	"filmoteka/config"
	"filmoteka/internal/controller/api"
	"filmoteka/internal/usecase"
	"filmoteka/internal/usecase/repo"
	"filmoteka/pkg/httpserver"
	"filmoteka/pkg/logger"
	"filmoteka/pkg/postgres"
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
	defer db.Close()

	// Creating usecase for actors
	actorsUseCase := usecase.NewActors(
		repo.NewActorsRepo(db),
		l,
	)

	// Creating usecase for movies
	moviesUseCase := usecase.NewMovies(
		repo.NewMoviesRepo(db),
		l,
	)

	// Creating usecase for many-to-many relationship between actors and films
	actorsMoviesUseCase := usecase.NewActorsMovies(
		repo.NewActorsMoviesRepo(db),
		l,
	)

	// HTTP Server
	r := chi.NewRouter()
	api.NewRouter(cfg, r, l, actorsUseCase, moviesUseCase, actorsMoviesUseCase)

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
