package usecase

import (
	"context"
	"main/internal/domain/entities"
	"main/internal/domain/repository"
)

// BookingUseCase defines business logic for Booking.
type BookingUseCase interface {
	GetByUser(ctx context.Context, userID string) ([]*entities.Booking, error)
	Create(ctx context.Context, b *entities.Booking) error
	Cancel(ctx context.Context, id string) error
}

type bookingInteractor struct {
	repo repository.BookingRepository
}

// NewBookingUseCase constructs a new BookingUseCase with the given repository.
func NewBookingUseCase(r repository.BookingRepository) BookingUseCase {
	return &bookingInteractor{repo: r}
}

func (u *bookingInteractor) GetByUser(ctx context.Context, userID string) ([]*entities.Booking, error) {
	return u.repo.FindAllByUser(ctx, userID)
}

func (u *bookingInteractor) Create(ctx context.Context, b *entities.Booking) error {
	b.Status = "active"
	b.CreatedAt = b.StartTime
	return u.repo.Create(ctx, b)
}

func (u *bookingInteractor) Cancel(ctx context.Context, id string) error {
	booking, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	booking.Status = "cancelled"
	return u.repo.Update(ctx, booking)
}
