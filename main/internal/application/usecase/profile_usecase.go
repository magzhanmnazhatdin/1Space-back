// internal/application/usecase/profile_usecase.go
package usecase

import (
	"context"
	"main/internal/domain/entities"
	"main/internal/domain/repository"
)

type ProfileUseCase interface {
	GetProfile(ctx context.Context, uid string) (*entities.Profile, error)
	UpdateProfile(ctx context.Context, p *entities.Profile) error
}

type profileInteractor struct {
	repo repository.ProfileRepository
}

func NewProfileUseCase(r repository.ProfileRepository) ProfileUseCase {
	return &profileInteractor{repo: r}
}

func (u *profileInteractor) GetProfile(ctx context.Context, uid string) (*entities.Profile, error) {
	return u.repo.GetByUser(ctx, uid)
}

func (u *profileInteractor) UpdateProfile(ctx context.Context, p *entities.Profile) error {
	return u.repo.Save(ctx, p)
}
