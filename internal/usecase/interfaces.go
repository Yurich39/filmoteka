package usecase

import (
	"context"

	"filmoteka/internal/entity"

	"golang.org/x/exp/slog"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	Actor interface {
		Save(ctx context.Context, data entity.ActorData) (entity.Actor, error)
		Update(ctx context.Context, updates entity.Actor) (entity.Actor, error)
		Delete(ctx context.Context, id int) (entity.Actor, error)
		Find(ctx context.Context, id int) (entity.Actor, error)
		List(ctx context.Context) ([]entity.Actor, error)
		Next(ctx context.Context) ([]entity.Actor, error)
	}

	Movie interface {
		Save(ctx context.Context, data entity.MovieData) (entity.Movie, error)
		Update(ctx context.Context, updates entity.Movie) (entity.Movie, error)
		Delete(ctx context.Context, id int) (entity.Movie, error)
		Find(ctx context.Context, id int) (entity.Movie, error)
		FindMovie(ctx context.Context) ([]entity.Movie, error)
		List(ctx context.Context) ([]entity.Movie, error)
		Next(ctx context.Context) ([]entity.Movie, error)
	}

	ActorMovie interface {
		Save(ctx context.Context, data entity.ActorMovie) error
		// Update(ctx context.Context, updates entity.ActorMovie) (entity.ActorMovie, error)
		// Delete(ctx context.Context, id int) (entity.ActorMovie, error)
		// Find(ctx context.Context, id int) (entity.ActorMovie, error)
		List(ctx context.Context) ([]entity.ActorMovieData, error)
		// Next(ctx context.Context) ([]entity.ActorMovie, error)
	}

	ActorsRepo interface {
		Save(ctx context.Context, data entity.ActorData) (int, error)
		Update(ctx context.Context, updates entity.Actor) (entity.Actor, error)
		Delete(ctx context.Context, id int) (entity.Actor, error)
		Get(ctx context.Context, id int) (entity.Actor, error)
		List(ctx context.Context) ([]entity.Actor, error)
		Next(ctx context.Context) ([]entity.Actor, error)
	}

	MoviesRepo interface {
		Save(ctx context.Context, data entity.MovieData) (int, error)
		Update(ctx context.Context, updates entity.Movie) (entity.Movie, error)
		Delete(ctx context.Context, id int) (entity.Movie, error)
		Get(ctx context.Context, id int) (entity.Movie, error)
		GetMovie(ctx context.Context) ([]entity.Movie, error)
		List(ctx context.Context) ([]entity.Movie, error)
		Next(ctx context.Context) ([]entity.Movie, error)
	}

	ActorsMoviesRepo interface {
		Save(ctx context.Context, data entity.ActorMovie) error
		// Update(ctx context.Context, updates entity.ActorMovie) (entity.ActorMovie, error)
		// Delete(ctx context.Context, id int) (entity.ActorMovie, error)
		// Get(ctx context.Context, id int) (entity.ActorMovie, error)
		List(ctx context.Context) ([]entity.ActorMovieData, error)
		// Next(ctx context.Context) ([]entity.ActorMovie, error)
	}

	Logger interface {
		Debug(msg string, args ...any)
		Info(msg string, args ...any)
		Warn(msg string, args ...any)
		Error(msg string, args ...any)
		Err(err error) slog.Attr
	}
)
