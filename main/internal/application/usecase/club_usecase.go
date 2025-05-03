package usecase

import (
	"context"
	"main/internal/domain/entities"
	"main/internal/domain/repository"
)

// ClubUseCase defines business logic for Club.
type ClubUseCase interface {
	GetAll(ctx context.Context) ([]*entities.Club, error)
	GetByID(ctx context.Context, id string) (*entities.Club, error)
	Create(ctx context.Context, c *entities.Club) error
	Update(ctx context.Context, c *entities.Club) error
	Delete(ctx context.Context, id string) error
}

type clubInteractor struct {
	repo repository.ClubRepository
}

// NewClubUseCase constructs a new ClubUseCase with the given repository.
func NewClubUseCase(r repository.ClubRepository) ClubUseCase {
	return &clubInteractor{repo: r}
}

func (i *clubInteractor) GetAll(ctx context.Context) ([]*entities.Club, error) {
	return i.repo.FindAll(ctx)
}

func (i *clubInteractor) GetByID(ctx context.Context, id string) (*entities.Club, error) {
	return i.repo.FindByID(ctx, id)
}

func (i *clubInteractor) Create(ctx context.Context, c *entities.Club) error {
	return i.repo.Create(ctx, c)
}

func (i *clubInteractor) Update(ctx context.Context, c *entities.Club) error {
	return i.repo.Update(ctx, c)
}

func (i *clubInteractor) Delete(ctx context.Context, id string) error {
	return i.repo.Delete(ctx, id)
}
