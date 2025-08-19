// Package services - Employee Service Unit Tests / Tests Unitarios del Servicio Employee
// Business logic testing for employee management with card number ID uniqueness validation
// Testing de lógica de negocio para gestión de empleados con validación de unicidad de ID de tarjeta
package services

import (
	"context"
	"errors"
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestEmployeeService_Create - Tests for Create method of EmployeeService
// TestEmployeeService_Create - Tests del método Create de EmployeeService
// Validates business logic for creation including card number ID uniqueness
// Valida lógica de negocio de creación incluyendo unicidad de card number ID
// Cases: successful creation, duplicate card number, repository error
// Casos: creación exitosa, card number duplicado, error de repositorio
func TestEmployeeService_Create(t *testing.T) {
	// create_ok: Valid data with unique card number creates employee successfully
	// create_ok: Datos válidos con card number único crea empleado exitosamente
	t.Run("create_ok", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		employee := models.Employee{
			CardNumberID: "CARD-001",
			FirstName:    "Carlos",
			LastName:     "Ruiz",
			WarehouseID:  1,
		}

		expectedEmployee := models.Employee{
			Id:           1,
			CardNumberID: "CARD-001",
			FirstName:    "Carlos",
			LastName:     "Ruiz",
			WarehouseID:  1,
		}

		repositoryMock.On("GetCardNumberIds").Return([]string{}, nil)
		repositoryMock.On("Create", mock.Anything, employee).Return(expectedEmployee, nil)

		result, err := service.Create(context.Background(), employee)

		assert.NoError(t, err)
		assert.Equal(t, expectedEmployee, result)
		repositoryMock.AssertExpectations(t)
	})

	// create_conflict: Duplicate card number returns conflict error
	// create_conflict: Card number duplicado retorna error de conflicto
	t.Run("create_conflict", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		employee := models.Employee{
			CardNumberID: "CARD-001",
			FirstName:    "Carlos",
			LastName:     "Ruiz",
			WarehouseID:  1,
		}

		repositoryMock.On("GetCardNumberIds").Return([]string{"CARD-001"}, nil)

		result, err := service.Create(context.Background(), employee)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrAlreadyExists))
		assert.Equal(t, models.Employee{}, result)
		repositoryMock.AssertExpectations(t)
	})

	// create_get_card_numbers_error: Error fetching existing card numbers during creation
	// create_get_card_numbers_error: Error al obtener card numbers existentes durante creación
	t.Run("create_get_card_numbers_error", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		employee := models.Employee{
			CardNumberID: "CARD-001",
			FirstName:    "Carlos",
			LastName:     "Ruiz",
			WarehouseID:  1,
		}

		repositoryMock.On("GetCardNumberIds").Return([]string{}, error_message.ErrInternalServerError)

		result, err := service.Create(context.Background(), employee)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
		assert.Equal(t, models.Employee{}, result)
		repositoryMock.AssertExpectations(t)
	})
}

// TestEmployeeService_GetAll - Tests for GetAll method of EmployeeService
// TestEmployeeService_GetAll - Tests del método GetAll de EmployeeService
// Validates retrieval of all employees from repository
// Valida obtención de todos los empleados del repositorio
// Cases: successful retrieval with data
// Casos: recuperación exitosa con datos
func TestEmployeeService_GetAll(t *testing.T) {
	// find_all: Repository returns data and service passes it without modification
	// find_all: Repositorio retorna datos y servicio los pasa sin modificar
	t.Run("find_all", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		expectedEmployees := map[int]models.Employee{
			1: {Id: 1, CardNumberID: "CARD-001", FirstName: "Carlos", LastName: "Ruiz", WarehouseID: 1},
			2: {Id: 2, CardNumberID: "CARD-002", FirstName: "Maria", LastName: "Lopez", WarehouseID: 2},
		}

		repositoryMock.On("GetAll", mock.Anything).Return(expectedEmployees, nil)

		result, err := service.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedEmployees, result)
		repositoryMock.AssertExpectations(t)
	})
}

