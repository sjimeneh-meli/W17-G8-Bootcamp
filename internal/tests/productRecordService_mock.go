package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

type (
	ProductRecordServiceMock struct {
		mock.Mock
	}
)

func GetNewProductRecordServiceMock() *ProductRecordServiceMock {
	return &ProductRecordServiceMock{}
}

func (prsm *ProductRecordServiceMock) CreateProductRecord(ctx context.Context, productRecord models.ProductRecord) (*models.ProductRecord, error) {
	args := prsm.Called(ctx, productRecord)
	return args.Get(0).(*models.ProductRecord), args.Error(1)
}

func (prsm *ProductRecordServiceMock) GetReportByIdProduct(ctx context.Context, id int64) (*models.ProductRecordReport, error) {
	args := prsm.Called(ctx, id)
	return args.Get(0).(*models.ProductRecordReport), args.Error(1)
}

func (prsm *ProductRecordServiceMock) GetReport(ctx context.Context) ([]*models.ProductRecordReport, error) {
	args := prsm.Called(ctx)
	return args.Get(0).([]*models.ProductRecordReport), args.Error(1)
}

func (prsm *ProductRecordServiceMock) ExistProductRecordByID(ctx context.Context, id int64) bool {
	args := prsm.Called(ctx, id)
	return args.Bool(0)
}
