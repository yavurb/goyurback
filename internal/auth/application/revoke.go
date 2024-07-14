package application

import "context"

func (uc *apiKeyUsecase) RevokeAPIKey(ctx context.Context, publicID string) error {
	err := uc.repository.RevokeAPIKey(ctx, publicID)
	if err != nil {
		return err
	}

	return nil
}
