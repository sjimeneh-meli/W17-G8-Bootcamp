package services

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

func GetPurchaseOrderService(repo repositories.PurchaseOrderRepositoryI) PurchaseOrderServiceI {
	return &PurchaseOrderService{
		repository: repo,
	}
}

type PurchaseOrderServiceI interface {
	GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error)
}

type PurchaseOrderService struct {
	repository repositories.PurchaseOrderRepositoryI
}

func (s *PurchaseOrderService) GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error) {
	return s.repository.GetAll(ctx)
}
