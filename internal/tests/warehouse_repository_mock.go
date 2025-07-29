package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type WarehouseRepositoryMock struct {
	GetWarehouseByIDFunc func(ctx context.Context, id int) (models.Warehouse, error)
	GetAllFunc           func(ctx context.Context) ([]models.Warehouse, error)
	CreateFunc           func(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error)
	UpdateFunc           func(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error)
	DeleteFunc           func(ctx context.Context, id int) error
	ExistsByCodeFunc     func(ctx context.Context, code string) (bool, error)
}

func (m *WarehouseRepositoryMock) GetWarehouseByID(ctx context.Context, id int) (models.Warehouse, error) {
	return m.GetWarehouseByIDFunc(ctx, id)
}

func (m *WarehouseRepositoryMock) GetById(ctx context.Context, id int) (models.Warehouse, error) {
	return m.GetWarehouseByID(ctx, id)
}

func (m *WarehouseRepositoryMock) GetAll(ctx context.Context) ([]models.Warehouse, error) {
	return m.GetAllFunc(ctx)
}

func (m *WarehouseRepositoryMock) Create(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error) {
	return m.CreateFunc(ctx, warehouse)
}

func (m *WarehouseRepositoryMock) Update(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error) {
	return m.UpdateFunc(ctx, id, warehouse)
}

func (m *WarehouseRepositoryMock) Delete(ctx context.Context, id int) error {
	return m.DeleteFunc(ctx, id)
}

func (m *WarehouseRepositoryMock) ExistsByCode(ctx context.Context, code string) (bool, error) {
	return m.ExistsByCodeFunc(ctx, code)
}
