package repositories

import (
	"context"
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
	"slices"
)

type EmployeeRepositoryI interface {
	GetAll(ctx context.Context) (map[int]models.Employee, error)
	GetById(ctx context.Context, id int) (models.Employee, error)
	DeleteById(ctx context.Context, id int) error
	Create(ctx context.Context, employee models.Employee) (models.Employee, error)
	Update(ctx context.Context, employeeId int, employee models.Employee) (models.Employee, error)
	GetCardNumberIds() ([]string, error)
	ExistEmployeeById(ctx context.Context, employeeId int) (bool, error)
}

type EmployeeRepository struct {
	storage map[int]models.Employee
	loader  loader.Storage[models.Employee]
}

func GetNewEmployeeRepository(loader loader.Storage[models.Employee]) (EmployeeRepositoryI, error) {
	storage, err := loader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w:%v", error_message.ErrInternalServerError, err)
	}

	return &EmployeeRepository{
		storage: storage,
		loader:  loader,
	}, nil
}

func (r *EmployeeRepository) GetAll(ctx context.Context) (map[int]models.Employee, error) {
	return r.storage, nil
}
func (r *EmployeeRepository) GetById(ctx context.Context, id int) (models.Employee, error) {
	_, exists := r.storage[id]
	if !exists {
		return models.Employee{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "employee with Id", id, "not exists")
	}

	return r.storage[id], nil
}
func (r *EmployeeRepository) Create(ctx context.Context, employee models.Employee) (models.Employee, error) {
	newId := r.GetNewId()
	employee.Id = newId

	_, exists := r.storage[employee.Id]
	if exists {
		return models.Employee{}, fmt.Errorf("%w. %s %d %s", error_message.ErrAlreadyExists, "employee with Id", employee.Id, "already exists.")
	}

	r.storage[employee.Id] = employee

	err := r.Save()
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err.Error())
	}

	return employee, nil
}
func (r *EmployeeRepository) Update(ctx context.Context, employeeId int, employee models.Employee) (models.Employee, error) {
	_, exists := r.storage[employeeId]
	if !exists {
		return models.Employee{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "employee with Id", employeeId, "not exists.")
	}

	updatedEmployee := r.storage[employeeId]
	if employee.CardNumberID != "" {
		existingCardNumberIds, _ := r.GetCardNumberIds()
		if slices.Contains(existingCardNumberIds, employee.CardNumberID) {
			return models.Employee{}, fmt.Errorf("%w. %s %s %s", error_message.ErrAlreadyExists, "Card number with Id", employee.CardNumberID, "already exists.")
		}
		updatedEmployee.CardNumberID = employee.CardNumberID
	}
	if employee.FirstName != "" {
		updatedEmployee.FirstName = employee.FirstName
	}
	if employee.LastName != "" {
		updatedEmployee.LastName = employee.LastName
	}

	r.storage[employeeId] = updatedEmployee
	err := r.Save()
	if err != nil {
		return models.Employee{}, err
	}

	return updatedEmployee, nil
}
func (r *EmployeeRepository) DeleteById(ctx context.Context, id int) error {
	_, exists := r.storage[id]
	if !exists {
		return fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", id, "doesn't exists.")
	}
	delete(r.storage, id)

	err := r.Save()
	if err != nil {
		return err
	}
	return nil
}
func (r *EmployeeRepository) GetCardNumberIds() ([]string, error) {
	cardNumbers := []string{}
	for _, employee := range r.storage {
		cardNumbers = append(cardNumbers, employee.CardNumberID)
	}
	return cardNumbers, nil
}
func (r *EmployeeRepository) ExistEmployeeById(ctx context.Context, employeeId int) (bool, error) {
	_, exists := r.storage[employeeId]
	if exists {
		return true, nil
	}
	return false, nil
}
func (r *EmployeeRepository) Save() error {
	employeeaArray := []models.Employee{}

	for _, employee := range r.storage {
		employeeaArray = append(employeeaArray, employee)
	}

	err := r.loader.WriteAll(employeeaArray)
	if err != nil {
		return fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err.Error())
	}
	return nil
}
func (r *EmployeeRepository) GetNewId() int {
	lastId := 0
	for _, employee := range r.storage {
		if employee.Id > lastId {
			lastId = employee.Id
		}
	}
	return (lastId + 1)
}
