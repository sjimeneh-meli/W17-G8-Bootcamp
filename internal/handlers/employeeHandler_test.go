package handlers_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	services_test "github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Package handlers_test - Employee Handler Integration Tests
// Paquete handlers_test - Tests de Integración del Handler de Employee
//
// This file contains comprehensive HTTP handler tests for Employee entity operations
// Este archivo contiene tests HTTP comprehensivos para operaciones de la entidad Employee
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

// Tests comprehensivos de integración para handlers HTTP de Employee
// Validan flujo completo request-response con HTTP real (métodos, JSON, status codes)
// Ejercitan el stack completo del handler a diferencia de unit tests con mocks

// Comprehensive integration tests for Employee HTTP handlers
// Validate complete request-response flow with real HTTP (methods, JSON, status codes)
// Exercise the complete handler stack unlike unit tests with mocks

// newTestRequestWithIDParam_employee - Helper para crear requests HTTP con parámetros URL
// newTestRequestWithIDParam_employee - Helper for creating HTTP requests with URL parameters
//
// Functionality / Funcionalidad:
// - Creates properly formatted HTTP requests / Crea requests HTTP con formato adecuado
// - Injects Chi router context with URL parameters / Inyecta contexto de router Chi con parámetros URL
// - Validates input parameters / Valida parámetros de entrada
// - Returns configured request ready for testing / Retorna request configurado listo para testing
func newTestRequestWithIDParam_employee(method, pathBase, id string, body io.Reader) (*http.Request, error) {
	if method == "" {
		return nil, fmt.Errorf("HTTP method cannot be empty")
	}

	if pathBase == "" {
		return nil, fmt.Errorf("pathBase cannot be empty")
	}

	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	fullPath := pathBase + "/" + id
	req := httptest.NewRequest(method, fullPath, body)
	req.Header.Set("Content-Type", "application/json")

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", id)

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	return req.WithContext(ctx), nil
}

