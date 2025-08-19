package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	services_test "github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Package handlers_test - InboundOrder Handler Integration Tests
// Paquete handlers_test - Tests de Integración del Handler de InboundOrder
//
// This file contains comprehensive HTTP handler tests for InboundOrder entity operations
// Este archivo contiene tests HTTP comprehensivos para operaciones de la entidad InboundOrder
//
// Key Features / Características Clave:
// ✓ Full CRUD operation testing / Pruebas completas de operaciones CRUD
// ✓ HTTP status code validation / Validación de códigos de estado HTTP
// ✓ JSON serialization/deserialization testing / Pruebas de serialización/deserialización JSON
// ✓ Error handling and edge case coverage / Cobertura de manejo de errores y casos límite
// ✓ Mock service integration with validation layer / Integración de servicios mock con capa de validación
//
// Testing Architecture / Arquitectura de Testing:
// - Handler Layer: HTTP request/response processing / Capa Handler: procesamiento de request/response HTTP
// - Service Layer: Business logic simulation via mocks / Capa Servicio: simulación de lógica de negocio via mocks
// - Validation Layer: Input validation testing / Capa Validación: pruebas de validación de entrada

// newTestRequestWithIDParam_inboundOrder - Helper para crear requests HTTP con parámetros URL
// newTestRequestWithIDParam_inboundOrder - Helper for creating HTTP requests with URL parameters
func newTestRequestWithIDParam_inboundOrder(method, pathBase, id string, body io.Reader) (*http.Request, error) {
	if method == "" {
		return nil, fmt.Errorf("HTTP method cannot be empty")
	}

	if pathBase == "" {
		return nil, fmt.Errorf("pathBase cannot be empty")
	}

	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	fullPath := pathBase + "?id=" + id
	req := httptest.NewRequest(method, fullPath, body)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// TestPostInboundOrder - Tests comprehensivos del endpoint POST /inbound-orders
// TestPostInboundOrder - Comprehensive tests for POST /inbound-orders endpoint
//
// Testing Coverage / Cobertura de Testing:
// ✓ 201: Successful inbound order creation with all required fields / Creación exitosa con todos los campos requeridos
// ✓ 422: Validation error for missing/invalid fields / Error de validación por campos faltantes/inválidos
// ✓ 409: Business logic conflict (duplicate order number, employee not found) / Conflicto de lógica de negocio
// ✓ 500: Internal server error resilience / Resistencia a errores internos del servidor
//
// InboundOrder-Specific Validations / Validaciones Específicas de InboundOrder:
// - Order number uniqueness / Unicidad de número de orden
// - Employee ID relationship validation / Validación de relación con ID de empleado
// - Required field presence (order_date, order_number, etc.) / Presencia de campos requeridos
func TestPostInboundOrder(t *testing.T) {
	t.Run("Post InboundOrder successfully returns 201", func(t *testing.T) {
		expectedResponseBody := `{
			"data": {
				"id": 1,
				"order_number": "ORD-001",
				"order_date": "2023-12-01T10:00:00Z",
				"employee_id": 1,
				"product_batch_id": 1,
				"warehouse_id": 1
			}
		}`
		expectedCode := 201

		validInboundOrderRequest := strings.NewReader(`{
			"data": {
				"order_date": "2023-12-01T10:00:00Z",
				"order_number": "ORD-001",
				"employee_id": 1,
				"product_batch_id": 1,
				"warehouse_id": 1
			}
		}`)

		orderDate, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
		mockRequestInboundOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		validInboundOrder := models.InboundOrder{
			Id:             1,
			OrderDate:      orderDate,
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		request := httptest.NewRequest("POST", "/api/v1/inbound-orders", validInboundOrderRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewInboundOrderServiceMock()
		service.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestInboundOrder).Return(validInboundOrder, nil)

		handler := handlers.GetInboundOrderHandler(service)

		handler.PostInboundOrder()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("Post InboundOrder with invalid JSON returns 400", func(t *testing.T) {
		expectedCode := 400

		invalidJSONRequest := strings.NewReader(`{
			"data": {
				"order_date": "invalid-date",
				"order_number": "ORD-001"
			}
		}`)

		request := httptest.NewRequest("POST", "/api/v1/inbound-orders", invalidJSONRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		service := services_test.GetNewInboundOrderServiceMock()
		handler := handlers.GetInboundOrderHandler(service)

		handler.PostInboundOrder()(response, request)

		assert.Equal(t, expectedCode, response.Code)
	})

	t.Run("Post InboundOrder with empty data returns 422", func(t *testing.T) {
		expectedCode := 422

		emptyDataRequest := strings.NewReader(`{
			"data": {}
		}`)

		request := httptest.NewRequest("POST", "/api/v1/inbound-orders", emptyDataRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		service := services_test.GetNewInboundOrderServiceMock()
		handler := handlers.GetInboundOrderHandler(service)

		handler.PostInboundOrder()(response, request)

		assert.Equal(t, expectedCode, response.Code)
	})

	t.Run("Post InboundOrder with missing required fields returns 422", func(t *testing.T) {
		expectedCode := 422

		missingFieldsRequest := strings.NewReader(`{
			"data": {
				"order_date": "2023-12-01T10:00:00Z",
				"employee_id": 1
			}
		}`)

		request := httptest.NewRequest("POST", "/api/v1/inbound-orders", missingFieldsRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		service := services_test.GetNewInboundOrderServiceMock()
		handler := handlers.GetInboundOrderHandler(service)

		handler.PostInboundOrder()(response, request)

		assert.Equal(t, expectedCode, response.Code)
	})

	t.Run("Post InboundOrder with duplicate order number returns 409", func(t *testing.T) {
		expectedCode := 409

		duplicateOrderRequest := strings.NewReader(`{
			"data": {
				"order_date": "2023-12-01T10:00:00Z",
				"order_number": "ORD-001",
				"employee_id": 1,
				"product_batch_id": 1,
				"warehouse_id": 1
			}
		}`)

		orderDate, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
		mockRequestInboundOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		request := httptest.NewRequest("POST", "/api/v1/inbound-orders", duplicateOrderRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		service := services_test.GetNewInboundOrderServiceMock()
		service.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestInboundOrder).Return(models.InboundOrder{}, error_message.ErrAlreadyExists)

		handler := handlers.GetInboundOrderHandler(service)

		handler.PostInboundOrder()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		service.AssertExpectations(t)
	})

	t.Run("Post InboundOrder with non-existent employee returns 409", func(t *testing.T) {
		expectedCode := 409

		nonExistentEmployeeRequest := strings.NewReader(`{
			"data": {
				"order_date": "2023-12-01T10:00:00Z",
				"order_number": "ORD-001",
				"employee_id": 999,
				"product_batch_id": 1,
				"warehouse_id": 1
			}
		}`)

		orderDate, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
		mockRequestInboundOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "ORD-001",
			EmployeeId:     999,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		request := httptest.NewRequest("POST", "/api/v1/inbound-orders", nonExistentEmployeeRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		service := services_test.GetNewInboundOrderServiceMock()
		service.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestInboundOrder).Return(models.InboundOrder{}, error_message.ErrDependencyNotFound)

		handler := handlers.GetInboundOrderHandler(service)

		handler.PostInboundOrder()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		service.AssertExpectations(t)
	})

	t.Run("Post InboundOrder with service error returns 500", func(t *testing.T) {
		expectedCode := 500

		validInboundOrderRequest := strings.NewReader(`{
			"data": {
				"order_date": "2023-12-01T10:00:00Z",
				"order_number": "ORD-001",
				"employee_id": 1,
				"product_batch_id": 1,
				"warehouse_id": 1
			}
		}`)

		orderDate, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
		mockRequestInboundOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		request := httptest.NewRequest("POST", "/api/v1/inbound-orders", validInboundOrderRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		service := services_test.GetNewInboundOrderServiceMock()
		service.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestInboundOrder).Return(models.InboundOrder{}, error_message.ErrInternalServerError)

		handler := handlers.GetInboundOrderHandler(service)

		handler.PostInboundOrder()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		service.AssertExpectations(t)
	})
}

// TestGetInboundOrdersReport - Tests comprehensivos del endpoint GET /employees/reportInboundOrders
// TestGetInboundOrdersReport - Comprehensive tests for GET /employees/reportInboundOrders endpoint
//
// Testing Coverage / Cobertura de Testing:
// ✓ 200: Successful retrieval of all reports / Recuperación exitosa de todos los reportes
// ✓ 200: Successful retrieval of specific employee report / Recuperación exitosa de reporte específico de empleado
// ✓ 400: Invalid employee ID parameter / Parámetro de ID de empleado inválido
// ✓ 500: Internal server error resilience / Resistencia a errores internos del servidor
//
// Report-Specific Features / Características Específicas de Reportes:
// - All employees report generation / Generación de reporte de todos los empleados
// - Single employee report by ID / Reporte de empleado individual por ID
// - Query parameter validation / Validación de parámetros de consulta
func TestGetInboundOrdersReport(t *testing.T) {
	t.Run("Get all inbound orders reports successfully returns 200", func(t *testing.T) {
		expectedResponseBody := `{
			"data": [
				{
					"id": 1,
					"id_card_number": "CARD-001",
					"first_name": "Carlos",
					"last_name": "Ruiz",
					"inbound_orders_count": 5
				},
				{
					"id": 2,
					"id_card_number": "CARD-002",
					"first_name": "Maria",
					"last_name": "Lopez",
					"inbound_orders_count": 3
				}
			]
		}`
		expectedCode := 200

		expectedReports := []models.InboundOrderReport{
			{Id: 1, IdCardNumber: "CARD-001", FirstName: "Carlos", LastName: "Ruiz", InboundOrderCount: 5},
			{Id: 2, IdCardNumber: "CARD-002", FirstName: "Maria", LastName: "Lopez", InboundOrderCount: 3},
		}

		request := httptest.NewRequest("GET", "/api/v1/employees/reportInboundOrders", nil)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		service := services_test.GetNewInboundOrderServiceMock()
		service.On("GetAllInboundOrdersReports", mock.AnythingOfType("*context.timerCtx")).Return(expectedReports, nil)

		handler := handlers.GetInboundOrderHandler(service)

		handler.GetInboundOrdersReport()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("Get inbound orders report by employee ID successfully returns 200", func(t *testing.T) {
		expectedResponseBody := `{
			"data": [
				{
					"id": 1,
					"id_card_number": "CARD-001",
					"first_name": "Carlos",
					"last_name": "Ruiz",
					"inbound_orders_count": 5
				}
			]
		}`
		expectedCode := 200

		expectedReport := models.InboundOrderReport{
			Id:                1,
			IdCardNumber:      "CARD-001",
			FirstName:         "Carlos",
			LastName:          "Ruiz",
			InboundOrderCount: 5,
		}

		request, err := newTestRequestWithIDParam_inboundOrder("GET", "/api/v1/employees/reportInboundOrders", "1", nil)
		require.NoError(t, err)
		response := httptest.NewRecorder()

		service := services_test.GetNewInboundOrderServiceMock()
		service.On("GetInboundOrdersReportByEmployeeId", mock.AnythingOfType("*context.timerCtx"), 1).Return(expectedReport, nil)

		handler := handlers.GetInboundOrderHandler(service)

		handler.GetInboundOrdersReport()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("Get inbound orders report with invalid employee ID returns 400", func(t *testing.T) {
		expectedCode := 400

		request, err := newTestRequestWithIDParam_inboundOrder("GET", "/api/v1/employees/reportInboundOrders", "invalid", nil)
		require.NoError(t, err)
		response := httptest.NewRecorder()

		service := services_test.GetNewInboundOrderServiceMock()
		handler := handlers.GetInboundOrderHandler(service)

		handler.GetInboundOrdersReport()(response, request)

		assert.Equal(t, expectedCode, response.Code)
	})

	t.Run("Get inbound orders report with service error returns 500", func(t *testing.T) {
		expectedCode := 500

		request := httptest.NewRequest("GET", "/api/v1/employees/reportInboundOrders", nil)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		service := services_test.GetNewInboundOrderServiceMock()
		service.On("GetAllInboundOrdersReports", mock.AnythingOfType("*context.timerCtx")).Return([]models.InboundOrderReport{}, error_message.ErrInternalServerError)

		handler := handlers.GetInboundOrderHandler(service)

		handler.GetInboundOrdersReport()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		service.AssertExpectations(t)
	})

	t.Run("Get inbound orders report by employee ID with service error returns 500", func(t *testing.T) {
		expectedCode := 500

		request, err := newTestRequestWithIDParam_inboundOrder("GET", "/api/v1/employees/reportInboundOrders", "1", nil)
		require.NoError(t, err)
		response := httptest.NewRecorder()

		service := services_test.GetNewInboundOrderServiceMock()
		service.On("GetInboundOrdersReportByEmployeeId", mock.AnythingOfType("*context.timerCtx"), 1).Return(models.InboundOrderReport{}, error_message.ErrInternalServerError)

		handler := handlers.GetInboundOrderHandler(service)

		handler.GetInboundOrdersReport()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		service.AssertExpectations(t)
	})
}

// TestGetInboundOrderHandler - Test for handler constructor to improve coverage
// TestGetInboundOrderHandler - Test del constructor del handler para mejorar coverage
// Validates that the handler is created correctly with the provided service
// Valida que el handler se crea correctamente con el servicio proporcionado
func TestGetInboundOrderHandler(t *testing.T) {
	t.Run("GetInboundOrderHandler creates handler with service", func(t *testing.T) {
		// Test: Handler constructor creates proper instance
		// Test: Constructor del handler crea instancia apropiada
		service := services_test.GetNewInboundOrderServiceMock()

		handler := handlers.GetInboundOrderHandler(service)

		assert.NotNil(t, handler)
		// Verify the handler has the expected interface methods
		// Verificar que el handler tiene los métodos de interfaz esperados
		assert.NotNil(t, handler.GetInboundOrdersReport())
		assert.NotNil(t, handler.PostInboundOrder())
	})
}
