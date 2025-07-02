package services

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

type WarehouseService interface {
	GetAll() (map[int]models.Warehouse, error)
	Create(warehouse models.Warehouse) (models.Warehouse, error)
}

type WarehouseServiceImpl struct {
	warehouseRepository repositories.WarehouseRepository
}

func NewWarehouseService(warehouseRepository repositories.WarehouseRepository) *WarehouseServiceImpl {
	return &WarehouseServiceImpl{warehouseRepository: warehouseRepository}
}

func (s *WarehouseServiceImpl) GetAll() (map[int]models.Warehouse, error) {
	return s.warehouseRepository.GetAll()
}
func (s *WarehouseServiceImpl) Create(warehouse models.Warehouse) (models.Warehouse, error) {
	return s.warehouseRepository.Create(warehouse)
}
