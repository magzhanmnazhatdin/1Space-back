// internal/infrastructure/firestore/profile_repo_fs.go
package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"main/internal/domain/entities"
	"main/internal/domain/repository"
)

type profileRepoFS struct{ client *firestore.Client }

func NewProfileRepoFS(c *firestore.Client) repository.ProfileRepository {
	return &profileRepoFS{client: c}
}

func (r *profileRepoFS) GetByUser(ctx context.Context, uid string) (*entities.Profile, error) {
	doc := r.client.Collection("profiles").Doc(uid)
	snap, err := doc.Get(ctx)
	if err != nil {
		return nil, err
	}
	var p entities.Profile
	snap.DataTo(&p)
	return &p, nil
}

func (r *profileRepoFS) Save(ctx context.Context, p *entities.Profile) error {
	_, err := r.client.Collection("profiles").Doc(p.UserID).Set(ctx, p)
	return err
}
