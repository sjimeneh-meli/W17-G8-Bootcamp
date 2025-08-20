// Package services - InboundOrder Service Unit Tests / Tests Unitarios del Servicio InboundOrder
// Business logic testing for inbound order management with order number uniqueness validation
// Testing de lógica de negocio para gestión de órdenes de entrada con validación de unicidad de número de orden
package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestInboundOrdersService_Create - Tests for Create method of InboundOrdersService
// TestInboundOrdersService_Create - Tests del método Create de InboundOrdersService
// Validates business logic for creation including order number uniqueness and employee validation
// Valida lógica de negocio de creación incluyendo unicidad de número de orden y validación de empleado
// Cases: successful creation, duplicate order number, invalid input, employee not found, repository error
// Casos: creación exitosa, número de orden duplicado, entrada inválida, empleado no encontrado, error de repositorio
func TestInboundOrdersService_Create(t *testing.T) {
	// create_ok: Valid data with unique order number creates inbound order successfully
	// create_ok: Datos válidos con número de orden único crea orden de entrada exitosamente
	t.Run("create_ok", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		orderDate := time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC)
		order := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		expectedOrder := models.InboundOrder{
			Id:             1,
			OrderDate:      orderDate,
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		inboundOrderRepositoryMock.On("ExistsByOrderNumber", mock.Anything, "ORD-001").Return(false, nil)
		employeeRepositoryMock.On("ExistEmployeeById", mock.Anything, 1).Return(true, nil)
		inboundOrderRepositoryMock.On("Create", mock.Anything, order).Return(expectedOrder, nil)

		result, err := service.Create(context.Background(), order)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrder, result)
		inboundOrderRepositoryMock.AssertExpectations(t)
		employeeRepositoryMock.AssertExpectations(t)
	})

	// create_invalid_input_empty_order_number: Empty order number returns invalid input error
	// create_invalid_input_empty_order_number: Número de orden vacío retorna error de entrada inválida
	t.Run("create_invalid_input_empty_order_number", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		order := models.InboundOrder{
			OrderDate:      time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC),
			OrderNumber:    "", // Empty order number
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		result, err := service.Create(context.Background(), order)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInvalidInput))
		assert.Equal(t, models.InboundOrder{}, result)
	})

	// create_invalid_input_zero_employee_id: Zero employee ID returns invalid input error
	// create_invalid_input_zero_employee_id: ID de empleado cero retorna error de entrada inválida
	t.Run("create_invalid_input_zero_employee_id", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		order := models.InboundOrder{
			OrderDate:      time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC),
			OrderNumber:    "ORD-001",
			EmployeeId:     0, // Zero employee ID
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		result, err := service.Create(context.Background(), order)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInvalidInput))
		assert.Equal(t, models.InboundOrder{}, result)
	})

	// create_invalid_input_zero_order_date: Zero order date returns invalid input error
	// create_invalid_input_zero_order_date: Fecha de orden cero retorna error de entrada inválida
	t.Run("create_invalid_input_zero_order_date", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		order := models.InboundOrder{
			OrderDate:      time.Time{}, // Zero date
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		result, err := service.Create(context.Background(), order)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInvalidInput))
		assert.Equal(t, models.InboundOrder{}, result)
	})

	// create_duplicate_order_number: Duplicate order number returns already exists error
	// create_duplicate_order_number: Número de orden duplicado retorna error de ya existe
	t.Run("create_duplicate_order_number", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		orderDate := time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC)
		order := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		inboundOrderRepositoryMock.On("ExistsByOrderNumber", mock.Anything, "ORD-001").Return(true, nil)

		result, err := service.Create(context.Background(), order)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrAlreadyExists))
		assert.Equal(t, models.InboundOrder{}, result)
		inboundOrderRepositoryMock.AssertExpectations(t)
	})

	// create_order_number_check_error: Error checking order number existence
	// create_order_number_check_error: Error al verificar existencia de número de orden
	t.Run("create_order_number_check_error", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		order := models.InboundOrder{
			OrderDate:      time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC),
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		inboundOrderRepositoryMock.On("ExistsByOrderNumber", mock.Anything, "ORD-001").Return(false, error_message.ErrInternalServerError)

		result, err := service.Create(context.Background(), order)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
		assert.Equal(t, models.InboundOrder{}, result)
		inboundOrderRepositoryMock.AssertExpectations(t)
	})

	// create_employee_not_found: Employee doesn't exist returns dependency not found error
	// create_employee_not_found: Empleado no existe retorna error de dependencia no encontrada
	t.Run("create_employee_not_found", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		order := models.InboundOrder{
			OrderDate:      time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC),
			OrderNumber:    "ORD-001",
			EmployeeId:     999,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		inboundOrderRepositoryMock.On("ExistsByOrderNumber", mock.Anything, "ORD-001").Return(false, nil)
		employeeRepositoryMock.On("ExistEmployeeById", mock.Anything, 999).Return(false, nil)

		result, err := service.Create(context.Background(), order)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrDependencyNotFound))
		assert.Equal(t, models.InboundOrder{}, result)
		inboundOrderRepositoryMock.AssertExpectations(t)
		employeeRepositoryMock.AssertExpectations(t)
	})

	// create_employee_check_error: Error checking employee existence
	// create_employee_check_error: Error al verificar existencia de empleado
	t.Run("create_employee_check_error", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		order := models.InboundOrder{
			OrderDate:      time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC),
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		inboundOrderRepositoryMock.On("ExistsByOrderNumber", mock.Anything, "ORD-001").Return(false, nil)
		employeeRepositoryMock.On("ExistEmployeeById", mock.Anything, 1).Return(false, error_message.ErrInternalServerError)

		result, err := service.Create(context.Background(), order)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
		assert.Equal(t, models.InboundOrder{}, result)
		inboundOrderRepositoryMock.AssertExpectations(t)
		employeeRepositoryMock.AssertExpectations(t)
	})

	// create_repository_error: Repository error during creation
	// create_repository_error: Error de repositorio durante creación
	t.Run("create_repository_error", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		order := models.InboundOrder{
			OrderDate:      time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC),
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		inboundOrderRepositoryMock.On("ExistsByOrderNumber", mock.Anything, "ORD-001").Return(false, nil)
		employeeRepositoryMock.On("ExistEmployeeById", mock.Anything, 1).Return(true, nil)
		inboundOrderRepositoryMock.On("Create", mock.Anything, order).Return(models.InboundOrder{}, error_message.ErrInternalServerError)

		result, err := service.Create(context.Background(), order)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
		assert.Equal(t, models.InboundOrder{}, result)
		inboundOrderRepositoryMock.AssertExpectations(t)
		employeeRepositoryMock.AssertExpectations(t)
	})
}

