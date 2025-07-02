package services

import (
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

type WarehouseService interface {
	GetAll() (map[int]models.Warehouse, error)
	Create(warehouse models.Warehouse) (models.Warehouse, error)
	ValidateCodeUniqueness(code string) error
	GetById(id int) (models.Warehouse, error)
	Delete(id int) error
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

func (s *WarehouseServiceImpl) ValidateCodeUniqueness(code string) error {
	exists, err := s.warehouseRepository.ExistsByCode(code)
	if err != nil {
		return fmt.Errorf("%w: %v", error_message.ErrDatabaseError, err)
	}

	if exists {
		return fmt.Errorf("%w: el código de almacén '%s' ya existe", error_message.ErrEntityExists, code)
	}

	return nil
}

func (s *WarehouseServiceImpl) GetById(id int) (models.Warehouse, error) {
	return s.warehouseRepository.GetById(id)
}

func (s *WarehouseServiceImpl) Delete(id int) error {
	return s.warehouseRepository.Delete(id)
}
