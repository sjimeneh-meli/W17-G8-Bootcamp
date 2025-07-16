package services

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var localityServiceInstance LocalityService

// NewSQLLocalityService - Creates and returns a new instance of SQLLocalityService with the required repository using singleton pattern
// NewSQLLocalityService - Crea y retorna una nueva instancia de SQLLocalityService con el repositorio requerido usando patrón singleton
func NewSQLLocalityService(repo repositories.LocalityRepository) LocalityService {
	if localityServiceInstance != nil {
		return localityServiceInstance
	}
	localityServiceInstance = &SQLLocalityService{repo: repo}
	return localityServiceInstance
}

// LocalityService - Interface defining the contract for locality service operations with business logic
// LocalityService - Interfaz que define el contrato para las operaciones del servicio de localidades con lógica de negocio
type LocalityService interface {
	// Save - Creates a new locality in the system with country and province management
	// Save - Crea una nueva localidad en el sistema con manejo de país y provincia
	Save(ctx context.Context, locality models.Locality) (models.Locality, error)

	// GetSellerReports - Retrieves seller reports for all localities or a specific locality by ID
	// GetSellerReports - Obtiene reportes de vendedores para todas las localidades o una localidad específica por ID
	GetSellerReports(ctx context.Context, id int) ([]responses.LocalitySellerReport, error)
}

// SQLLocalityService - Implementation of LocalityService containing business logic for locality operations
// SQLLocalityService - Implementación de LocalityService que contiene la lógica de negocio para operaciones de localidades
type SQLLocalityService struct {
	repo repositories.LocalityRepository // Repository dependency for data access / Dependencia del repositorio para acceso a datos
}

// Save - Delegates saving a locality to the repository
// Save - Delega el guardado de una localidad al repositorio
func (s *SQLLocalityService) Save(ctx context.Context, locality models.Locality) (models.Locality, error) {
	return s.repo.Save(ctx, locality)
}

// GetSellerReports - Delegates retrieving seller reports to the repository
// GetSellerReports - Delega la obtención de reportes de vendedores al repositorio
func (s *SQLLocalityService) GetSellerReports(ctx context.Context, id int) ([]responses.LocalitySellerReport, error) {
	return s.repo.GetSellerReports(ctx, id)
}
