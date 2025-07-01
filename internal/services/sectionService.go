package services

import (
	"errors"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

func GetSectionService() SectionServiceI {
	return &sectionService{
		repository: repositories.GetSectionRepository(),
	}
}

type SectionServiceI interface {
	GetAll() []*models.Section
	GetByID(id int) (*models.Section, error)
	Create(model *models.Section) error
	DeleteByID(id int) error
}

type sectionService struct {
	repository repositories.SectionRepositoryI
}

func (s *sectionService) GetAll() []*models.Section {
	return s.repository.GetAll()
}

func (s *sectionService) GetByID(id int) (*models.Section, error) {
	if model := s.repository.GetByID(id); model != nil {
		return model, nil
	}
	return nil, error_message.ErrNotFound
}

func (s *sectionService) Create(model *models.Section) error {
	if s.repository.ExistsWithSectionNumber(model.SectionNumber) {
		return errors.New("already exist a section with the same number")
	}

	s.repository.Create(model)
	return nil
}

func (s *sectionService) DeleteByID(id int) error {
	if s.repository.DeleteByID(id) {
		return nil
	}
	return error_message.ErrNotFound
}
