package services

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var sectionServiceInstance SectionServiceI

// GetSectionService creates and returns a singleton instance of sectionService with the required repository
// GetSectionService crea y retorna una instancia singleton de sectionService con el repositorio requerido
func GetSectionService(repository repositories.SectionRepositoryI) SectionServiceI {
	if sectionServiceInstance != nil {
		return sectionServiceInstance
	}

	sectionServiceInstance = &sectionService{
		repository: repository,
	}
	return sectionServiceInstance
}

// SectionServiceI defines the contract for section service operations with business logic and validation
// SectionServiceI define el contrato para las operaciones de servicio de secciones con lógica de negocio y validación
type SectionServiceI interface {
	GetAll(ctx context.Context) ([]*models.Section, error)
	GetByID(ctx context.Context, id int) (*models.Section, error)
	Create(ctx context.Context, model *models.Section) error
	Update(ctx context.Context, model *models.Section) error
	DeleteByID(ctx context.Context, id int) error
	ExistWithID(ctx context.Context, id int) bool
	ExistsWithSectionNumber(ctx context.Context, id int, sectionNumber string) bool
}

// sectionService implements SectionServiceI and contains business logic for section operations
// sectionService implementa SectionServiceI y contiene la lógica de negocio para operaciones de secciones
type sectionService struct {
	repository repositories.SectionRepositoryI // Repository for section data access / Repositorio para acceso a datos de secciones
}

// GetAll retrieves all sections from the repository
// GetAll recupera todas las secciones del repositorio
func (s *sectionService) GetAll(ctx context.Context) ([]*models.Section, error) {
	return s.repository.GetAll(ctx)
}

// GetByID retrieves a section by its ID with error handling for non-existent sections
// GetByID recupera una sección por su ID con manejo de errores para secciones no existentes
func (s *sectionService) GetByID(ctx context.Context, id int) (*models.Section, error) {
	if model, err := s.repository.GetByID(ctx, id); err != nil {
		return model, error_message.ErrNotFound
	} else {
		return model, nil
	}
}

// Create creates a new section in the repository
// Create crea una nueva sección en el repositorio
func (s *sectionService) Create(ctx context.Context, model *models.Section) error {
	return s.repository.Create(ctx, model)
}

// Update modifies an existing section in the repository
// Update modifica una sección existente en el repositorio
func (s *sectionService) Update(ctx context.Context, model *models.Section) error {
	return s.repository.Update(ctx, model)
}

// DeleteByID removes a section by its ID with error handling for non-existent sections
// DeleteByID elimina una sección por su ID con manejo de errores para secciones no existentes
func (s *sectionService) DeleteByID(ctx context.Context, id int) error {
	if err := s.repository.DeleteByID(ctx, id); err != nil {
		return error_message.ErrNotFound
	}
	return nil
}

// ExistWithID checks if a section exists by its ID
// ExistWithID verifica si una sección existe por su ID
func (s *sectionService) ExistWithID(ctx context.Context, id int) bool {
	return s.repository.ExistWithID(ctx, id)
}

// ExistsWithSectionNumber checks if a section number exists, excluding a specific ID for updates
// ExistsWithSectionNumber verifica si un número de sección existe, excluyendo un ID específico para actualizaciones
func (s *sectionService) ExistsWithSectionNumber(ctx context.Context, id int, sectionNumber string) bool {
	return s.repository.ExistsWithSectionNumber(ctx, id, sectionNumber)
}
