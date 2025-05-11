package usecase

import (
	"context"
	"main/internal/domain/entities"
	"main/internal/domain/repository"
)

// ClubUseCase defines business logic for Club, включая подсчёт свободных ПК.
type ClubUseCase interface {
	GetAll(ctx context.Context) ([]*entities.Club, error)
	GetByID(ctx context.Context, id string) (*entities.Club, error)
	Create(ctx context.Context, c *entities.Club) error
	Update(ctx context.Context, c *entities.Club) error
	Delete(ctx context.Context, id string) error
}

type clubInteractor struct {
	repo     repository.ClubRepository
	compRepo repository.ComputerRepository
}

// NewClubUseCase constructs a new ClubUseCase with the given repository.
func NewClubUseCase(r repository.ClubRepository, c repository.ComputerRepository) ClubUseCase {
	return &clubInteractor{repo: r, compRepo: c}
}

func (i *clubInteractor) GetAll(ctx context.Context) ([]*entities.Club, error) {
	clubs, err := i.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, club := range clubs {
		// динамически считаем доступные ПК
		comps, _ := i.compRepo.FindByClub(ctx, club.ID)
		count := 0
		for _, pc := range comps {
			if pc.IsAvailable {
				count++
			}
		}
		club.AvailablePCs = count
	}
	return clubs, nil
}

func (i *clubInteractor) GetByID(ctx context.Context, id string) (*entities.Club, error) {
	club, err := i.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	comps, _ := i.compRepo.FindByClub(ctx, id)
	count := 0
	for _, pc := range comps {
		if pc.IsAvailable {
			count++
		}
	}
	club.AvailablePCs = count
	return club, nil
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
