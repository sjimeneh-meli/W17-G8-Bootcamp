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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPost(t *testing.T) {
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

		request := httptest.NewRequest("POST", "http://localhost:8080/api/v1/buyers", validBuyerRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewBuyerServiceMock()
		service.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestBuyer).Return(validBuyer, nil)

		handler := handlers.GetBuyerHandler(service)

		handler.PostBuyer()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Post Buyer fails because request Buyer isn't valid returns 400", func(t *testing.T) {
		expectedCode := 400
		expectedResponseBody := `{
		"message":"first_name: cannot be blank; id_card_number: cannot be blank; last_name: cannot be blank.", 
		"status":"Bad Request"
		}`

		invalidBuyerRequest := strings.NewReader(`{
			"idcard_number": "CARD-1001",
			"firstname": "Juan",
			"lastname": "Pérez"
		}`)

		request := httptest.NewRequest("POST", "http://localhost:8080/api/v1/buyers", invalidBuyerRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewBuyerServiceMock()
		handler := handlers.GetBuyerHandler(service)

		handler.PostBuyer()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Post Buyer fails because request Buyer has an card number id that already exists returns 409", func(t *testing.T) {
		expectedCode := 409
		expectedResponseBody := `{
		"message":"error: resource with the provided identifier already exists", 
		"status":"Conflict"
		}`

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

		serviceMock := services_test.GetNewBuyerServiceMock()
		serviceMock.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestBuyer).Return(models.Buyer{}, error_message.ErrAlreadyExists)
		handler := handlers.GetBuyerHandler(serviceMock)

		request := httptest.NewRequest("POST", "http://localhost:8080/api/v1/buyers", repeatedCardNumberBuyerRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.PostBuyer()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Post Buyer fails because of an internal server error returns 500", func(t *testing.T) {
		expectedCode := 500
		expectedResponseBody := `{
		"message":"error: an unexpected internal server error occurred", 
		"status":"Internal Server Error"
		}`
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

		serviceMock := services_test.GetNewBuyerServiceMock()
		serviceMock.On("Create", mock.AnythingOfType("*context.timerCtx"), mockRequestBuyer).Return(models.Buyer{}, error_message.ErrInternalServerError)
		handler := handlers.GetBuyerHandler(serviceMock)

		request := httptest.NewRequest("POST", "http://localhost:8080/api/v1/buyers", internalServerErrorBuyerRequest)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.PostBuyer()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

func TestGetAll(t *testing.T) {

	t.Run("error on service returns 500", func(t *testing.T) {
		expectedCode := 500
		expectedResponseBody := `{
									"message":"error: an unexpected internal server error occurred", 
									"status":"Internal Server Error"
								 }`

		service := services_test.GetNewBuyerServiceMock()
		service.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(map[int]models.Buyer{}, error_message.ErrInternalServerError)
		handler := handlers.GetBuyerHandler(service)

		request := httptest.NewRequest("GET", "http://localhost:8080/api/v1/buyers", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetAll()(response, request)

		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Successfully get all buyers", func(t *testing.T) {
		expectedCode := 200

		mockBuyers := map[int]models.Buyer{
			1: {Id: 1, CardNumberId: "CARD-1001", FirstName: "Juan", LastName: "Pérez"},
			2: {Id: 2, CardNumberId: "CARD-1002", FirstName: "María", LastName: "Gómez"},
			3: {Id: 3, CardNumberId: "CARD-1003", FirstName: "Carlos", LastName: "López"},
		}
		service := services_test.GetNewBuyerServiceMock()
		service.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(mockBuyers, nil)
		handler := handlers.GetBuyerHandler(service)

		request := httptest.NewRequest("GET", "http://localhost:8080/api/v1/buyers", nil)
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetAll()(response, request)

		require.Equal(t, expectedCode, response.Code)
	})
}

func TestGetById(t *testing.T) {
	t.Run("Get By Id fails because request buyer id doesn't exists returns 404", func(t *testing.T) {
		id := "100"
		numberId := 100
		expectedCode := 404
		expectedResponseBody := `{
									"message":"error: the requested resource was not found", 
									"status":"Not Found"
								 }`

		serviceMock := services_test.GetNewBuyerServiceMock()
		serviceMock.On("GetById", mock.AnythingOfType("*context.timerCtx"), numberId).Return(models.Buyer{}, error_message.ErrNotFound)

		handler := handlers.GetBuyerHandler(serviceMock)

		request, err := newTestRequestWithIDParam("GET", "/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetById()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())

	})
	t.Run("Get By Id fails because request id parameter isn't a number returns 400", func(t *testing.T) {
		id := "100a"

		expectedCode := 400
		expectedResponseBody := `{
									"message":"error: the provided input is invalid or missing required fields", 
									"status":"Bad Request"
								 }`

		serviceMock := services_test.GetNewBuyerServiceMock()
		handler := handlers.GetBuyerHandler(serviceMock)

		request, err := newTestRequestWithIDParam("GET", "/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetById()(response, request)

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

		serviceMock := services_test.GetNewBuyerServiceMock()
		serviceMock.On("GetById", mock.AnythingOfType("*context.timerCtx"), numberId).Return(models.Buyer{}, error_message.ErrInternalServerError)

		handler := handlers.GetBuyerHandler(serviceMock)

		request, err := newTestRequestWithIDParam("GET", "/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.GetById()(response, request)

		assert.Equal(t, expectedCode, response.Code)
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

		serviceMock := services_test.GetNewBuyerServiceMock()
		serviceMock.On("GetById", mock.AnythingOfType("*context.timerCtx"), numberId).Return(mockBuyer, nil)

		handler := handlers.GetBuyerHandler(serviceMock)

		request, err := newTestRequestWithIDParam("GET", "/buyers", id, nil)
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

func TestPatch(t *testing.T) {
	t.Run("Patch Buyer failes because request id parameter isn't a number returns 400", func(t *testing.T) {
		id := "100a"

		expectedCode := 400
		expectedResponseBody := `{
									"message":"error: the provided input is invalid or missing required fields", 
									"status":"Bad Request"
								 }`

		serviceMock := services_test.GetNewBuyerServiceMock()
		handler := handlers.GetBuyerHandler(serviceMock)

		request, err := newTestRequestWithIDParam("PATCH", "/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		handler.PatchBuyer()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Patch Buyer fails because request buyer doesn't have any valid buyer field to update returns 400", func(t *testing.T) {
		id := "1"
		expectedCode := 400
		expectedResponseBody := `{
		"message":"at least one of id_card_number, first_name, or last_name is required", 
		"status":"Bad Request"
		}`

		invalidBuyerRequest := strings.NewReader(`{
			"idcard_number": "CARD-1001",
			"firstname": "Juan",
			"lastname": "Pérez"
		}`)

		request, err := newTestRequestWithIDParam("PATCH", "/buyers", id, invalidBuyerRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewBuyerServiceMock()
		handler := handlers.GetBuyerHandler(service)

		handler.PatchBuyer()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Patch Buyer failes because request buyer id doesn't exists returns 404", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 404
		expectedResponseBody := `{
		"message":"error: the requested resource was not found", 
		"status":"Not Found"
		}`

		PatchBuyerRequest := strings.NewReader(`{
			"id_card_number": "CARD-10012"
		}`)
		mockPatchBuyerRequest := models.Buyer{
			CardNumberId: "CARD-10012",
		}

		request, err := newTestRequestWithIDParam("PATCH", "/buyers", id, PatchBuyerRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewBuyerServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchBuyerRequest).Return(models.Buyer{}, error_message.ErrNotFound)
		handler := handlers.GetBuyerHandler(service)

		handler.PatchBuyer()(response, request)

		assert.Equal(t, expectedCode, response.Code)
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

		PatchBuyerRequest := strings.NewReader(`{
			"id_card_number": "CARD-1001"
		}`)
		mockPatchBuyerRequest := models.Buyer{
			CardNumberId: "CARD-1001",
		}

		request, err := newTestRequestWithIDParam("PATCH", "/buyers", id, PatchBuyerRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewBuyerServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchBuyerRequest).Return(models.Buyer{}, error_message.ErrAlreadyExists)
		handler := handlers.GetBuyerHandler(service)

		handler.PatchBuyer()(response, request)

		assert.Equal(t, expectedCode, response.Code)
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

		PatchBuyerRequest := strings.NewReader(`{
			"id_card_number": "CARD-1001"
		}`)
		mockPatchBuyerRequest := models.Buyer{
			CardNumberId: "CARD-1001",
		}

		request, err := newTestRequestWithIDParam("PATCH", "/buyers", id, PatchBuyerRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewBuyerServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchBuyerRequest).Return(models.Buyer{}, error_message.ErrInternalServerError)
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

		request, err := newTestRequestWithIDParam("PATCH", "/buyers", id, PatchBuyerRequest)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewBuyerServiceMock()
		service.On("Update", mock.AnythingOfType("*context.timerCtx"), idNumber, mockPatchBuyerRequest).Return(mockReturnBuyer, nil)
		handler := handlers.GetBuyerHandler(service)

		handler.PatchBuyer()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

func TestDeleteById(t *testing.T) {
	t.Run("Delete By Id fails because request buyer id parameter isn't a number returns 400", func(t *testing.T) {
		id := "100a"

		expectedCode := 400
		expectedResponseBody := `{
									"message":"error: the provided input is invalid or missing required fields", 
									"status":"Bad Request"
								 }`

		serviceMock := services_test.GetNewBuyerServiceMock()
		handler := handlers.GetBuyerHandler(serviceMock)

		request, err := newTestRequestWithIDParam("DELETE", "/buyers", id, nil)
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

		request, err := newTestRequestWithIDParam("DELETE", "/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewBuyerServiceMock()
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

		request, err := newTestRequestWithIDParam("DELETE", "/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewBuyerServiceMock()
		service.On("DeleteById", mock.AnythingOfType("*context.timerCtx"), idNumber).Return(error_message.ErrInternalServerError)
		handler := handlers.GetBuyerHandler(service)

		handler.DeleteById()(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Delete By Id successfully deletes a buyer returns 204", func(t *testing.T) {
		id := "100"
		idNumber := 100

		expectedCode := 204

		request, err := newTestRequestWithIDParam("DELETE", "/buyers", id, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		response := httptest.NewRecorder()
		response.Header().Set("Content-type", "application/json")

		service := services_test.GetNewBuyerServiceMock()
		service.On("DeleteById", mock.AnythingOfType("*context.timerCtx"), idNumber).Return(nil)

		handler := handlers.GetBuyerHandler(service)

		handler.DeleteById()(response, request)

		assert.Equal(t, expectedCode, response.Code)

	})
}

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
