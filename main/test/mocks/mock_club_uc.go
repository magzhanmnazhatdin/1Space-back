package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"main/internal/domain/entities"
)

type MockClubUC struct {
	mock.Mock
}

func (m *MockClubUC) GetAll(ctx context.Context) ([]*entities.Club, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entities.Club), args.Error(1)
}

func (m *MockClubUC) GetByID(ctx context.Context, id string) (*entities.Club, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Club), args.Error(1)
}

func (m *MockClubUC) Create(ctx context.Context, c *entities.Club) error {
	return m.Called(ctx, c).Error(0)
}

func (m *MockClubUC) Update(ctx context.Context, c *entities.Club) error {
	return m.Called(ctx, c).Error(0)
}

func (m *MockClubUC) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
