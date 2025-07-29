// Package handlers_test contains comprehensive HTTP handler integration tests for the Buyer entity
// Este paquete contiene tests de integración HTTP comprehensivos para la entidad Buyer
//
// Key Testing Aspects / Aspectos Clave de Testing:
// - HTTP method testing (POST, GET, PATCH, DELETE) / Pruebas de métodos HTTP
// - Status code validation / Validación de códigos de estado
// - JSON request/response handling / Manejo de requests/responses JSON
// - Error scenarios and edge cases / Escenarios de error y casos límite
// - Mock service integration / Integración con servicios mock
//
// Testing Pattern / Patrón de Testing:
// - Arrange: Setup mock services, requests, expected responses / Configurar servicios mock, requests, respuestas esperadas
// - Act: Execute handler function / Ejecutar función del handler
// - Assert: Validate response codes and bodies / Validar códigos y cuerpos de respuesta
package handlers_test

import (
	"context"
	"encoding/json"
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
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test_Buyer_Post - Comprehensive testing suite for POST /api/v1/buyers endpoint
// Test_Buyer_Post - Suite de pruebas comprehensiva para el endpoint POST /api/v1/buyers
//
// Test Coverage / Cobertura de Pruebas:
// ✓ 201: Successful buyer creation / Creación exitosa de comprador
// ✓ 422: Invalid input validation / Validación de entrada inválida
// ✓ 409: Duplicate card number conflict / Conflicto por número de tarjeta duplicado
// ✓ 500: Internal server error handling / Manejo de errores internos del servidor
//
// HTTP Testing Patterns / Patrones de Testing HTTP:
// - JSON payload validation / Validación de payload JSON
// - Business logic error mapping / Mapeo de errores de lógica de negocio
// - Service layer mock interaction / Interacción con mocks de capa de servicio
func Test_Buyer_Post(t *testing.T) {
	// Test Case 1: Success scenario - validates complete happy path flow
	// Caso de Prueba 1: Escenario exitoso - valida flujo completo de caso feliz
	t.Run("Post Buyer successfully returns 201", func(t *testing.T) {
		expectedResponseBody := `{
			"data": {
				"id": 100,
				"id_card_number": "CARD-1001",
				"first_name": "Juan",
				"last_name": "Pérez"
			}
		}`
		expectedCode := 201

		validBuyerRequest := strings.NewReader(`{
			"id_card_number": "CARD-1001",
			"first_name": "Juan",
			"last_name": "Pérez"
		}`)

		mockRequestBuyer := models.Buyer{
			Id:           0,
			CardNumberId: "CARD-1001",
			FirstName:    "Juan",
			LastName:     "Pérez",
		}

		validBuyer := models.Buyer{
			Id:           100,
			CardNumberId: "CARD-1001",
			FirstName:    "Juan",
			LastName:     "Pérez",
		}

		request := httptest.NewRequest("POST", "/api/v1/buyers", validBuyerRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewBuyerServiceMock()
		service.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestBuyer).Return(validBuyer, nil)

		handler := handlers.GetBuyerHandler(service)

		handler.PostBuyer()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Post Buyer fails because request Buyer isn't valid returns 422", func(t *testing.T) {
		expectedCode := 422
		expectedResponseBody := `{
		"message":"first_name: cannot be blank; id_card_number: cannot be blank; last_name: cannot be blank.", 
		"status":"Unprocessable Entity"
		}`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "first_name: cannot be blank; id_card_number: cannot be blank; last_name: cannot be blank.",
			Status:  "Unprocessable Entity",
		}

		invalidBuyerRequest := strings.NewReader(`{
			"idcard_number": "CARD-1001",
			"firstname": "Juan",
			"lastname": "Pérez"
		}`)

		request := httptest.NewRequest("POST", "/api/v1/buyers", invalidBuyerRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewBuyerServiceMock()
		handler := handlers.GetBuyerHandler(service)

		handler.PostBuyer()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err := json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())

	})

	t.Run("Post Buyer fails because request Buyer has an card number id that already exists returns 409", func(t *testing.T) {
		expectedCode := 409
		expectedResponseBody := `{
		"message":"error: resource with the provided identifier already exists", 
		"status":"Conflict"
		}`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: resource with the provided identifier already exists",
			Status:  "Conflict",
		}

		repeatedCardNumberBuyerRequest := strings.NewReader(`{
			"id_card_number": "CARD1001",
			"first_name": "Juan",
			"last_name": "Pérez"
		}`)

		mockRequestBuyer := models.Buyer{
			Id:           0,
			CardNumberId: "CARD1001",
			FirstName:    "Juan",
			LastName:     "Pérez",
		}

		serviceMock := tests.GetNewBuyerServiceMock()
		serviceMock.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestBuyer).Return(models.Buyer{}, error_message.ErrAlreadyExists)
		handler := handlers.GetBuyerHandler(serviceMock)

		request := httptest.NewRequest("POST", "/api/v1/buyers", repeatedCardNumberBuyerRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.PostBuyer()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err := json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Post Buyer fails because of an internal server error returns 500", func(t *testing.T) {
		expectedCode := 500
		expectedResponseBody := `{
		"message":"error: an unexpected internal server error occurred", 
		"status":"Internal Server Error"
		}`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: an unexpected internal server error occurred",
			Status:  "Internal Server Error",
		}

		internalServerErrorBuyerRequest := strings.NewReader(`{
			"id_card_number": "CARD-1001",
			"first_name": "Juan",
			"last_name": "Pérez"
		}`)

		mockRequestBuyer := models.Buyer{
			Id:           0,
			CardNumberId: "CARD-1001",
			FirstName:    "Juan",
			LastName:     "Pérez",
		}

		serviceMock := tests.GetNewBuyerServiceMock()
		serviceMock.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestBuyer).Return(models.Buyer{}, error_message.ErrInternalServerError)
		handler := handlers.GetBuyerHandler(serviceMock)

		request := httptest.NewRequest("POST", "/api/v1/buyers", internalServerErrorBuyerRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.PostBuyer()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err := json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// Test_Buyer_GetAll - Complete testing for GET /api/v1/buyers endpoint (retrieve all buyers)
// Test_Buyer_GetAll - Pruebas completas para endpoint GET /api/v1/buyers (obtener todos los compradores)
//
// Test Scenarios / Escenarios de Prueba:
// ✓ 500: Service layer error handling / Manejo de errores de capa de servicio
// ✓ 200: Successful data retrieval / Recuperación exitosa de datos
//
// Key Validations / Validaciones Clave:
// - Error propagation from service to handler / Propagación de errores de servicio a handler
// - Response format consistency / Consistencia de formato de respuesta
// - Mock service behavior verification / Verificación de comportamiento de servicio mock
func Test_Buyer_GetAll(t *testing.T) {

	t.Run("error on service returns 500", func(t *testing.T) {
		expectedCode := 500
		expectedResponseBody := `{
									"message":"error: an unexpected internal server error occurred", 
									"status":"Internal Server Error"
								 }`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: an unexpected internal server error occurred",
			Status:  "Internal Server Error",
		}

		service := tests.GetNewBuyerServiceMock()
		service.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(map[int]models.Buyer{}, error_message.ErrInternalServerError)
		handler := handlers.GetBuyerHandler(service)

		request := httptest.NewRequest("GET", "/api/v1/buyers", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetAll()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err := json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		require.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		require.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Successfully get all buyers", func(t *testing.T) {
		expectedCode := 200

		mockBuyers := map[int]models.Buyer{
			1: {Id: 1, CardNumberId: "CARD-1001", FirstName: "Juan", LastName: "Pérez"},
			2: {Id: 2, CardNumberId: "CARD-1002", FirstName: "María", LastName: "Gómez"},
			3: {Id: 3, CardNumberId: "CARD-1003", FirstName: "Carlos", LastName: "López"},
		}
		service := tests.GetNewBuyerServiceMock()
		service.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(mockBuyers, nil)
		handler := handlers.GetBuyerHandler(service)

		request := httptest.NewRequest("GET", "/api/v1/buyers", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetAll()(response, request)

		require.Equal(t, expectedCode, response.Code)
	})
}

// Test_Buyer_GetById - Comprehensive testing for GET /api/v1/buyers/{id} endpoint
// Test_Buyer_GetById - Pruebas comprehensivas para endpoint GET /api/v1/buyers/{id}
//
// URL Parameter Testing / Pruebas de Parámetros URL:
// - Valid numeric ID handling / Manejo de ID numérico válido
// - Invalid ID format validation / Validación de formato de ID inválido
// - Non-existent ID error handling / Manejo de errores de ID inexistente
//
// HTTP Status Code Coverage / Cobertura de Códigos de Estado HTTP:
// ✓ 404: Resource not found / Recurso no encontrado
// ✓ 400: Invalid ID parameter / Parámetro ID inválido
// ✓ 500: Internal server error / Error interno del servidor
// ✓ 200: Successful resource retrieval / Recuperación exitosa del recurso
//
// Testing Infrastructure / Infraestructura de Testing:
// - Chi router context simulation / Simulación de contexto de router Chi
// - URL parameter injection / Inyección de parámetros URL
// - Helper function usage (newTestRequestWithIDParam) / Uso de funciones helper
func Test_Buyer_GetById(t *testing.T) {
	t.Run("Get By Id fails because request buyer id doesn't exists returns 404", func(t *testing.T) {
		id := "100"
		numberId := 100
		expectedCode := 404
		expectedResponseBody := `{
									"message":"error: the requested resource was not found", 
									"status":"Not Found"
								 }`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: the requested resource was not found",
			Status:  "Not Found",
		}

		serviceMock := tests.GetNewBuyerServiceMock()
		serviceMock.On("GetById", mock.AnythingOfType("*context.timerCtx"), numberId).Return(models.Buyer{}, error_message.ErrNotFound).Once()

		handler := handlers.GetBuyerHandler(serviceMock)

		request, err := newTestRequestWithIDParam("GET", "/api/v1/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetById()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())

	})
	t.Run("Get By Id fails because request id parameter isn't a number returns 400", func(t *testing.T) {
		id := "100a"

		expectedCode := 400
		expectedResponseBody := `{
									"message":"error: the provided input is invalid or missing required fields", 
									"status":"Bad Request"
								 }`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: the provided input is invalid or missing required fields",
			Status:  "Bad Request",
		}

		serviceMock := tests.GetNewBuyerServiceMock()
		handler := handlers.GetBuyerHandler(serviceMock)

		request, err := newTestRequestWithIDParam("GET", "/api/v1/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetById()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
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

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: an unexpected internal server error occurred",
			Status:  "Internal Server Error",
		}

		serviceMock := tests.GetNewBuyerServiceMock()
		serviceMock.On("GetById", mock.AnythingOfType("*context.timerCtx"), numberId).Return(models.Buyer{}, error_message.ErrInternalServerError)

		handler := handlers.GetBuyerHandler(serviceMock)

		request, err := newTestRequestWithIDParam("GET", "/api/v1/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetById()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Get By Id successfully returns requested buyer information returns 200", func(t *testing.T) {
		expectedCode := 200
		expectedResponseBody := `{
			"data": {
				"id": 1,
				"id_card_number": "CARD-1001",
				"first_name": "Juan",
				"last_name": "Pérez"
			}
		}`

		id := "100"
		numberId := 100
		mockBuyer := models.Buyer{
			Id:           1,
			CardNumberId: "CARD-1001",
			FirstName:    "Juan",
			LastName:     "Pérez",
		}

		serviceMock := tests.GetNewBuyerServiceMock()
		serviceMock.On("GetById", mock.AnythingOfType("*context.timerCtx"), numberId).Return(mockBuyer, nil)

		handler := handlers.GetBuyerHandler(serviceMock)

		request, err := newTestRequestWithIDParam("GET", "/api/v1/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetById()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// Test_Buyer_Patch - Extensive testing suite for PATCH /api/v1/buyers/{id} endpoint
// Test_Buyer_Patch - Suite de pruebas extensiva para endpoint PATCH /api/v1/buyers/{id}
//
// PATCH Operation Testing / Pruebas de Operación PATCH:
// - Partial update validation / Validación de actualización parcial
// - Required field presence checking / Verificación de presencia de campos requeridos
// - Data persistence verification / Verificación de persistencia de datos
//
// Business Logic Validations / Validaciones de Lógica de Negocio:
// ✓ ID parameter validation / Validación de parámetro ID
// ✓ Request body field requirements / Requisitos de campos del cuerpo de request
// ✓ Resource existence verification / Verificación de existencia de recurso
// ✓ Unique constraint handling / Manejo de restricciones de unicidad
// ✓ Server error resilience / Resistencia a errores del servidor
// ✓ Multiple field update success / Éxito en actualización de múltiples campos
//
// HTTP Response Testing / Pruebas de Respuesta HTTP:
// - Error message format consistency / Consistencia de formato de mensajes de error
// - Success response structure validation / Validación de estructura de respuesta exitosa
func Test_Buyer_Patch(t *testing.T) {
	t.Run("Patch Buyer fails because request id parameter isn't a number returns 400", func(t *testing.T) {
		id := "100a"

		expectedCode := 400
		expectedResponseBody := `{
								"message":"error: the provided input is invalid or missing required fields", 
								"status":"Bad Request"
								}`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: the provided input is invalid or missing required fields",
			Status:  "Bad Request",
		}

		serviceMock := tests.GetNewBuyerServiceMock()
		handler := handlers.GetBuyerHandler(serviceMock)

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.PatchBuyer()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		assert.NotEmpty(t, response.Body.String(), "Body should not be empty on error")
	})

	t.Run("Patch Buyer fails because request buyer doesn't have any valid buyer field to update returns 400", func(t *testing.T) {
		id := "1"
		expectedCode := 400
		expectedResponseBody := `{
		"message":"at least one of id_card_number, first_name, or last_name is required", 
		"status":"Bad Request"
		}`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "at least one of id_card_number, first_name, or last_name is required",
			Status:  "Bad Request",
		}

		invalidBuyerRequest := strings.NewReader(`{
			"idcard_number": "CARD-1001",
			"firstname": "Juan",
			"lastname": "Pérez"
		}`)

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/buyers", id, invalidBuyerRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewBuyerServiceMock()
		handler := handlers.GetBuyerHandler(service)

		handler.PatchBuyer()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Patch Buyer fails because request buyer id doesn't exists returns 404", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 404
		expectedResponseBody := `{
		"message":"error: the requested resource was not found", 
		"status":"Not Found"
		}`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: the requested resource was not found",
			Status:  "Not Found",
		}

		PatchBuyerRequest := strings.NewReader(`{
			"id_card_number": "CARD-10012"
		}`)
		mockPatchBuyerRequest := models.Buyer{
			CardNumberId: "CARD-10012",
		}

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/buyers", id, PatchBuyerRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewBuyerServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchBuyerRequest).Return(models.Buyer{}, error_message.ErrNotFound)
		handler := handlers.GetBuyerHandler(service)

		handler.PatchBuyer()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Patch buyer fails because request buyer id tries to update a buyer card number id to an already existing one returns 409", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 409
		expectedResponseBody := `{
		"message":"error: resource with the provided identifier already exists", 
		"status":"Conflict"
		}`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: resource with the provided identifier already exists",
			Status:  "Conflict",
		}

		PatchBuyerRequest := strings.NewReader(`{
			"id_card_number": "CARD-1001"
		}`)
		mockPatchBuyerRequest := models.Buyer{
			CardNumberId: "CARD-1001",
		}

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/buyers", id, PatchBuyerRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewBuyerServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchBuyerRequest).Return(models.Buyer{}, error_message.ErrAlreadyExists)
		handler := handlers.GetBuyerHandler(service)

		handler.PatchBuyer()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())

	})

	t.Run("Patch Buyer fails because of internal server error returns 500", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 500
		expectedResponseBody := `{
		"message":"error: an unexpected internal server error occurred", 
		"status":"Internal Server Error"
		}`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: an unexpected internal server error occurred",
			Status:  "Internal Server Error",
		}

		PatchBuyerRequest := strings.NewReader(`{
			"id_card_number": "CARD-1001"
		}`)
		mockPatchBuyerRequest := models.Buyer{
			CardNumberId: "CARD-1001",
		}

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/buyers", id, PatchBuyerRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewBuyerServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchBuyerRequest).Return(models.Buyer{}, error_message.ErrInternalServerError)
		handler := handlers.GetBuyerHandler(service)

		handler.PatchBuyer()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("patch Buyer successfully updates multiple fields returns 200", func(t *testing.T) {
		id := "101"
		idNumber := 101

		expectedCode := 200
		expectedResponseBody := `{
		"data": {
			"id": 101,
			"id_card_number": "CARD-101101",
			"first_name": "Ignacio",
			"last_name": "Gomez"
		}
	}`

		PatchBuyerRequest := strings.NewReader(`{
		"id_card_number": "CARD-101101",
		"first_name": "Ignacio",
		"last_name": "Gomez"
	}`)
		mockPatchBuyerRequest := models.Buyer{
			CardNumberId: "CARD-101101",
			FirstName:    "Ignacio",
			LastName:     "Gomez",
		}
		mockReturnBuyer := models.Buyer{
			Id:           101,
			CardNumberId: "CARD-101101",
			FirstName:    "Ignacio",
			LastName:     "Gomez",
		}

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/buyers", id, PatchBuyerRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewBuyerServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchBuyerRequest).Return(mockReturnBuyer, nil)
		handler := handlers.GetBuyerHandler(service)

		handler.PatchBuyer()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("patch Buyer successfully updates a buyer returns 200", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 200
		expectedResponseBody := `{
			"data": {
				"id": 100,
				"id_card_number": "CARD-100100",
				"first_name": "Ignacio",
				"last_name": "Garcia"
			}
		}`

		PatchBuyerRequest := strings.NewReader(`{
			"id_card_number": "CARD-100100"
		}`)
		mockPatchBuyerRequest := models.Buyer{
			CardNumberId: "CARD-100100",
		}
		mockReturnBuyer := models.Buyer{
			Id:           100,
			CardNumberId: "CARD-100100",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/buyers", id, PatchBuyerRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewBuyerServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchBuyerRequest).Return(mockReturnBuyer, nil)
		handler := handlers.GetBuyerHandler(service)

		handler.PatchBuyer()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// Test_Buyer_DeleteById - Complete testing coverage for DELETE /api/v1/buyers/{id} endpoint
// Test_Buyer_DeleteById - Cobertura completa de pruebas para endpoint DELETE /api/v1/buyers/{id}
//
// Delete Operation Scenarios / Escenarios de Operación DELETE:
// ✓ 400: Invalid ID parameter format / Formato de parámetro ID inválido
// ✓ 404: Resource not found for deletion / Recurso no encontrado para eliminación
// ✓ 500: Internal server error during deletion / Error interno durante eliminación
// ✓ 204: Successful deletion (No Content) / Eliminación exitosa (Sin Contenido)
//
// Key Testing Aspects / Aspectos Clave de Testing:
// - ID parameter validation and conversion / Validación y conversión de parámetro ID
// - Service layer error propagation / Propagación de errores de capa de servicio
// - HTTP status code accuracy / Precisión de códigos de estado HTTP
// - Empty response body validation for 204 / Validación de cuerpo de respuesta vacío para 204
func Test_Buyer_DeleteById(t *testing.T) {
	t.Run("Delete By Id fails because request buyer id parameter isn't a number returns 400", func(t *testing.T) {
		id := "100a"

		expectedCode := 400
		expectedResponseBody := `{
									"message":"error: the provided input is invalid or missing required fields", 
									"status":"Bad Request"
								 }`

		serviceMock := tests.GetNewBuyerServiceMock()
		handler := handlers.GetBuyerHandler(serviceMock)

		request, err := newTestRequestWithIDParam("DELETE", "/api/v1/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.DeleteById()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Delete By Id fails because request buyer id doesn't exists returns 404", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 404
		expectedResponseBody := `{
		"message":"error: the requested resource was not found", 
		"status":"Not Found"
		}`

		request, err := newTestRequestWithIDParam("DELETE", "/api/v1/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewBuyerServiceMock()
		service.On("DeleteById", mock.AnythingOfType("*context.timerCtx"), idNumber).Return(error_message.ErrNotFound)
		handler := handlers.GetBuyerHandler(service)

		handler.DeleteById()(response, request)

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

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: an unexpected internal server error occurred",
			Status:  "Internal Server Error",
		}

		request, err := newTestRequestWithIDParam("DELETE", "/api/v1/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewBuyerServiceMock()
		service.On("DeleteById", mock.AnythingOfType("*context.timerCtx"), idNumber).Return(error_message.ErrInternalServerError)
		handler := handlers.GetBuyerHandler(service)

		handler.DeleteById()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err = json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Delete By Id successfully deletes a buyer returns 204", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 204

		request, err := newTestRequestWithIDParam("DELETE", "/api/v1/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewBuyerServiceMock()
		service.On("DeleteById", mock.AnythingOfType("*context.timerCtx"), idNumber).Return(nil)

		handler := handlers.GetBuyerHandler(service)

		handler.DeleteById()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.Empty(t, response.Body.String(), "Body should be empty on 204 No Content")
	})
}

// newTestRequestWithIDParam - Utility helper function for creating HTTP requests with URL parameters
// newTestRequestWithIDParam - Función utilitaria helper para crear requests HTTP con parámetros URL
//
// Purpose / Propósito:
// - Simulates Chi router URL parameter extraction / Simula extracción de parámetros URL del router Chi
// - Reduces test code duplication / Reduce duplicación de código de test
// - Provides consistent request setup / Proporciona configuración consistente de requests
//
// Parameters / Parámetros:
// - method: HTTP method (GET, POST, PATCH, DELETE) / Método HTTP
// - pathBase: Base URL path / Ruta base URL
// - id: URL parameter value / Valor del parámetro URL
// - body: Request body reader / Reader del cuerpo de request
//
// Returns / Retorna:
// - Configured HTTP request with Chi route context / Request HTTP configurado con contexto de ruta Chi
// - Error if invalid parameters provided / Error si se proporcionan parámetros inválidos
func newTestRequestWithIDParam(method, pathBase, id string, body io.Reader) (*http.Request, error) {
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
