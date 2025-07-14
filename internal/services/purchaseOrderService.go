package services

import (
	"context"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

func GetPurchaseOrderService(purchaseOrderRepository repositories.PurchaseOrderRepositoryI, buyerRepository repositories.BuyerRepositoryI) PurchaseOrderServiceI {
	return &PurchaseOrderService{
		PurchaseOrderRepository: purchaseOrderRepository,
		BuyerRepository:         buyerRepository,
		//ProductRecordRepository: IProductRecordRepository
	}
}

type PurchaseOrderServiceI interface {
	GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error)
	GetPurchaseOrdersReport(ctx context.Context, id *int) ([]models.PurchaseOrderReport, error)
	Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error)
}

type PurchaseOrderService struct {
	PurchaseOrderRepository repositories.PurchaseOrderRepositoryI
	BuyerRepository         repositories.BuyerRepositoryI
}

func (s *PurchaseOrderService) GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error) {
	return s.PurchaseOrderRepository.GetAll(ctx)
}

func (s *PurchaseOrderService) GetPurchaseOrdersReport(ctx context.Context, id *int) ([]models.PurchaseOrderReport, error) {
	if id != nil {
		reports := []models.PurchaseOrderReport{}
		report, err := s.PurchaseOrderRepository.GetPurchaseOrdersReportByBuyerId(ctx, *id)
		if err != nil {
			return reports, err
		}
		reports = append(reports, report)
		return reports, nil
	}
	return s.PurchaseOrderRepository.GetAllPurchaseOrdersReports(ctx)
}

func (s *PurchaseOrderService) Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error) {
	//Valido que el order number no exista.

	exists, err := s.PurchaseOrderRepository.ExistPurchaseOrderByOrderNumber(ctx, order.OrderNumber)
	if err != nil {
		return models.PurchaseOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	if exists {
		return models.PurchaseOrder{}, fmt.Errorf("%w. %s %s %s", error_message.ErrAlreadyExists, "Order number ", order.OrderNumber, "already exists.")
	}

	//Valido que el buyer exista:
	exists, err = s.BuyerRepository.ExistBuyerById(ctx, order.BuyerId)
	if err != nil {
		return models.PurchaseOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	if !exists {
		return models.PurchaseOrder{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", order.BuyerId, "doesn't exists.")
	}

	/*
		Valido que el product record exista:
		Pendiente de implementar porque necesitaría productRecordRepository
		exists, err = s.ProductRecordRepository.ExistsProductRecordById(order.ProductRecordId)
		if err != nil {
			return models.PurchaseOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
		}
		if !exists {
			return models.PurchaseOrder{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Product record with Id", order.ProductRecordId, "doesn't exists.")
		}

	*/

	return s.PurchaseOrderRepository.Create(ctx, order)
}
