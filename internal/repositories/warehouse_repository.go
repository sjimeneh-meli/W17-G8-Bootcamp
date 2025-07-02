package repositories

import (
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

type WarehouseRepository interface {
	GetAll() (map[int]models.Warehouse, error)
	Create(warehouse models.Warehouse) (models.Warehouse, error)
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
		return nil, fmt.Errorf("error al leer warehouses: %w", err)
	}

	return warehouses, nil
}

func (r *WarehouseRepositoryImpl) Create(warehouse models.Warehouse) (models.Warehouse, error) {
	warehouses, err := r.loader.ReadAll()
	if err != nil {
		return models.Warehouse{}, fmt.Errorf("Error al crear el Id : %w", err)
	}

	warehouse.Id = len(warehouses) + 1

	warehouses[warehouse.Id] = warehouse

	warehousesSlice := r.loader.MapToSlice(warehouses)

	if err := r.loader.WriteAll(warehousesSlice); err != nil {
		return models.Warehouse{}, fmt.Errorf("error al guardar el warehouse: %w", err)
	}

	return warehouse, nil
}
