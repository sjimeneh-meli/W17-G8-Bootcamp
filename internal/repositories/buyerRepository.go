package repositories

import (
	"context"
	"fmt"
	"slices"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

// BuyerRepositoryI defines the interface for buyer repository operations
type BuyerRepositoryI interface {
	GetAll(ctx context.Context) (map[int]models.Buyer, error)
	GetById(ctx context.Context, id int) (models.Buyer, error)
	DeleteById(ctx context.Context, id int) error
	Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error)
	Update(ctx context.Context, buyerId int, buyer models.Buyer) (models.Buyer, error)

	GetCardNumberIds() ([]string, error)
	ExistBuyerById(ctx context.Context, buyerId int) (bool, error)
}

// BuyerRepository implements BuyerRepositoryI for in-memory storage with file persistence
type BuyerRepository struct {
	storage map[int]models.Buyer
	loader  loader.Storage[models.Buyer]
}

// Update updates an existing buyer in the repository
// Validates that the buyer exists and checks for duplicate card numbers
func (r *BuyerRepository) Update(ctx context.Context, id int, buyer models.Buyer) (models.Buyer, error) {
	_, exists := r.storage[id]
	if !exists {
		return models.Buyer{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", id, "doesn't exists.")
	}

	updatedBuyer := r.storage[id]
	if buyer.CardNumberId != "" {
		existingCardNumberIds, _ := r.GetCardNumberIds()
		if slices.Contains(existingCardNumberIds, buyer.CardNumberId) {
			return models.Buyer{}, fmt.Errorf("%w. %s %s %s", error_message.ErrAlreadyExists, "Card number with Id", buyer.CardNumberId, "already exists.")
		}
		updatedBuyer.CardNumberId = buyer.CardNumberId
	}
	if buyer.FirstName != "" {
		updatedBuyer.FirstName = buyer.FirstName
	}
	if buyer.LastName != "" {
		updatedBuyer.LastName = buyer.LastName
	}

	r.storage[id] = updatedBuyer
	err := r.Save()
	if err != nil {
		return models.Buyer{}, err
	}

	return updatedBuyer, nil
}

// GetAll retrieves all buyers from the repository
func (r *BuyerRepository) GetAll(ctx context.Context) (map[int]models.Buyer, error) {
	return r.storage, nil
}

// GetById retrieves a buyer by their ID
// Returns an error if the buyer doesn't exist
func (r *BuyerRepository) GetById(ctx context.Context, id int) (models.Buyer, error) {
	_, exists := r.storage[id]
	if !exists {
		return models.Buyer{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", id, "doesn't exists.")
	}

	return r.storage[id], nil
}

// DeleteById removes a buyer from the repository by their ID
// Returns an error if the buyer doesn't exist
func (r *BuyerRepository) DeleteById(ctx context.Context, id int) error {
	_, exists := r.storage[id]
	if !exists {
		return fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", id, "doesn't exists.")
	}
	delete(r.storage, id)

	err := r.Save()
	if err != nil {
		return err
	}
	return nil
}

// Create adds a new buyer to the repository
// Generates a new ID and validates that the buyer doesn't already exist
func (r *BuyerRepository) Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error) {
	newId := r.GetNewId()
	buyer.Id = newId

	_, exists := r.storage[buyer.Id]
	if exists {
		return models.Buyer{}, fmt.Errorf("%w. %s %d %s", error_message.ErrAlreadyExists, "Buyer with Id", buyer.Id, "already exists.")
	}

	r.storage[buyer.Id] = buyer

	err := r.Save()
	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err.Error())
	}

	return buyer, nil
}

// GetNewId generates a new unique ID for buyers
// Finds the highest existing ID and returns the next number
func (r *BuyerRepository) GetNewId() int {
	lastId := 0
	for _, buyer := range r.storage {
		if buyer.Id > lastId {
			lastId = buyer.Id
		}
	}
	return (lastId + 1)
}

// Save persists all buyers to the storage using the loader
func (r *BuyerRepository) Save() error {
	buyerArray := []models.Buyer{}

	for _, buyer := range r.storage {
		buyerArray = append(buyerArray, buyer)
	}

	err := r.loader.WriteAll(buyerArray)
	if err != nil {
		return fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err.Error())
	}
	return nil
}

// GetNewBuyerRepository creates and returns a new instance of BuyerRepository
// Loads existing buyers from the storage using the provided loader
func GetNewBuyerRepository(loader loader.Storage[models.Buyer]) (BuyerRepositoryI, error) {
	storage, err := loader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w:%v", error_message.ErrInternalServerError, err)
	}

	return &BuyerRepository{
		storage: storage,
		loader:  loader,
	}, nil
}

// GetCardNumberIds retrieves all card number IDs from the repository
// Returns a slice of all existing card number IDs
func (r *BuyerRepository) GetCardNumberIds() ([]string, error) {
	cardNumbers := []string{}
	for _, buyer := range r.storage {
		cardNumbers = append(cardNumbers, buyer.CardNumberId)
	}
	return cardNumbers, nil
}

// ExistBuyerById checks if a buyer with the given ID exists in the repository
// Returns true if the buyer exists, false otherwise
func (r *BuyerRepository) ExistBuyerById(ctx context.Context, buyerId int) (bool, error) {
	_, exists := r.storage[buyerId]
	if exists {
		return true, nil
	}
	return false, nil
}
