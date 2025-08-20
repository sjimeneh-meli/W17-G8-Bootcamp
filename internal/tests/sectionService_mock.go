package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

type SectionServiceMock struct {
	mock.Mock
}

func GetSectionServiceMock() *SectionServiceMock {
	return &SectionServiceMock{}
}

func (r *SectionServiceMock) GetAll(ctx context.Context) ([]*models.Section, error) {
	args := r.Called(ctx)
	return args.Get(0).([]*models.Section), args.Error(1)
}

func (r *SectionServiceMock) GetByID(ctx context.Context, id int) (*models.Section, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*models.Section), args.Error(1)
}

func (r *SectionServiceMock) Create(ctx context.Context, section *models.Section) error {
	args := r.Called(ctx, section)
	section.Id = 1
	return args.Error(0)
}

func (r *SectionServiceMock) Update(ctx context.Context, section *models.Section) error {
	args := r.Called(ctx, section)
	return args.Error(0)
}

func (r *SectionServiceMock) DeleteByID(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func (r *SectionServiceMock) ExistWithID(ctx context.Context, id int) bool {
	r.Called(ctx, id)
	return id == 5
}

func (r *SectionServiceMock) ExistsWithSectionNumber(ctx context.Context, id int, sectionNumber string) bool {
	r.Called(ctx, id, sectionNumber)
	return sectionNumber == "B-01"
}
