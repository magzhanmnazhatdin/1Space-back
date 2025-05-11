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
	GetUsersByClub(ctx context.Context, clubID string) ([]string, error)
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

func (u *bookingInteractor) GetUsersByClub(ctx context.Context, clubID string) ([]string, error) {
	bookings, err := u.bookingRepo.FindAllByClub(ctx, clubID)
	if err != nil {
		return nil, err
	}
	uniq := make(map[string]struct{})
	for _, b := range bookings {
		uniq[b.UserID] = struct{}{}
	}
	var users []string
	for uid := range uniq {
		users = append(users, uid)
	}
	return users, nil
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
