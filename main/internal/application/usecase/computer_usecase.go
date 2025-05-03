package usecase

import (
	"context"
	"main/internal/domain/entities"
	"main/internal/domain/repository"
)

// ComputerUseCase defines business logic for Computer.
type ComputerUseCase interface {
	GetAll(ctx context.Context) ([]*entities.Computer, error)
	GetByClub(ctx context.Context, clubID string) ([]*entities.Computer, error)
	Create(ctx context.Context, comp *entities.Computer) error
}

type computerInteractor struct {
	repo repository.ComputerRepository
}

// NewComputerUseCase constructs a new ComputerUseCase with the given repository.
func NewComputerUseCase(r repository.ComputerRepository) ComputerUseCase {
	return &computerInteractor{repo: r}
}

func (u *computerInteractor) GetAll(ctx context.Context) ([]*entities.Computer, error) {
	return u.repo.FindAll(ctx)
}

func (u *computerInteractor) GetByClub(ctx context.Context, clubID string) ([]*entities.Computer, error) {
	return u.repo.FindByClub(ctx, clubID)
}

func (u *computerInteractor) Create(ctx context.Context, comp *entities.Computer) error {
	return u.repo.Create(ctx, comp)
}
