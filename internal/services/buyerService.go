package services

import (
	"fmt"
	"slices"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
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
	Create(buyer models.Buyer) (models.Buyer, error)
	Update(buyerId int, buyer models.Buyer) (models.Buyer, error)
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

func (s *BuyerService) Create(buyer models.Buyer) (models.Buyer, error) {
	newId := s.repository.GetNewId()
	buyer.Id = newId

	existingCardNumbers := s.repository.GetCardNumberIds()
	if slices.Contains(existingCardNumbers, buyer.CardNumberId) {
		return models.Buyer{}, fmt.Errorf("%w - %s %s", error_message.ErrAlreadyExists, "card number with id:", buyer.CardNumberId)
	}

	return s.repository.Create(buyer)
}

func (s *BuyerService) Update(id int, buyer models.Buyer) (models.Buyer, error) {
	return s.repository.Update(id, buyer)
}
