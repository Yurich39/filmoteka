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
	log    Logger
}

func New(r PersonRepo, w EnrichWebAPI, l Logger) *PersonUseCase {
	return &PersonUseCase{
		repo:   r,
		webAPI: w,
		log:    l,
	}
}

func (uc *PersonUseCase) Find(ctx context.Context, id int) (entity.Person, error) {
	res, err := uc.repo.Find(ctx, id)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Find returned error: %w", op, err)
	}

	return res, nil
}

func (uc *PersonUseCase) Save(ctx context.Context, data entity.Data) (entity.Person, error) {

	uc.enrich(ctx, &data)

	person, err := uc.repo.Save(ctx, data)

	if err != nil {
		return person, fmt.Errorf("%s: repo.Save returned error: %w", op, err)
	}

	return person, nil
}

func (uc *PersonUseCase) enrich(ctx context.Context, data *entity.Data) {

	// Request Age
	age, err := uc.webAPI.EnrichAge(*data.Name)

	if err != nil {
		uc.log.Debug("API fail:", err)
	} else {
		data.Age = &age
	}

	// Request Gender
	gender, err := uc.webAPI.EnrichGender(*data.Name)

	if err != nil {
		uc.log.Debug("API fail:", err)
	} else {
		data.Gender = &gender
	}

	// Request Nationality
	nationality, err := uc.webAPI.EnrichNationality(*data.Name)

	if err != nil {
		uc.log.Debug("API fail:", err)
	} else {
		data.Nationality = &nationality
	}

}

func (uc *PersonUseCase) Update(ctx context.Context, updates entity.Person) (entity.Person, error) {
	res, err := uc.repo.Update(ctx, updates)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Update returned error: %w", op, err)
	}

	return res, nil
}

func (uc *PersonUseCase) Delete(ctx context.Context, id int) (entity.Person, error) {
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
