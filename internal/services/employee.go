package services

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

func GetEmployeeService() EmployeeServiceI {
	return &employeeService{
		repository: repositories.GetEmployeeRepository(),
	}
}

type EmployeeServiceI interface {
	GetAll() []*models.Employee
	Create(e *models.Employee) error
	GetById(id int) (*models.Employee, error)
	DeleteById(id int) error
	ExistsWhCardNumber(id int, cardNumber string) bool
}
type employeeService struct {
	repository repositories.EmployeeRepositoryI
}

func (s employeeService) GetAll() []*models.Employee {
	return s.repository.GetAll()
}

func (s employeeService) Create(e *models.Employee) error {
	s.repository.Create(e)
	return nil
}

func (s employeeService) GetById(id int) (*models.Employee, error) {
	if model := s.repository.GetById(id); model != nil {
		return model, nil
	}
	return nil, error_message.ErrNotFound
}

func (s employeeService) DeleteById(id int) error {
	if s.repository.DeleteById(id) {
		return nil
	}
	return error_message.ErrNotFound
}

func (s employeeService) ExistsWhCardNumber(id int, cardNumber string) bool {
	return s.repository.ExistsWhCardNumber(id, cardNumber)

}