// TestEmployeeService_GetById - Tests for GetById method of EmployeeService
// TestEmployeeService_GetById - Tests del método GetById de EmployeeService
// Covers successful search and employee not found scenarios
// Cubre búsqueda exitosa y empleado no encontrado
// Cases: existing employee, non-existent employee
// Casos: empleado existente, empleado inexistente
func TestEmployeeService_GetById(t *testing.T) {
	// find_by_id_existent: Employee exists, returns employee data
	// find_by_id_existent: Empleado existe, retorna datos del empleado
	t.Run("find_by_id_existent", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		expectedEmployee := models.Employee{
			Id:           1,
			CardNumberID: "CARD-001",
			FirstName:    "Carlos",
			LastName:     "Ruiz",
			WarehouseID:  1,
		}

		repositoryMock.On("GetById", mock.Anything, 1).Return(expectedEmployee, nil)

		result, err := service.GetById(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedEmployee, result)
		repositoryMock.AssertExpectations(t)
	})

	// find_by_id_non_existent: Employee doesn't exist, returns not found error
	// find_by_id_non_existent: Empleado no existe, retorna error not found
	t.Run("find_by_id_non_existent", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		repositoryMock.On("GetById", mock.Anything, 999).Return(models.Employee{}, error_message.ErrNotFound)

		result, err := service.GetById(context.Background(), 999)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrNotFound))
		assert.Equal(t, models.Employee{}, result)
		repositoryMock.AssertExpectations(t)
	})
}

// TestEmployeeService_Update - Tests for Update method of EmployeeService
// TestEmployeeService_Update - Tests del método Update de EmployeeService
// Validates update with card number ID uniqueness verification
// Valida actualización con verificación de unicidad de card number ID
// Cases: successful update, non-existent employee, repository error
// Casos: actualización exitosa, empleado inexistente, error de repositorio
func TestEmployeeService_Update(t *testing.T) {
	// update_existent: Employee exists, updates and returns updated employee
	// update_existent: Empleado existe, actualiza y retorna empleado actualizado
	t.Run("update_existent", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		updateEmployee := models.Employee{
			CardNumberID: "CARD-002",
			FirstName:    "Updated Name",
			LastName:     "Updated LastName",
		}

		expectedEmployee := models.Employee{
			Id:           1,
			CardNumberID: "CARD-002",
			FirstName:    "Updated Name",
			LastName:     "Updated LastName",
			WarehouseID:  1,
		}

		repositoryMock.On("GetCardNumberIds").Return([]string{}, nil)
		repositoryMock.On("Update", mock.Anything, 1, updateEmployee).Return(expectedEmployee, nil)

		result, err := service.Update(context.Background(), 1, updateEmployee)

		assert.NoError(t, err)
		assert.Equal(t, expectedEmployee, result)
		repositoryMock.AssertExpectations(t)
	})

	// update_non_existent: Employee doesn't exist, returns not found error
	// update_non_existent: Empleado no existe, retorna error not found
	t.Run("update_non_existent", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		updateEmployee := models.Employee{
			FirstName: "Updated Name",
		}

		repositoryMock.On("GetCardNumberIds").Return([]string{}, nil)
		repositoryMock.On("Update", mock.Anything, 999, updateEmployee).Return(models.Employee{}, error_message.ErrNotFound)

		result, err := service.Update(context.Background(), 999, updateEmployee)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrNotFound))
		assert.Equal(t, models.Employee{}, result)
		repositoryMock.AssertExpectations(t)
	})

	// update_get_card_numbers_error: Error fetching existing card numbers
	// update_get_card_numbers_error: Error al obtener card numbers existentes
	t.Run("update_get_card_numbers_error", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		updateEmployee := models.Employee{
			FirstName: "Updated Name",
		}

		repositoryMock.On("GetCardNumberIds").Return([]string{}, error_message.ErrInternalServerError)

		result, err := service.Update(context.Background(), 1, updateEmployee)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
		assert.Equal(t, models.Employee{}, result)
		repositoryMock.AssertExpectations(t)
	})

	// update_conflict: Card number already exists, returns conflict error
	// update_conflict: Card number ya existe, retorna error de conflicto
	t.Run("update_conflict", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		updateEmployee := models.Employee{
			CardNumberID: "CARD-EXISTING",
			FirstName:    "Updated Name",
		}

		// Mock card numbers that include the one we're trying to use
		repositoryMock.On("GetCardNumberIds").Return([]string{"CARD-001", "CARD-EXISTING", "CARD-003"}, nil)

		result, err := service.Update(context.Background(), 1, updateEmployee)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrAlreadyExists))
		assert.Equal(t, models.Employee{}, result)
		repositoryMock.AssertExpectations(t)
	})

	// update_with_empty_card_number: Update with empty card number (no conflict check needed)
	// update_with_empty_card_number: Actualizar con card number vacío (no necesita verificación de conflicto)
	t.Run("update_with_empty_card_number", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		updateEmployee := models.Employee{
			CardNumberID: "", // Empty card number should not trigger conflict check
			FirstName:    "Updated Name",
			LastName:     "Updated LastName",
		}

		expectedEmployee := models.Employee{
			Id:           1,
			CardNumberID: "",
			FirstName:    "Updated Name",
			LastName:     "Updated LastName",
			WarehouseID:  1,
		}

		repositoryMock.On("GetCardNumberIds").Return([]string{"CARD-001", "CARD-002"}, nil)
		repositoryMock.On("Update", mock.Anything, 1, updateEmployee).Return(expectedEmployee, nil)

		result, err := service.Update(context.Background(), 1, updateEmployee)

		assert.NoError(t, err)
		assert.Equal(t, expectedEmployee, result)
		repositoryMock.AssertExpectations(t)
	})
}