// TestInboundOrdersService_GetAllInboundOrdersReports - Tests for GetAllInboundOrdersReports method
// TestInboundOrdersService_GetAllInboundOrdersReports - Tests del método GetAllInboundOrdersReports
// Validates retrieval of all inbound order reports from repository
// Valida obtención de todos los reportes de órdenes de entrada del repositorio
// Cases: successful retrieval with data, repository error
// Casos: recuperación exitosa con datos, error de repositorio
func TestInboundOrdersService_GetAllInboundOrdersReports(t *testing.T) {
	// get_all_reports_ok: Repository returns data and service passes it without modification
	// get_all_reports_ok: Repositorio retorna datos y servicio los pasa sin modificar
	t.Run("get_all_reports_ok", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		expectedReports := []models.InboundOrderReport{
			{Id: 1, IdCardNumber: "CARD-001", FirstName: "Carlos", LastName: "Ruiz", InboundOrderCount: 5},
			{Id: 2, IdCardNumber: "CARD-002", FirstName: "Maria", LastName: "Lopez", InboundOrderCount: 3},
		}

		inboundOrderRepositoryMock.On("GetAllInboundOrdersReports", mock.Anything).Return(expectedReports, nil)

		result, err := service.GetAllInboundOrdersReports(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedReports, result)
		inboundOrderRepositoryMock.AssertExpectations(t)
	})

	// get_all_reports_error: Repository error during retrieval
	// get_all_reports_error: Error de repositorio durante recuperación
	t.Run("get_all_reports_error", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		inboundOrderRepositoryMock.On("GetAllInboundOrdersReports", mock.Anything).Return([]models.InboundOrderReport{}, error_message.ErrInternalServerError)

		result, err := service.GetAllInboundOrdersReports(context.Background())

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
		assert.Empty(t, result)
		inboundOrderRepositoryMock.AssertExpectations(t)
	})
}

