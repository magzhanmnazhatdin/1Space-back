package repository

import (
	"context"
	"main/internal/domain/entities"
)

// ComputerRepository defines persistence operations for Computer.
type ComputerRepository interface {
	FindAll(ctx context.Context) ([]*entities.Computer, error)
	FindByID(ctx context.Context, id string) (*entities.Computer, error)
	FindByClub(ctx context.Context, clubID string) ([]*entities.Computer, error)
	Create(ctx context.Context, comp *entities.Computer) error
	Update(ctx context.Context, comp *entities.Computer) error
	Delete(ctx context.Context, id string) error
}
