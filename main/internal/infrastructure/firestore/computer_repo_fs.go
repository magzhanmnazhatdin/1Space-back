package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"main/internal/domain/entities"
	"main/internal/domain/repository"
)

// computerRepoFS implements ComputerRepository using Firestore as backend.
type computerRepoFS struct {
	client *firestore.Client
}

// NewComputerRepoFS creates a Firestore-based implementation of ComputerRepository.
func NewComputerRepoFS(c *firestore.Client) repository.ComputerRepository {
	return &computerRepoFS{client: c}
}

func (r *computerRepoFS) FindAll(ctx context.Context) ([]*entities.Computer, error) {
	docs, err := r.client.Collection("computers").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var out []*entities.Computer
	for _, doc := range docs {
		var c entities.Computer
		doc.DataTo(&c)
		c.ID = doc.Ref.ID
		out = append(out, &c)
	}
	return out, nil
}

func (r *computerRepoFS) FindByID(ctx context.Context, id string) (*entities.Computer, error) {
	doc, err := r.client.Collection("computers").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var c entities.Computer
	doc.DataTo(&c)
	c.ID = doc.Ref.ID
	return &c, nil
}

func (r *computerRepoFS) FindByClub(ctx context.Context, clubID string) ([]*entities.Computer, error) {
	iter := r.client.Collection("computers").Where("club_id", "==", clubID).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, err
	}
	var out []*entities.Computer
	for _, doc := range docs {
		var c entities.Computer
		doc.DataTo(&c)
		c.ID = doc.Ref.ID
		out = append(out, &c)
	}
	return out, nil
}

func (r *computerRepoFS) Create(ctx context.Context, c *entities.Computer) error {
	ref := r.client.Collection("computers").NewDoc()
	c.ID = ref.ID
	_, err := ref.Set(ctx, c)
	return err
}

func (r *computerRepoFS) Update(ctx context.Context, c *entities.Computer) error {
	_, err := r.client.Collection("computers").Doc(c.ID).Set(ctx, c)
	return err
}

func (r *computerRepoFS) Delete(ctx context.Context, id string) error {
	_, err := r.client.Collection("computers").Doc(id).Delete(ctx)
	return err
}
