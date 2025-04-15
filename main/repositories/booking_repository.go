package repositories

import (
	"context"
	"main/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookingRepository interface {
	GetUserBookings(ctx context.Context, userID string) ([]models.Booking, error)
	GetBookingByID(ctx context.Context, id string) (*models.Booking, error)
	CreateBooking(ctx context.Context, booking *models.Booking) error
	CancelBooking(ctx context.Context, id string) error
	GetActiveBookingsByComputer(ctx context.Context, clubID string, pcNumber int) ([]models.Booking, error)
}

type bookingRepository struct {
	client *firestore.Client
}

func NewBookingRepository(client *firestore.Client) BookingRepository {
	return &bookingRepository{client: client}
}

func (r *bookingRepository) GetUserBookings(ctx context.Context, userID string) ([]models.Booking, error) {
	docs, err := r.client.Collection("bookings").
		Where("user_id", "==", userID).
		Where("status", "==", "active").
		Where("end_time", ">", time.Now()).
		OrderBy("start_time", firestore.Asc).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	bookings := make([]models.Booking, 0, len(docs))
	for _, doc := range docs {
		var booking models.Booking
		if err := doc.DataTo(&booking); err != nil {
			continue
		}
		booking.ID = doc.Ref.ID
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (r *bookingRepository) GetBookingByID(ctx context.Context, id string) (*models.Booking, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Booking ID is required")
	}
	doc, err := r.client.Collection("bookings").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.NotFound, "Booking not found")
		}
		return nil, err
	}
	var booking models.Booking
	if err := doc.DataTo(&booking); err != nil {
		return nil, err
	}
	booking.ID = doc.Ref.ID
	return &booking, nil
}

func (r *bookingRepository) CreateBooking(ctx context.Context, booking *models.Booking) error {
	if booking.ID == "" {
		return status.Error(codes.InvalidArgument, "Booking ID is required")
	}
	_, err := r.client.Collection("bookings").Doc(booking.ID).Set(ctx, booking)
	return err
}

func (r *bookingRepository) CancelBooking(ctx context.Context, id string) error {
	if id == "" {
		return status.Error(codes.InvalidArgument, "Booking ID is required")
	}
	_, err := r.client.Collection("bookings").Doc(id).Update(ctx, []firestore.Update{
		{Path: "status", Value: "cancelled"},
	})
	return err
}

func (r *bookingRepository) GetActiveBookingsByComputer(ctx context.Context, clubID string, pcNumber int) ([]models.Booking, error) {
	docs, err := r.client.Collection("bookings").
		Where("club_id", "==", clubID).
		Where("pc_number", "==", pcNumber).
		Where("status", "==", "active").
		Where("end_time", ">", time.Now()).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	bookings := make([]models.Booking, 0, len(docs))
	for _, doc := range docs {
		var booking models.Booking
		if err := doc.DataTo(&booking); err != nil {
			continue
		}
		booking.ID = doc.Ref.ID
		bookings = append(bookings, booking)
	}
	return bookings, nil
}
