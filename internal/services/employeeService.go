package services

import (
	"context"
	"fmt"
	"slices"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var employeeServiceInstance EmployeeServiceI

func GetEmployeeService(repository repositories.EmployeeRepositoryI) EmployeeServiceI {
	if employeeServiceInstance != nil {
		return employeeServiceInstance
	}
	employeeServiceInstance = &EmployeeService{
		repository: repository,
	}
	return employeeServiceInstance
}

type EmployeeServiceI interface {
	GetAll(ctx context.Context) (map[int]models.Employee, error)
	GetById(ctx context.Context, id int) (models.Employee, error)
	DeleteById(ctx context.Context, id int) error
	Create(ctx context.Context, employee models.Employee) (models.Employee, error)
	Update(ctx context.Context, employeeId int, employee models.Employee) (models.Employee, error)
}
type EmployeeService struct {
	repository repositories.EmployeeRepositoryI
}

func (s *EmployeeService) GetAll(ctx context.Context) (map[int]models.Employee, error) {
	return s.repository.GetAll(ctx)

}
func (s *EmployeeService) GetById(ctx context.Context, id int) (models.Employee, error) {
	return s.repository.GetById(ctx, id)
}

func (s *EmployeeService) DeleteById(ctx context.Context, id int) error {
	return s.repository.DeleteById(ctx, id)
}

func (s *EmployeeService) Create(ctx context.Context, employee models.Employee) (models.Employee, error) {
	existingCardNumbers, err := s.repository.GetCardNumberIds()
	if err != nil {
		return models.Employee{}, err
	}
	if slices.Contains(existingCardNumbers, employee.CardNumberID) {
		return models.Employee{}, fmt.Errorf("%w - %s %s %s", error_message.ErrAlreadyExists, "card number with id:", employee.CardNumberID, "already exists.")
	}
	return s.repository.Create(ctx, employee)
}

func (s *EmployeeService) Update(ctx context.Context, employeeId int, employee models.Employee) (models.Employee, error) {
	existingCardNumbers, err := s.repository.GetCardNumberIds()
	if err != nil {
		return models.Employee{}, err
	}
	if slices.Contains(existingCardNumbers, employee.CardNumberID) {
		return models.Employee{}, fmt.Errorf("%w - %s %s %s", error_message.ErrAlreadyExists, "card number with id:", employee.CardNumberID, "already exists.")
	}
	return s.repository.Update(ctx, employeeId, employee)
}
