package services

import (
	"context"

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
	GetAll(ctx context.Context) ([]*models.Section, error)
	GetByID(ctx context.Context, id int) (*models.Section, error)
	Create(ctx context.Context, model *models.Section) error
	Update(ctx context.Context, model *models.Section) error
	DeleteByID(ctx context.Context, id int) error
	ExistWithID(ctx context.Context, id int) bool
	ExistsWithSectionNumber(ctx context.Context, id int, sectionNumber string) bool
}

type sectionService struct {
	repository repositories.SectionRepositoryI
}

func (s *sectionService) GetAll(ctx context.Context) ([]*models.Section, error) {
	return s.repository.GetAll(ctx)
}

func (s *sectionService) GetByID(ctx context.Context, id int) (*models.Section, error) {
	if model, err := s.repository.GetByID(ctx, id); err != nil {
		return model, error_message.ErrNotFound
	} else {
		return model, nil
	}

}

func (s *sectionService) Create(ctx context.Context, model *models.Section) error {
	return s.repository.Create(ctx, model)
}

func (s *sectionService) Update(ctx context.Context, model *models.Section) error {
	return s.repository.Update(ctx, model)
}

func (s *sectionService) DeleteByID(ctx context.Context, id int) error {
	if err := s.repository.DeleteByID(ctx, id); err != nil {
		return error_message.ErrNotFound
	}
	return nil
}

func (s *sectionService) ExistWithID(ctx context.Context, id int) bool {
	return s.repository.ExistWithID(ctx, id)
}

func (s *sectionService) ExistsWithSectionNumber(ctx context.Context, id int, sectionNumber string) bool {
	return s.repository.ExistsWithSectionNumber(ctx, id, sectionNumber)
}
