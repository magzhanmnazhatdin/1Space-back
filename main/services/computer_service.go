package services

import (
	"context"
	"main/models"
	"main/repositories"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ComputerService interface {
	GetAllComputers(ctx context.Context) ([]models.Computer, error)
	GetComputersByClubID(ctx context.Context, clubID string) ([]models.Computer, error)
	CreateComputers(ctx context.Context, clubID string, computers []models.Computer) ([]models.Computer, error)
}

type computerService struct {
	repo     repositories.ComputerRepository
	clubRepo repositories.ClubRepository
}

func NewComputerService(repo repositories.ComputerRepository, clubRepo repositories.ClubRepository) ComputerService {
	return &computerService{repo: repo, clubRepo: clubRepo}
}

func (s *computerService) GetAllComputers(ctx context.Context) ([]models.Computer, error) {
	return s.repo.GetAllComputers(ctx)
}

func (s *computerService) GetComputersByClubID(ctx context.Context, clubID string) ([]models.Computer, error) {
	if clubID == "" {
		return nil, status.Error(codes.InvalidArgument, "Club ID is required")
	}
	return s.repo.GetComputersByClubID(ctx, clubID)
}

func (s *computerService) CreateComputers(ctx context.Context, clubID string, computers []models.Computer) ([]models.Computer, error) {
	if clubID == "" {
		return nil, status.Error(codes.InvalidArgument, "Club ID is required")
	}
	if len(computers) == 0 {
		return nil, status.Error(codes.InvalidArgument, "At least one computer is required")
	}

	_, err := s.clubRepo.GetClubByID(ctx, clubID)
	if err != nil {
		return nil, err
	}

	for i, comp := range computers {
		if comp.PCNumber <= 0 || strings.TrimSpace(comp.Description) == "" {
			return nil, status.Error(codes.InvalidArgument, "Invalid computer data")
		}
		computers[i].ID = uuid.New().String()
		computers[i].ClubID = clubID
		computers[i].IsAvailable = true
	}

	if err := s.repo.CreateComputers(ctx, computers); err != nil {
		return nil, err
	}
	return computers, nil
}