// TestInboundOrdersService_GetInboundOrdersReportByEmployeeId - Tests for GetInboundOrdersReportByEmployeeId method
// TestInboundOrdersService_GetInboundOrdersReportByEmployeeId - Tests del método GetInboundOrdersReportByEmployeeId
// Covers successful search and employee not found scenarios
// Cubre búsqueda exitosa y empleado no encontrado
// Cases: existing employee, non-existent employee, repository error
// Casos: empleado existente, empleado inexistente, error de repositorio
func TestInboundOrdersService_GetInboundOrdersReportByEmployeeId(t *testing.T) {
	// get_report_by_employee_id_ok: Employee exists, returns employee report data
	// get_report_by_employee_id_ok: Empleado existe, retorna datos del reporte del empleado
	t.Run("get_report_by_employee_id_ok", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		expectedReport := models.InboundOrderReport{
			Id:                1,
			IdCardNumber:      "CARD-001",
			FirstName:         "Carlos",
			LastName:          "Ruiz",
			InboundOrderCount: 5,
		}

		inboundOrderRepositoryMock.On("GetInboundOrdersReportByEmployeeId", mock.Anything, 1).Return(expectedReport, nil)

		result, err := service.GetInboundOrdersReportByEmployeeId(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedReport, result)
		inboundOrderRepositoryMock.AssertExpectations(t)
	})

	// get_report_by_employee_id_not_found: Employee doesn't exist, returns not found error
	// get_report_by_employee_id_not_found: Empleado no existe, retorna error not found
	t.Run("get_report_by_employee_id_not_found", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		inboundOrderRepositoryMock.On("GetInboundOrdersReportByEmployeeId", mock.Anything, 999).Return(models.InboundOrderReport{}, error_message.ErrNotFound)

		result, err := service.GetInboundOrdersReportByEmployeeId(context.Background(), 999)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrNotFound))
		assert.Equal(t, models.InboundOrderReport{}, result)
		inboundOrderRepositoryMock.AssertExpectations(t)
	})

	// get_report_by_employee_id_error: Repository error during retrieval
	// get_report_by_employee_id_error: Error de repositorio durante recuperación
	t.Run("get_report_by_employee_id_error", func(t *testing.T) {
		ResetInboundOrdersServiceInstance()

		inboundOrderRepositoryMock := tests.GetNewInboundOrderRepositoryMock()
		employeeRepositoryMock := tests.GetNewEmployeeRepositoryMock()
		service := GetInboundOrdersService(inboundOrderRepositoryMock, employeeRepositoryMock)

		inboundOrderRepositoryMock.On("GetInboundOrdersReportByEmployeeId", mock.Anything, 1).Return(models.InboundOrderReport{}, error_message.ErrInternalServerError)

		result, err := service.GetInboundOrdersReportByEmployeeId(context.Background(), 1)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
		assert.Equal(t, models.InboundOrderReport{}, result)
		inboundOrderRepositoryMock.AssertExpectations(t)
	})
}

// TestGetInboundOrdersService_Singleton - Test for singleton pattern to improve coverage
// TestGetInboundOrdersService_Singleton - Test del patrón singleton para mejorar coverage
// Validates that the service follows singleton pattern correctly
// Valida que el servicio sigue el patrón singleton correctamente
func TestGetInboundOrdersService_Singleton(t *testing.T) {
	t.Run("GetInboundOrdersService_returns_same_instance_when_already_exists", func(t *testing.T) {
		// Test: Singleton pattern ensures same instance is returned
		// Test: Patrón singleton asegura que se retorna la misma instancia

		// Reset singleton instance to ensure clean test
		// Resetea instancia singleton para asegurar test limpio
		inboundOrdersServiceInstance = nil

		mockInboundOrderRepository1 := tests.GetNewInboundOrderRepositoryMock()
		mockEmployeeRepository1 := tests.GetNewEmployeeRepositoryMock()
		mockInboundOrderRepository2 := tests.GetNewInboundOrderRepositoryMock()
		mockEmployeeRepository2 := tests.GetNewEmployeeRepositoryMock()

		// First call - creates new instance / Primera llamada - crea nueva instancia
		service1 := GetInboundOrdersService(mockInboundOrderRepository1, mockEmployeeRepository1)

		// Second call - should return same instance (ignores new repositories)
		// Segunda llamada - debe retornar la misma instancia (ignora nuevos repositories)
		service2 := GetInboundOrdersService(mockInboundOrderRepository2, mockEmployeeRepository2)

		// Verify they are the same instance / Verifica que son la misma instancia
		assert.Equal(t, service1, service2)

		// Reset for next tests / Resetea para próximos tests
		inboundOrdersServiceInstance = nil
	})
}
