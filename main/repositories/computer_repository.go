package repositories

import (
	"context"
	"main/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ComputerRepository interface {
	GetAllComputers(ctx context.Context) ([]models.Computer, error)
	GetComputersByClubID(ctx context.Context, clubID string) ([]models.Computer, error)
	CreateComputers(ctx context.Context, computers []models.Computer) error
	GetComputerByClubIDAndNumber(ctx context.Context, clubID string, pcNumber int) (*models.Computer, error)
	UpdateComputerAvailability(ctx context.Context, id string, isAvailable bool) error
}

type computerRepository struct {
	client *firestore.Client
}

func NewComputerRepository(client *firestore.Client) ComputerRepository {
	return &computerRepository{client: client}
}

func (r *computerRepository) GetAllComputers(ctx context.Context) ([]models.Computer, error) {
	docs, err := r.client.Collection("computers").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	computers := make([]models.Computer, 0, len(docs))
	for _, doc := range docs {
		var comp models.Computer
		if err := doc.DataTo(&comp); err != nil {
			continue
		}
		comp.ID = doc.Ref.ID
		computers = append(computers, comp)
	}
	return computers, nil
}

func (r *computerRepository) GetComputersByClubID(ctx context.Context, clubID string) ([]models.Computer, error) {
	docs, err := r.client.Collection("computers").
		Where("club_id", "==", clubID).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	computers := make([]models.Computer, 0, len(docs))
	for _, doc := range docs {
		var comp models.Computer
		if err := doc.DataTo(&comp); err != nil {
			continue
		}
		comp.ID = doc.Ref.ID
		computers = append(computers, comp)
	}
	return computers, nil
}

func (r *computerRepository) CreateComputers(ctx context.Context, computers []models.Computer) error {
	batch := r.client.Batch()
	for _, comp := range computers {
		if comp.ID == "" {
			return status.Error(codes.InvalidArgument, "Computer ID is required")
		}
		batch.Set(r.client.Collection("computers").Doc(comp.ID), comp)
	}
	_, err := batch.Commit(ctx)
	return err
}

func (r *computerRepository) GetComputerByClubIDAndNumber(ctx context.Context, clubID string, pcNumber int) (*models.Computer, error) {
	docs, err := r.client.Collection("computers").
		Where("club_id", "==", clubID).
		Where("pc_number", "==", pcNumber).
		Limit(1).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, status.Error(codes.NotFound, "Computer not found")
	}
	var comp models.Computer
	if err := docs[0].DataTo(&comp); err != nil {
		return nil, err
	}
	comp.ID = docs[0].Ref.ID
	return &comp, nil
}

func (r *computerRepository) UpdateComputerAvailability(ctx context.Context, id string, isAvailable bool) error {
	if id == "" {
		return status.Error(codes.InvalidArgument, "Computer ID is required")
	}
	_, err := r.client.Collection("computers").Doc(id).Update(ctx, []firestore.Update{
		{Path: "is_available", Value: isAvailable},
	})
	return err
}
