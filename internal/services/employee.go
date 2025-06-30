package services

import (
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
	Add(e *models.Employee) (*models.Employee, error)
}
type employeeService struct {
	repository repositories.EmployeeRepositoryI
}

func (s *employeeService) Add(e *models.Employee) (*models.Employee, error) {
	return s.repository.Add(e)
}

func (s *employeeService) GetAll() []*models.Employee {
	return s.repository.GetAll()
}
