package repositories

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

func GetSectionRepository(loader *loader.StorageJSON[models.Section]) (SectionRepositoryI, error) {
	storage, err := loader.ReadAll()
	if err != nil {
		return nil, err
	}

	return &sectionRepository{
		storage: storage,
		loader:  loader,
	}, nil
}

type SectionRepositoryI interface {
	GetAll() []*models.Section
	GetByID(id int) *models.Section
	Create(model *models.Section)
	ExistsWithSectionNumber(id int, sectionNumber string) bool
	DeleteByID(id int) bool
}

type sectionRepository struct {
	storage map[int]models.Section
	loader  *loader.StorageJSON[models.Section]
}

func (r *sectionRepository) GetAll() []*models.Section {
	var list []*models.Section
	for _, m := range r.loader.MapToSlice(r.storage) {
		list = append(list, &m)
	}

	return list
}

func (r *sectionRepository) GetByID(id int) *models.Section {
	for _, m := range r.storage {
		if m.Id == id {
			return &m
		}
	}
	return nil
}

func (r *sectionRepository) Create(model *models.Section) {
	model.Id = len(r.storage) + 1
	r.storage[model.Id] = *model
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
			delete(r.storage, i)
			return true
		}
	}
	return false
}
