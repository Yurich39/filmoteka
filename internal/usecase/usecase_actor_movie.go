package usecase

import (
	"context"
	"fmt"

	"filmoteka/internal/entity"
)

type ActorMovieUseCase struct {
	repo ActorsMoviesRepo
	log  Logger
}

func NewActorsMovies(repoActorsMovies ActorsMoviesRepo, l Logger) *ActorMovieUseCase {
	return &ActorMovieUseCase{
		repo: repoActorsMovies,
		log:  l,
	}
}

func (uc *ActorMovieUseCase) Save(ctx context.Context, data entity.ActorMovie) error {

	err := uc.repo.Save(ctx, data)

	if err != nil {
		return fmt.Errorf("%s: repo.Save returned error: %w", op, err)
	}

	return nil
}

func (uc *ActorMovieUseCase) List(ctx context.Context) ([]entity.ActorMovieData, error) {

	res, err := uc.repo.List(ctx)

	if err != nil {
		return res, fmt.Errorf("%s: repo.List returned error: %w", op, err)
	}

	return res, nil
}
