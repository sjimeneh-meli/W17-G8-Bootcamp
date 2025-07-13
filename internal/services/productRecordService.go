package services

import (
	"context"
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

type ProductRecordServiceI interface {
	CreateProductRecord(ctx context.Context, productRecord models.ProductRecord) (*models.ProductRecord, error)
	GetReportByIdProduct(ctx context.Context, id int64) (*models.ProductRecordReport, error)
}
type productRecordService struct {
	Repository repositories.IProductRecordRepository
}

func NewProductRecordService(repository repositories.IProductRecordRepository) ProductRecordServiceI {
	return &productRecordService{Repository: repository}
}

func (prs *productRecordService) CreateProductRecord(ctx context.Context, productRecord models.ProductRecord) (*models.ProductRecord, error) {
	exist, err := prs.Repository.ExistProductByID(ctx, productRecord.ProductID)

	if err != nil {
		return &models.ProductRecord{}, err
	}
	if !exist {
		return &models.ProductRecord{}, fmt.Errorf("error: product by id : %d does not exist. %w", productRecord.ProductID, error_message.ErrDependencyNotFound)
	}

	return prs.Repository.Create(ctx, &productRecord)
}
func (prs *productRecordService) GetReportByIdProduct(ctx context.Context, id int64) (*models.ProductRecordReport, error) {
	exist, err := prs.Repository.ExistProductByID(ctx, id)

	if err != nil {
		return &models.ProductRecordReport{}, err
	}
	if !exist {
		return &models.ProductRecordReport{}, fmt.Errorf("error: product by id : %d does not exist %w", id, error_message.ErrNotFound)
	}

	return prs.Repository.GetReportByIdProduct(ctx, id)
}
