package usecase

import (
	"context"

	"people-finder/internal/entity"

	"golang.org/x/exp/slog"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	Person interface {
		Save(ctx context.Context, data entity.Data) (entity.Person, error)
		Update(ctx context.Context, updates entity.Person) (entity.Person, error)
		Delete(ctx context.Context, id int) (entity.Person, error)
		Find(ctx context.Context, id int) (entity.Person, error)
		List(ctx context.Context) ([]entity.Person, error)
		Next(ctx context.Context) ([]entity.Person, error)
	}

	PersonRepo interface {
		Save(ctx context.Context, data entity.Data) (int, error)
		Update(ctx context.Context, updates entity.Person) (entity.Person, error)
		Delete(ctx context.Context, id int) (entity.Person, error)
		Get(ctx context.Context, id int) (entity.Person, error)
		List(ctx context.Context) ([]entity.Person, error)
		Next(ctx context.Context) ([]entity.Person, error)
	}

	EnrichWebAPI interface {
		EnrichAge(name string) (int, error)
		EnrichGender(name string) (string, error)
		EnrichNationality(name string) (string, error)
	}

	Logger interface {
		Debug(msg string, args ...any)
		Info(msg string, args ...any)
		Warn(msg string, args ...any)
		Error(msg string, args ...any)
		Err(err error) slog.Attr
	}
)
