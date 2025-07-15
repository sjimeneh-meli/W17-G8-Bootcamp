package services

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

type LocalityService interface {
	Save(ctx context.Context, locality models.Locality) (models.Locality, error)
	GetSellerReports(ctx context.Context, id int) ([]responses.LocalitySellerReport, error)
}

type SQLLocalityService struct {
	repo repositories.LocalityRepository
}

func NewSQLLocalityService(repo repositories.LocalityRepository) *SQLLocalityService {
	return &SQLLocalityService{repo: repo}
}

func (s *SQLLocalityService) Save(ctx context.Context, locality models.Locality) (models.Locality, error) {
	return s.repo.Save(ctx, locality)
}

func (s *SQLLocalityService) GetSellerReports(ctx context.Context, id int) ([]responses.LocalitySellerReport, error) {
	return s.repo.GetSellerReports(ctx, id)
}
