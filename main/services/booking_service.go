package services

import (
	"context"
	"main/models"
	"main/repositories"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookingService interface {
	GetUserBookings(ctx context.Context, userID string) ([]models.Booking, error)
	CreateBooking(ctx context.Context, userID string, bookingInput *BookingInput) (*models.Booking, error)
	CancelBooking(ctx context.Context, userID, bookingID string) error
}

type BookingInput struct {
	ClubID    string    `json:"club_id"`
	PCNumber  int       `json:"pc_number"`
	StartTime time.Time `json:"start_time"`
	Hours     int       `json:"hours"`
}

type bookingService struct {
	bookingRepo  repositories.BookingRepository
	clubRepo     repositories.ClubRepository
	computerRepo repositories.ComputerRepository
	client       *firestore.Client
}

func NewBookingService(bookingRepo repositories.BookingRepository, clubRepo repositories.ClubRepository, computerRepo repositories.ComputerRepository, client *firestore.Client) BookingService {
	return &bookingService{
		bookingRepo:  bookingRepo,
		clubRepo:     clubRepo,
		computerRepo: computerRepo,
		client:       client,
	}
}

func (s *bookingService) GetUserBookings(ctx context.Context, userID string) ([]models.Booking, error) {
	bookings, err := s.bookingRepo.GetUserBookings(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i, booking := range bookings {
		club, err := s.clubRepo.GetClubByID(ctx, booking.ClubID)
		if err == nil {
			bookings[i].ClubName = club.Name
		}
	}
	return bookings, nil
}

func (s *bookingService) CreateBooking(ctx context.Context, userID string, input *BookingInput) (*models.Booking, error) {
	if input.ClubID == "" || input.PCNumber <= 0 || input.Hours <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid booking data")
	}
	if input.StartTime.Before(time.Now()) {
		return nil, status.Error(codes.InvalidArgument, "Start time must be in the future")
	}

	club, err := s.clubRepo.GetClubByID(ctx, input.ClubID)
	if err != nil {
		return nil, err
	}

	computer, err := s.computerRepo.GetComputerByClubIDAndNumber(ctx, input.ClubID, input.PCNumber)
	if err != nil {
		return nil, err
	}
	if !computer.IsAvailable {
		return nil, status.Error(codes.FailedPrecondition, "Computer is already booked")
	}

	existingBookings, err := s.bookingRepo.GetActiveBookingsByComputer(ctx, input.ClubID, input.PCNumber)
	if err != nil {
		return nil, err
	}
	if len(existingBookings) > 0 {
		return nil, status.Error(codes.FailedPrecondition, "Computer is booked for the requested time")
	}

	booking := &models.Booking{
		ID:         uuid.New().String(),
		ClubID:     input.ClubID,
		UserID:     userID,
		PCNumber:   input.PCNumber,
		StartTime:  input.StartTime,
		EndTime:    input.StartTime.Add(time.Duration(input.Hours) * time.Hour),
		TotalPrice: club.PricePerHour * float64(input.Hours),
		Status:     "active",
		CreatedAt:  time.Now(),
	}

	err = s.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		if err := tx.Set(s.client.Collection("bookings").Doc(booking.ID), booking); err != nil {
			return err
		}
		return tx.Update(s.client.Collection("computers").Doc(computer.ID), []firestore.Update{
			{Path: "is_available", Value: false},
		})
	})
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *bookingService) CancelBooking(ctx context.Context, userID, bookingID string) error {
	if bookingID == "" {
		return status.Error(codes.InvalidArgument, "Booking ID is required")
	}

	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if booking.UserID != userID {
		return status.Error(codes.PermissionDenied, "Cannot cancel someone else's booking")
	}
	if booking.Status != "active" {
		return status.Error(codes.FailedPrecondition, "Booking is already cancelled or completed")
	}
	if time.Until(booking.StartTime) < time.Hour {
		return status.Error(codes.FailedPrecondition, "Can only cancel at least 1 hour before start time")
	}

	computer, err := s.computerRepo.GetComputerByClubIDAndNumber(ctx, booking.ClubID, booking.PCNumber)
	if err != nil {
		return err
	}

	err = s.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		if err := tx.Update(s.client.Collection("bookings").Doc(bookingID), []firestore.Update{
			{Path: "status", Value: "cancelled"},
		}); err != nil {
			return err
		}
		return tx.Update(s.client.Collection("computers").Doc(computer.ID), []firestore.Update{
			{Path: "is_available", Value: true},
		})
	})
	return err
}
