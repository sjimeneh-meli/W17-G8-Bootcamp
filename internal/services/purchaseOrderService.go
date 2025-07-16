package services

import (
	"context"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var purchaseOrderServiceInstance PurchaseOrderServiceI

// GetPurchaseOrderService creates and returns a singleton instance of PurchaseOrderService with the required repositories
// GetPurchaseOrderService crea y retorna una instancia singleton de PurchaseOrderService con los repositorios requeridos
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

// PurchaseOrderServiceI defines the contract for purchase order service operations with business logic
// PurchaseOrderServiceI define el contrato para las operaciones de servicio de órdenes de compra con lógica de negocio
type PurchaseOrderServiceI interface {
	GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error)
	GetPurchaseOrdersReport(ctx context.Context, id *int) ([]models.PurchaseOrderReport, error)
	Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error)
}

// PurchaseOrderService implements PurchaseOrderServiceI and contains business logic for purchase order operations
// PurchaseOrderService implementa PurchaseOrderServiceI y contiene la lógica de negocio para operaciones de órdenes de compra
type PurchaseOrderService struct {
	PurchaseOrderRepository repositories.PurchaseOrderRepositoryI // Repository for purchase order data access / Repositorio para acceso a datos de órdenes de compra
	BuyerRepository         repositories.BuyerRepositoryI         // Repository for buyer validation / Repositorio para validación de compradores
	ProductRecordRepository repositories.IProductRecordRepository // Repository for product record validation / Repositorio para validación de registros de productos
}

// GetAll retrieves all purchase orders from the repository
// GetAll recupera todas las órdenes de compra del repositorio
func (s *PurchaseOrderService) GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error) {
	return s.PurchaseOrderRepository.GetAll(ctx)
}

// GetPurchaseOrdersReport retrieves purchase order reports with optional filtering by buyer ID
// If id is provided, returns report for that specific buyer, otherwise returns reports for all buyers
// GetPurchaseOrdersReport recupera reportes de órdenes de compra con filtrado opcional por ID de comprador
// Si se proporciona id, retorna el reporte para ese comprador específico, de lo contrario retorna reportes para todos los compradores
func (s *PurchaseOrderService) GetPurchaseOrdersReport(ctx context.Context, id *int) ([]models.PurchaseOrderReport, error) {
	if id != nil {
		// Get report for specific buyer / Obtener reporte para comprador específico
		reports := []models.PurchaseOrderReport{}
		report, err := s.PurchaseOrderRepository.GetPurchaseOrdersReportByBuyerId(ctx, *id)
		if err != nil {
			return reports, err
		}
		reports = append(reports, report)
		return reports, nil
	}
	// Get reports for all buyers / Obtener reportes para todos los compradores
	return s.PurchaseOrderRepository.GetAllPurchaseOrdersReports(ctx)
}

// Create creates a new purchase order with comprehensive business validation
// Validates that the order number doesn't already exist, the buyer exists, and the product record exists
// Create crea una nueva orden de compra con validación de negocio comprensiva
// Valida que el número de orden no exista, que el comprador exista, y que el registro de producto exista
func (s *PurchaseOrderService) Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error) {
	// Validate that order number doesn't exist / Validar que el número de orden no exista
	exists, err := s.PurchaseOrderRepository.ExistPurchaseOrderByOrderNumber(ctx, order.OrderNumber)
	if err != nil {
		return models.PurchaseOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	if exists {
		return models.PurchaseOrder{}, fmt.Errorf("%w. %s %s %s", error_message.ErrAlreadyExists, "Order number ", order.OrderNumber, "already exists.")
	}

	// Validate that buyer ID exists / Validar que el ID del comprador exista
	exists, err = s.BuyerRepository.ExistBuyerById(ctx, order.BuyerId)
	if err != nil {
		return models.PurchaseOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	if !exists {
		return models.PurchaseOrder{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", order.BuyerId, "doesn't exists.")
	}

	// Validate that product record ID exists / Validar que el ID del registro de producto exista
	exists = s.ProductRecordRepository.ExistProductRecordByID(ctx, int64(order.ProductRecordId))
	if !exists {
		return models.PurchaseOrder{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Product record with Id", order.ProductRecordId, "doesn't exists.")
	}

	// Create the purchase order after all validations pass / Crear la orden de compra después de que todas las validaciones pasen
	return s.PurchaseOrderRepository.Create(ctx, order)
}
