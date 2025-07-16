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

func NewCarryService(r repositories.CarryRepository, lr repositories.LocalityRepository) CarryService {
	if carryServiceInstance != nil {
		return carryServiceInstance
	}
	carryServiceInstance = &CarryServiceImpl{carryRepository: r, localityRepository: lr}
	return carryServiceInstance
}

type CarryService interface {
	CreateCarry(ctx context.Context, carry models.Carry) (models.Carry, error)
	GetCarryReportByLocality(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error)
}

type CarryServiceImpl struct {
	carryRepository    repositories.CarryRepository
	localityRepository repositories.LocalityRepository
}

// Falta validar que locality id exista en la base de datos

func (s *CarryServiceImpl) CreateCarry(ctx context.Context, carry models.Carry) (models.Carry, error) {
	// Validar que locality existe
	localityExists, err := s.localityRepository.ExistById(ctx, carry.LocalityId)
	if err != nil {
		return models.Carry{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	if !localityExists {
		return models.Carry{}, fmt.Errorf("%w: locality with id %d", error_message.ErrNotFound, carry.LocalityId)
	}

	// Validar que CID sea único
	exists, err := s.carryRepository.ExistsByCid(ctx, carry.Cid)
	if err != nil {
		return models.Carry{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	if exists {
		return models.Carry{}, fmt.Errorf("%w: resource with the provided identifier already exists", error_message.ErrAlreadyExists)
	}

	// Crear carry
	carry, err = s.carryRepository.Create(ctx, carry)
	if err != nil {
		return models.Carry{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	return carry, nil
}

func (s *CarryServiceImpl) GetCarryReportByLocality(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
	// Solo validar existencia si se especifica un localityID específico
	if localityID != 0 {
		exists, err := s.localityRepository.ExistById(ctx, localityID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
		}
		if !exists {
			return nil, fmt.Errorf("%w: locality with id %d", error_message.ErrNotFound, localityID)
		}
	}

	// Obtener reporte (si localityID = 0, traerá todas las localidades)
	reports, err := s.carryRepository.GetCarryReportsByLocality(ctx, localityID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	return reports, nil
}
