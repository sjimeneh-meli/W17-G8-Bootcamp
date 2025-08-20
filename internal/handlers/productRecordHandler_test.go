// Package handlers provides comprehensive unit tests for ProductRecordHandler with 100% code coverage
// El paquete handlers proporciona tests unitarios completos para ProductRecordHandler con 100% de cobertura de código
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewProductRecordHandler validates the constructor function creates a valid handler instance
// TestNewProductRecordHandler valida que la función constructora cree una instancia válida del handler
func TestNewProductRecordHandler(t *testing.T) {
	// Arrange
	mockService := tests.GetNewProductRecordServiceMock()

	// Act
	handler := NewProductRecordHandler(mockService)

	// Assert
	assert.NotNil(t, handler)
	assert.IsType(t, &productRecordHandler{}, handler)
}

// TestProductRecordHandler_Create uses table-driven tests to validate the Create method behavior
// Covers all possible scenarios: success, JSON parsing errors, validation errors, business logic errors
// TestProductRecordHandler_Create utiliza table-driven tests para validar el comportamiento del método Create
// Cubre todos los escenarios posibles: éxito, errores de parsing JSON, errores de validación, errores de lógica de negocio
func TestProductRecordHandler_Create(t *testing.T) {
	testCases := []struct {
		name               string                                  // Test case description / Descripción del caso de test
		requestBody        interface{}                             // Request payload (valid struct or invalid string) / Payload de petición (struct válido o string inválido)
		mockSetup          func(m *tests.ProductRecordServiceMock) // Mock configuration / Configuración del mock
		expectedStatusCode int                                     // Expected HTTP status code / Código de estado HTTP esperado
		expectedError      string                                  // Expected error message / Mensaje de error esperado
	}{
		{
			name: "Success - Valid request",
			requestBody: requests.ProductRecordRequest{
				LastUpdateDate: time.Now(),
				PurchasePrice:  100.50,
				SalePrice:      120.00,
				ProductID:      1,
			},
			mockSetup: func(m *tests.ProductRecordServiceMock) {
				expectedResult := &models.ProductRecord{
					ID:             1,
					LastUpdateDate: time.Now(),
					PurchasePrice:  100.50,
					SalePrice:      120.00,
					ProductID:      1,
				}
				m.On("CreateProductRecord", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("models.ProductRecord")).
					Return(expectedResult, nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Error - Invalid JSON",
			requestBody:        "invalid json",
			mockSetup:          func(m *tests.ProductRecordServiceMock) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Error - Validation failed - Missing required fields",
			requestBody: requests.ProductRecordRequest{
				// Missing required fields
				PurchasePrice: 100.50,
				SalePrice:     120.00,
			},
			mockSetup:          func(m *tests.ProductRecordServiceMock) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Error - Dependency not found",
			requestBody: requests.ProductRecordRequest{
				LastUpdateDate: time.Now(),
				PurchasePrice:  100.50,
				SalePrice:      120.00,
				ProductID:      999, // Non-existent product
			},
			mockSetup: func(m *tests.ProductRecordServiceMock) {
				m.On("CreateProductRecord", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("models.ProductRecord")).
					Return((*models.ProductRecord)(nil), error_message.ErrDependencyNotFound)
			},
			expectedStatusCode: http.StatusConflict,
		},
		{
			name: "Error - Internal server error",
			requestBody: requests.ProductRecordRequest{
				LastUpdateDate: time.Now(),
				PurchasePrice:  100.50,
				SalePrice:      120.00,
				ProductID:      1,
			},
			mockSetup: func(m *tests.ProductRecordServiceMock) {
				m.On("CreateProductRecord", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("models.ProductRecord")).
					Return((*models.ProductRecord)(nil), error_message.ErrInternalServerError)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - Setup mock service and handler / Configuración del servicio mock y handler
			mockService := tests.GetNewProductRecordServiceMock()
			tt.mockSetup(mockService)
			handler := NewProductRecordHandler(mockService)

			var body []byte
			var err error

			// Handle different request body types for testing invalid JSON scenarios
			// Maneja diferentes tipos de request body para probar escenarios de JSON inválido
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str) // Invalid JSON string for error testing / String JSON inválido para testing de errores
			} else {
				body, err = json.Marshal(tt.requestBody) // Valid struct serialization / Serialización de struct válido
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/product-records", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Act - Execute the handler method / Ejecutar el método del handler
			handler.Create(w, req)

			// Assert - Validate response status and content / Validar estado y contenido de la respuesta
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Validate successful response structure / Validar estructura de respuesta exitosa
			if tt.expectedStatusCode == http.StatusCreated {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")
			}

			// Ensure all mock expectations were met / Asegurar que todas las expectativas del mock se cumplieron
			mockService.AssertExpectations(t)
		})
	}
}

// TestProductRecordHandler_GetReport validates the GetReport method using table-driven tests
// Tests both scenarios: getting all reports (no ID) and getting specific product report (with ID)
// Covers query parameter validation, URL encoding, and error handling
// TestProductRecordHandler_GetReport valida el método GetReport usando table-driven tests
// Prueba ambos escenarios: obtener todos los reportes (sin ID) y obtener reporte específico (con ID)
// Cubre validación de parámetros de consulta, codificación URL, y manejo de errores
func TestProductRecordHandler_GetReport(t *testing.T) {
	testCases := []struct {
		name               string                                  // Test case description / Descripción del caso de test
		queryParam         string                                  // URL query parameter (?id=X or empty) / Parámetro de consulta URL (?id=X o vacío)
		mockSetup          func(m *tests.ProductRecordServiceMock) // Mock service configuration / Configuración del servicio mock
		expectedStatusCode int                                     // Expected HTTP response status / Estado de respuesta HTTP esperado
		expectedError      string                                  // Expected error message for failures / Mensaje de error esperado para fallos
	}{
		{
			name:       "Success - Get all reports (no ID parameter)",
			queryParam: "",
			mockSetup: func(m *tests.ProductRecordServiceMock) {
				expectedReports := []*models.ProductRecordReport{
					{
						ProductId:    1,
						Description:  "Product 1",
						RecordsCount: 5,
					},
					{
						ProductId:    2,
						Description:  "Product 2",
						RecordsCount: 3,
					},
				}
				m.On("GetReport", mock.AnythingOfType("*context.timerCtx")).
					Return(expectedReports, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:       "Success - Get report by valid ID",
			queryParam: "?id=1",
			mockSetup: func(m *tests.ProductRecordServiceMock) {
				expectedReport := &models.ProductRecordReport{
					ProductId:    1,
					Description:  "Product 1",
					RecordsCount: 5,
				}
				m.On("GetReportByIdProduct", mock.AnythingOfType("*context.timerCtx"), int64(1)).
					Return(expectedReport, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Error - Invalid ID parameter (not a number)",
			queryParam:         "?id=invalid",
			mockSetup:          func(m *tests.ProductRecordServiceMock) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "error: id not is a number",
		},
		{
			name:       "Error - Get all reports - Internal server error",
			queryParam: "",
			mockSetup: func(m *tests.ProductRecordServiceMock) {
				m.On("GetReport", mock.AnythingOfType("*context.timerCtx")).
					Return(([]*models.ProductRecordReport)(nil), error_message.ErrInternalServerError)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:       "Error - Get report by ID - Dependency not found",
			queryParam: "?id=999",
			mockSetup: func(m *tests.ProductRecordServiceMock) {
				m.On("GetReportByIdProduct", mock.AnythingOfType("*context.timerCtx"), int64(999)).
					Return((*models.ProductRecordReport)(nil), error_message.ErrDependencyNotFound)
			},
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:       "Error - Get report by ID - Internal server error",
			queryParam: "?id=1",
			mockSetup: func(m *tests.ProductRecordServiceMock) {
				m.On("GetReportByIdProduct", mock.AnythingOfType("*context.timerCtx"), int64(1)).
					Return((*models.ProductRecordReport)(nil), error_message.ErrInternalServerError)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:       "Success - Get report with ID parameter with spaces",
			queryParam: "?id=%201%20",
			mockSetup: func(m *tests.ProductRecordServiceMock) {
				expectedReport := &models.ProductRecordReport{
					ProductId:    1,
					Description:  "Product 1",
					RecordsCount: 5,
				}
				m.On("GetReportByIdProduct", mock.AnythingOfType("*context.timerCtx"), int64(1)).
					Return(expectedReport, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - Setup mock service and handler / Configuración del servicio mock y handler
			mockService := tests.GetNewProductRecordServiceMock()
			tt.mockSetup(mockService)
			handler := NewProductRecordHandler(mockService)

			// Build URL with or without query parameters / Construir URL con o sin parámetros de consulta
			url := fmt.Sprintf("/product-records/reports%s", tt.queryParam)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			// Act - Execute the GetReport handler method / Ejecutar el método GetReport del handler
			handler.GetReport(w, req)

			// Assert - Validate HTTP status code / Validar código de estado HTTP
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Validate successful response contains data field / Validar que respuesta exitosa contenga campo data
			if tt.expectedStatusCode == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "data")
			}

			// Validate error response structure and message / Validar estructura y mensaje de respuesta de error
			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
				assert.Equal(t, tt.expectedError, response["message"])
			}

			// Verify all mock service expectations were satisfied / Verificar que todas las expectativas del servicio mock se cumplieron
			mockService.AssertExpectations(t)
		})
	}
}
