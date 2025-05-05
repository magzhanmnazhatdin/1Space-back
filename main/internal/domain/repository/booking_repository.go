package repository

import (
	"context"
	"main/internal/domain/entities"
)

// BookingRepository defines persistence operations for Booking.
type BookingRepository interface {
	FindAllByUser(ctx context.Context, userID string) ([]*entities.Booking, error)
	FindByID(ctx context.Context, id string) (*entities.Booking, error)
	Create(ctx context.Context, b *entities.Booking) error
	Update(ctx context.Context, b *entities.Booking) error
}
