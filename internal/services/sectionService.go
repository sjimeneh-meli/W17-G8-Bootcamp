package services

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var sectionServiceInstance SectionServiceI

func GetSectionService(repository repositories.SectionRepositoryI) SectionServiceI {
	if sectionServiceInstance != nil {
		return sectionServiceInstance
	}

	sectionServiceInstance = &sectionService{
		repository: repository,
	}
	return sectionServiceInstance
}

type SectionServiceI interface {
	GetAll() ([]*models.Section, error)
	GetByID(id int) (*models.Section, error)
	Create(model *models.Section) error
	Update(model *models.Section) error
	DeleteByID(id int) error
	ExistWithID(id int) bool
	ExistsWithSectionNumber(id int, sectionNumber string) bool
}

type sectionService struct {
	repository repositories.SectionRepositoryI
}

func (s *sectionService) GetAll() ([]*models.Section, error) {
	return s.repository.GetAll()
}

func (s *sectionService) GetByID(id int) (*models.Section, error) {
	if model, err := s.repository.GetByID(id); err != nil {
		return model, error_message.ErrNotFound
	} else {
		return model, nil
	}

}

func (s *sectionService) Create(model *models.Section) error {
	return s.repository.Create(model)
}

func (s *sectionService) Update(model *models.Section) error {
	return s.repository.Update(model)
}

func (s *sectionService) DeleteByID(id int) error {
	if err := s.repository.DeleteByID(id); err != nil {
		return error_message.ErrNotFound
	}
	return nil
}

func (s *sectionService) ExistWithID(id int) bool {
	return s.repository.ExistWithID(id)
}

func (s *sectionService) ExistsWithSectionNumber(id int, sectionNumber string) bool {
	return s.repository.ExistsWithSectionNumber(id, sectionNumber)
}
