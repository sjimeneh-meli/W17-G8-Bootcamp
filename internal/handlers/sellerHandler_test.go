// Package handlers_test - Tests de integración para Seller Handler
// Tests HTTP comprehensivos para operaciones CRUD de la entidad Seller
package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"errors"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestPostSeller - Tests para POST /sellers
// Casos: 201 (creación exitosa), 400 (JSON inválido), 422 (validación), 409 (CID duplicado)
func TestPostSeller(t *testing.T) {
	t.Run("should return 201 and a seller", func(t *testing.T) {
		expectedResponseBody := `{
			"data": {
				"id": 1000,
				"cid": "12345678",
				"company_name": "Test Company",
				"address": "Test Address 123",
				"telephone": "555-1234",
				"locality_id": 1
			}
		}`
		expectedCode := 201

		validSellerRequest := strings.NewReader(`{
			"cid": "12345678",
			"company_name": "Test Company",
			"address": "Test Address 123",
			"telephone": "555-1234",
			"locality_id": 1
		}`)

		validSeller := models.Seller{
			Id:          1000,
			CID:         "12345678",
			CompanyName: "Test Company",
			Address:     "Test Address 123",
			Telephone:   "555-1234",
			LocalityID:  1,
		}

		request := httptest.NewRequest(http.MethodPost, "/sellers", validSellerRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		service := tests.NewMockSellerService()
		service.On("Save", mock.AnythingOfType("models.Seller")).Return([]models.Seller{validSeller}, nil)

		handler := handlers.NewSellerHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 400 when JSON has incorrect field", func(t *testing.T) {
		// Test error JSON: tipo incorrecto para locality_id
		expectedCode := 400
		expectedResponseBody := `{
			"status": "Bad Request",
    		"message": "json: cannot unmarshal string into Go struct field SellerRequest.locality_id of type int"
		}`

		invalidSellerRequest := strings.NewReader(`{
			"cid": "12345678",
			"company_name": "Test Company",
			"address": "Test Address 123",
			"telephone": "555-1234",
			"locality_id": "invalid_id_type"
		}`)

		request := httptest.NewRequest("POST", "/api/v1/sellers", invalidSellerRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		service := tests.NewMockSellerService()

		handler := handlers.NewSellerHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 422 when JSON missing required fields", func(t *testing.T) {
		// Test validación: campo cid faltante
		expectedCode := 422
		expectedResponseBody := `{
			"status": "Unprocessable Entity",
    		"message": "cid: cannot be blank."
		}`
		invalidSellerRequest := strings.NewReader(`{
			"company_name": "Test Company",
			"address": "Test Address 123",
			"telephone": "555-1234",
			"locality_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/sellers", invalidSellerRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		service := tests.NewMockSellerService()

		handler := handlers.NewSellerHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 409 when CID already exists", func(t *testing.T) {
		// Test conflicto: CID duplicado
		expectedCode := 409
		expectedResponseBody := `{
			"status": "Conflict",
			"message": "error: resource with the provided identifier already exists"
		}`

		existingCIDRequest := strings.NewReader(`{
			"cid": "SEL-001",
			"company_name": "Test Company",
			"address": "Test Address 123",
			"telephone": "555-1234",
			"locality_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/sellers", existingCIDRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		service := tests.NewMockSellerService()
		service.On("Save", mock.AnythingOfType("models.Seller")).Return([]models.Seller{}, error_message.ErrAlreadyExists)

		handler := handlers.NewSellerHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// TestGet - Tests para GET /sellers y GET /sellers/{id}
// Casos: 200 (colección y individual), 404 (no encontrado)
func TestGet(t *testing.T) {
	t.Run("should return 200 and list of sellers", func(t *testing.T) {
		expectedCode := 200
		expectedResponseBody := `{
			"data": [
				{
					"id": 1,
					"cid": "SEL-001",
					"company_name": "Company One",
					"address": "Address One 123",
					"telephone": "555-0001",
					"locality_id": 1
				},
				{
					"id": 2,
					"cid": "SEL-002",
					"company_name": "Company Two",
					"address": "Address Two 456",
					"telephone": "555-0002",
					"locality_id": 2
				}
			]
		}`

		mockSellers := []models.Seller{
			{
				Id:          1,
				CID:         "SEL-001",
				CompanyName: "Company One",
				Address:     "Address One 123",
				Telephone:   "555-0001",
				LocalityID:  1,
			},
			{
				Id:          2,
				CID:         "SEL-002",
				CompanyName: "Company Two",
				Address:     "Address Two 456",
				Telephone:   "555-0002",
				LocalityID:  2,
			},
		}

		request := httptest.NewRequest(http.MethodGet, "/sellers", nil)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		service := tests.NewMockSellerService()
		service.On("GetAll").Return(mockSellers, nil)

		handler := handlers.NewSellerHandler(service)
		handler.GetAll(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 404 when seller does not exist", func(t *testing.T) {
		id := "100"
		numberId := 100
		expectedCode := 404
		expectedResponseBody := `{
									"message":"error: the requested resource was not found", 
									"status":"Not Found"
								 }`

		serviceMock := tests.NewMockSellerService()
		serviceMock.On("GetById", numberId).Return(models.Seller{}, error_message.ErrNotFound)

		handler := handlers.NewSellerHandler(serviceMock)

		request, err := newTestRequestWithIDParam(http.MethodGet, "/sellers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetById(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 200 and a seller", func(t *testing.T) {
		expectedCode := 200
		id := "1"
		numberId := 1
		expectedResponseBody := `{
				"data": {
					"id": 1,
					"cid": "SEL-001",
					"company_name": "Distribuidora Lácteos del Campo",
					"address": "Calle Falsa 123",
					"telephone": "555-0101",
					"locality_id": 1
				}
			}`

		expectedSeller := models.Seller{
			Id:          1,
			CID:         "SEL-001",
			CompanyName: "Distribuidora Lácteos del Campo",
			Address:     "Calle Falsa 123",
			Telephone:   "555-0101",
			LocalityID:  1,
		}

		serviceMock := tests.NewMockSellerService()
		serviceMock.On("GetById", numberId).Return(expectedSeller, nil)

		handler := handlers.NewSellerHandler(serviceMock)

		request, err := newTestRequestWithIDParam(http.MethodGet, "/sellers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetById(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// TestPut - Tests para PUT /sellers/{id}
// Casos: 200 (actualización exitosa), 404 (no encontrado)
func TestPut(t *testing.T) {
	t.Run("should return 200 and update a seller", func(t *testing.T) {
		expectedCode := 200
		id := "1"
		numberId := 1
		expectedResponseBody := `{
			"data": [
				{
					"id": 1,
					"cid": "SEL-001",
					"company_name": "Alkemy V2",
					"address": "Monroe 861",
					"telephone": "47470009",
					"locality_id": 1
				}
			]
		}`

		requestBody := strings.NewReader(`{
			"cid": "SEL-001",
			"company_name": "Alkemy V2",
			"address": "Monroe 861",
			"telephone": "47470009",
			"locality_id": 1
		}`)

		updatedSeller := models.Seller{
			Id:          1,
			CID:         "SEL-001",
			CompanyName: "Alkemy V2",
			Address:     "Monroe 861",
			Telephone:   "47470009",
			LocalityID:  1,
		}

		serviceMock := tests.NewMockSellerService()
		serviceMock.On("Update", numberId, mock.AnythingOfType("models.Seller")).Return([]models.Seller{updatedSeller}, nil)

		handler := handlers.NewSellerHandler(serviceMock)

		request, err := newTestRequestWithIDParam(http.MethodPatch, "/sellers", id, requestBody)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
	t.Run("should return 404 when seller does not exist", func(t *testing.T) {
		expectedCode := 404
		id := "100"
		numberId := 100
		expectedResponseBody := `{
			"status": "Not Found",
			"message": "error: the requested resource was not found"
		}`

		requestBody := strings.NewReader(`{
			"cid": "SEL-100",
			"company_name": "Test Company",
			"address": "Test Address",
			"telephone": "555-0000",
			"locality_id": 1
		}`)

		serviceMock := tests.NewMockSellerService()
		serviceMock.On("Update", numberId, mock.AnythingOfType("models.Seller")).Return([]models.Seller{}, error_message.ErrNotFound)

		handler := handlers.NewSellerHandler(serviceMock)
		request, err := newTestRequestWithIDParam(http.MethodPatch, "/sellers", id, requestBody)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// TestDelete - Tests para DELETE /sellers/{id}
// Casos: 204 (eliminación exitosa), 404 (no encontrado)
func TestDelete(t *testing.T) {
	t.Run("should return 404 when seller does not exist", func(t *testing.T) {
		expectedCode := 404
		id := "999"
		numberId := 999
		expectedResponseBody := `{
			"status": "Not Found",
			"message": "error: the requested resource was not found"
		}`

		serviceMock := tests.NewMockSellerService()
		serviceMock.On("Delete", numberId).Return(error_message.ErrNotFound)

		handler := handlers.NewSellerHandler(serviceMock)

		request, err := newTestRequestWithIDParam(http.MethodDelete, "/sellers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		handler.Delete(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 204 when seller is deleted successfully", func(t *testing.T) {
		expectedCode := 204
		id := "1"
		numberId := 1

		serviceMock := tests.NewMockSellerService()
		serviceMock.On("Delete", numberId).Return(nil)

		handler := handlers.NewSellerHandler(serviceMock)

		request, err := newTestRequestWithIDParam(http.MethodDelete, "/sellers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		handler.Delete(response, request)

		assert.Equal(t, expectedCode, response.Code)
	})
}

// TestGetAllError - Test casos de error para GetAll
func TestGetAllError(t *testing.T) {
	t.Run("should return 404 when service returns error", func(t *testing.T) {
		expectedCode := 404
		expectedResponseBody := `{
			"status": "Not Found",
			"message": "error: the requested resource was not found"
		}`

		request := httptest.NewRequest(http.MethodGet, "/sellers", nil)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		service := tests.NewMockSellerService()
		service.On("GetAll").Return([]models.Seller{}, error_message.ErrNotFound)

		handler := handlers.NewSellerHandler(service)
		handler.GetAll(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// TestGetByIdInvalidID - Test ID inválido para GetById
func TestGetByIdInvalidID(t *testing.T) {
	t.Run("should return 400 when ID is not a valid number", func(t *testing.T) {
		expectedCode := 400
		id := "invalid_id"

		serviceMock := tests.NewMockSellerService()

		handler := handlers.NewSellerHandler(serviceMock)

		request, err := newTestRequestWithIDParam(http.MethodGet, "/sellers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetById(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.Contains(t, response.Body.String(), "strconv.Atoi")
	})
}

// TestSaveAdditionalErrors - Tests adicionales de error para Save
func TestSaveAdditionalErrors(t *testing.T) {
	t.Run("should return 422 when dependency not found", func(t *testing.T) {
		// Test error: locality_id no existe
		expectedCode := 422
		expectedResponseBody := `{
			"status": "Unprocessable Entity",
			"message": "error: a required dependent entity was not found"
		}`

		validSellerRequest := strings.NewReader(`{
			"cid": "12345678",
			"company_name": "Test Company",
			"address": "Test Address 123",
			"telephone": "555-1234",
			"locality_id": 999
		}`)

		request := httptest.NewRequest(http.MethodPost, "/sellers", validSellerRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		service := tests.NewMockSellerService()
		service.On("Save", mock.AnythingOfType("models.Seller")).Return([]models.Seller{}, error_message.ErrDependencyNotFound)

		handler := handlers.NewSellerHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 500 when internal server error occurs", func(t *testing.T) {
		expectedCode := 500
		expectedResponseBody := `{
			"status": "Internal Server Error",
			"message": "internal server error"
		}`

		validSellerRequest := strings.NewReader(`{
			"cid": "12345678",
			"company_name": "Test Company",
			"address": "Test Address 123",
			"telephone": "555-1234",
			"locality_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/sellers", validSellerRequest)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		service := tests.NewMockSellerService()
		service.On("Save", mock.AnythingOfType("models.Seller")).Return([]models.Seller{}, errors.New("internal server error"))

		handler := handlers.NewSellerHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// TestUpdateAdditionalErrors - Tests adicionales de error para Update
func TestUpdateAdditionalErrors(t *testing.T) {
	t.Run("should return 400 when ID is empty", func(t *testing.T) {
		expectedCode := 400
		expectedResponseBody := `{
			"status": "Bad Request",
			"message": "id is required"
		}`

		requestBody := strings.NewReader(`{
			"cid": "SEL-001",
			"company_name": "Test Company",
			"address": "Test Address",
			"telephone": "555-0000",
			"locality_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPatch, "/sellers/", requestBody)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		serviceMock := tests.NewMockSellerService()
		handler := handlers.NewSellerHandler(serviceMock)

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 400 when ID is not a valid number", func(t *testing.T) {
		// Test ID inválido
		expectedCode := 400
		id := "invalid_id"

		requestBody := strings.NewReader(`{
			"cid": "SEL-001",
			"company_name": "Test Company",
			"address": "Test Address",
			"telephone": "555-0000",
			"locality_id": 1
		}`)

		serviceMock := tests.NewMockSellerService()
		handler := handlers.NewSellerHandler(serviceMock)

		request, err := newTestRequestWithIDParam(http.MethodPatch, "/sellers", id, requestBody)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.Contains(t, response.Body.String(), "strconv.Atoi")
	})

	t.Run("should return 400 when JSON body is invalid", func(t *testing.T) {
		// Test JSON inválido
		expectedCode := 400
		id := "1"

		invalidJsonBody := strings.NewReader(`{
			"cid": "SEL-001",
			"company_name": "Test Company",
			"address": "Test Address",
			"telephone": "555-0000",
			"locality_id": "invalid_id_type"
		}`)

		serviceMock := tests.NewMockSellerService()
		handler := handlers.NewSellerHandler(serviceMock)

		request, err := newTestRequestWithIDParam(http.MethodPatch, "/sellers", id, invalidJsonBody)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.Contains(t, response.Body.String(), "cannot unmarshal")
	})

	t.Run("should return 500 when internal server error occurs", func(t *testing.T) {
		expectedCode := 500
		expectedResponseBody := `{
			"status": "Internal Server Error",
			"message": "internal server error"
		}`
		id := "1"
		numberId := 1

		requestBody := strings.NewReader(`{
			"cid": "SEL-001",
			"company_name": "Test Company",
			"address": "Test Address",
			"telephone": "555-0000",
			"locality_id": 1
		}`)

		serviceMock := tests.NewMockSellerService()
		serviceMock.On("Update", numberId, mock.AnythingOfType("models.Seller")).Return([]models.Seller{}, errors.New("internal server error"))

		handler := handlers.NewSellerHandler(serviceMock)
		request, err := newTestRequestWithIDParam(http.MethodPatch, "/sellers", id, requestBody)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// TestDeleteAdditionalErrors - Tests adicionales de error para Delete
func TestDeleteAdditionalErrors(t *testing.T) {
	t.Run("should return 400 when ID is empty", func(t *testing.T) {
		expectedCode := 400
		expectedResponseBody := `{
			"status": "Bad Request",
			"message": "id is required"
		}`

		request := httptest.NewRequest(http.MethodDelete, "/sellers/", nil)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		serviceMock := tests.NewMockSellerService()
		handler := handlers.NewSellerHandler(serviceMock)

		handler.Delete(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("should return 400 when ID is not a valid number", func(t *testing.T) {
		// Test ID inválido
		expectedCode := 400
		id := "invalid_id"

		serviceMock := tests.NewMockSellerService()
		handler := handlers.NewSellerHandler(serviceMock)

		request, err := newTestRequestWithIDParam(http.MethodDelete, "/sellers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		handler.Delete(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.Contains(t, response.Body.String(), "strconv.Atoi")
	})

	t.Run("should return 500 when internal server error occurs", func(t *testing.T) {
		expectedCode := 500
		expectedResponseBody := `{
			"status": "Internal Server Error",
			"message": "internal server error"
		}`
		id := "1"
		numberId := 1

		serviceMock := tests.NewMockSellerService()
		serviceMock.On("Delete", numberId).Return(errors.New("internal server error"))

		handler := handlers.NewSellerHandler(serviceMock)

		request, err := newTestRequestWithIDParam(http.MethodDelete, "/sellers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")

		handler.Delete(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}
