package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"main/internal/domain/entities"
)

type MockComputerUC struct {
	mock.Mock
}

func (m *MockComputerUC) GetAll(ctx context.Context) ([]*entities.Computer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Computer), args.Error(1)
}

func (m *MockComputerUC) GetByClub(ctx context.Context, clubID string) ([]*entities.Computer, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]*entities.Computer), args.Error(1)
}

func (m *MockComputerUC) GetByID(ctx context.Context, id string) (*entities.Computer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Computer), args.Error(1)
}

func (m *MockComputerUC) Create(ctx context.Context, comp *entities.Computer) error {
	return m.Called(ctx, comp).Error(0)
}

func (m *MockComputerUC) Update(ctx context.Context, comp *entities.Computer) error {
	return m.Called(ctx, comp).Error(0)
}

func (m *MockComputerUC) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
