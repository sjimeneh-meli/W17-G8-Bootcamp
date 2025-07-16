package services

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var productBatchServiceInstance ProductBatchServiceI

// GetProductBatchService - Creates and returns a new instance of productBatchService with the required repository using singleton pattern
// GetProductBatchService - Crea y retorna una nueva instancia de productBatchService con el repositorio requerido usando patrón singleton
func GetProductBatchService(repository repositories.ProductBatchRepositoryI) ProductBatchServiceI {
	if productBatchServiceInstance != nil {
		return productBatchServiceInstance
	}

	productBatchServiceInstance = &productBatchService{
		repository: repository,
	}
	return productBatchServiceInstance
}

// ProductBatchServiceI - Interface defining the contract for product batch service operations with business logic
// ProductBatchServiceI - Interfaz que define el contrato para las operaciones del servicio de lotes de productos con lógica de negocio
type ProductBatchServiceI interface {
	// Create - Creates a new product batch in the system
	// Create - Crea un nuevo lote de producto en el sistema
	Create(ctx context.Context, model *models.ProductBatch) error

	// GetProductQuantityBySectionId - Retrieves the total quantity of products in a specific section
	// GetProductQuantityBySectionId - Obtiene la cantidad total de productos en una sección específica
	GetProductQuantityBySectionId(ctx context.Context, id int) int

	// ExistsWithBatchNumber - Checks if a product batch exists with the given batch number, excluding a specific ID
	// ExistsWithBatchNumber - Verifica si existe un lote de producto con el número de lote dado, excluyendo un ID específico
	ExistsWithBatchNumber(ctx context.Context, id int, batchNumber string) bool
}

// productBatchService - Implementation of ProductBatchServiceI containing business logic for product batch operations
// productBatchService - Implementación de ProductBatchServiceI que contiene la lógica de negocio para operaciones de lotes de productos
type productBatchService struct {
	repository repositories.ProductBatchRepositoryI // Repository dependency for data access / Dependencia del repositorio para acceso a datos
}

// Create - Delegates creating a product batch to the repository
// Create - Delega la creación de un lote de producto al repositorio
func (s *productBatchService) Create(ctx context.Context, model *models.ProductBatch) error {
	return s.repository.Create(ctx, model)
}

// GetProductQuantityBySectionId - Delegates retrieving product quantity by section ID to the repository
// GetProductQuantityBySectionId - Delega la obtención de cantidad de producto por ID de sección al repositorio
func (s *productBatchService) GetProductQuantityBySectionId(ctx context.Context, id int) int {
	return s.repository.GetProductQuantityBySectionId(ctx, id)
}

// ExistsWithBatchNumber - Delegates checking product batch existence by batch number to the repository
// ExistsWithBatchNumber - Delega la verificación de existencia de lote de producto por número de lote al repositorio
func (s *productBatchService) ExistsWithBatchNumber(ctx context.Context, id int, batchNumber string) bool {
	return s.repository.ExistsWithBatchNumber(ctx, id, batchNumber)
}
