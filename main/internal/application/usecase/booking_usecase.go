package usecase

import (
	"context"
	"fmt"
	"main/internal/domain/entities"
	"main/internal/domain/repository"
	"time"
)

type BookingUseCase interface {
	GetByUser(ctx context.Context, userID string) ([]*entities.Booking, error)
	Create(ctx context.Context, b *entities.Booking) error
	Cancel(ctx context.Context, id string) error
}

// booking_usecase.go

type bookingInteractor struct {
	bookingRepo repository.BookingRepository
	compRepo    repository.ComputerRepository
}

func NewBookingUseCase(
	bRepo repository.BookingRepository,
	cRepo repository.ComputerRepository,
) BookingUseCase {
	return &bookingInteractor{bookingRepo: bRepo, compRepo: cRepo}
}

func (u *bookingInteractor) GetByUser(ctx context.Context, userID string) ([]*entities.Booking, error) {
	return u.bookingRepo.FindAllByUser(ctx, userID)
}

func (u *bookingInteractor) Create(ctx context.Context, b *entities.Booking) error {
	b.Status = "active"
	b.CreatedAt = time.Now()
	if err := u.bookingRepo.Create(ctx, b); err != nil {
		return err
	}
	// now mark that computer as unavailable
	comps, err := u.compRepo.FindByClub(ctx, b.ClubID)
	if err != nil {
		return err
	}
	for _, comp := range comps {
		if comp.PCNumber == b.PCNumber {
			comp.IsAvailable = false
			return u.compRepo.Update(ctx, comp)
		}
	}
	return fmt.Errorf("computer not found for update availability")
}

func (u *bookingInteractor) Cancel(ctx context.Context, id string) error {
	b, err := u.bookingRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	b.Status = "cancelled"
	if err := u.bookingRepo.Update(ctx, b); err != nil {
		return err
	}
	// restore availability
	comps, err := u.compRepo.FindByClub(ctx, b.ClubID)
	if err != nil {
		return err
	}
	for _, comp := range comps {
		if comp.PCNumber == b.PCNumber {
			comp.IsAvailable = true
			return u.compRepo.Update(ctx, comp)
		}
	}
	return fmt.Errorf("computer not found for restore availability")
}
