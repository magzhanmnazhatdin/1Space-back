// internal/application/usecase/computer_usecase.go
package usecase

import (
	"context"
	"main/internal/domain/entities"
	"main/internal/domain/repository"
)

type ComputerUseCase interface {
	GetAll(ctx context.Context) ([]*entities.Computer, error)
	GetByClub(ctx context.Context, clubID string) ([]*entities.Computer, error)
	GetByID(ctx context.Context, id string) (*entities.Computer, error) // ← new
	Create(ctx context.Context, comp *entities.Computer) error
	Update(ctx context.Context, comp *entities.Computer) error // ← new
	Delete(ctx context.Context, id string) error               // ← new
}

type computerInteractor struct {
	repo repository.ComputerRepository
}

func NewComputerUseCase(r repository.ComputerRepository) ComputerUseCase {
	return &computerInteractor{repo: r}
}

func (u *computerInteractor) GetAll(ctx context.Context) ([]*entities.Computer, error) {
	return u.repo.FindAll(ctx)
}

func (u *computerInteractor) GetByClub(ctx context.Context, clubID string) ([]*entities.Computer, error) {
	return u.repo.FindByClub(ctx, clubID)
}

func (u *computerInteractor) GetByID(ctx context.Context, id string) (*entities.Computer, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *computerInteractor) Create(ctx context.Context, comp *entities.Computer) error {
	return u.repo.Create(ctx, comp)
}

func (u *computerInteractor) Update(ctx context.Context, comp *entities.Computer) error {
	return u.repo.Update(ctx, comp)
}

func (u *computerInteractor) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