// TestEmployeeService_DeleteById - Tests for DeleteById method of EmployeeService
// TestEmployeeService_DeleteById - Tests del método DeleteById de EmployeeService
// Validates deletion for existing and non-existing employees
// Valida eliminación para empleados existentes y no existentes
// Cases: successful deletion, non-existent employee
// Casos: eliminación exitosa, empleado inexistente
func TestEmployeeService_DeleteById(t *testing.T) {
	// delete_ok: Employee exists, deletes successfully without error
	// delete_ok: Empleado existe, elimina exitosamente sin error
	t.Run("delete_ok", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		repositoryMock.On("DeleteById", mock.Anything, 1).Return(nil)

		err := service.DeleteById(context.Background(), 1)

		assert.NoError(t, err)
		repositoryMock.AssertExpectations(t)
	})

	// delete_non_existent: Employee doesn't exist, returns not found error
	// delete_non_existent: Empleado no existe, retorna error not found
	t.Run("delete_non_existent", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		repositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetEmployeeService(repositoryMock)

		repositoryMock.On("DeleteById", mock.Anything, 999).Return(error_message.ErrNotFound)

		err := service.DeleteById(context.Background(), 999)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrNotFound))
		repositoryMock.AssertExpectations(t)
	})
}

// TestGetEmployeeService_Singleton - Test for singleton pattern to improve coverage
// TestGetEmployeeService_Singleton - Test del patrón singleton para mejorar coverage
// Validates that the service follows singleton pattern correctly
// Valida que el servicio sigue el patrón singleton correctamente
func TestGetEmployeeService_Singleton(t *testing.T) {
	t.Run("GetEmployeeService_returns_same_instance_when_already_exists", func(t *testing.T) {
		// Test: Singleton pattern ensures same instance is returned
		// Test: Patrón singleton asegura que se retorna la misma instancia
		ResetEmployeeServiceInstance()

		mockRepository1 := tests.GetNewEmployeeRepositoryMock()
		mockRepository2 := tests.GetNewEmployeeRepositoryMock()

		// First call - creates new instance / Primera llamada - crea nueva instancia
		service1 := GetEmployeeService(mockRepository1)

		// Second call - should return same instance (ignores new repository)
		// Segunda llamada - debe retornar la misma instancia (ignora nuevo repository)
		service2 := GetEmployeeService(mockRepository2)

		// Verify they are the same instance / Verifica que son la misma instancia
		assert.Equal(t, service1, service2)
	})
}
