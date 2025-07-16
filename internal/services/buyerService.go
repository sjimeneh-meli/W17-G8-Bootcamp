package services

import (
	"context"
	"fmt"
	"slices"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var buyerServiceInstance BuyerServiceI

// GetBuyerService creates and returns a new instance of BuyerService with the required repository
func GetBuyerService(repo repositories.BuyerRepositoryI) BuyerServiceI {
	if buyerServiceInstance != nil {
		return buyerServiceInstance
	}

	buyerServiceInstance = &BuyerService{
		repository: repo,
	}
	return buyerServiceInstance
}

// BuyerServiceI defines the interface for buyer service operations
type BuyerServiceI interface {
	GetAll(ctx context.Context) (map[int]models.Buyer, error)
	GetById(ctx context.Context, id int) (models.Buyer, error)
	DeleteById(ictx context.Context, d int) error
	Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error)
	Update(ctx context.Context, buyerId int, buyer models.Buyer) (models.Buyer, error)
}

// BuyerService implements BuyerServiceI and contains business logic for buyer operations
type BuyerService struct {
	repository repositories.BuyerRepositoryI
}

// GetAll retrieves all buyers from the repository
func (s *BuyerService) GetAll(ctx context.Context) (map[int]models.Buyer, error) {
	return s.repository.GetAll(ctx)
}

// GetById retrieves a buyer by their ID from the repository
func (s *BuyerService) GetById(ctx context.Context, id int) (models.Buyer, error) {
	return s.repository.GetById(ctx, id)
}

// DeleteById removes a buyer from the repository by their ID
func (s *BuyerService) DeleteById(ctx context.Context, id int) error {
	return s.repository.DeleteById(ctx, id)
}

// Create creates a new buyer with validation
// Validates that the card number doesn't already exist
func (s *BuyerService) Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error) {

	existingCardNumbers, err := s.repository.GetCardNumberIds()
	if err != nil {
		return models.Buyer{}, err
	}
	if slices.Contains(existingCardNumbers, buyer.CardNumberId) {
		return models.Buyer{}, fmt.Errorf("%w - %s %s %s", error_message.ErrAlreadyExists, "card number with id:", buyer.CardNumberId, "already exists.")
	}

	return s.repository.Create(ctx, buyer)
}

// Update updates an existing buyer with validation
// Validates that the card number doesn't already exist if it's being updated
func (s *BuyerService) Update(ctx context.Context, id int, buyer models.Buyer) (models.Buyer, error) {
	existingCardNumbers, err := s.repository.GetCardNumberIds()
	if err != nil {
		return models.Buyer{}, err
	}
	if slices.Contains(existingCardNumbers, buyer.CardNumberId) {
		return models.Buyer{}, fmt.Errorf("%w - %s %s %s", error_message.ErrAlreadyExists, "card number with id:", buyer.CardNumberId, "already exists.")
	}
	return s.repository.Update(ctx, id, buyer)
}
