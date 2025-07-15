package services

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

func GetProductBatchService(repository repositories.ProductBatchRepositoryI) ProductBatchServiceI {
	return &productBatchService{
		repository: repository,
	}
}

type ProductBatchServiceI interface {
	Create(model *models.ProductBatch) error
	ExistsWithBatchNumber(id int, batchNumber string) bool
}

type productBatchService struct {
	repository repositories.ProductBatchRepositoryI
}

func (s *productBatchService) Create(model *models.ProductBatch) error {
	return s.repository.Create(model)
}

func (s *productBatchService) ExistsWithBatchNumber(id int, batchNumber string) bool {
	return s.repository.ExistsWithBatchNumber(id, batchNumber)
}
