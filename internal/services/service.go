package services

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

type SellerService interface {
	GetAll() (map[int]models.Seller, error)
}

type JsonSellerService struct {
	repo repositories.SellerRepository
}

func NewJSONSellerService(repo repositories.SellerRepository) *JsonSellerService {
	return &JsonSellerService{
		repo: repo,
	}
}

func (s *JsonSellerService) GetAll() (map[int]models.Seller, error) {
	sellers, err := s.repo.GetAll()
	return sellers, err
}
