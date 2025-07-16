package services

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var productBatchServiceInstance ProductBatchServiceI

func GetProductBatchService(repository repositories.ProductBatchRepositoryI) ProductBatchServiceI {
	if productBatchServiceInstance != nil {
		return productBatchServiceInstance
	}

	productBatchServiceInstance = &productBatchService{
		repository: repository,
	}
	return productBatchServiceInstance
}

type ProductBatchServiceI interface {
	Create(ctx context.Context, model *models.ProductBatch) error
	GetProductQuantityBySectionId(ctx context.Context, id int) int
	ExistsWithBatchNumber(ctx context.Context, id int, batchNumber string) bool
}

type productBatchService struct {
	repository repositories.ProductBatchRepositoryI
}

func (s *productBatchService) Create(ctx context.Context, model *models.ProductBatch) error {
	return s.repository.Create(ctx, model)
}

func (s *productBatchService) GetProductQuantityBySectionId(ctx context.Context, id int) int {
	return s.repository.GetProductQuantityBySectionId(ctx, id)
}

func (s *productBatchService) ExistsWithBatchNumber(ctx context.Context, id int, batchNumber string) bool {
	return s.repository.ExistsWithBatchNumber(ctx, id, batchNumber)
}
