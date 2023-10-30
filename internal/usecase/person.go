package usecase

import (
	"context"
	"fmt"
	"sync"

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
	res, err := uc.repo.Get(ctx, id)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Find returned error: %w", op, err)
	}

	return res, nil
}

func (uc *PersonUseCase) Save(ctx context.Context, data entity.Data) (entity.Person, error) {

	uc.enrich(ctx, &data)

	id, err := uc.repo.Save(ctx, data)

	if err != nil {
		return entity.Person{}, fmt.Errorf("%s: usecase.Save returned error: %w", op, err)
	}

	res := entity.Person{
		Id:   &id,
		Data: data,
	}

	return res, nil
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

func (uc *PersonUseCase) enrich(ctx context.Context, data *entity.Data) {

	wg := new(sync.WaitGroup)
	wg.Add(3)

	// Request Age
	go func() {
		defer wg.Done()

		age, err := uc.webAPI.EnrichAge(*data.Name)

		if err != nil {
			uc.log.Debug("API fail:", err)
		} else {
			data.Age = &age
		}
	}()

	// Request Gender
	go func() {
		defer wg.Done()

		gender, err := uc.webAPI.EnrichGender(*data.Name)

		if err != nil {
			uc.log.Debug("API fail:", err)
		} else {
			data.Gender = &gender
		}
	}()

	// Request Nationality
	go func() {
		defer wg.Done()

		nationality, err := uc.webAPI.EnrichNationality(*data.Name)

		if err != nil {
			uc.log.Debug("API fail:", err)
		} else {
			data.Nationality = &nationality
		}
	}()

	wg.Wait()
}
