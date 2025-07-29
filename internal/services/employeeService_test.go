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

// TestEmployeeService_Create - Tests del método Create de EmployeeService
// Valida lógica de negocio de creación incluyendo unicidad de card number ID
func TestEmployeeService_Create(t *testing.T) {
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

// TestEmployeeService_GetAll - Tests del método GetAll de EmployeeService
// Valida obtención de todos los empleados del repositorio
func TestEmployeeService_GetAll(t *testing.T) {
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

// TestEmployeeService_GetById - Tests del método GetById de EmployeeService
// Cubre búsqueda exitosa y empleado no encontrado
func TestEmployeeService_GetById(t *testing.T) {
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

// TestEmployeeService_Update - Tests del método Update de EmployeeService
// Valida actualización con verificación de unicidad de card number ID
func TestEmployeeService_Update(t *testing.T) {
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
}

// TestEmployeeService_DeleteById - Tests del método DeleteById de EmployeeService
// Valida eliminación para empleados existentes y no existentes
func TestEmployeeService_DeleteById(t *testing.T) {
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

// TestGetEmployeeService_Singleton - Test del patrón singleton para mejorar coverage
func TestGetEmployeeService_Singleton(t *testing.T) {
	t.Run("GetEmployeeService_returns_same_instance_when_already_exists", func(t *testing.T) {
		ResetEmployeeServiceInstance()

		mockRepository1 := tests.GetNewEmployeeRepositoryMock()
		mockRepository2 := tests.GetNewEmployeeRepositoryMock()

		// Primera llamada - crea nueva instancia
		service1 := GetEmployeeService(mockRepository1)

		// Segunda llamada - debe retornar la misma instancia (ignora nuevo repository)
		service2 := GetEmployeeService(mockRepository2)

		// Verifica que son la misma instancia
		assert.Equal(t, service1, service2)
	})
}
