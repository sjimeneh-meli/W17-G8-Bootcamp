package repositories

import (
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
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
	GetAll() []models.Employee
	Create(e models.Employee)
	GetById(id int) (models.Employee, error)
	DeleteById(id int) bool
	Update(e models.Employee) bool
	ExistsWhCardNumber(id int, cardNumber string) bool
}

type employeeRepository struct {
	storage map[int]models.Employee
	loader  *loader.StorageJSON[models.Employee]
}

func (r *employeeRepository) GetAll() []models.Employee {
	return r.loader.MapToSlice(r.storage)
}

func (r *employeeRepository) GetById(id int) (models.Employee, error) {
	_, exists := r.storage[id]
	if !exists {
		return models.Employee{}, fmt.Errorf("%w. %s %d", error_message.ErrNotFound, "employee with Id", id)
	}
	return r.storage[id], nil
}

func (r *employeeRepository) Create(e models.Employee) {
	e.Id = len(r.storage) + 1
	r.storage[e.Id] = e
}

func (r *employeeRepository) ExistsWhCardNumber(id int, cardNumber string) bool {
	for _, e := range r.storage {
		if e.CardNumberID == cardNumber && e.Id != id {
			return true
		}
	}
	return false
}

func (r *employeeRepository) Update(e models.Employee) bool {
	if _, exists := r.storage[e.Id]; !exists {
		return false
	}
	r.storage[e.Id] = e
	return true
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
