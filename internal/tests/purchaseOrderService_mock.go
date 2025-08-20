package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

type PurchaseOrderServiceMock struct {
	mock.Mock
}

func GetNewPurchaseOrderServiceMock() *PurchaseOrderServiceMock {
	return &PurchaseOrderServiceMock{}
}

func (r *PurchaseOrderServiceMock) GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error) {
	args := r.Called(ctx)
	return args.Get(0).(map[int]models.PurchaseOrder), args.Error(1)
}

func (r *PurchaseOrderServiceMock) GetPurchaseOrdersReport(ctx context.Context, id *int) ([]models.PurchaseOrderReport, error) {
	args := r.Called(ctx, id)
	return args.Get(0).([]models.PurchaseOrderReport), args.Error(1)
}

func (r *PurchaseOrderServiceMock) Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error) {
	args := r.Called(ctx, order)
	return args.Get(0).(models.PurchaseOrder), args.Error(1)
}
