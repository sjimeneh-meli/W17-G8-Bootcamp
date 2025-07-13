package services

import (
	"context"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

// ProductRecordServiceI - Interface defining the contract for product record business logic
// ProductRecordServiceI - Interfaz que define el contrato para la lógica de negocio de registros de productos
type ProductRecordServiceI interface {
	// CreateProductRecord - Creates a new product record with business validation
	// CreateProductRecord - Crea un nuevo registro de producto con validación de negocio
	CreateProductRecord(ctx context.Context, productRecord models.ProductRecord) (*models.ProductRecord, error)

	// GetReportByIdProduct - Generates a report for a specific product with business validation
	// GetReportByIdProduct - Genera un reporte para un producto específico con validación de negocio
	GetReportByIdProduct(ctx context.Context, id int64) (*models.ProductRecordReport, error)
}

// productRecordService - Service layer implementation that handles business logic for product records
// productRecordService - Implementación de la capa de servicio que maneja la lógica de negocio para registros de productos
type productRecordService struct {
	Repository repositories.IProductRecordRepository // Repository dependency for data access / Dependencia del repositorio para acceso a datos
}

// NewProductRecordService - Constructor function that creates a new service instance with dependency injection
// NewProductRecordService - Función constructora que crea una nueva instancia del servicio con inyección de dependencias
func NewProductRecordService(repository repositories.IProductRecordRepository) ProductRecordServiceI {
	return &productRecordService{Repository: repository}
}

// CreateProductRecord - Business logic for creating product records with validation
// CreateProductRecord - Lógica de negocio para crear registros de productos con validación
func (prs *productRecordService) CreateProductRecord(ctx context.Context, productRecord models.ProductRecord) (*models.ProductRecord, error) {
	// BUSINESS VALIDATION: Check if the referenced product exists before creating a record
	// VALIDACIÓN DE NEGOCIO: Verificar si el producto referenciado existe antes de crear un registro
	exist, err := prs.Repository.ExistProductByID(ctx, productRecord.ProductID)

	if err != nil {
		// Return empty struct on database error / Retorna estructura vacía en caso de error de base de datos
		return &models.ProductRecord{}, err
	}

	// BUSINESS RULE: Cannot create a record for a non-existent product
	// REGLA DE NEGOCIO: No se puede crear un registro para un producto que no existe
	if !exist {
		return &models.ProductRecord{}, fmt.Errorf("error: product by id : %d does not exist. %w", productRecord.ProductID, error_message.ErrDependencyNotFound)
	}

	// If validation passes, delegate to repository for data persistence
	// Si la validación pasa, delega al repositorio para la persistencia de datos
	return prs.Repository.Create(ctx, &productRecord)
}

// GetReportByIdProduct - Business logic for generating product reports with validation
// GetReportByIdProduct - Lógica de negocio para generar reportes de productos con validación
func (prs *productRecordService) GetReportByIdProduct(ctx context.Context, id int64) (*models.ProductRecordReport, error) {
	// BUSINESS VALIDATION: Verify product exists before generating report
	// VALIDACIÓN DE NEGOCIO: Verificar que el producto existe antes de generar el reporte
	exist, err := prs.Repository.ExistProductByID(ctx, id)

	if err != nil {
		// Return empty struct on database error / Retorna estructura vacía en caso de error de base de datos
		return &models.ProductRecordReport{}, err
	}

	// BUSINESS RULE: Cannot generate report for non-existent product
	// REGLA DE NEGOCIO: No se puede generar reporte para un producto que no existe
	if !exist {
		return &models.ProductRecordReport{}, fmt.Errorf("error: product by id : %d does not exist %w", id, error_message.ErrNotFound)
	}

	// If validation passes, delegate to repository for data retrieval
	// Si la validación pasa, delega al repositorio para la recuperación de datos
	return prs.Repository.GetReportByIdProduct(ctx, id)
}
