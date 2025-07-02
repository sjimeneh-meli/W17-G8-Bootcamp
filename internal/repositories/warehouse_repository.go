package repositories

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

type WarehouseRepository interface {
	GetAll() (map[int]models.Warehouse, error)
}

type WarehouseRepositoryImpl struct {
	loader loader.StorageJSON[models.Warehouse]
}

func NewWarehouseRepository(loader loader.StorageJSON[models.Warehouse]) *WarehouseRepositoryImpl {
	return &WarehouseRepositoryImpl{loader: loader}
}

func (r *WarehouseRepositoryImpl) GetAll() (map[int]models.Warehouse, error) {
	warehouses, err := r.loader.ReadAll()
	if err != nil {
		return nil, err
	}

	return warehouses, nil
}
