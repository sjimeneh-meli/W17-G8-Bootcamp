package repositories

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

func GetEmployeeRepository(loader *loader.StorageJSON[models.Employee]) (EmployeeRepositoryI, error) {
	storage, err := loader.ReadAll()
	if err != nil {
		return nil, err
	}

	return &employeeRepository{
		storage: storage,
		loader:  loader,
	}, nil
}

type EmployeeRepositoryI interface {
	GetAll() []*models.Employee
	Create(e *models.Employee)
	GetById(id int) *models.Employee
	DeleteById(id int) bool
	ExistsWhCardNumber(id int, cardNumber string) bool
}

type employeeRepository struct {
	storage map[int]models.Employee
	loader  *loader.StorageJSON[models.Employee]
}

func (r *employeeRepository) GetAll() []*models.Employee {
	var list []*models.Employee
	for _, m := range r.loader.MapToSlice(r.storage) {
		list = append(list, &m)
	}

	return list
}

func (r *employeeRepository) GetById(id int) *models.Employee {
	for _, m := range r.storage {
		if m.Id == id {
			return &m
		}
	}
	return nil
}

func (r *employeeRepository) Create(e *models.Employee) {
	e.Id = len(r.storage) + 1
	r.storage[e.Id] = *e
}

func (r *employeeRepository) ExistsWhCardNumber(id int, cardNumber string) bool {
	for _, e := range r.storage {
		if e.CardNumberID == cardNumber && e.Id != id {
			return true
		}
	}
	return false
}

func (r *employeeRepository) DeleteById(id int) bool {
	for i, m := range r.storage {
		if m.Id == id {
			delete(r.storage, i)
			return true
		}
	}
	return false
}