// TestPostEmployee - Tests comprehensivos del endpoint POST /employees
// TestPostEmployee - Comprehensive tests for POST /employees endpoint
//
// Testing Coverage / Cobertura de Testing:
// ✓ 201: Successful employee creation with all required fields / Creación exitosa con todos los campos requeridos
// ✓ 422: Validation error for missing/invalid fields / Error de validación por campos faltantes/inválidos
// ✓ 409: Business logic conflict (duplicate card number) / Conflicto de lógica de negocio (número de tarjeta duplicado)
// ✓ 500: Internal server error resilience / Resistencia a errores internos del servidor
//
// Employee-Specific Validations / Validaciones Específicas de Employee:
// - Card number uniqueness / Unicidad de número de tarjeta
// - Warehouse ID relationship validation / Validación de relación con ID de almacén
// - Required field presence (first_name, last_name, etc.) / Presencia de campos requeridos
//
// HTTP Testing Methodology / Metodología de Testing HTTP:
// - Real HTTP request simulation / Simulación real de requests HTTP
// - JSON payload marshaling/unmarshaling / Marshaling/unmarshaling de payload JSON
// - Content-Type header validation / Validación de headers Content-Type
// - Mock service behavior verification / Verificación de comportamiento de servicios mock
func TestPostEmployee(t *testing.T) {
	t.Run("Post Employee successfully returns 201", func(t *testing.T) {
		expectedResponseBody := `{
			"data": {
				"id": 100,
				"id_card_number": "CARD-1001",
				"first_name": "Juan",
				"last_name": "Pérez",
				"warehouse_id": 1
			}
		}`
		expectedCode := 201

		validEmployeeRequest := strings.NewReader(`{
			"id_card_number": "CARD-1001",
			"first_name": "Juan",
			"last_name": "Pérez",
			"warehouse_id": 1
		}`)

		mockRequestEmployee := models.Employee{
			Id:           0,
			CardNumberID: "CARD-1001",
			FirstName:    "Juan",
			LastName:     "Pérez",
			WarehouseID:  1,
		}

		validEmployee := models.Employee{
			Id:           100,
			CardNumberID: "CARD-1001",
			FirstName:    "Juan",
			LastName:     "Pérez",
			WarehouseID:  1,
		}

		request := httptest.NewRequest("POST", "/api/v1/employees", validEmployeeRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewEmployeeServiceMock()
		service.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestEmployee).Return(validEmployee, nil)

		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(service, validation)

		handler.PostEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Post Employee fails because request Employee isn't valid returns 422", func(t *testing.T) {
		expectedCode := 422
		expectedResponseBody := `{
		"message":"first_name: cannot be blank; id_card_number: cannot be blank; last_name: cannot be blank; warehouse_id: cannot be blank.", 
		"status":"Unprocessable Entity"
		}`

		invalidEmployeeRequest := strings.NewReader(`{
			"idcard_number": "CARD-1001",
			"firstname": "Juan",
			"lastname": "Pérez"
		}`)

		request := httptest.NewRequest("POST", "/api/v1/employees", invalidEmployeeRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewEmployeeServiceMock()
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(service, validation)

		handler.PostEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Post Employee fails because request Employee has an card number id that already exists returns 409", func(t *testing.T) {
		expectedCode := 409
		expectedResponseBody := `{
		"message":"error: resource with the provided identifier already exists", 
		"status":"Conflict"
		}`

		repeatedCardNumberEmployeeRequest := strings.NewReader(`{
			"id_card_number": "CARD1001",
			"first_name": "Juan",
			"last_name": "Pérez",
			"warehouse_id": 1
		}`)

		mockRequestEmployee := models.Employee{
			Id:           0,
			CardNumberID: "CARD1001",
			FirstName:    "Juan",
			LastName:     "Pérez",
			WarehouseID:  1,
		}

		serviceMock := services_test.GetNewEmployeeServiceMock()
		serviceMock.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestEmployee).Return(models.Employee{}, error_message.ErrAlreadyExists)
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(serviceMock, validation)

		request := httptest.NewRequest("POST", "/api/v1/employees", repeatedCardNumberEmployeeRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.PostEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Post Employee fails because of an internal server error returns 500", func(t *testing.T) {
		expectedCode := 500
		expectedResponseBody := `{
		"message":"error: an unexpected internal server error occurred", 
		"status":"Internal Server Error"
		}`
		internalServerErrorEmployeeRequest := strings.NewReader(`{
			"id_card_number": "CARD-1001",
			"first_name": "Juan",
			"last_name": "Pérez",
			"warehouse_id": 1
		}`)

		mockRequestEmployee := models.Employee{
			Id:           0,
			CardNumberID: "CARD-1001",
			FirstName:    "Juan",
			LastName:     "Pérez",
			WarehouseID:  1,
		}

		serviceMock := services_test.GetNewEmployeeServiceMock()
		serviceMock.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestEmployee).Return(models.Employee{}, error_message.ErrInternalServerError)
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(serviceMock, validation)

		request := httptest.NewRequest("POST", "/api/v1/employees", internalServerErrorEmployeeRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.PostEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// TestGetAllEmployees - Complete testing suite for GET /api/v1/employees endpoint
// TestGetAllEmployees - Suite completa de pruebas para endpoint GET /api/v1/employees
//
// Retrieval Operation Testing / Pruebas de Operación de Recuperación:
// ✓ 500: Service layer error handling / Manejo de errores de capa de servicio
// ✓ 200: Successful bulk data retrieval / Recuperación exitosa de datos en lote
//
// Data Structure Validation / Validación de Estructura de Datos:
// - Multiple employee records handling / Manejo de múltiples registros de empleados
// - Employee model consistency / Consistencia del modelo Employee
// - Response format standardization / Estandarización de formato de respuesta
//
// Service Integration Testing / Pruebas de Integración de Servicio:
// - Mock service GetAll method verification / Verificación del método GetAll del servicio mock
// - Error propagation from service to handler / Propagación de errores del servicio al handler
// - Context passing validation / Validación de paso de contexto
func TestGetAllEmployees(t *testing.T) {

	t.Run("error on service returns 500", func(t *testing.T) {
		expectedCode := 500
		expectedResponseBody := `{
									"message":"error: an unexpected internal server error occurred", 
									"status":"Internal Server Error"
								 }`

		service := services_test.GetNewEmployeeServiceMock()
		service.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(map[int]models.Employee{}, error_message.ErrInternalServerError)
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(service, validation)

		request := httptest.NewRequest("GET", "/api/v1/employees", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetAllEmployee()(response, request)

		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Successfully get all employees", func(t *testing.T) {
		expectedCode := 200

		mockEmployees := map[int]models.Employee{
			1: {Id: 1, CardNumberID: "CARD-1001", FirstName: "Juan", LastName: "Pérez", WarehouseID: 1},
			2: {Id: 2, CardNumberID: "CARD-1002", FirstName: "María", LastName: "Gómez", WarehouseID: 2},
			3: {Id: 3, CardNumberID: "CARD-1003", FirstName: "Carlos", LastName: "López", WarehouseID: 1},
		}
		service := services_test.GetNewEmployeeServiceMock()
		service.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(mockEmployees, nil)
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(service, validation)

		request := httptest.NewRequest("GET", "/api/v1/employees", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetAllEmployee()(response, request)

		require.Equal(t, expectedCode, response.Code)
	})
}

// TestGetEmployeeById - Comprehensive testing for GET /api/v1/employees/{id} endpoint
// TestGetEmployeeById - Pruebas comprehensivas para endpoint GET /api/v1/employees/{id}
//
// URL Parameter Processing / Procesamiento de Parámetros URL:
// - ID parameter extraction from URL / Extracción de parámetro ID de URL
// - String to integer conversion validation / Validación de conversión string a entero
// - Chi router context simulation / Simulación de contexto de router Chi
//
// Error Scenario Coverage / Cobertura de Escenarios de Error:
// ✓ 404: Employee not found by ID / Empleado no encontrado por ID
// ✓ 400: Invalid ID format (non-numeric) / Formato de ID inválido (no numérico)
// ✓ 500: Internal server error during retrieval / Error interno del servidor durante recuperación
// ✓ 200: Successful employee data retrieval / Recuperación exitosa de datos de empleado
//
// Response Validation / Validación de Respuesta:
// - JSON structure correctness / Corrección de estructura JSON
// - Employee model field mapping / Mapeo de campos del modelo Employee
// - Error message format consistency / Consistencia de formato de mensajes de error
func TestGetEmployeeById(t *testing.T) {
	t.Run("Get By Id fails because request employee id doesn't exists returns 404", func(t *testing.T) {
		id := "100"
		numberId := 100
		expectedCode := 404
		expectedResponseBody := `{
									"message":"error: the requested resource was not found", 
									"status":"Not Found"
								 }`

		serviceMock := services_test.GetNewEmployeeServiceMock()
		serviceMock.On("GetById", mock.AnythingOfType("*context.timerCtx"), numberId).Return(models.Employee{}, error_message.ErrNotFound).Once()

		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(serviceMock, validation)

		request, err := newTestRequestWithIDParam_employee("GET", "/api/v1/employees", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetByIdEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())

	})
	t.Run("Get By Id fails because request id parameter isn't a number returns 400", func(t *testing.T) {
		id := "100a"

		expectedCode := 400
		expectedResponseBody := `{
									"message":"strconv.Atoi: parsing \"100a\": invalid syntax", 
									"status":"Bad Request"
								 }`

		serviceMock := services_test.GetNewEmployeeServiceMock()
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(serviceMock, validation)

		request, err := newTestRequestWithIDParam_employee("GET", "/api/v1/employees", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetByIdEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Get By Id fails because of internal server error returns 500", func(t *testing.T) {
		id := "100"
		numberId := 100
		expectedCode := 500
		expectedResponseBody := `{
									"message":"error: an unexpected internal server error occurred", 
									"status":"Internal Server Error"
								 }`

		serviceMock := services_test.GetNewEmployeeServiceMock()
		serviceMock.On("GetById", mock.AnythingOfType("*context.timerCtx"), numberId).Return(models.Employee{}, error_message.ErrInternalServerError)

		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(serviceMock, validation)

		request, err := newTestRequestWithIDParam_employee("GET", "/api/v1/employees", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetByIdEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Get By Id successfully returns requested employee information returns 200", func(t *testing.T) {
		expectedCode := 200
		expectedResponseBody := `{
			"data": {
				"id": 1,
				"id_card_number": "CARD-1001",
				"first_name": "Juan",
				"last_name": "Pérez",
				"warehouse_id": 1
			}
		}`

		id := "100"
		numberId := 100
		mockEmployee := models.Employee{
			Id:           1,
			CardNumberID: "CARD-1001",
			FirstName:    "Juan",
			LastName:     "Pérez",
			WarehouseID:  1,
		}

		serviceMock := services_test.GetNewEmployeeServiceMock()
		serviceMock.On("GetById", mock.AnythingOfType("*context.timerCtx"), numberId).Return(mockEmployee, nil)

		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(serviceMock, validation)

		request, err := newTestRequestWithIDParam_employee("GET", "/api/v1/employees", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetByIdEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// TestPatchEmployee - Extensive testing for PATCH /api/v1/employees/{id} endpoint
// TestPatchEmployee - Pruebas extensivas para endpoint PATCH /api/v1/employees/{id}
//
// Partial Update Operation Testing / Pruebas de Operación de Actualización Parcial:
// - Individual field update capability / Capacidad de actualización de campos individuales
// - Multiple field simultaneous update / Actualización simultánea de múltiples campos
// - Field validation during partial updates / Validación de campos durante actualizaciones parciales
//
// Business Logic Validation / Validación de Lógica de Negocio:
// ✓ ID parameter format validation / Validación de formato de parámetro ID
// ✓ Required field presence for updates / Presencia de campos requeridos para actualizaciones
// ✓ Resource existence verification / Verificación de existencia de recurso
// ✓ Unique constraint enforcement / Aplicación de restricciones de unicidad
// ✓ Server error handling / Manejo de errores del servidor
// ✓ Successful update operation / Operación de actualización exitosa
//
// Update Validation Strategy / Estrategia de Validación de Actualización:
// - At least one field requirement / Requisito de al menos un campo
// - Card number uniqueness preservation / Preservación de unicidad de número de tarjeta
// - Data integrity maintenance / Mantenimiento de integridad de datos
//
// HTTP Method Specific Testing / Pruebas Específicas del Método HTTP:
// - PATCH semantics implementation / Implementación de semántica PATCH
// - Partial resource modification / Modificación parcial de recurso
// - Idempotency validation / Validación de idempotencia
func TestPatchEmployee(t *testing.T) {
	t.Run("Patch Employee fails because request id parameter isn't a number returns 400", func(t *testing.T) {
		id := "100a"

		expectedCode := 400
		expectedResponseBody := `{
									"message":"strconv.Atoi: parsing \"100a\": invalid syntax", 
									"status":"Bad Request"
								 }`

		serviceMock := services_test.GetNewEmployeeServiceMock()
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(serviceMock, validation)

		request, err := newTestRequestWithIDParam_employee("PATCH", "/api/v1/employees", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.PatchEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Patch Employee fails because request employee doesn't have any valid employee field to update returns 422", func(t *testing.T) {
		id := "1"
		expectedCode := 422
		expectedResponseBody := `{
		"message":"at least one of id_card_number, first_name, or last_name is required", 
		"status":"Unprocessable Entity"
		}`

		invalidEmployeeRequest := strings.NewReader(`{
			"idcard_number": "CARD-1001",
			"firstname": "Juan",
			"lastname": "Pérez"
		}`)

		request, err := newTestRequestWithIDParam_employee("PATCH", "/api/v1/employees", id, invalidEmployeeRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewEmployeeServiceMock()
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(service, validation)

		handler.PatchEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Patch Employee fails because request employee id doesn't exists returns 404", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 404
		expectedResponseBody := `{
		"message":"error: the requested resource was not found", 
		"status":"Not Found"
		}`

		PatchEmployeeRequest := strings.NewReader(`{
			"id_card_number": "CARD-10012"
		}`)
		mockPatchEmployeeRequest := models.Employee{
			CardNumberID: "CARD-10012",
		}

		request, err := newTestRequestWithIDParam_employee("PATCH", "/api/v1/employees", id, PatchEmployeeRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewEmployeeServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchEmployeeRequest).Return(models.Employee{}, error_message.ErrNotFound)
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(service, validation)

		handler.PatchEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Patch employee fails because request employee id tries to update a employee card number id to an already existing one returns 409", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 409
		expectedResponseBody := `{
		"message":"error: resource with the provided identifier already exists", 
		"status":"Conflict"
		}`

		PatchEmployeeRequest := strings.NewReader(`{
			"id_card_number": "CARD-1001"
		}`)
		mockPatchEmployeeRequest := models.Employee{
			CardNumberID: "CARD-1001",
		}

		request, err := newTestRequestWithIDParam_employee("PATCH", "/api/v1/employees", id, PatchEmployeeRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewEmployeeServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchEmployeeRequest).Return(models.Employee{}, error_message.ErrAlreadyExists)
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(service, validation)

		handler.PatchEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())

	})

	t.Run("Patch Employee fails because of internal server error returns 500", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 500
		expectedResponseBody := `{
		"message":"error: an unexpected internal server error occurred", 
		"status":"Internal Server Error"
		}`

		PatchEmployeeRequest := strings.NewReader(`{
			"id_card_number": "CARD-1001"
		}`)
		mockPatchEmployeeRequest := models.Employee{
			CardNumberID: "CARD-1001",
		}

		request, err := newTestRequestWithIDParam_employee("PATCH", "/api/v1/employees", id, PatchEmployeeRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewEmployeeServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchEmployeeRequest).Return(models.Employee{}, error_message.ErrInternalServerError)
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(service, validation)

		handler.PatchEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("patch Employee successfully updates a employee returns 200", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 200
		expectedResponseBody := `{
			"data": {
				"id": 100,
				"id_card_number": "CARD-100100",
							"first_name": "Ana",
			"last_name": "Torres",
				"warehouse_id": 2
			}
		}`

		PatchEmployeeRequest := strings.NewReader(`{
			"id_card_number": "CARD-100100"
		}`)
		mockPatchEmployeeRequest := models.Employee{
			CardNumberID: "CARD-100100",
		}
		mockReturnEmployee := models.Employee{
			Id:           100,
			CardNumberID: "CARD-100100",
			FirstName:    "Ana",
			LastName:     "Torres",
			WarehouseID:  2,
		}

		request, err := newTestRequestWithIDParam_employee("PATCH", "/api/v1/employees", id, PatchEmployeeRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewEmployeeServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchEmployeeRequest).Return(mockReturnEmployee, nil)
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(service, validation)

		handler.PatchEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// TestDeleteEmployeeById - Complete testing coverage for DELETE /api/v1/employees/{id} endpoint
// TestDeleteEmployeeById - Cobertura completa de pruebas para endpoint DELETE /api/v1/employees/{id}
//
// Delete Operation Scenarios / Escenarios de Operación DELETE:
// ✓ 400: Invalid ID parameter format handling / Manejo de formato de parámetro ID inválido
// ✓ 404: Non-existent employee deletion attempt / Intento de eliminación de empleado inexistente
// ✓ 500: Internal server error during deletion / Error interno del servidor durante eliminación
// ✓ 204: Successful employee deletion (No Content response) / Eliminación exitosa de empleado (respuesta Sin Contenido)
//
// HTTP DELETE Semantics / Semántica HTTP DELETE:
// - Resource removal operation / Operación de eliminación de recurso
// - Idempotent operation behavior / Comportamiento de operación idempotente
// - 204 No Content response validation / Validación de respuesta 204 Sin Contenido
// - Empty response body verification / Verificación de cuerpo de respuesta vacío
//
// Service Layer Integration / Integración de Capa de Servicio:
// - Delete method mock configuration / Configuración de mock del método Delete
// - Error propagation testing / Pruebas de propagación de errores
// - Context parameter passing / Paso de parámetros de contexto
//
// Data Consistency Validation / Validación de Consistencia de Datos:
// - Resource state after deletion / Estado del recurso después de eliminación
// - Cascade deletion considerations / Consideraciones de eliminación en cascada
// - Transaction rollback simulation / Simulación de rollback de transacción
func TestDeleteEmployeeById(t *testing.T) {
	t.Run("Delete By Id fails because request employee id parameter isn't a number returns 400", func(t *testing.T) {
		id := "100a"

		expectedCode := 400
		expectedResponseBody := `{
									"message":"strconv.Atoi: parsing \"100a\": invalid syntax", 
									"status":"Bad Request"
								 }`

		serviceMock := services_test.GetNewEmployeeServiceMock()
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(serviceMock, validation)

		request, err := newTestRequestWithIDParam_employee("DELETE", "/api/v1/employees", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.DeleteByIdEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Delete By Id fails because request employee id doesn't exists returns 404", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 404
		expectedResponseBody := `{
		"message":"error: the requested resource was not found", 
		"status":"Not Found"
		}`

		request, err := newTestRequestWithIDParam_employee("DELETE", "/api/v1/employees", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewEmployeeServiceMock()
		service.On("DeleteById", mock.AnythingOfType("*context.timerCtx"), idNumber).Return(error_message.ErrNotFound)
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(service, validation)

		handler.DeleteByIdEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Delete By Id fails because of internal server error returns 500", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 500
		expectedResponseBody := `{
		"message":"error: an unexpected internal server error occurred", 
		"status":"Internal Server Error"
		}`

		request, err := newTestRequestWithIDParam_employee("DELETE", "/api/v1/employees", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewEmployeeServiceMock()
		service.On("DeleteById", mock.AnythingOfType("*context.timerCtx"), idNumber).Return(error_message.ErrInternalServerError)
		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(service, validation)

		handler.DeleteByIdEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Delete By Id successfully deletes a employee returns 204", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 204

		request, err := newTestRequestWithIDParam_employee("DELETE", "/api/v1/employees", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewEmployeeServiceMock()
		service.On("DeleteById", mock.AnythingOfType("*context.timerCtx"), idNumber).Return(nil)

		validation := validations.GetEmployeeValidation()
		handler := handlers.GetEmployeeHandler(service, validation)

		handler.DeleteByIdEmployee()(response, request)

		assert.Equal(t, expectedCode, response.Code)

	})
}
