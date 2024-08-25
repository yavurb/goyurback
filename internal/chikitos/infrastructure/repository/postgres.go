package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yavurb/goyurback/internal/chikitos/domain"
	"github.com/yavurb/goyurback/internal/database/postgres"
)

type Repository struct {
	db *postgres.Queries
}

func NewRepo(connpool *pgxpool.Pool) domain.ChikitoRepository {
	db := postgres.New(connpool)
	return &Repository{db}
}

func (r *Repository) CreateChikito(ctx context.Context, chikito *domain.ChikitoCreate) (*domain.Chikito, error) {
	chikito_, err := r.db.CreateChikito(ctx, postgres.CreateChikitoParams{
		PublicID:    chikito.PublicID,
		Url:         chikito.URL,
		Description: chikito.Description,
	})
	if err != nil {
		log.Printf("DB Error creating chikito: %v\n", err)

		return nil, err
	}

	chikitoCreated := &domain.Chikito{
		ID:          chikito_.ID,
		PublicID:    chikito_.PublicID,
		URL:         chikito_.Url,
		Description: chikito_.Description,
		CreatedAt:   chikito_.CreatedAt.Time,
		UpdatedAt:   chikito_.UpdatedAt.Time,
	}

	return chikitoCreated, nil
}

func (r *Repository) GetChikito(ctx context.Context, id string) (*domain.Chikito, error) {
	chikito_, err := r.db.GetChikito(ctx, id)
	if err != nil {
		log.Printf("DB Error getting chikito: %v\n", err)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrChikitoNotFound
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
