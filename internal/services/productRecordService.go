package services

import (
	"context"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

// ProductRecordServiceI - Interfaz que define el contrato para la lógica de negocio de registros de productos
// ProductRecordServiceI - Interface defining the contract for product record business logic
type ProductRecordServiceI interface {
	// CreateProductRecord - Crea un nuevo registro de producto con validación de negocio
	// CreateProductRecord - Creates a new product record with business validation
	CreateProductRecord(ctx context.Context, productRecord models.ProductRecord) (*models.ProductRecord, error)

	// GetReportByIdProduct - Genera un reporte para un producto específico con validación de negocio
	// GetReportByIdProduct - Generates a report for a specific product with business validation
	GetReportByIdProduct(ctx context.Context, id int64) (*models.ProductRecordReport, error)

	// GetReport - Obtiene el reporte completo de todos los registros de productos
	// GetReport - Retrieves the complete report of all product records
	GetReport(ctx context.Context) ([]*models.ProductRecordReport, error)

	// ExistProductRecordByID - Verifica si existe un registro de producto por su ID
	// ExistProductRecordByID - Checks if a product record exists by its ID
	ExistProductRecordByID(ctx context.Context, id int64) bool
}

// productRecordService - Implementación de la capa de servicio que maneja la lógica de negocio para registros de productos
// productRecordService - Service layer implementation that handles business logic for product records
type productRecordService struct {
	Repository     repositories.IProductRecordRepository // Dependencia del repositorio para acceso a datos / Repository dependency for data access
	ProductService ProductService                        // Dependencia del servicio de productos para validaciones / Product service dependency for validations
}

// NewProductRecordService - Función constructora que crea una nueva instancia del servicio con inyección de dependencias
// NewProductRecordService - Constructor function that creates a new service instance with dependency injection
func NewProductRecordService(repository repositories.IProductRecordRepository, ProductService ProductService) ProductRecordServiceI {
	return &productRecordService{Repository: repository, ProductService: ProductService}
}

// CreateProductRecord - Lógica de negocio para crear registros de productos con validación
// CreateProductRecord - Business logic for creating product records with validation
func (prs *productRecordService) CreateProductRecord(ctx context.Context, productRecord models.ProductRecord) (*models.ProductRecord, error) {
	// VALIDACIÓN DE NEGOCIO: Verificar si el producto referenciado existe antes de crear un registro
	// BUSINESS VALIDATION: Check if the referenced product exists before creating a record
	exist, err := prs.ProductService.ExistById(ctx, productRecord.ProductID)

	if err != nil {
		// Retorna estructura vacía en caso de error de base de datos
		// Return empty struct on database error
		return &models.ProductRecord{}, err
	}

	// REGLA DE NEGOCIO: No se puede crear un registro para un producto que no existe
	// BUSINESS RULE: Cannot create a record for a non-existent product
	if !exist {
		return &models.ProductRecord{}, fmt.Errorf("error: product by id : %d does not exist. %w", productRecord.ProductID, error_message.ErrDependencyNotFound)
	}

	// Si la validación pasa, delega al repositorio para la persistencia de datos
	// If validation passes, delegate to repository for data persistence
	return prs.Repository.Create(ctx, &productRecord)
}

// GetReportByIdProduct - Lógica de negocio para generar reportes de productos con validación
// GetReportByIdProduct - Business logic for generating product reports with validation
func (prs *productRecordService) GetReportByIdProduct(ctx context.Context, id int64) (*models.ProductRecordReport, error) {
	// VALIDACIÓN DE NEGOCIO: Verificar que el producto existe antes de generar el reporte
	// BUSINESS VALIDATION: Verify product exists before generating report
	exist, err := prs.ProductService.ExistById(ctx, id)

	if err != nil {
		// Retorna estructura vacía en caso de error de base de datos
		// Return empty struct on database error
		return &models.ProductRecordReport{}, err
	}

	// REGLA DE NEGOCIO: No se puede generar reporte para un producto que no existe
	// BUSINESS RULE: Cannot generate report for non-existent product
	if !exist {
		return &models.ProductRecordReport{}, fmt.Errorf("error: product by id : %d does not exist %w", id, error_message.ErrDependencyNotFound)
	}

	// Si la validación pasa, delega al repositorio para la recuperación de datos
	// If validation passes, delegate to repository for data retrieval
	return prs.Repository.GetReportByIdProduct(ctx, id)
}

// GetReport - Obtiene el reporte completo de todos los registros de productos
// GetReport - Retrieves the complete report of all product records
func (prs *productRecordService) GetReport(ctx context.Context) ([]*models.ProductRecordReport, error) {
	return prs.Repository.GetReport(ctx)
}

// ExistProductRecordByID - Verifica si existe un registro de producto por su ID
// ExistProductRecordByID - Checks if a product record exists by its ID
func (prs *productRecordService) ExistProductRecordByID(ctx context.Context, id int64) bool {
	return prs.Repository.ExistProductRecordByID(ctx, id)
}
