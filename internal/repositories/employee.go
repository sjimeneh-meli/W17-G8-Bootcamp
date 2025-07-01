package repositories

import (
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
	"os"
)

func GetEmployeeRepository() EmployeeRepositoryI {
	jsonLoader := loader.NewJSONStorage[models.Employee](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), "employees.json"))
	storage, err := jsonLoader.ReadAll()
	if err != nil {
		fmt.Println("error loading employees:", err)
	}
	return &employeeRepository{
		storage: jsonLoader.MapToSlice(storage),
	}
}

type EmployeeRepositoryI interface {
	GetAll() []*models.Employee
	Create(e *models.Employee)
	GetById(id int) *models.Employee
	DeleteById(id int) bool
	ExistsWhCardNumber(id int, cardNumber string) bool
}

type employeeRepository struct {
	storage []*models.Employee
}

func (r *employeeRepository) GetAll() []*models.Employee {
	return r.storage
}
func (r *employeeRepository) GetById(id int) *models.Employee {
	for _, e := range r.storage {
		if e.Id == id {
			return e
		}
	}
	return nil
}

func (r *employeeRepository) Create(e *models.Employee) {
	e.Id = len(r.storage) + 1
	r.storage = append(r.storage, e)
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
	for i, e := range r.storage {
		if e.Id == id {
			r.storage = append(r.storage[:i], r.storage[i+1:]...)
			return true
		}
	}
	return false
}
