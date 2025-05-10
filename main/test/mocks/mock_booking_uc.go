package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"main/internal/domain/entities"
)

type MockBookingUC struct {
	mock.Mock
}

func (m *MockBookingUC) GetByUser(ctx context.Context, userID string) ([]*entities.Booking, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entities.Booking), args.Error(1)
}

func (m *MockBookingUC) Create(ctx context.Context, b *entities.Booking) error {
	return m.Called(ctx, b).Error(0)
}

func (m *MockBookingUC) Cancel(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
