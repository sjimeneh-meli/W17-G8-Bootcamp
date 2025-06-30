package services

import (
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
