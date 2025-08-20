// Package handlers_test - Locality Handler Unit Tests / Tests Unitarios del Handler Locality
// HTTP endpoint testing for locality operations with JSON validation and service integration
// Testing de endpoints HTTP para operaciones de localidades con validación JSON e integración de servicios
package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestLocalityHandler_Save - Tests for POST /localities endpoint
// Cases: 200 (successful creation), 400 (invalid JSON), 422 (validation error), 409 (already exists), 500 (internal error), 504 (timeout)
// Tests para endpoint POST /localities
// Casos: 200 (creación exitosa), 400 (JSON inválido), 422 (error de validación), 409 (ya existe), 500 (error interno), 504 (timeout)
func TestLocalityHandler_Save(t *testing.T) {
	t.Run("should return 200 and locality when creation is successful", func(t *testing.T) {
		// Test: Valid locality request creates locality successfully / Solicitud válida de localidad crea localidad exitosamente
		expectedCode := 200
		expectedResponseBody := `{
			"data": {
				"id": 1,
				"locality_name": "Test City",
				"province_name": "Test Province",
				"country_name": "Test Country"
			}
		}`

		validLocalityRequest := strings.NewReader(`{
			"data": {
				"locality_name": "Test City",
				"province_name": "Test Province",
				"country_name": "Test Country"
			}
		}`)

		expectedLocality := models.Locality{
			Id:           1,
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		request := httptest.NewRequest(http.MethodPost, "/localities", validLocalityRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		service.On("Save", mock.Anything, mock.AnythingOfType("models.Locality")).Return(expectedLocality, nil)

		handler := handlers.NewLocalityHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("should return 400 when JSON is invalid", func(t *testing.T) {
		// Test: Malformed JSON returns bad request error / JSON malformado retorna error de solicitud incorrecta
		expectedCode := 400
		expectedResponseBody := `{
			"status": "Bad Request",
			"message": "invalid character '}' looking for beginning of object key string"
		}`

		invalidJSONRequest := strings.NewReader(`{
			"data": {
				"locality_name": "Test City",
				"province_name": "Test Province",
				"country_name": "Test Country",
			}
		}`)

		request := httptest.NewRequest(http.MethodPost, "/localities", invalidJSONRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		handler := handlers.NewLocalityHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 422 when locality_name is empty", func(t *testing.T) {
		// Test: Empty locality name fails validation / Nombre de localidad vacío falla validación
		expectedCode := 422
		expectedResponseBody := `{
			"status": "Unprocessable Entity",
			"message": "locality_name: cannot be blank."
		}`

		invalidLocalityRequest := strings.NewReader(`{
			"data": {
				"locality_name": "",
				"province_name": "Test Province",
				"country_name": "Test Country"
			}
		}`)

		request := httptest.NewRequest(http.MethodPost, "/localities", invalidLocalityRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		handler := handlers.NewLocalityHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 422 when province_name is empty", func(t *testing.T) {
		// Test: Empty province name fails validation / Nombre de provincia vacío falla validación
		expectedCode := 422
		expectedResponseBody := `{
			"status": "Unprocessable Entity",
			"message": "province_name: cannot be blank."
		}`

		invalidLocalityRequest := strings.NewReader(`{
			"data": {
				"locality_name": "Test City",
				"province_name": "",
				"country_name": "Test Country"
			}
		}`)

		request := httptest.NewRequest(http.MethodPost, "/localities", invalidLocalityRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		handler := handlers.NewLocalityHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 422 when country_name is empty", func(t *testing.T) {
		// Test: Empty country name fails validation / Nombre de país vacío falla validación
		expectedCode := 422
		expectedResponseBody := `{
			"status": "Unprocessable Entity",
			"message": "country_name: cannot be blank."
		}`

		invalidLocalityRequest := strings.NewReader(`{
			"data": {
				"locality_name": "Test City",
				"province_name": "Test Province",
				"country_name": ""
			}
		}`)

		request := httptest.NewRequest(http.MethodPost, "/localities", invalidLocalityRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		handler := handlers.NewLocalityHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 409 when locality already exists", func(t *testing.T) {
		// Test: Duplicate locality returns conflict error / Localidad duplicada retorna error de conflicto
		expectedCode := 409
		expectedResponseBody := `{
			"status": "Conflict",
			"message": "error: resource with the provided identifier already exists"
		}`

		validLocalityRequest := strings.NewReader(`{
			"data": {
				"locality_name": "Existing City",
				"province_name": "Test Province",
				"country_name": "Test Country"
			}
		}`)

		request := httptest.NewRequest(http.MethodPost, "/localities", validLocalityRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		service.On("Save", mock.Anything, mock.AnythingOfType("models.Locality")).Return(models.Locality{}, error_message.ErrAlreadyExists)

		handler := handlers.NewLocalityHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("should return 500 when service returns query error", func(t *testing.T) {
		// Test: Service query error returns internal server error / Error de consulta del servicio retorna error interno del servidor
		expectedCode := 500
		expectedResponseBody := `{
			"status": "Internal Server Error",
			"message": "error: failed to query or insert"
		}`

		validLocalityRequest := strings.NewReader(`{
			"data": {
				"locality_name": "Test City",
				"province_name": "Test Province",
				"country_name": "Test Country"
			}
		}`)

		request := httptest.NewRequest(http.MethodPost, "/localities", validLocalityRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		service.On("Save", mock.Anything, mock.AnythingOfType("models.Locality")).Return(models.Locality{}, error_message.ErrQuery)

		handler := handlers.NewLocalityHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("should return 504 when request times out", func(t *testing.T) {
		// Test: Context timeout returns gateway timeout error / Timeout de contexto retorna error de timeout de gateway
		expectedCode := 504
		expectedResponseBody := `{
			"status": "Gateway Timeout",
			"message": "context deadline exceeded"
		}`

		validLocalityRequest := strings.NewReader(`{
			"data": {
				"locality_name": "Test City",
				"province_name": "Test Province",
				"country_name": "Test Country"
			}
		}`)

		request := httptest.NewRequest(http.MethodPost, "/localities", validLocalityRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		service.On("Save", mock.Anything, mock.AnythingOfType("models.Locality")).Return(models.Locality{}, context.DeadlineExceeded)

		handler := handlers.NewLocalityHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		service.AssertExpectations(t)
	})
}

// TestLocalityHandler_GetSellerReportByLocality - Tests for GET /localities/reportSellers endpoint
// Cases: 200 (successful retrieval all), 200 (successful retrieval specific), 400 (invalid query), 404 (not found), 500 (internal error), 504 (timeout)
// Tests para endpoint GET /localities/reportSellers
// Casos: 200 (recuperación exitosa todas), 200 (recuperación exitosa específica), 400 (consulta inválida), 404 (no encontrada), 500 (error interno), 504 (timeout)
func TestLocalityHandler_GetSellerReportByLocality(t *testing.T) {
	t.Run("should return 200 and all locality seller reports when no id parameter", func(t *testing.T) {
		// Test: Request without ID parameter returns all locality reports / Solicitud sin parámetro ID retorna todos los reportes de localidades
		expectedCode := 200
		expectedReports := []responses.LocalitySellerReport{
			{
				LocalityID:   1,
				LocalityName: "City One",
				SellerCount:  5,
			},
			{
				LocalityID:   2,
				LocalityName: "City Two",
				SellerCount:  3,
			},
		}

		expectedResponseBody, _ := json.Marshal(responses.DataResponse{Data: expectedReports})

		request := httptest.NewRequest(http.MethodGet, "/localities/reportSellers", nil)
		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		service.On("GetSellerReports", mock.Anything, 0).Return(expectedReports, nil)

		handler := handlers.NewLocalityHandler(service)
		handler.GetSellerReportByLocality(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, string(expectedResponseBody), response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("should return 200 and specific locality seller report when valid id parameter", func(t *testing.T) {
		// Test: Request with valid ID parameter returns specific locality report / Solicitud con parámetro ID válido retorna reporte específico de localidad
		expectedCode := 200
		expectedReports := []responses.LocalitySellerReport{
			{
				LocalityID:   1,
				LocalityName: "Specific City",
				SellerCount:  2,
			},
		}

		expectedResponseBody, _ := json.Marshal(responses.DataResponse{Data: expectedReports})

		request := httptest.NewRequest(http.MethodGet, "/localities/reportSellers?id=1", nil)
		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		service.On("GetSellerReports", mock.Anything, 1).Return(expectedReports, nil)

		handler := handlers.NewLocalityHandler(service)
		handler.GetSellerReportByLocality(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, string(expectedResponseBody), response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("should return 500 when invalid query parameter format", func(t *testing.T) {
		// Test: Invalid query parameter format returns internal server error / Formato de parámetro de consulta inválido retorna error interno del servidor
		expectedCode := 500
		expectedResponseBody := `{
			"status": "Internal Server Error",
			"message": "Invalid query parameter"
		}`

		request := httptest.NewRequest(http.MethodGet, "/localities/reportSellers?invalid=param", nil)
		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		handler := handlers.NewLocalityHandler(service)
		handler.GetSellerReportByLocality(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 400 when id parameter is not a valid integer", func(t *testing.T) {
		// Test: Non-integer ID parameter returns bad request error / Parámetro ID no entero retorna error de solicitud incorrecta
		expectedCode := 400
		expectedResponseBody := `{
			"status": "Bad Request",
			"message": "strconv.Atoi: parsing \"invalid\": invalid syntax"
		}`

		request := httptest.NewRequest(http.MethodGet, "/localities/reportSellers?id=invalid", nil)
		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		handler := handlers.NewLocalityHandler(service)
		handler.GetSellerReportByLocality(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 404 when locality not found", func(t *testing.T) {
		// Test: Non-existent locality returns not found error / Localidad inexistente retorna error not found
		expectedCode := 404
		expectedResponseBody := `{
			"status": "Not Found",
			"message": "error: the requested resource was not found"
		}`

		request := httptest.NewRequest(http.MethodGet, "/localities/reportSellers?id=999", nil)
		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		service.On("GetSellerReports", mock.Anything, 999).Return([]responses.LocalitySellerReport(nil), error_message.ErrNotFound)

		handler := handlers.NewLocalityHandler(service)
		handler.GetSellerReportByLocality(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("should return 500 when service returns existence check error", func(t *testing.T) {
		// Test: Service existence check error returns internal server error / Error de verificación de existencia del servicio retorna error interno del servidor
		expectedCode := 500
		expectedResponseBody := `{
			"status": "Internal Server Error",
			"message": "error: failed checking existence"
		}`

		request := httptest.NewRequest(http.MethodGet, "/localities/reportSellers?id=1", nil)
		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		service.On("GetSellerReports", mock.Anything, 1).Return([]responses.LocalitySellerReport(nil), error_message.ErrFailedCheckingExistence)

		handler := handlers.NewLocalityHandler(service)
		handler.GetSellerReportByLocality(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("should return 500 when service returns query report error", func(t *testing.T) {
		// Test: Service query report error returns internal server error / Error de consulta de reporte del servicio retorna error interno del servidor
		expectedCode := 500
		expectedResponseBody := `{
			"status": "Internal Server Error",
			"message": "error: querying report failed"
		}`

		request := httptest.NewRequest(http.MethodGet, "/localities/reportSellers", nil)
		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		service.On("GetSellerReports", mock.Anything, 0).Return([]responses.LocalitySellerReport(nil), error_message.ErrQueryingReport)

		handler := handlers.NewLocalityHandler(service)
		handler.GetSellerReportByLocality(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("should return 500 when service returns scan error", func(t *testing.T) {
		// Test: Service scan error returns internal server error / Error de escaneo del servicio retorna error interno del servidor
		expectedCode := 500
		expectedResponseBody := `{
			"status": "Internal Server Error",
			"message": "error: failed to scan record row"
		}`

		request := httptest.NewRequest(http.MethodGet, "/localities/reportSellers", nil)
		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		service.On("GetSellerReports", mock.Anything, 0).Return([]responses.LocalitySellerReport(nil), error_message.ErrFailedToScan)

		handler := handlers.NewLocalityHandler(service)
		handler.GetSellerReportByLocality(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("should return 504 when request times out", func(t *testing.T) {
		// Test: Context timeout returns gateway timeout error / Timeout de contexto retorna error de timeout de gateway
		expectedCode := 504
		expectedResponseBody := `{
			"status": "Gateway Timeout",
			"message": "context deadline exceeded"
		}`

		request := httptest.NewRequest(http.MethodGet, "/localities/reportSellers?id=1", nil)
		response := httptest.NewRecorder()

		service := tests.NewMockLocalityService()
		service.On("GetSellerReports", mock.Anything, 1).Return([]responses.LocalitySellerReport(nil), context.DeadlineExceeded)

		handler := handlers.NewLocalityHandler(service)
		handler.GetSellerReportByLocality(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
		service.AssertExpectations(t)
	})
}
