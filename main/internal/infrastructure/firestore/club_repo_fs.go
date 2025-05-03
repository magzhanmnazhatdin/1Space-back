package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"main/internal/domain/entities"
	"main/internal/domain/repository"
)

// clubRepoFS implements ClubRepository using Firestore as backend.
type clubRepoFS struct {
	client *firestore.Client
}

// NewClubRepoFS creates a Firestore-based implementation of ClubRepository.
func NewClubRepoFS(c *firestore.Client) repository.ClubRepository {
	return &clubRepoFS{client: c}
}

func (r *clubRepoFS) FindAll(ctx context.Context) ([]*entities.Club, error) {
	docs, err := r.client.Collection("clubs").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var out []*entities.Club
	for _, doc := range docs {
		var c entities.Club
		doc.DataTo(&c)
		c.ID = doc.Ref.ID
		out = append(out, &c)
	}
	return out, nil
}

func (r *clubRepoFS) FindByID(ctx context.Context, id string) (*entities.Club, error) {
	doc, err := r.client.Collection("clubs").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var c entities.Club
	doc.DataTo(&c)
	c.ID = doc.Ref.ID
	return &c, nil
}

func (r *clubRepoFS) Create(ctx context.Context, c *entities.Club) error {
	ref := r.client.Collection("clubs").NewDoc()
	c.ID = ref.ID
	_, err := ref.Set(ctx, c)
	return err
}

func (r *clubRepoFS) Update(ctx context.Context, c *entities.Club) error {
	_, err := r.client.Collection("clubs").Doc(c.ID).Set(ctx, c)
	return err
}

func (r *clubRepoFS) Delete(ctx context.Context, id string) error {
	_, err := r.client.Collection("clubs").Doc(id).Delete(ctx)
	return err
}
