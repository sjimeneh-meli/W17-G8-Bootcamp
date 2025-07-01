package repositories

import (
	"fmt"
	"os"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

type BuyerRepositoryI interface {
	GetAll() (map[int]models.Buyer, error)
	GetById(id int) (models.Buyer, error)
	DeleteById(id int) error
	Create(buyer models.Buyer) (models.Buyer, error)
	GetNewId() int
	Save() error
	GetCardNumberIds() []string
}

type BuyerRepository struct {
	storage map[int]models.Buyer
}

func (r *BuyerRepository) GetAll() (map[int]models.Buyer, error) {
	return r.storage, nil
}

func (r *BuyerRepository) GetById(id int) (models.Buyer, error) {
	_, exists := r.storage[id]
	if !exists {
		return models.Buyer{}, fmt.Errorf("%w. %s %d", error_message.ErrNotFound, "Buyer with Id", id)
	}

	return r.storage[id], nil
}

func (r *BuyerRepository) DeleteById(id int) error {
	_, exists := r.storage[id]
	if !exists {
		return fmt.Errorf("%w. %s %d", error_message.ErrNotFound, "Buyer with Id", id)
	}
	delete(r.storage, id)

	err := r.Save()
	if err != nil {
		return err
	}
	return nil
}

func (r *BuyerRepository) Create(buyer models.Buyer) (models.Buyer, error) {
	_, exists := r.storage[buyer.Id]
	if exists {
		return models.Buyer{}, fmt.Errorf("%w. %s %d", error_message.ErrAlreadyExists, "Buyer with Id", buyer.Id)
	}

	r.storage[buyer.Id] = buyer

	err := r.Save()
	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err.Error())
	}

	return buyer, nil
}

func (r *BuyerRepository) GetNewId() int {
	lastId := 0
	for _, buyer := range r.storage {
		if buyer.Id > lastId {
			lastId = buyer.Id
		}
	}
	return (lastId + 1)
}

func (r *BuyerRepository) Save() error {
	jsonLoader := loader.NewJSONStorage[models.Buyer](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "buyers.json"))
	buyerArray := []models.Buyer{}

	for _, buyer := range r.storage {
		buyerArray = append(buyerArray, buyer)
	}

	err := jsonLoader.WriteAll(buyerArray)
	if err != nil {
		return fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err.Error())
	}
	return nil
}

func GetJsonBuyerRepository() (BuyerRepositoryI, error) {
	jsonLoader := loader.NewJSONStorage[models.Buyer](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "buyers.json"))
	storage, err := jsonLoader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w:%v", error_message.ErrInternalServerError, err)
	}

	return &BuyerRepository{
		storage: storage,
	}, nil
}

func (r *BuyerRepository) GetCardNumberIds() []string {
	cardNumbers := []string{}
	for _, buyer := range r.storage {
		cardNumbers = append(cardNumbers, buyer.CardNumberId)
	}
	return cardNumbers
}
