package application

import (
	"context"
	"log"

	"github.com/yavurb/goyurback/internal/chikitos/domain"
	"github.com/yavurb/goyurback/internal/pgk/ids"
)

const prefix = "ch"

func (uc *ChikitoUsecase) Create(ctx context.Context, url, description string) (*domain.Chikito, error) {
	publicID, err := ids.NewPublicID(prefix)
	if err != nil {
		log.Printf("Error creating public id for chikito: %v\n", err)

		return nil, err
	}

	chikito := &domain.ChikitoCreate{
		PublicID:    publicID,
		URL:         url,
		Description: description,
	}

	chikitoCreated, err := uc.repository.CreateChikito(ctx, chikito)
	if err != nil {
		log.Printf("Unable to create chikito. Got: %v\n", err)

		return nil, err
	}

	return chikitoCreated, nil
}
