package services

import (
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
	Create(model *models.ProductBatch) error
	GetProductQuantityBySectionId(id int) int
	ExistsWithBatchNumber(id int, batchNumber string) bool
}

type productBatchService struct {
	repository repositories.ProductBatchRepositoryI
}

func (s *productBatchService) Create(model *models.ProductBatch) error {
	return s.repository.Create(model)
}

func (s *productBatchService) GetProductQuantityBySectionId(id int) int {
	return s.repository.GetProductQuantityBySectionId(id)
}

func (s *productBatchService) ExistsWithBatchNumber(id int, batchNumber string) bool {
	return s.repository.ExistsWithBatchNumber(id, batchNumber)
}
