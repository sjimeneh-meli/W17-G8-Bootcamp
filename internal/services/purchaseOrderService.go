package services

import (
	"context"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var purchaseOrderServiceInstance PurchaseOrderServiceI

// GetPurchaseOrderService creates and returns a new instance of PurchaseOrderService with the required repositories
func GetPurchaseOrderService(purchaseOrderRepository repositories.PurchaseOrderRepositoryI, buyerRepository repositories.BuyerRepositoryI, productRecordRepository repositories.IProductRecordRepository) PurchaseOrderServiceI {
	if purchaseOrderServiceInstance != nil {
		return purchaseOrderServiceInstance
	}
	purchaseOrderServiceInstance = &PurchaseOrderService{
		PurchaseOrderRepository: purchaseOrderRepository,
		BuyerRepository:         buyerRepository,
		ProductRecordRepository: productRecordRepository,
	}
	return purchaseOrderServiceInstance
}

// PurchaseOrderServiceI defines the interface for purchase order service operations
type PurchaseOrderServiceI interface {
	GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error)
	GetPurchaseOrdersReport(ctx context.Context, id *int) ([]models.PurchaseOrderReport, error)
	Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error)
}

// PurchaseOrderService implements PurchaseOrderServiceI and contains business logic for purchase order operations
type PurchaseOrderService struct {
	PurchaseOrderRepository repositories.PurchaseOrderRepositoryI
	BuyerRepository         repositories.BuyerRepositoryI
	ProductRecordRepository repositories.IProductRecordRepository
}

// GetAll retrieves all purchase orders from the repository
func (s *PurchaseOrderService) GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error) {
	return s.PurchaseOrderRepository.GetAll(ctx)
}

// GetPurchaseOrdersReport retrieves purchase order reports
// If id is provided, returns report for that specific buyer, otherwise returns reports for all buyers
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

// Create creates a new purchase order with validation
// Validates that the order number doesn't already exist and that the buyer exists
func (s *PurchaseOrderService) Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error) {
	//validate that order number doesn't exists.

	exists, err := s.PurchaseOrderRepository.ExistPurchaseOrderByOrderNumber(ctx, order.OrderNumber)
	if err != nil {
		return models.PurchaseOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	if exists {
		return models.PurchaseOrder{}, fmt.Errorf("%w. %s %s %s", error_message.ErrAlreadyExists, "Order number ", order.OrderNumber, "already exists.")
	}

	//validate that buyer id exists.
	exists, err = s.BuyerRepository.ExistBuyerById(ctx, order.BuyerId)
	if err != nil {
		return models.PurchaseOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	if !exists {
		return models.PurchaseOrder{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", order.BuyerId, "doesn't exists.")
	}

	//validate that product record id exists.
	exists = s.ProductRecordRepository.ExistProductRecordByID(ctx, int64(order.ProductRecordId))
	if !exists {
		return models.PurchaseOrder{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Product record with Id", order.ProductRecordId, "doesn't exists.")
	}

	return s.PurchaseOrderRepository.Create(ctx, order)
}
