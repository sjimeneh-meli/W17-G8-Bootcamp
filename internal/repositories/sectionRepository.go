package repositories

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func GetSectionRepository() SectionRepositoryI {
	return &sectionRepository{}
}

type SectionRepositoryI interface{}

type sectionRepository struct {
	Storage *models.Section
}
