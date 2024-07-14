package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavurb/goyurback/internal/auth/domain"
	"github.com/yavurb/goyurback/internal/database/postgres"
	"github.com/yavurb/goyurback/internal/pgk/ids"
)

const prefix = "ak"

type Repository struct {
	db *postgres.Queries
}

func NewRepo(connpool *pgxpool.Pool) domain.APIKeyRepository {
	db := postgres.New(connpool)

	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateAPIKey(ctx context.Context, apiKey *domain.APIKeyCreate) (*domain.APIKey, error) {
	id, err := ids.NewPublicID(prefix) // TODO: handle errors and validate if the id already exists
	if err != nil {
		return nil, err
	}

	apiKey_, err := r.db.CreateAPIKey(ctx, postgres.CreateAPIKeyParams{
		PublicID: id,
		Key:      apiKey.Key,
		Name:     apiKey.Name,
	})
	// TODO: print log
	if err != nil {
		return nil, err
	}

	newApiKey := &domain.APIKey{
		ID:        apiKey_.ID,
		PublicID:  apiKey_.PublicID,
		CreatedAt: apiKey_.CreatedAt.Time,
		UpdatedAt: apiKey_.UpdatedAt.Time,
	}

	return newApiKey, nil
}

func (r *Repository) RevokeAPIKey(ctx context.Context, publicID string) error {
	err := r.db.RevokeAPIKey(ctx, publicID)
	if err != nil {
		log.Printf("DB Error revoking key: %v\n", err)

		return err
	}

	return nil
}
