// internal/domain/repository/profile_repository.go
package repository

import (
	"context"
	"main/internal/domain/entities"
)

type ProfileRepository interface {
	GetByUser(ctx context.Context, uid string) (*entities.Profile, error)
	Save(ctx context.Context, p *entities.Profile) error
}
