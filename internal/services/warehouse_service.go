package services

import (
	"context"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

type WarehouseService interface {
	GetAll(ctx context.Context) ([]models.Warehouse, error)
	Create(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error)
	ValidateCodeUniqueness(ctx context.Context, code string) error
	GetById(ctx context.Context, id int) (models.Warehouse, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error)
}

type WarehouseServiceImpl struct {
	warehouseRepository repositories.WarehouseRepository
}

func NewWarehouseService(warehouseRepository repositories.WarehouseRepository) *WarehouseServiceImpl {
	return &WarehouseServiceImpl{warehouseRepository: warehouseRepository}
}

func (s *WarehouseServiceImpl) GetAll(ctx context.Context) ([]models.Warehouse, error) {
	return s.warehouseRepository.GetAll(ctx)
}

func (s *WarehouseServiceImpl) Create(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error) {
	return s.warehouseRepository.Create(ctx, warehouse)
}

func (s *WarehouseServiceImpl) ValidateCodeUniqueness(ctx context.Context, code string) error {
	exists, err := s.warehouseRepository.ExistsByCode(ctx, code)
	if err != nil {
		return fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	if exists {
		return fmt.Errorf("%w: warehouse code '%s' already exists", error_message.ErrAlreadyExists, code)
	}

	return nil
}

func (s *WarehouseServiceImpl) GetById(ctx context.Context, id int) (models.Warehouse, error) {
	return s.warehouseRepository.GetById(ctx, id)
}

func (s *WarehouseServiceImpl) Delete(ctx context.Context, id int) error {
	return s.warehouseRepository.Delete(ctx, id)
}

func (s *WarehouseServiceImpl) Update(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error) {
	currentWarehouse, err := s.warehouseRepository.GetById(ctx, id)
	if err != nil {
		return models.Warehouse{}, err
	}

	if currentWarehouse.WareHouseCode != warehouse.WareHouseCode {
		if err := s.ValidateCodeUniqueness(ctx, warehouse.WareHouseCode); err != nil {
			return models.Warehouse{}, err
		}
	}

	return s.warehouseRepository.Update(ctx, id, warehouse)
}
