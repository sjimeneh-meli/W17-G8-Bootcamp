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
