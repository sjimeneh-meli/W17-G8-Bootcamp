package services

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

type LocalityService interface {
	Save(locality models.Locality) (models.Locality, error)
	GetSellerReports(id int) ([]responses.LocalitySellerReport, error)
}

type SQLLocalityService struct {
	repo repositories.LocalityRepository
}

func NewSQLLocalityService(repo repositories.LocalityRepository) *SQLLocalityService {
	return &SQLLocalityService{repo: repo}
}

func (s *SQLLocalityService) Save(locality models.Locality) (models.Locality, error) {
	return s.repo.Save(locality)
}

func (s *SQLLocalityService) GetSellerReports(id int) ([]responses.LocalitySellerReport, error) {
	return s.repo.GetSellerReports(id)
}
