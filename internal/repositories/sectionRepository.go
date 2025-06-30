package repositories

import (
	"fmt"
	"os"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

func GetSectionRepository() SectionRepositoryI {
	jsonLoader := loader.NewJSONStorage[models.Section](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "sections.json"))
	storage, err := jsonLoader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	return &sectionRepository{
		storage: jsonLoader.MapToSlice(storage),
	}
}

type SectionRepositoryI interface {
	GetAll() []*models.Section
	GetByID(id int) *models.Section
}

type sectionRepository struct {
	storage []*models.Section
}

func (r sectionRepository) GetAll() []*models.Section {
	return r.storage
}

func (r sectionRepository) GetByID(id int) *models.Section {
	for _, m := range r.storage {
		if m.Id == id {
			return m
		}
	}
	return nil
}
