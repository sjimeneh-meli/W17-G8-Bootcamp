package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

type ProductRecordRepositoryMock struct {
	mock.Mock
}

func GetNewProductRecordRepositoryMock() *ProductRecordRepositoryMock {
	return &ProductRecordRepositoryMock{}
}

func (r *ProductRecordRepositoryMock) Create(ctx context.Context, pr *models.ProductRecord) (*models.ProductRecord, error) {
	args := r.Called(ctx, pr)
	return args.Get(0).(*models.ProductRecord), args.Error(1)
}

func (r *ProductRecordRepositoryMock) GetReportByIdProduct(ctx context.Context, id int64) (*models.ProductRecordReport, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*models.ProductRecordReport), args.Error(1)
}

func (r *ProductRecordRepositoryMock) GetReport(ctx context.Context) ([]*models.ProductRecordReport, error) {
	args := r.Called(ctx)
	return args.Get(0).([]*models.ProductRecordReport), args.Error(1)
}

func (r *ProductRecordRepositoryMock) ExistProductRecordByID(ctx context.Context, id int64) bool {
	args := r.Called(ctx, id)
	return args.Bool(0)
}
