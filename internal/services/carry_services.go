package services

import (
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

type CarryService interface {
	CreateCarry(carry models.Carry) (models.Carry, error)
}

type CarryServiceImpl struct {
	carryRepository repositories.CarryRepository
}

func NewCarryService(r repositories.CarryRepository) *CarryServiceImpl {
	return &CarryServiceImpl{carryRepository: r}
}

// Falta validar que locality id exista en la base de datos

func (s *CarryServiceImpl) CreateCarry(carry models.Carry) (models.Carry, error) {
	exists, err := s.carryRepository.ExistsByCid(carry.Cid)
	if err != nil {
		return models.Carry{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	if exists {
		return models.Carry{}, fmt.Errorf("%w: resource with the provided identifier already exists", error_message.ErrAlreadyExists)
	}
	carry, err = s.carryRepository.Create(carry)
	if err != nil {
		return models.Carry{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	return carry, nil
}
