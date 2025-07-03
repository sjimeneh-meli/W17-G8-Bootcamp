package services

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

func GetEmployeeService(repository repositories.EmployeeRepositoryI) EmployeeServiceI {
	return &employeeService{
		repository: repository,
	}
}

type EmployeeServiceI interface {
	GetAll() []models.Employee
	Create(e models.Employee) error
	GetById(id int) (models.Employee, error)
	DeleteById(id int) error
	ExistsWhCardNumber(id int, cardNumber string) bool
	Update(e models.Employee) error
}
type employeeService struct {
	repository repositories.EmployeeRepositoryI
}

func (s employeeService) GetAll() []models.Employee {
	return s.repository.GetAll()
}

func (s employeeService) Create(e models.Employee) error {
	s.repository.Create(e)
	return nil
}

func (s employeeService) GetById(id int) (models.Employee, error) {
	return s.repository.GetById(id)
}

func (s employeeService) DeleteById(id int) error {
	if s.repository.DeleteById(id) {
		return nil
	}
	return error_message.ErrNotFound
}

func (s employeeService) Update(e models.Employee) error {
	if !s.repository.Update(e) {
		return error_message.ErrNotFound
	}
	return nil
}

func (s employeeService) ExistsWhCardNumber(id int, cardNumber string) bool {
	return s.repository.ExistsWhCardNumber(id, cardNumber)

}
