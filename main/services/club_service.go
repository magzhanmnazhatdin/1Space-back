package services

import (
	"context"
	"main/models"
	"main/repositories"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClubService interface {
	GetAllClubs(ctx context.Context) ([]models.ComputerClub, error)
	GetClubByID(ctx context.Context, id string) (*models.ComputerClub, error)
	CreateClub(ctx context.Context, club *models.ComputerClub) error
	UpdateClub(ctx context.Context, id string, club *models.ComputerClub) error
	DeleteClub(ctx context.Context, id string) error
}

type clubService struct {
	repo repositories.ClubRepository
}

func NewClubService(repo repositories.ClubRepository) ClubService {
	return &clubService{repo: repo}
}

func (s *clubService) GetAllClubs(ctx context.Context) ([]models.ComputerClub, error) {
	return s.repo.GetAllClubs(ctx)
}

func (s *clubService) GetClubByID(ctx context.Context, id string) (*models.ComputerClub, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Club ID is required")
	}
	return s.repo.GetClubByID(ctx, id)
}

func (s *clubService) CreateClub(ctx context.Context, club *models.ComputerClub) error {
	if club.ID == "" {
		return status.Error(codes.InvalidArgument, "Club ID is required")
	}
	if club.Name == "" || club.Address == "" || club.PricePerHour <= 0 {
		return status.Error(codes.InvalidArgument, "Invalid club data")
	}
	return s.repo.CreateClub(ctx, club)
}

func (s *clubService) UpdateClub(ctx context.Context, id string, club *models.ComputerClub) error {
	if id == "" {
		return status.Error(codes.InvalidArgument, "Club ID is required")
	}

	// Проверяем, что ID из тела совпадает с параметром маршрута
	if club.ID != id {
		return status.Error(codes.InvalidArgument, "Club ID in body must match the URL parameter")
	}

	// Валидация данных клуба
	if club.Name == "" {
		return status.Error(codes.InvalidArgument, "Club name is required")
	}
	if club.Address == "" {
		return status.Error(codes.InvalidArgument, "Club address is required")
	}
	if club.PricePerHour <= 0 {
		return status.Error(codes.InvalidArgument, "Price per hour must be greater than 0")
	}
	if club.AvailablePCs < 0 {
		return status.Error(codes.InvalidArgument, "Available PCs cannot be negative")
	}

	// Проверяем, существует ли клуб
	existingClub, err := s.repo.GetClubByID(ctx, id)
	if err != nil {
		return err
	}
	if existingClub == nil {
		return status.Error(codes.NotFound, "Club not found")
	}

	// Обновляем клуб в репозитории
	return s.repo.UpdateClub(ctx, id, club)
}

func (s *clubService) DeleteClub(ctx context.Context, id string) error {
	if id == "" {
		return status.Error(codes.InvalidArgument, "Club ID is required")
	}
	return s.repo.DeleteClub(ctx, id)
}
