package repositories

import (
	"context"
	"main/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClubRepository interface {
	GetAllClubs(ctx context.Context) ([]models.ComputerClub, error)
	GetClubByID(ctx context.Context, id string) (*models.ComputerClub, error)
	CreateClub(ctx context.Context, club *models.ComputerClub) error
	UpdateClub(ctx context.Context, id string, club *models.ComputerClub) error
	DeleteClub(ctx context.Context, id string) error
}

type clubRepository struct {
	client *firestore.Client
}

func NewClubRepository(client *firestore.Client) ClubRepository {
	return &clubRepository{client: client}
}

func (r *clubRepository) GetAllClubs(ctx context.Context) ([]models.ComputerClub, error) {
	docs, err := r.client.Collection("clubs").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	clubs := make([]models.ComputerClub, 0, len(docs))
	for _, doc := range docs {
		var club models.ComputerClub
		if err := doc.DataTo(&club); err != nil {
			continue
		}
		club.ID = doc.Ref.ID
		clubs = append(clubs, club)
	}
	return clubs, nil
}

func (r *clubRepository) GetClubByID(ctx context.Context, id string) (*models.ComputerClub, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Club ID is required")
	}
	doc, err := r.client.Collection("clubs").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.NotFound, "Club not found")
		}
		return nil, err
	}
	var club models.ComputerClub
	if err := doc.DataTo(&club); err != nil {
		return nil, err
	}
	club.ID = doc.Ref.ID
	return &club, nil
}

func (r *clubRepository) CreateClub(ctx context.Context, club *models.ComputerClub) error {
	if club.ID == "" {
		return status.Error(codes.InvalidArgument, "Club ID is required")
	}
	_, err := r.client.Collection("clubs").Doc(club.ID).Set(ctx, club)
	return err
}

func (r *clubRepository) UpdateClub(ctx context.Context, id string, club *models.ComputerClub) error {
	updates := []firestore.Update{}
	if club.Name != "" {
		updates = append(updates, firestore.Update{Path: "name", Value: club.Name})
	}
	if club.Address != "" {
		updates = append(updates, firestore.Update{Path: "address", Value: club.Address})
	}
	if club.PricePerHour > 0 {
		updates = append(updates, firestore.Update{Path: "price_per_hour", Value: club.PricePerHour})
	}
	if club.AvailablePCs >= 0 {
		updates = append(updates, firestore.Update{Path: "available_pcs", Value: club.AvailablePCs})
	}

	if len(updates) == 0 {
		return status.Error(codes.InvalidArgument, "No fields to update")
	}

	_, err := r.client.Collection("clubs").Doc(id).Update(ctx, updates)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return status.Error(codes.NotFound, "Club not found")
		}
		return err
	}
	return nil
}

func (r *clubRepository) DeleteClub(ctx context.Context, id string) error {
	if id == "" {
		return status.Error(codes.InvalidArgument, "Club ID is required")
	}
	_, err := r.client.Collection("clubs").Doc(id).Delete(ctx)
	return err
}
