package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type WarehouseServiceMock struct {
	GetWarehouseByIDFunc       func(ctx context.Context, id int) (models.Warehouse, error)
	GetByIdFunc                func(ctx context.Context, id int) (models.Warehouse, error)
	GetAllFunc                 func(ctx context.Context) ([]models.Warehouse, error)
	CreateFunc                 func(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error)
	ValidateCodeUniquenessFunc func(ctx context.Context, code string) error
	DeleteFunc                 func(ctx context.Context, id int) error
	UpdateFunc                 func(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error)
}

func (m *WarehouseServiceMock) GetById(ctx context.Context, id int) (models.Warehouse, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id)
	}
	return m.GetWarehouseByIDFunc(ctx, id)
}

func (m *WarehouseServiceMock) GetAll(ctx context.Context) ([]models.Warehouse, error) {
	return m.GetAllFunc(ctx)
}

func (m *WarehouseServiceMock) Create(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error) {
	return m.CreateFunc(ctx, warehouse)
}

func (m *WarehouseServiceMock) ValidateCodeUniqueness(ctx context.Context, code string) error {
	return m.ValidateCodeUniquenessFunc(ctx, code)
}

func (m *WarehouseServiceMock) Delete(ctx context.Context, id int) error {
	return m.DeleteFunc(ctx, id)
}

func (m *WarehouseServiceMock) Update(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error) {
	return m.UpdateFunc(ctx, id, warehouse)
}
