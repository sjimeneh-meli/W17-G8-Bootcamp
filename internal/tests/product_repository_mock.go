package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

type (
	ProductRepositoryMock struct {
		mock.Mock
	}
)

func (psm *ProductRepositoryMock) GetAll(ctx context.Context) ([]models.Product, error) {
	arg := psm.Called(ctx)

	if arg.Get(0) == nil {
		return nil, arg.Error(1)
	}

	return arg.Get(0).([]models.Product), arg.Error(1)
}

func (psm *ProductRepositoryMock) GetByID(ctx context.Context, id int64) (models.Product, error) {
	arg := psm.Called(ctx, id)

	return arg.Get(0).(models.Product), arg.Error(1)
}

func (psm *ProductRepositoryMock) Create(ctx context.Context, newProduct models.Product) (models.Product, error) {
	arg := psm.Called(ctx, newProduct)

	return arg.Get(0).(models.Product), arg.Error(1)
}

func (psm *ProductRepositoryMock) CreateByBatch(ctx context.Context, products []models.Product) ([]models.Product, error) {
	arg := psm.Called(ctx, products)

	if arg.Get(0) == nil {
		return nil, arg.Error(1)
	}

	return arg.Get(0).([]models.Product), arg.Error(1)
}

func (psm *ProductRepositoryMock) Update(ctx context.Context, id int64, product models.Product) (models.Product, error) {
	arg := psm.Called(ctx, id, product)

	return arg.Get(0).(models.Product), arg.Error(1)
}

func (psm *ProductRepositoryMock) Delete(ctx context.Context, id int64) error {
	arg := psm.Called(ctx, id)

	return arg.Error(0)
}

func (psm *ProductRepositoryMock) Exists(ctx context.Context, id int64) (bool, error) {
	arg := psm.Called(ctx, id)

	return arg.Bool(0), arg.Error(1)
}

func (psm *ProductRepositoryMock) ExistsByProductCode(ctx context.Context, productCode string) (bool, error) {
	arg := psm.Called(ctx, productCode)

	return arg.Bool(0), arg.Error(1)
}
