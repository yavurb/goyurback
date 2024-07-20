package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavurb/goyurback/internal/auth/domain"
	"github.com/yavurb/goyurback/internal/database/postgres"
	"github.com/yavurb/goyurback/internal/pgk/ids"
)

const prefix = "ak"

type Repository struct {
	db *postgres.Queries
}

func NewAPIKeyRepo(connpool *pgxpool.Pool) domain.APIKeyRepository {
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
		Name:      apiKey_.Name,
		Key:       apiKey_.Key,
		Revoked:   apiKey_.Revoked,
		CreatedAt: apiKey_.CreatedAt.Time,
		UpdatedAt: apiKey_.UpdatedAt.Time,
		RevokedAt: apiKey_.RevokedAt.Time,
	}

	return newApiKey, nil
}

func (r *Repository) GetAPIKeyByValue(ctx context.Context, key string) (*domain.APIKey, error) {
	apiKey_, err := r.db.GetAPIKeyByValue(ctx, key)
	if err != nil {
		log.Printf("Error getting APIKey. Got: %v", err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAPIKeyNotFound
		}

		log.Panicf("DB Error obtaining APIKey: %v", err)
	}

	apiKey := &domain.APIKey{
		ID:        apiKey_.ID,
		PublicID:  apiKey_.PublicID,
		Name:      apiKey_.Name,
		Key:       apiKey_.Key,
		Revoked:   apiKey_.Revoked,
		CreatedAt: apiKey_.CreatedAt.Time,
		UpdatedAt: apiKey_.UpdatedAt.Time,
		RevokedAt: apiKey_.RevokedAt.Time,
	}

	return apiKey, nil
}

func (r *Repository) RevokeAPIKey(ctx context.Context, publicID string) error {
	err := r.db.RevokeAPIKey(ctx, publicID)
	if err != nil {
		log.Printf("DB Error revoking key: %v\n", err)

		return err
	}

	return nil
}
