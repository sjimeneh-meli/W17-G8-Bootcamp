package repositories

import (
	"fmt"
	"os"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

type BuyerRepositoryI interface {
	GetAll() []*models.Buyer
}

type BuyerRepository struct {
	storage []*models.Buyer
}

func (r BuyerRepository) GetAll() []*models.Buyer {
	return r.storage
}

func GetJsonBuyerRepository() BuyerRepositoryI {
	jsonLoader := loader.NewJSONStorage[models.Buyer](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "buyers.json"))
	storage, err := jsonLoader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	return &BuyerRepository{
		storage: jsonLoader.MapToSlice(storage),
	}
}
