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

// Tests comprehensivos de integración para handlers HTTP de Employee
// Validan flujo completo request-response con HTTP real (métodos, JSON, status codes)
// Ejercitan el stack completo del handler a diferencia de unit tests con mocks

// newTestRequestWithIDParam - Helper para crear requests HTTP con parámetros URL
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

// TestPostEmployee - Tests comprehensivos del endpoint POST /employees
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

		request, err := newTestRequestWithIDParam("GET", "/api/v1/employees", id, nil)
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

		request, err := newTestRequestWithIDParam("GET", "/api/v1/employees", id, nil)
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

		request, err := newTestRequestWithIDParam("GET", "/api/v1/employees", id, nil)
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

		request, err := newTestRequestWithIDParam("GET", "/api/v1/employees", id, nil)
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

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/employees", id, nil)
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

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/employees", id, invalidEmployeeRequest)
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

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/employees", id, PatchEmployeeRequest)
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

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/employees", id, PatchEmployeeRequest)
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

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/employees", id, PatchEmployeeRequest)
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

		request, err := newTestRequestWithIDParam("PATCH", "/api/v1/employees", id, PatchEmployeeRequest)
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

		request, err := newTestRequestWithIDParam("DELETE", "/api/v1/employees", id, nil)
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

		request, err := newTestRequestWithIDParam("DELETE", "/api/v1/employees", id, nil)
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

		request, err := newTestRequestWithIDParam("DELETE", "/api/v1/employees", id, nil)
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

		request, err := newTestRequestWithIDParam("DELETE", "/api/v1/employees", id, nil)
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
