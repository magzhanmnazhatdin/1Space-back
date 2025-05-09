package repository

import (
	"context"
	"main/internal/domain/entities"
)

// ClubRepository defines persistence operations for Club.
type ClubRepository interface {
	FindAll(ctx context.Context) ([]*entities.Club, error)
	FindByID(ctx context.Context, id string) (*entities.Club, error)
	Create(ctx context.Context, club *entities.Club) error
	Update(ctx context.Context, club *entities.Club) error
	Delete(ctx context.Context, id string) error
}
