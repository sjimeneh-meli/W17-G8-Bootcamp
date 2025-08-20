package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

type ProductBatchServiceMock struct {
	mock.Mock
}

func GetProductBatchServiceMock() *ProductBatchServiceMock {
	return &ProductBatchServiceMock{}
}

func (sm *ProductBatchServiceMock) Create(ctx context.Context, model *models.ProductBatch) error {
	args := sm.Called(ctx, model)
	return args.Error(0)
}

func (sm *ProductBatchServiceMock) GetProductQuantityBySectionId(ctx context.Context, id int) int {
	args := sm.Called(ctx, id)
	return args.Get(0).(int)
}

func (sm *ProductBatchServiceMock) ExistsWithBatchNumber(ctx context.Context, id int, batchNumber string) bool {
	args := sm.Called(ctx, id, batchNumber)
	return args.Get(0).(bool)
}
