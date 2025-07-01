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
	Create(e *models.Employee) (*models.Employee, error)
	GetById(id int) *models.Employee
	DeleteById(id int) bool
}
type employeeService struct {
	repository repositories.EmployeeRepositoryI
}

func (s *employeeService) GetAll() []*models.Employee {
	return s.repository.GetAll()
}

func (s *employeeService) GetById(id int) *models.Employee {
	return s.repository.GetById(id)
}

func (s *employeeService) Create(e *models.Employee) (*models.Employee, error) {
	return s.repository.Create(e)
}

func (s *employeeService) DeleteById(id int) bool {
	return s.repository.DeleteById(id)
}
