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
	GetAll() (map[int]models.Buyer, error)
	GetById(id int) (models.Buyer, error)
	DeleteById(id int) error
}

type BuyerService struct {
	repository repositories.BuyerRepositoryI
}

func (s *BuyerService) GetAll() (map[int]models.Buyer, error) {
	return s.repository.GetAll()
}

func (s *BuyerService) GetById(id int) (models.Buyer, error) {
	return s.repository.GetById(id)
}

func (s *BuyerService) DeleteById(id int) error {
	return s.repository.DeleteById(id)
}
