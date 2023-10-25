package usecase

import (
	"context"
	"fmt"

	"people-finder/internal/entity"
)

const op = "internal.usecase"

type PersonUseCase struct {
	repo   PersonRepo
	webAPI EnrichWebAPI
}

func New(r PersonRepo, w EnrichWebAPI) *PersonUseCase {
	return &PersonUseCase{
		repo:   r,
		webAPI: w,
	}
}

func (uc *PersonUseCase) Find(ctx context.Context, id entity.Id) ([]entity.Person, error) {
	res, err := uc.repo.Find(ctx, id)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Find returned error: %w", op, err)
	}

	return res, nil
}

func (uc *PersonUseCase) Save(ctx context.Context, person entity.Person) ([]entity.Person, error) {

	res, err := uc.repo.Save(ctx, person)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Save returned error: %w", op, err)
	}

	return res, nil
}

func (uc *PersonUseCase) Update(ctx context.Context, updates entity.Person) ([]entity.Person, error) {
	res, err := uc.repo.Update(ctx, updates)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Update returned error: %w", op, err)
	}

	return res, nil
}

func (uc *PersonUseCase) Delete(ctx context.Context, id entity.Id) ([]entity.Person, error) {
	res, err := uc.repo.Delete(ctx, id)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Delete returned error: %w", op, err)
	}

	return res, nil
}

func (uc *PersonUseCase) List(ctx context.Context) ([]entity.Person, error) {
	res, err := uc.repo.List(ctx)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Find returned error: %w", op, err)
	}

	return res, nil
}

func (uc *PersonUseCase) Next(ctx context.Context) ([]entity.Person, error) {
	res, err := uc.repo.Next(ctx)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Find returned error: %w", op, err)
	}

	return res, nil
}

func (uc *PersonUseCase) GetAge(name string) (int, error) {

	res, err := uc.webAPI.EnrichAge(name)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc *PersonUseCase) GetGender(name string) (string, error) {

	res, err := uc.webAPI.EnrichGender(name)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc *PersonUseCase) GetNationality(name string) (string, error) {

	res, err := uc.webAPI.EnrichNationality(name)
	if err != nil {
		return res, err
	}

	return res, nil
}
