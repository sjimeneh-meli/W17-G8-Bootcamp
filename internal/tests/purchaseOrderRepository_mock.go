package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

type PurchaseOrderRepositoryMock struct {
	mock.Mock
}

func GetNewPurchaseOrderRepositoryMock() *PurchaseOrderRepositoryMock {
	return &PurchaseOrderRepositoryMock{}
}

func (r *PurchaseOrderRepositoryMock) GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error) {
	args := r.Called(ctx)
	return args.Get(0).(map[int]models.PurchaseOrder), args.Error(1)
}

func (r *PurchaseOrderRepositoryMock) Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error) {
	args := r.Called(ctx, order)
	return args.Get(0).(models.PurchaseOrder), args.Error(1)
}

func (r *PurchaseOrderRepositoryMock) ExistPurchaseOrderByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	args := r.Called(ctx, orderNumber)
	return args.Bool(0), args.Error(1)
}

func (r *PurchaseOrderRepositoryMock) GetPurchaseOrdersReportByBuyerId(ctx context.Context, buyerId int) (models.PurchaseOrderReport, error) {
	args := r.Called(ctx, buyerId)
	return args.Get(0).(models.PurchaseOrderReport), args.Error(1)
}

func (r *PurchaseOrderRepositoryMock) GetAllPurchaseOrdersReports(ctx context.Context) ([]models.PurchaseOrderReport, error) {
	args := r.Called(ctx)
	return args.Get(0).([]models.PurchaseOrderReport), args.Error(1)
}
