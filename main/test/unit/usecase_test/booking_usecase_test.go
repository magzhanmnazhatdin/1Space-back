package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"main/internal/application/usecase"
	"main/internal/domain/entities"
)

// --- Mocks ---

type mockBookingRepo struct {
	mock.Mock
}

func (m *mockBookingRepo) FindAllByUser(ctx context.Context, userID string) ([]*entities.Booking, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entities.Booking), args.Error(1)
}

func (m *mockBookingRepo) FindByID(ctx context.Context, id string) (*entities.Booking, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockBookingRepo) Create(ctx context.Context, b *entities.Booking) error {
	return m.Called(ctx, b).Error(0)
}

func (m *mockBookingRepo) Update(ctx context.Context, b *entities.Booking) error {
	return m.Called(ctx, b).Error(0)
}

type mockComputerRepo struct {
	mock.Mock
}

func (m *mockComputerRepo) FindByClub(ctx context.Context, clubID string) ([]*entities.Computer, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]*entities.Computer), args.Error(1)
}

func (m *mockComputerRepo) FindByID(ctx context.Context, id string) (*entities.Computer, error) {
	return nil, nil
}
func (m *mockComputerRepo) FindAll(ctx context.Context) ([]*entities.Computer, error) {
	return nil, nil
}
func (m *mockComputerRepo) Create(ctx context.Context, c *entities.Computer) error { return nil }
func (m *mockComputerRepo) Update(ctx context.Context, c *entities.Computer) error { return nil }
func (m *mockComputerRepo) Delete(ctx context.Context, id string) error            { return nil }

// --- Tests ---

func TestCreateBooking_Success(t *testing.T) {
	ctx := context.Background()

	bookingRepo := new(mockBookingRepo)
	compRepo := new(mockComputerRepo)
	usecase := usecase.NewBookingUseCase(bookingRepo, compRepo)

	booking := &entities.Booking{
		ClubID:    "club1",
		UserID:    "user1",
		PCNumber:  1,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}

	computers := []*entities.Computer{
		{PCNumber: 1, IsAvailable: true},
		{PCNumber: 2, IsAvailable: true},
	}

	bookingRepo.On("Create", ctx, booking).Return(nil)
	compRepo.On("FindByClub", ctx, "club1").Return(computers, nil)
	compRepo.On("Update", ctx, mock.Anything).Return(nil)

	err := usecase.Create(ctx, booking)

	assert.NoError(t, err)
	assert.Equal(t, "active", booking.Status)
	assert.False(t, computers[0].IsAvailable)
}

func TestCreateBooking_ComputerNotFound(t *testing.T) {
	ctx := context.Background()

	bookingRepo := new(mockBookingRepo)
	compRepo := new(mockComputerRepo)
	usecase := usecase.NewBookingUseCase(bookingRepo, compRepo)

	booking := &entities.Booking{
		ClubID:    "club2",
		UserID:    "user2",
		PCNumber:  99,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}

	bookingRepo.On("Create", ctx, booking).Return(nil)
	compRepo.On("FindByClub", ctx, "club2").Return([]*entities.Computer{}, nil)

	err := usecase.Create(ctx, booking)

	assert.EqualError(t, err, "computer not found for update availability")
}
