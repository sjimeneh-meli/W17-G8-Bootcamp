package repositories

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func GetServiceRepository() ServiceRepositoryI {
	return &ServiceRepository{}
}

type ServiceRepositoryI interface{}

type ServiceRepository struct {
	Storage *models.Section
}
