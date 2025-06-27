package services

import (
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
}

type sectionService struct {
	repository repositories.SectionRepositoryI
}

func (s *sectionService) GetAll() []*models.Section {
	return s.repository.GetAll()
}
