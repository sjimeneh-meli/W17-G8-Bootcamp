package services

import (
	"context"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var carryServiceInstance CarryService

// NewCarryService - Creates and returns a new instance of CarryServiceImpl with required repositories using singleton pattern
// NewCarryService - Crea y retorna una nueva instancia de CarryServiceImpl con los repositorios requeridos usando patrón singleton
func NewCarryService(r repositories.CarryRepository, lr repositories.LocalityRepository) CarryService {
	if carryServiceInstance != nil {
		return carryServiceInstance
	}
	carryServiceInstance = &CarryServiceImpl{carryRepository: r, localityRepository: lr}
	return carryServiceInstance
}

// CarryService - Interface defining the contract for carry service operations with business logic
// CarryService - Interfaz que define el contrato para las operaciones del servicio de transportistas con lógica de negocio
type CarryService interface {
	// CreateCarry - Creates a new carry (carrier) with business validation (locality existence and CID uniqueness)
	// CreateCarry - Crea un nuevo transportista con validación de negocio (existencia de localidad y unicidad de CID)
	CreateCarry(ctx context.Context, carry models.Carry) (models.Carry, error)

	// GetCarryReportByLocality - Retrieves carry reports for all localities or a specific locality with validation
	// GetCarryReportByLocality - Obtiene reportes de transportistas para todas las localidades o una localidad específica con validación
	GetCarryReportByLocality(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error)
}

// CarryServiceImpl - Implementation of CarryService containing business logic for carry operations
// CarryServiceImpl - Implementación de CarryService que contiene la lógica de negocio para operaciones de transportistas
type CarryServiceImpl struct {
	carryRepository    repositories.CarryRepository    // Repository dependency for carry data access / Dependencia del repositorio para acceso a datos de transportistas
	localityRepository repositories.LocalityRepository // Repository dependency for locality validation / Dependencia del repositorio para validación de localidades
}

// CreateCarry - Creates a new carry with comprehensive business validation
// CreateCarry - Crea un nuevo transportista con validación integral de negocio
func (s *CarryServiceImpl) CreateCarry(ctx context.Context, carry models.Carry) (models.Carry, error) {
	// Business validation: Verify that the locality exists before creating the carry
	// Validación de negocio: Verificar que la localidad existe antes de crear el transportista
	localityExists, err := s.localityRepository.ExistById(ctx, carry.LocalityId)
	if err != nil {
		return models.Carry{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	if !localityExists {
		return models.Carry{}, fmt.Errorf("%w: locality with id %d", error_message.ErrNotFound, carry.LocalityId)
	}

	// Business rule: CID must be unique across all carries
	// Regla de negocio: El CID debe ser único entre todos los transportistas
	exists, err := s.carryRepository.ExistsByCid(ctx, carry.Cid)
	if err != nil {
		return models.Carry{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	if exists {
		return models.Carry{}, fmt.Errorf("%w: resource with the provided identifier already exists", error_message.ErrAlreadyExists)
	}

	// If all validations pass, delegate to repository for persistence
	// Si todas las validaciones pasan, delegar al repositorio para la persistencia
	carry, err = s.carryRepository.Create(ctx, carry)
	if err != nil {
		return models.Carry{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	return carry, nil
}

// GetCarryReportByLocality - Retrieves carry reports with optional locality validation
// GetCarryReportByLocality - Obtiene reportes de transportistas con validación opcional de localidad
func (s *CarryServiceImpl) GetCarryReportByLocality(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
	// Business validation: Only validate locality existence if a specific locality ID is provided
	// Validación de negocio: Solo validar la existencia de la localidad si se proporciona un ID de localidad específico
	if localityID != 0 {
		exists, err := s.localityRepository.ExistById(ctx, localityID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
		}
		if !exists {
			return nil, fmt.Errorf("%w: locality with id %d", error_message.ErrNotFound, localityID)
		}
	}

	// Delegate to repository for data retrieval (if localityID = 0, returns all localities)
	// Delegar al repositorio para la recuperación de datos (si localityID = 0, retorna todas las localidades)
	reports, err := s.carryRepository.GetCarryReportsByLocality(ctx, localityID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	return reports, nil
}
