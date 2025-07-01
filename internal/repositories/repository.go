package repositories

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

type SellerRepository interface {
	GetAll() (map[int]models.Seller, error)
}

type JsonSellerRepository struct {
	storage loader.Storage[models.Seller]
}

func NewJSONSellerRepository(storage loader.Storage[models.Seller]) *JsonSellerRepository {

	return &JsonSellerRepository{
		storage: storage,
	}
}

func (r *JsonSellerRepository) GetAll() (map[int]models.Seller, error) {
	data, err := r.storage.ReadAll()
	//fmt.Println(data)
	if err != nil {
		return nil, err
	}
	return data, err
}
