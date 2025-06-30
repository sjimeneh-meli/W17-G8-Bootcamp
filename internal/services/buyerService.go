package services

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

func GetBuyerService(repo repositories.BuyerRepositoryI) BuyerServiceI {
	return &BuyerService{
		repository: repo,
	}
}

type BuyerServiceI interface {
	GetAll() []*models.Buyer
}

type BuyerService struct {
	repository repositories.BuyerRepositoryI
}

func (s *BuyerService) GetAll() []*models.Buyer {
	return s.repository.GetAll()
}
