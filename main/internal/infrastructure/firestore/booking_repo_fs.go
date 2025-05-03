package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"main/internal/domain/entities"
	"main/internal/domain/repository"
)

// bookingRepoFS implements BookingRepository using Firestore as backend.
type bookingRepoFS struct {
	client *firestore.Client
}

// NewBookingRepoFS creates a Firestore-based implementation of BookingRepository.
func NewBookingRepoFS(c *firestore.Client) repository.BookingRepository {
	return &bookingRepoFS{client: c}
}

func (r *bookingRepoFS) FindAllByUser(ctx context.Context, userID string) ([]*entities.Booking, error) {
	iter := r.client.Collection("bookings").Where("user_id", "==", userID).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, err
	}
	var out []*entities.Booking
	for _, doc := range docs {
		var b entities.Booking
		doc.DataTo(&b)
		b.ID = doc.Ref.ID
		out = append(out, &b)
	}
	return out, nil
}

func (r *bookingRepoFS) FindByID(ctx context.Context, id string) (*entities.Booking, error) {
	doc, err := r.client.Collection("bookings").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var b entities.Booking
	doc.DataTo(&b)
	b.ID = doc.Ref.ID
	return &b, nil
}

func (r *bookingRepoFS) Create(ctx context.Context, b *entities.Booking) error {
	ref := r.client.Collection("bookings").NewDoc()
	b.ID = ref.ID
	_, err := ref.Set(ctx, b)
	return err
}

func (r *bookingRepoFS) Update(ctx context.Context, b *entities.Booking) error {
	_, err := r.client.Collection("bookings").Doc(b.ID).Set(ctx, b)
	return err
}
