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
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

		// Valid seller request JSON
		validSellerRequest := strings.NewReader(`{
			"cid": "12345678",
			"company_name": "Test Company",
			"address": "Test Address 123",
			"telephone": "555-1234",
			"locality_id": 1
		}`)

		// Valid seller with ID (as returned by service)
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
		expectedCode := 400
		expectedResponseBody := `{
			"status": "Bad Request",
    		"message": "json: cannot unmarshal string into Go struct field SellerRequest.locality_id of type int"
		}` //necesario?

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
		assert.JSONEq(t, expectedResponseBody, response.Body.String()) //necesario?

	})

	t.Run("should return 422 when JSON missing required fields", func(t *testing.T) {
		expectedCode := 422
		expectedResponseBody := `{
			"status": "Unprocessable Entity",
    		"message": "cid: cannot be blank."
		}` //necesario?
		// Invalid seller request JSON missing required field "cid"
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
		// No need to configure mock expectations since the handler should return early with 422

		handler := handlers.NewSellerHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String()) //necesario?

	})

	t.Run("should return 409 when CID already exists", func(t *testing.T) {
		expectedCode := 409
		expectedResponseBody := `{
			"status": "Conflict",
			"message": "error: resource with the provided identifier already exists"
		}` //necesario?

		// Valid seller request JSON but with existing CID
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
		// Configure mock to return conflict error when trying to save
		service.On("Save", mock.AnythingOfType("models.Seller")).Return([]models.Seller{}, error_message.ErrAlreadyExists)

		handler := handlers.NewSellerHandler(service)
		handler.Save(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String()) //necesario?
	})

}

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

		// Mock sellers list
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

		// Mock seller response
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

		// Request body with updated seller data
		requestBody := strings.NewReader(`{
			"cid": "SEL-001",
			"company_name": "Alkemy V2",
			"address": "Monroe 861",
			"telephone": "47470009",
			"locality_id": 1
		}`)

		// Expected updated seller
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

		// Request body for update attempt
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

func newTestRequestWithIDParamSeller(method, pathBase, id string, body io.Reader) (*http.Request, error) {
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
