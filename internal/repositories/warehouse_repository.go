package repositories

import (
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

type WarehouseRepository interface {
	GetAll() (map[int]models.Warehouse, error)
	Create(warehouse models.Warehouse) (models.Warehouse, error)
	ExistsByCode(code string) (bool, error)
	GetById(id int) (models.Warehouse, error)
	Delete(id int) error
	Update(id int, warehouse models.Warehouse) (models.Warehouse, error)
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
		return nil, fmt.Errorf("%w: %v", error_message.ErrDatabaseError, err)
	}

	return warehouses, nil
}

func (r *WarehouseRepositoryImpl) Create(warehouse models.Warehouse) (models.Warehouse, error) {
	warehouses, err := r.loader.ReadAll()
	if err != nil {
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrDatabaseError, err)
	}

	warehouse.Id = len(warehouses) + 1

	warehouses[warehouse.Id] = warehouse

	warehousesSlice := r.loader.MapToSlice(warehouses)

	if err := r.loader.WriteAll(warehousesSlice); err != nil {
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrDatabaseError, err)
	}

	return warehouse, nil
}

func (r *WarehouseRepositoryImpl) ExistsByCode(code string) (bool, error) {
	warehouses, err := r.loader.ReadAll()
	if err != nil {
		return false, fmt.Errorf("%w: %v", error_message.ErrDatabaseError, err)
	}

	for _, warehouse := range warehouses {
		if warehouse.WareHouseCode == code {
			return true, nil
		}
	}

	return false, nil
}

func (r *WarehouseRepositoryImpl) GetById(id int) (models.Warehouse, error) {
	warehouses, err := r.loader.ReadAll()
	if err != nil {
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrDatabaseError, err)
	}

	warehouse, exists := warehouses[id]
	if !exists {
		return models.Warehouse{}, fmt.Errorf("%w: warehouse con id %d", error_message.ErrEntityNotFound, id)
	}

	return warehouse, nil
}

func (r *WarehouseRepositoryImpl) Delete(id int) error {
	warehouses, err := r.loader.ReadAll()
	if err != nil {
		return fmt.Errorf("%w: %v", error_message.ErrDatabaseError, err)
	}

	// Verificar que el warehouse existe
	if _, exists := warehouses[id]; !exists {
		return fmt.Errorf("%w: warehouse con id %d", error_message.ErrEntityNotFound, id)
	}

	delete(warehouses, id)

	warehousesSlice := r.loader.MapToSlice(warehouses)

	if err := r.loader.WriteAll(warehousesSlice); err != nil {
		return fmt.Errorf("%w: %v", error_message.ErrDatabaseError, err)
	}

	return nil
}

func (r *WarehouseRepositoryImpl) Update(id int, warehouse models.Warehouse) (models.Warehouse, error) {
	warehouses, err := r.loader.ReadAll()
	if err != nil {
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrDatabaseError, err)
	}

	// Verificar que el warehouse existe
	if _, exists := warehouses[id]; !exists {
		return models.Warehouse{}, fmt.Errorf("%w: warehouse con id %d", error_message.ErrEntityNotFound, id)
	}

	// Asignar el ID al warehouse y actualizarlo
	warehouse.Id = id
	warehouses[id] = warehouse

	warehousesSlice := r.loader.MapToSlice(warehouses)

	if err := r.loader.WriteAll(warehousesSlice); err != nil {
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrDatabaseError, err)
	}

	return warehouse, nil
}
