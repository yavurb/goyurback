package application

import (
	"github.com/yavurb/goyurback/internal/chikitos/domain"
)

type ChikitoUsecase struct {
	repository domain.ChikitoRepository
}

func NewChikitoUsecase(repository domain.ChikitoRepository) domain.ChikitoUsecase {
	return &ChikitoUsecase{repository}
}
