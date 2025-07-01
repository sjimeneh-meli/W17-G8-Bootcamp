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
	Create(model *models.Section)
	ExistsWithSectionNumber(id int, sectionNumber string) bool
	DeleteByID(id int) bool
}

type sectionRepository struct {
	storage []*models.Section
}

func (r *sectionRepository) GetAll() []*models.Section {
	return r.storage
}

func (r *sectionRepository) GetByID(id int) *models.Section {
	for _, m := range r.storage {
		if m.Id == id {
			return m
		}
	}
	return nil
}

func (r *sectionRepository) Create(model *models.Section) {
	model.Id = len(r.storage) + 1
	r.storage = append(r.storage, model)
}

func (r *sectionRepository) ExistsWithSectionNumber(id int, sectionNumber string) bool {
	for _, m := range r.storage {
		if m.SectionNumber == sectionNumber && m.Id != id {
			return true
		}
	}
	return false
}

func (r *sectionRepository) DeleteByID(id int) bool {
	for i, m := range r.storage {
		if m.Id == id {
			r.storage = append(r.storage[:i], r.storage[i+1:]...)
			return true
		}
	}
	return false
}
