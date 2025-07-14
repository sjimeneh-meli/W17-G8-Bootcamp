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
	GetPurchaseOrdersReport(ctx context.Context, id *int) ([]models.PurchaseOrderReport, error)
}

type PurchaseOrderService struct {
	repository repositories.PurchaseOrderRepositoryI
}

func (s *PurchaseOrderService) GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error) {
	return s.repository.GetAll(ctx)
}

func (s *PurchaseOrderService) GetPurchaseOrdersReport(ctx context.Context, id *int) ([]models.PurchaseOrderReport, error) {
	if id != nil {
		reports := []models.PurchaseOrderReport{}
		report, err := s.repository.GetPurchaseOrdersReportByBuyerId(ctx, *id)
		if err != nil {
			return reports, err
		}
		reports = append(reports, report)
		return reports, nil
	}
	return s.repository.GetAllPurchaseOrdersReports(ctx)
}
