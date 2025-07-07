package services

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

type SellerService interface {
	GetAll() ([]models.Seller, error)
	GetById(id int) (models.Seller, error)
	Save(seller models.Seller) ([]models.Seller, error)
	Update(id int, seller models.Seller) ([]models.Seller, error)
	Delete(id int) error
}

type JsonSellerService struct {
	repo repositories.SellerRepository
}

func NewJSONSellerService(repo repositories.SellerRepository) *JsonSellerService {
	return &JsonSellerService{
		repo: repo,
	}
}

func (s *JsonSellerService) GetAll() ([]models.Seller, error) {
	sellers, err := s.repo.GetAll()
	return sellers, err
}

func (s *JsonSellerService) GetById(id int) (models.Seller, error) {
	sellers, err := s.repo.GetAll()
	if err != nil {
		return models.Seller{}, error_message.ErrNotFound
	}

	for _, seller := range sellers {
		if seller.Id == id {
			return seller, nil
		}
	}

	return models.Seller{}, error_message.ErrNotFound
}

func (s *JsonSellerService) Save(seller models.Seller) ([]models.Seller, error) {
	sellerCreated, err := s.repo.Save(seller)
	return sellerCreated, err
}

func (s *JsonSellerService) Update(id int, seller models.Seller) ([]models.Seller, error) {
	sellerFounded, err := s.repo.Update(id, seller)
	return sellerFounded, err
}

func (s *JsonSellerService) Delete(id int) error {
	err := s.repo.Delete(id)
	return err
}
