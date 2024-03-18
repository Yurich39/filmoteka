package usecase

import (
	"context"
	"fmt"

	"filmoteka/internal/entity"
)

type MovieUseCase struct {
	repo MoviesRepo
	log  Logger
}

func NewMovies(repoMovies MoviesRepo, l Logger) *MovieUseCase {
	return &MovieUseCase{
		repo: repoMovies,
		log:  l,
	}
}

func (uc *MovieUseCase) Find(ctx context.Context, id int) (entity.Movie, error) {
	res, err := uc.repo.Get(ctx, id)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Find returned error: %w", op, err)
	}

	return res, nil
}

func (uc *MovieUseCase) FindMovie(ctx context.Context) ([]entity.Movie, error) {
	res, err := uc.repo.GetMovie(ctx)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Find returned error: %w", op, err)
	}

	return res, nil
}

func (uc *MovieUseCase) Save(ctx context.Context, data entity.MovieData) (entity.Movie, error) {

	id, err := uc.repo.Save(ctx, data)

	if err != nil {
		return entity.Movie{}, fmt.Errorf("%s: usecase.Save returned error: %w", op, err)
	}

	res := entity.Movie{
		Id:        &id,
		MovieData: data,
	}

	return res, nil
}

func (uc *MovieUseCase) Update(ctx context.Context, updates entity.Movie) (entity.Movie, error) {
	res, err := uc.repo.Update(ctx, updates)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Update returned error: %w", op, err)
	}

	return res, nil
}

func (uc *MovieUseCase) Delete(ctx context.Context, id int) (entity.Movie, error) {
	res, err := uc.repo.Delete(ctx, id)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Delete returned error: %w", op, err)
	}

	return res, nil
}

func (uc *MovieUseCase) List(ctx context.Context) ([]entity.Movie, error) {
	res, err := uc.repo.List(ctx)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Find returned error: %w", op, err)
	}

	return res, nil
}

func (uc *MovieUseCase) Next(ctx context.Context) ([]entity.Movie, error) {
	res, err := uc.repo.Next(ctx)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Find returned error: %w", op, err)
	}

	return res, nil
}
