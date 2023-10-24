// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"people-finder/internal/entity"
)

type (
	Person interface {
		Save(ctx context.Context, person entity.Person) ([]entity.Person, error)
		Update(ctx context.Context, updates entity.UpdateRequest) ([]entity.Person, error)
		Delete(ctx context.Context, deleter entity.DelRequest) ([]entity.Person, error)
		Find(ctx context.Context, getter entity.Options) ([]entity.Person, error)
		ListPeople(ctx context.Context, getter entity.Options, personID int) ([]entity.Person, error)

		GetAge(name string) (int, error)
		GetGender(name string) (string, error)
		GetNationality(name string) (string, error)
	}

	PersonRepo interface {
		Save(ctx context.Context, person entity.Person) ([]entity.Person, error)
		Update(ctx context.Context, updates entity.UpdateRequest) ([]entity.Person, error)
		Delete(ctx context.Context, deleter entity.DelRequest) ([]entity.Person, error)
		Find(ctx context.Context, getter entity.Options) ([]entity.Person, error)
		ListPeople(ctx context.Context, getter entity.Options, personID int) ([]entity.Person, error)
	}

	EnrichWebAPI interface {
		EnrichAge(name string) (int, error)
		EnrichGender(name string) (string, error)
		EnrichNationality(name string) (string, error)
	}
)
