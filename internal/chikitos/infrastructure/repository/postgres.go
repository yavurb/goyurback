package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavurb/goyurback/internal/chikitos/domain"
	"github.com/yavurb/goyurback/internal/database/postgres"
	"github.com/yavurb/goyurback/internal/pgk/ids"
)

const prefix = "ch"

type Repository struct {
	db *postgres.Queries
}

func NewRepo(connpool *pgxpool.Pool) domain.ChikitoRepository {
	db := postgres.New(connpool)
	return &Repository{db}
}

func (r *Repository) CreateChikito(ctx context.Context, url, description string) (*domain.Chikito, error) {
	publicID, err := ids.NewPublicID(prefix)
	if err != nil {
		log.Printf("Error creating public id for chikito: %v\n", err)

		return nil, err
	}

	chikito_, err := r.db.CreateChikito(ctx, postgres.CreateChikitoParams{
		PublicID:    publicID,
		Url:         url,
		Description: description,
	})
	if err != nil {
		log.Printf("DB Error creating chikito: %v\n", err)

		return nil, err
	}

	chikito := &domain.Chikito{
		ID:          chikito_.ID,
		PublicID:    chikito_.PublicID,
		URL:         chikito_.Url,
		Description: chikito_.Description,
		CreatedAt:   chikito_.CreatedAt.Time,
		UpdatedAt:   chikito_.UpdatedAt.Time,
	}

	return chikito, nil
}

func (r *Repository) GetChikito(ctx context.Context, id string) (*domain.Chikito, error) {
	chikito_, err := r.db.GetChikito(ctx, id)
	if err != nil {
		log.Printf("DB Error getting chikito: %v\n", err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotChikitoFound
		}

		return nil, err
	}

	chikito := &domain.Chikito{
		ID:          chikito_.ID,
		PublicID:    chikito_.PublicID,
		URL:         chikito_.Url,
		Description: chikito_.Description,
		CreatedAt:   chikito_.CreatedAt.Time,
		UpdatedAt:   chikito_.UpdatedAt.Time,
	}

	return chikito, nil
}
