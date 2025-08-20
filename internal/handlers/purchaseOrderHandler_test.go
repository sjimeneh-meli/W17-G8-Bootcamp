// Package handlers_test contains comprehensive HTTP handler integration tests for the PurchaseOrder entity
// Este paquete contiene tests de integración HTTP comprehensivos para la entidad PurchaseOrder
//
// Key Testing Aspects / Aspectos Clave de Testing:
// - HTTP method testing (POST, GET) / Pruebas de métodos HTTP
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
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test_PurchaseOrder_Post - Comprehensive testing suite for POST /api/v1/purchaseOrders endpoint
// Test_PurchaseOrder_Post - Suite de pruebas comprehensiva para el endpoint POST /api/v1/purchaseOrders
//
// Test Coverage / Cobertura de Pruebas:
// ✓ 201: Successful purchase order creation / Creación exitosa de orden de compra
// ✓ 422: Invalid input validation / Validación de entrada inválida
// ✓ 409: Business logic conflicts (order number exists, buyer/product record not found) / Conflictos de lógica de negocio
// ✓ 500: Internal server error handling / Manejo de errores internos del servidor
//
// HTTP Testing Patterns / Patrones de Testing HTTP:
// - JSON payload validation / Validación de payload JSON
// - Business logic error mapping / Mapeo de errores de lógica de negocio
// - Service layer mock interaction / Interacción con mocks de capa de servicio
func Test_PurchaseOrder_Post(t *testing.T) {
	// Test Case 1: Success scenario - validates complete happy path flow
	// Caso de Prueba 1: Escenario exitoso - valida flujo completo de caso feliz
	t.Run("Post PurchaseOrder successfully returns 201", func(t *testing.T) {
		expectedResponseBody := `{
			"data": {
				"id": 17,
				"order_number": "ORDER-001",
				"order_date": "2024-01-15T10:30:00Z",
				"tracking_code": "TRACK-001",
				"buyer_id": 1,
				"product_record_id": 1
			}
		}`
		expectedCode := 201

		validPurchaseOrderRequest := strings.NewReader(`{
			"data": {
				"order_number": "ORDER-001",
				"order_date": "2024-01-15T10:30:00Z",
				"tracking_code": "TRACK-001",
				"buyer_id": 1,
				"product_record_id": 1
			}
		}`)

		mockRequestOrder := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		validOrder := models.PurchaseOrder{
			Id:              17,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		request := httptest.NewRequest("POST", "/api/v1/purchaseOrders", validPurchaseOrderRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewPurchaseOrderServiceMock()
		service.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestOrder).Return(validOrder, nil)

		handler := handlers.GetPurchaseOrderHandler(service)

		handler.PostPurchaseOrder()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Post PurchaseOrder fails because request data isn't valid returns 422", func(t *testing.T) {
		expectedCode := 422
		expectedResponseBody := `{
		"message":"data: cannot be blank. fields order_number, order_date, tracking_code, buyer_id, product_record_id are neccesary inside of data", 
		"status":"Unprocessable Entity"
		}`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "data: cannot be blank. fields order_number, order_date, tracking_code, buyer_id, product_record_id are neccesary inside of data",
			Status:  "Unprocessable Entity",
		}

		invalidPurchaseOrderRequest := strings.NewReader(`{
			"data": {}
		}`)

		request := httptest.NewRequest("POST", "/api/v1/purchaseOrders", invalidPurchaseOrderRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := tests.GetNewPurchaseOrderServiceMock()
		handler := handlers.GetPurchaseOrderHandler(service)

		handler.PostPurchaseOrder()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err := json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Post PurchaseOrder fails because order number already exists returns 409", func(t *testing.T) {
		expectedCode := 409
		expectedResponseBody := `{
		"message":"error: resource with the provided identifier already exists", 
		"status":"Conflict"
		}`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: resource with the provided identifier already exists",
			Status:  "Conflict",
		}

		duplicateOrderNumberRequest := strings.NewReader(`{
			"data": {
				"order_number": "ORDER-001",
				"order_date": "2024-01-15T10:30:00Z",
				"tracking_code": "TRACK-001",
				"buyer_id": 1,
				"product_record_id": 1
			}
		}`)

		mockRequestOrder := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		serviceMock := tests.GetNewPurchaseOrderServiceMock()
		serviceMock.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestOrder).Return(models.PurchaseOrder{}, error_message.ErrAlreadyExists)
		handler := handlers.GetPurchaseOrderHandler(serviceMock)

		request := httptest.NewRequest("POST", "/api/v1/purchaseOrders", duplicateOrderNumberRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.PostPurchaseOrder()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err := json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Post PurchaseOrder fails because buyer doesn't exist returns 409", func(t *testing.T) {
		expectedCode := 409
		expectedResponseBody := `{
		"message":"error: the requested resource was not found", 
		"status":"Conflict"
		}`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: the requested resource was not found",
			Status:  "Conflict",
		}

		buyerNotFoundRequest := strings.NewReader(`{
			"data": {
				"order_number": "ORDER-002",
				"order_date": "2024-01-15T10:30:00Z",
				"tracking_code": "TRACK-002",
				"buyer_id": 999,
				"product_record_id": 1
			}
		}`)

		mockRequestOrder := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-002",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-002",
			BuyerId:         999,
			ProductRecordId: 1,
		}

		serviceMock := tests.GetNewPurchaseOrderServiceMock()
		serviceMock.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestOrder).Return(models.PurchaseOrder{}, error_message.ErrNotFound)
		handler := handlers.GetPurchaseOrderHandler(serviceMock)

		request := httptest.NewRequest("POST", "/api/v1/purchaseOrders", buyerNotFoundRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.PostPurchaseOrder()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err := json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Post PurchaseOrder fails because of an internal server error returns 500", func(t *testing.T) {
		expectedCode := 500
		expectedResponseBody := `{
		"message":"error: an unexpected internal server error occurred", 
		"status":"Internal Server Error"
		}`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: an unexpected internal server error occurred",
			Status:  "Internal Server Error",
		}

		internalServerErrorRequest := strings.NewReader(`{
			"data": {
				"order_number": "ORDER-003",
				"order_date": "2024-01-15T10:30:00Z",
				"tracking_code": "TRACK-003",
				"buyer_id": 1,
				"product_record_id": 1
			}
		}`)

		mockRequestOrder := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-003",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-003",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		serviceMock := tests.GetNewPurchaseOrderServiceMock()
		serviceMock.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestOrder).Return(models.PurchaseOrder{}, error_message.ErrInternalServerError)
		handler := handlers.GetPurchaseOrderHandler(serviceMock)

		request := httptest.NewRequest("POST", "/api/v1/purchaseOrders", internalServerErrorRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.PostPurchaseOrder()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err := json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// Test_PurchaseOrder_GetAll - Complete testing for GET /api/v1/purchaseOrders endpoint (retrieve all purchase orders)
// Test_PurchaseOrder_GetAll - Pruebas completas para endpoint GET /api/v1/purchaseOrders (obtener todas las órdenes de compra)
//
// Test Scenarios / Escenarios de Prueba:
// ✓ 500: Service layer error handling / Manejo de errores de capa de servicio
// ✓ 200: Successful data retrieval / Recuperación exitosa de datos
//
// Key Validations / Validaciones Clave:
// - Error propagation from service to handler / Propagación de errores de servicio a handler
// - Response format consistency / Consistencia de formato de respuesta
// - Mock service behavior verification / Verificación de comportamiento de servicio mock
func Test_PurchaseOrder_GetAll(t *testing.T) {

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

		service := tests.GetNewPurchaseOrderServiceMock()
		service.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(map[int]models.PurchaseOrder{}, error_message.ErrInternalServerError)
		handler := handlers.GetPurchaseOrderHandler(service)

		request := httptest.NewRequest("GET", "/api/v1/purchaseOrders", nil)
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

	t.Run("Successfully get all purchase orders", func(t *testing.T) {
		expectedCode := 200

		mockPurchaseOrders := map[int]models.PurchaseOrder{
			1: {
				Id:              1,
				OrderNumber:     "ORDER-001",
				OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				TrackingCode:    "TRACK-001",
				BuyerId:         1,
				ProductRecordId: 1,
			},
			2: {
				Id:              2,
				OrderNumber:     "ORDER-002",
				OrderDate:       time.Date(2024, 1, 16, 14, 45, 0, 0, time.UTC),
				TrackingCode:    "TRACK-002",
				BuyerId:         2,
				ProductRecordId: 2,
			},
		}
		service := tests.GetNewPurchaseOrderServiceMock()
		service.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(mockPurchaseOrders, nil)
		handler := handlers.GetPurchaseOrderHandler(service)

		request := httptest.NewRequest("GET", "/api/v1/purchaseOrders", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetAll()(response, request)

		require.Equal(t, expectedCode, response.Code)
	})
}

// Test_PurchaseOrder_GetPurchaseOrdersReport - Comprehensive testing for GET /api/v1/purchaseOrders/reportByBuyers endpoint
// Test_PurchaseOrder_GetPurchaseOrdersReport - Pruebas comprehensivas para endpoint GET /api/v1/purchaseOrders/reportByBuyers
//
// Query Parameter Testing / Pruebas de Parámetros de Consulta:
// - Valid numeric ID handling / Manejo de ID numérico válido
// - Invalid ID format validation / Validación de formato de ID inválido
// - No ID parameter handling (all buyers) / Manejo sin parámetro ID (todos los compradores)
// - Non-existent ID error handling / Manejo de errores de ID inexistente
//
// HTTP Status Code Coverage / Cobertura de Códigos de Estado HTTP:
// ✓ 404: Resource not found / Recurso no encontrado
// ✓ 400: Invalid ID parameter / Parámetro ID inválido
// ✓ 500: Internal server error / Error interno del servidor
// ✓ 200: Successful resource retrieval for specific buyer / Recuperación exitosa del recurso para comprador específico
// ✓ 200: Successful resource retrieval for all buyers / Recuperación exitosa del recurso para todos los compradores
func Test_PurchaseOrder_GetPurchaseOrdersReport(t *testing.T) {
	t.Run("GetPurchaseOrdersReport successfully returns report for specific buyer returns 200", func(t *testing.T) {
		expectedCode := 200
		buyerId := 1

		expectedReports := []models.PurchaseOrderReport{
			{
				Id:                 1,
				IdCardNumber:       "CARD-001",
				FirstName:          "Ignacio",
				LastName:           "Garcia",
				PurchaseOrderCount: 3,
			},
		}

		serviceMock := tests.GetNewPurchaseOrderServiceMock()
		serviceMock.On("GetPurchaseOrdersReport", mock.AnythingOfType("*context.timerCtx"), &buyerId).Return(expectedReports, nil)

		handler := handlers.GetPurchaseOrderHandler(serviceMock)

		request := httptest.NewRequest("GET", "/api/v1/purchaseOrders/reportByBuyers?id=1", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetPurchaseOrdersReport()(response, request)

		assert.Equal(t, expectedCode, response.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &responseBody)
		require.NoError(t, err, "failed to unmarshal response body")

		data := responseBody["data"].([]interface{})
		assert.Equal(t, 1, len(data))
	})

	t.Run("GetPurchaseOrdersReport successfully returns reports for all buyers returns 200", func(t *testing.T) {
		expectedCode := 200

		expectedReports := []models.PurchaseOrderReport{
			{
				Id:                 1,
				IdCardNumber:       "CARD-001",
				FirstName:          "Ignacio",
				LastName:           "Garcia",
				PurchaseOrderCount: 3,
			},
			{
				Id:                 2,
				IdCardNumber:       "CARD-002",
				FirstName:          "Jesus",
				LastName:           "Ortega",
				PurchaseOrderCount: 2,
			},
		}

		serviceMock := tests.GetNewPurchaseOrderServiceMock()
		serviceMock.On("GetPurchaseOrdersReport", mock.AnythingOfType("*context.timerCtx"), (*int)(nil)).Return(expectedReports, nil)

		handler := handlers.GetPurchaseOrderHandler(serviceMock)

		request := httptest.NewRequest("GET", "/api/v1/purchaseOrders/reportByBuyers", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetPurchaseOrdersReport()(response, request)

		assert.Equal(t, expectedCode, response.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &responseBody)
		require.NoError(t, err, "failed to unmarshal response body")

		data := responseBody["data"].([]interface{})
		assert.Equal(t, 2, len(data))
	})

	t.Run("GetPurchaseOrdersReport fails because request id parameter isn't a number returns 400", func(t *testing.T) {
		expectedCode := 400
		expectedResponseBody := `{
									"message":"strconv.Atoi: parsing \"100a\": invalid syntax", 
									"status":"Bad Request"
								 }`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "strconv.Atoi: parsing \"100a\": invalid syntax",
			Status:  "Bad Request",
		}

		serviceMock := tests.GetNewPurchaseOrderServiceMock()
		handler := handlers.GetPurchaseOrderHandler(serviceMock)

		request := httptest.NewRequest("GET", "/api/v1/purchaseOrders/reportByBuyers?id=100a", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetPurchaseOrdersReport()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err := json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("GetPurchaseOrdersReport fails because buyer doesn't exist returns 404", func(t *testing.T) {
		buyerId := 999
		expectedCode := 404
		expectedResponseBody := `{
									"message":"error: the requested resource was not found", 
									"status":"Not Found"
								 }`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: the requested resource was not found",
			Status:  "Not Found",
		}

		serviceMock := tests.GetNewPurchaseOrderServiceMock()
		serviceMock.On("GetPurchaseOrdersReport", mock.AnythingOfType("*context.timerCtx"), &buyerId).Return([]models.PurchaseOrderReport{}, error_message.ErrNotFound).Once()

		handler := handlers.GetPurchaseOrderHandler(serviceMock)

		request := httptest.NewRequest("GET", "/api/v1/purchaseOrders/reportByBuyers?id=999", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetPurchaseOrdersReport()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err := json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("GetPurchaseOrdersReport fails because of internal server error returns 500", func(t *testing.T) {
		buyerId := 1
		expectedCode := 500
		expectedResponseBody := `{
									"message":"error: an unexpected internal server error occurred", 
									"status":"Internal Server Error"
								 }`

		expectedErrorResponse := tests.ErrorResponse{
			Message: "error: an unexpected internal server error occurred",
			Status:  "Internal Server Error",
		}

		serviceMock := tests.GetNewPurchaseOrderServiceMock()
		serviceMock.On("GetPurchaseOrdersReport", mock.AnythingOfType("*context.timerCtx"), &buyerId).Return([]models.PurchaseOrderReport{}, error_message.ErrInternalServerError)

		handler := handlers.GetPurchaseOrderHandler(serviceMock)

		request := httptest.NewRequest("GET", "/api/v1/purchaseOrders/reportByBuyers?id=1", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetPurchaseOrdersReport()(response, request)

		var actualErrorResponse tests.ErrorResponse
		err := json.Unmarshal(response.Body.Bytes(), &actualErrorResponse)
		require.NoError(t, err, "failed to unmarshal response body")

		assert.Equal(t, expectedCode, response.Code)
		assert.Equal(t, expectedErrorResponse.Message, actualErrorResponse.Message)
		assert.Equal(t, expectedErrorResponse.Status, actualErrorResponse.Status)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}
