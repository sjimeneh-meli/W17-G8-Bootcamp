package services

import (
	"context"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var warehouseServiceInstance WarehouseService

// NewWarehouseService creates and returns a singleton instance of WarehouseServiceImpl with the required repository
// NewWarehouseService crea y retorna una instancia singleton de WarehouseServiceImpl con el repositorio requerido
func NewWarehouseService(warehouseRepository repositories.WarehouseRepository) WarehouseService {
	if warehouseServiceInstance != nil {
		return warehouseServiceInstance
	}
	warehouseServiceInstance = &WarehouseServiceImpl{warehouseRepository: warehouseRepository}
	return warehouseServiceInstance
}

// WarehouseService defines the contract for warehouse service operations with business logic and validation
// WarehouseService define el contrato para las operaciones de servicio de almacenes con lógica de negocio y validación
type WarehouseService interface {
	GetAll(ctx context.Context) ([]models.Warehouse, error)
	Create(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error)
	ValidateCodeUniqueness(ctx context.Context, code string) error
	GetById(ctx context.Context, id int) (models.Warehouse, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error)
}

// WarehouseServiceImpl implements WarehouseService and contains business logic for warehouse operations
// WarehouseServiceImpl implementa WarehouseService y contiene la lógica de negocio para operaciones de almacenes
type WarehouseServiceImpl struct {
	warehouseRepository repositories.WarehouseRepository // Repository for warehouse data access / Repositorio para acceso a datos de almacenes
}

// GetAll retrieves all warehouses from the repository
// GetAll recupera todos los almacenes del repositorio
func (s *WarehouseServiceImpl) GetAll(ctx context.Context) ([]models.Warehouse, error) {
	return s.warehouseRepository.GetAll(ctx)
}

// Create creates a new warehouse in the repository
// Create crea un nuevo almacén en el repositorio
func (s *WarehouseServiceImpl) Create(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error) {
	return s.warehouseRepository.Create(ctx, warehouse)
}

// ValidateCodeUniqueness validates that a warehouse code is unique in the system
// Returns an error if the code already exists or if there's an internal server error
// ValidateCodeUniqueness valida que un código de almacén sea único en el sistema
// Retorna un error si el código ya existe o si hay un error interno del servidor
func (s *WarehouseServiceImpl) ValidateCodeUniqueness(ctx context.Context, code string) error {
	// Check if warehouse code already exists / Verificar si el código de almacén ya existe
	exists, err := s.warehouseRepository.ExistsByCode(ctx, code)
	if err != nil {
		return fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	// Return error if code already exists / Retornar error si el código ya existe
	if exists {
		return fmt.Errorf("%w: warehouse code '%s' already exists", error_message.ErrAlreadyExists, code)
	}

	return nil
}

// GetById retrieves a warehouse by its ID from the repository
// GetById recupera un almacén por su ID del repositorio
func (s *WarehouseServiceImpl) GetById(ctx context.Context, id int) (models.Warehouse, error) {
	return s.warehouseRepository.GetById(ctx, id)
}

// Delete removes a warehouse by its ID from the repository
// Delete elimina un almacén por su ID del repositorio
func (s *WarehouseServiceImpl) Delete(ctx context.Context, id int) error {
	return s.warehouseRepository.Delete(ctx, id)
}

// Update modifies an existing warehouse with business validation for code uniqueness
// Validates code uniqueness only if the warehouse code has changed
// Update modifica un almacén existente con validación de negocio para unicidad de código
// Valida la unicidad del código solo si el código del almacén ha cambiado
func (s *WarehouseServiceImpl) Update(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error) {
	// Get current warehouse to compare codes / Obtener almacén actual para comparar códigos
	currentWarehouse, err := s.warehouseRepository.GetById(ctx, id)
	if err != nil {
		return models.Warehouse{}, err
	}

	// Validate code uniqueness only if code has changed / Validar unicidad del código solo si el código ha cambiado
	if currentWarehouse.WareHouseCode != warehouse.WareHouseCode {
		if err := s.ValidateCodeUniqueness(ctx, warehouse.WareHouseCode); err != nil {
			return models.Warehouse{}, err
		}
	}

	// Update warehouse after validation / Actualizar almacén después de la validación
	return s.warehouseRepository.Update(ctx, id, warehouse)
}
