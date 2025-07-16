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

// GetEmployeeService - Creates and returns a new instance of EmployeeService with the required repository using singleton pattern
// GetEmployeeService - Crea y retorna una nueva instancia de EmployeeService con el repositorio requerido usando patrón singleton
func GetEmployeeService(repository repositories.EmployeeRepositoryI) EmployeeServiceI {
	if employeeServiceInstance != nil {
		return employeeServiceInstance
	}
	employeeServiceInstance = &EmployeeService{
		repository: repository,
	}
	return employeeServiceInstance
}

// EmployeeServiceI - Interface defining the contract for employee service operations with business logic
// EmployeeServiceI - Interfaz que define el contrato para las operaciones del servicio de empleados con lógica de negocio
type EmployeeServiceI interface {
	// GetAll - Retrieves all employees from the system
	// GetAll - Obtiene todos los empleados del sistema
	GetAll(ctx context.Context) (map[int]models.Employee, error)

	// GetById - Retrieves a specific employee by their ID
	// GetById - Obtiene un empleado específico por su ID
	GetById(ctx context.Context, id int) (models.Employee, error)

	// DeleteById - Removes an employee from the system by their ID
	// DeleteById - Elimina un empleado del sistema por su ID
	DeleteById(ctx context.Context, id int) error

	// Create - Creates a new employee with business validation (card number uniqueness)
	// Create - Crea un nuevo empleado con validación de negocio (unicidad del número de tarjeta)
	Create(ctx context.Context, employee models.Employee) (models.Employee, error)

	// Update - Updates an existing employee with business validation (card number uniqueness)
	// Update - Actualiza un empleado existente con validación de negocio (unicidad del número de tarjeta)
	Update(ctx context.Context, employeeId int, employee models.Employee) (models.Employee, error)
}

// EmployeeService - Implementation of EmployeeServiceI containing business logic for employee operations
// EmployeeService - Implementación de EmployeeServiceI que contiene la lógica de negocio para operaciones de empleados
type EmployeeService struct {
	repository repositories.EmployeeRepositoryI // Repository dependency for data access / Dependencia del repositorio para acceso a datos
}

// GetAll - Delegates retrieving all employees to the repository
// GetAll - Delega la obtención de todos los empleados al repositorio
func (s *EmployeeService) GetAll(ctx context.Context) (map[int]models.Employee, error) {
	return s.repository.GetAll(ctx)

}

// GetById - Delegates retrieving an employee by their ID to the repository
// GetById - Delega la obtención de un empleado por su ID al repositorio
func (s *EmployeeService) GetById(ctx context.Context, id int) (models.Employee, error) {
	return s.repository.GetById(ctx, id)
}

// DeleteById - Delegates removing an employee from the repository by their ID
// DeleteById - Delega la eliminación de un empleado del repositorio por su ID
func (s *EmployeeService) DeleteById(ctx context.Context, id int) error {
	return s.repository.DeleteById(ctx, id)
}

// Create - Creates a new employee with business validation to ensure card number uniqueness
// Create - Crea un nuevo empleado con validación de negocio para asegurar la unicidad del número de tarjeta
func (s *EmployeeService) Create(ctx context.Context, employee models.Employee) (models.Employee, error) {
	// Business validation: Get all existing card numbers to check for duplicates
	// Validación de negocio: Obtener todos los números de tarjeta existentes para verificar duplicados
	existingCardNumbers, err := s.repository.GetCardNumberIds()
	if err != nil {
		return models.Employee{}, err
	}

	// Business rule: Card number must be unique across all employees
	// Regla de negocio: El número de tarjeta debe ser único entre todos los empleados
	if slices.Contains(existingCardNumbers, employee.CardNumberID) {
		return models.Employee{}, fmt.Errorf("%w - %s %s %s", error_message.ErrAlreadyExists, "card number with id:", employee.CardNumberID, "already exists.")
	}

	// If validation passes, delegate to repository for persistence
	// Si la validación pasa, delegar al repositorio para la persistencia
	return s.repository.Create(ctx, employee)
}

// Update - Updates an existing employee with business validation to ensure card number uniqueness
// Update - Actualiza un empleado existente con validación de negocio para asegurar la unicidad del número de tarjeta
func (s *EmployeeService) Update(ctx context.Context, employeeId int, employee models.Employee) (models.Employee, error) {
	// Business validation: Get all existing card numbers to check for duplicates
	// Validación de negocio: Obtener todos los números de tarjeta existentes para verificar duplicados
	existingCardNumbers, err := s.repository.GetCardNumberIds()
	if err != nil {
		return models.Employee{}, err
	}

	// Business rule: Card number must be unique across all employees (excluding current employee)
	// Regla de negocio: El número de tarjeta debe ser único entre todos los empleados (excluyendo el empleado actual)
	if slices.Contains(existingCardNumbers, employee.CardNumberID) {
		return models.Employee{}, fmt.Errorf("%w - %s %s %s", error_message.ErrAlreadyExists, "card number with id:", employee.CardNumberID, "already exists.")
	}

	// If validation passes, delegate to repository for persistence
	// Si la validación pasa, delegar al repositorio para la persistencia
	return s.repository.Update(ctx, employeeId, employee)
}
