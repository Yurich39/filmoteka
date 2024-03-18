package usecase

import (
	"context"
	"fmt"

	"filmoteka/internal/entity"
)

const op = "internal.usecase"

type ActorUseCase struct {
	repo ActorsRepo
	log  Logger
}

func NewActors(repoActors ActorsRepo, l Logger) *ActorUseCase {
	return &ActorUseCase{
		repo: repoActors,
		log:  l,
	}
}

func (uc *ActorUseCase) Find(ctx context.Context, id int) (entity.Actor, error) {
	res, err := uc.repo.Get(ctx, id)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Find returned error: %w", op, err)
	}

	return res, nil
}

func (uc *ActorUseCase) Save(ctx context.Context, data entity.ActorData) (entity.Actor, error) {

	id, err := uc.repo.Save(ctx, data)

	if err != nil {
		return entity.Actor{}, err
	}

	res := entity.Actor{
		Id:        &id,
		ActorData: data,
	}

	return res, nil
}

func (uc *ActorUseCase) Update(ctx context.Context, updates entity.Actor) (entity.Actor, error) {
	res, err := uc.repo.Update(ctx, updates)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Update returned error: %w", op, err)
	}

	return res, nil
}

func (uc *ActorUseCase) Delete(ctx context.Context, id int) (entity.Actor, error) {
	res, err := uc.repo.Delete(ctx, id)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Delete returned error: %w", op, err)
	}

	return res, nil
}

func (uc *ActorUseCase) List(ctx context.Context) ([]entity.Actor, error) {
	res, err := uc.repo.List(ctx)
	if err != nil {
		return res, fmt.Errorf("%s: repo.List returned error: %w", op, err)
	}

	return res, nil
}

func (uc *ActorUseCase) Next(ctx context.Context) ([]entity.Actor, error) {
	res, err := uc.repo.Next(ctx)
	if err != nil {
		return res, fmt.Errorf("%s: repo.Next returned error: %w", op, err)
	}

	return res, nil
}
