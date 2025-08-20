// Package handlers_test - Tests de integración para Section Handler
// Tests HTTP comprehensivos para operaciones CRUD de la entidad Section
package handlers_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/seeders"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestPostSection - Tests para POST /sections
// Casos: 201 (creación exitosa), 422 (validación), 417 (errores de negocio/JSON)
func TestPostSection(t *testing.T) {
	t.Run("Entity: Section, Method: POST, Code: 201", func(t *testing.T) {
		expectedResponseCode := 201
		expectedResponseBody := `{
			"data": {
				"id": 1,  
				"section_number": "A-01",
				"current_capacity": 2,
				"current_temperature": 3.43,
				"maximum_capacity": 2,
				"minimum_capacity": 2,
				"minimum_temperature": 2,
				"product_type_id": 2,
				"warehouse_id": 1
			}
		}`

		actualRequest := strings.NewReader(`{
			"section_number": "A-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/sections", actualRequest)
		response := httptest.NewRecorder()
		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{
			GetByIdFunc: func(ctx context.Context, id int) (models.Warehouse, error) {
				return models.Warehouse{}, nil
			},
		}

		sectionMock.On("Create", mock.AnythingOfType("*context.timerCtx"), seeders.NewSectionModel).Return(nil)
		sectionMock.On("ExistsWithSectionNumber", mock.AnythingOfType("*context.timerCtx"), seeders.NewSectionModel.Id, seeders.NewSectionModel.SectionNumber).Return(false)
		vldMock.On("ValidateSectionRequestStruct", seeders.NewSectionRequest).Return(nil)

		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)
		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: POST, Code: 417", func(t *testing.T) {
		// Test error JSON: tipo incorrecto para section_number
		expectedResponseCode := 417
		expectedResponseBody := `{
    		"status": "Expectation Failed",
    		"message": "json: cannot unmarshal number into Go struct field SectionRequest.section_number of type string"
		}`

		actualRequest := strings.NewReader(`{
			"section_number": 1,
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/sections", actualRequest)
		response := httptest.NewRecorder()
		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{
			GetByIdFunc: func(ctx context.Context, id int) (models.Warehouse, error) {
				return models.Warehouse{}, nil
			},
		}

		sectionMock.On("Create", mock.AnythingOfType("*context.timerCtx"), seeders.NewSectionModel).Return(nil)
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)
		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: POST, Code: 409", func(t *testing.T) {
		// Test conflicto: section_number duplicado
		expectedResponseCode := 409
		expectedResponseBody := `{
    		"status": "Conflict",
    		"message": "already exist a section with the same number"
		}`

		actualRequest := strings.NewReader(`{
			"section_number": "B-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/sections", actualRequest)
		response := httptest.NewRecorder()
		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{
			GetByIdFunc: func(ctx context.Context, id int) (models.Warehouse, error) {
				return models.Warehouse{}, nil
			},
		}

		sectionMock.On("ExistsWithSectionNumber", mock.AnythingOfType("*context.timerCtx"), 0, "B-01").Return(true)
		vldMock.On("ValidateSectionRequestStruct", seeders.NewSectionRequestTwo).Return(nil)
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)
		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: POST, Code: 404", func(t *testing.T) {
		expectedResponseCode := 404
		expectedResponseBody := `{
    		"status": "Not Found",
    		"message": "error: a required dependent entity was not found"
		}`

		actualRequest := strings.NewReader(`{
			"section_number": "A-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/sections", actualRequest)
		response := httptest.NewRecorder()
		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{
			GetByIdFunc: func(ctx context.Context, id int) (models.Warehouse, error) {
				return models.Warehouse{}, error_message.ErrDependencyNotFound
			},
		}

		sectionMock.On("Create", mock.AnythingOfType("*context.timerCtx"), seeders.NewSectionModel).Return(nil)
		vldMock.On("ValidateSectionRequestStruct", seeders.NewSectionRequest).Return(nil)

		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)
		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: POST, Code: 422", func(t *testing.T) {
		// Test validación: campo section_number faltante
		expectedResponseCode := 422
		expectedResponseBody := `{
    		"status": "Unprocessable Entity",
    		"message": "section_number: cannot be blank."
		}`

		actualRequest := strings.NewReader(`{
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/sections", actualRequest)
		response := httptest.NewRecorder()
		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{
			GetByIdFunc: func(ctx context.Context, id int) (models.Warehouse, error) {
				return models.Warehouse{}, nil
			},
		}

		sectionMock.On("Create", mock.AnythingOfType("*context.timerCtx"), actualRequest).Return(nil)
		vldMock.On("ValidateSectionRequestStruct", seeders.NewSectionRequestWithoutNumber).Return(errors.New("section_number: cannot be blank."))
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)
		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: POST, Code: 417", func(t *testing.T) {
		// Test validación: campo section_number faltante
		expectedResponseCode := 417
		expectedResponseBody := `{
    		"status": "Expectation Failed",
    		"message": "invalid character 'a' looking for beginning of value"
		}`

		actualRequest := strings.NewReader(`{
			"section_number": "A-01",
			"current_capacity": a,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/sections", actualRequest)
		response := httptest.NewRecorder()
		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{
			GetByIdFunc: func(ctx context.Context, id int) (models.Warehouse, error) {
				return models.Warehouse{}, nil
			},
		}

		sectionMock.On("Create", mock.AnythingOfType("*context.timerCtx"), actualRequest).Return(nil)
		vldMock.On("ValidateSectionRequestStruct", seeders.NewSectionModel).Return(nil)
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)
		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// TestGetAllSections - Tests para GET /sections
// Casos: 200 (recuperación exitosa), 404 (sin resultados)
func TestGetAllSections(t *testing.T) {
	t.Run("Entity: Section, Method: Get, Code: 200", func(t *testing.T) {
		expectedCode := 200

		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(seeders.Sections, nil)
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodGet, "/sections", nil)
		response := httptest.NewRecorder()

		handler.GetAll(response, request)

		require.Equal(t, expectedCode, response.Code)
	})

	t.Run("Entity: Section, Method: Get, Code: 404", func(t *testing.T) {
		expectedCode := 404

		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(seeders.Sections, error_message.ErrNotFound)
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodGet, "/sections", nil)
		response := httptest.NewRecorder()

		handler.GetAll(response, request)

		require.Equal(t, expectedCode, response.Code)
	})
}

// TestGetByIdSection - Tests para GET /sections/{id}
// Casos: 200 (encontrado), 404 (no encontrado), 417 (ID inválido)
func TestGetByIdSection(t *testing.T) {
	t.Run("Entity: Section, Method: Get, Code: 200", func(t *testing.T) {
		expectedCode := 200
		expectedResponseBody := `{
			"data": {
				"id": 1, 
				"section_number": "A-01",
				"current_capacity": 1,
				"current_temperature": 3.43,
				"maximum_capacity": 1,
				"minimum_capacity": 1,
				"minimum_temperature": 1,
				"product_type_id": 1,
				"warehouse_id": 1
			}
		}`

		id := "1"

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), seeders.Sections[0].Id).Return(seeders.Sections[0], nil)
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sections/%s", id), nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.GetByID(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: Get, Code: 404", func(t *testing.T) {
		expectedCode := 404
		id := "100"
		sectionId := 100
		expectedResponseBody := `{
			"message":"error: the requested resource was not found", 
			"status":"Not Found"
		}`

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(&models.Section{}, error_message.ErrNotFound)
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sections/%s", id), nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.GetByID(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: Get, Code: 417", func(t *testing.T) {
		// Test ID inválido (formato no numérico)
		expectedCode := 417
		id := "100a"
		sectionId := "100s"
		expectedResponseBody := `{
			"message":"strconv.Atoi: parsing \"100a\": invalid syntax", 
			"status":"Expectation Failed"
		}`

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(&models.Section{}, error_message.ErrNotFound)
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sections/%s", id), nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.GetByID(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// TestUpdate - Tests para PUT /sections/{id}
// Casos: 200 (actualización exitosa), 404 (no encontrado), 417 (errores), 422 (validación)
func TestUpdate(t *testing.T) {
	t.Run("Entity: Section, Method: Update, Code: 200", func(t *testing.T) {
		id := "1"

		expectedCode := 200
		expectedResponseBody := `{
			"data": {
				"id": 1, 
				"section_number": "A-01",
				"current_capacity": 2,
				"current_temperature": 3.43,
				"maximum_capacity": 2,
				"minimum_capacity": 2,
				"minimum_temperature": 2,
				"product_type_id": 2,
				"warehouse_id": 1
			}
		}`

		sectionRequest := strings.NewReader(`{
			"section_number": "A-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("ExistsWithSectionNumber", mock.AnythingOfType("*context.timerCtx"), seeders.Sections[0].Id, seeders.Sections[0].SectionNumber).Return(false)
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), seeders.Sections[0].Id).Return(seeders.Sections[0], nil)
		sectionMock.On("Update", mock.AnythingOfType("*context.timerCtx"), seeders.Sections[0]).Return(nil)
		vldMock.On("ValidateSectionRequestStruct", seeders.NewSectionRequest).Return(nil)

		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/sections/%s", id), sectionRequest)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: Update, Code: 417 (Update)", func(t *testing.T) {
		id := "1"

		expectedCode := 417
		expectedResponseBody := `{
			"message":"service error", 
			"status":"Expectation Failed"
		}`

		sectionRequest := strings.NewReader(`{
			"section_number": "A-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("ExistsWithSectionNumber", mock.AnythingOfType("*context.timerCtx"), seeders.Sections[0].Id, seeders.Sections[0].SectionNumber).Return(false)
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), seeders.Sections[0].Id).Return(seeders.Sections[0], nil)
		sectionMock.On("Update", mock.AnythingOfType("*context.timerCtx"), seeders.Sections[0]).Return(errors.New("service error"))
		vldMock.On("ValidateSectionRequestStruct", seeders.NewSectionRequest).Return(nil)

		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/sections/%s", id), sectionRequest)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: Update, Code: 404", func(t *testing.T) {
		id := "1000"
		sectionId := 1000

		expectedCode := 404
		expectedResponseBody := `{
			"message":"error: the requested resource was not found", 
			"status":"Not Found"
		}`

		sectionRequest := strings.NewReader(`{
			"section_number": "K-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(&models.Section{}, error_message.ErrNotFound)

		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/sections/%s", id), sectionRequest)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: Update, Code: 417", func(t *testing.T) {
		// Test ID inválido
		id := "1000a"
		sectionId := "1000a"

		expectedCode := 417
		expectedResponseBody := `{
			"message":"strconv.Atoi: parsing \"1000a\": invalid syntax", 
			"status":"Expectation Failed"
		}`

		sectionRequest := strings.NewReader(`{
			"section_number": "K-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(&models.Section{}, error_message.ErrNotFound)

		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/sections/%s", id), sectionRequest)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: Update, Code: 422", func(t *testing.T) {
		// Test validación: campo current_capacity faltante
		id := "1"
		sectionId := 1

		expectedCode := 422
		expectedResponseBody := `{
			"message":"section_number: cannot be blank.", 
			"status":"Unprocessable Entity"
		}`

		sectionRequest := strings.NewReader(`{
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(seeders.Sections[0], nil)
		vldMock.On("ValidateSectionRequestStruct", seeders.NewSectionRequestWithoutNumber).Return(errors.New("section_number: cannot be blank."))

		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/sections/%s", id), sectionRequest)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: Update, Code: 409", func(t *testing.T) {
		// Test conflicto: section_number duplicado
		id := "1"

		expectedCode := 409
		expectedResponseBody := `{
			"message":"already exist a section with the same number", 
			"status":"Conflict"
		}`

		sectionRequest := strings.NewReader(`{
			"section_number": "B-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), seeders.Sections[0].Id).Return(seeders.Sections[0], nil)
		sectionMock.On("ExistsWithSectionNumber", mock.AnythingOfType("*context.timerCtx"), seeders.Sections[0].Id, "B-01").Return(false)
		vldMock.On("ValidateSectionRequestStruct", seeders.NewSectionRequestTwo).Return(nil)
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/sections/%s", id), sectionRequest)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: Update, Code: 417 (Bad Request)", func(t *testing.T) {
		// Test error JSON: tipo incorrecto
		id := "1"
		sectionId := 1

		expectedCode := 417
		expectedResponseBody := `{
			"message":"json: cannot unmarshal number into Go struct field SectionRequest.section_number of type string", 
			"status":"Expectation Failed"
		}`

		sectionRequest := strings.NewReader(`{
			"section_number": 1,
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		section := &models.Section{
			Id: 1,
		}

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(section, nil)

		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/sections/%s", id), sectionRequest)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.Update(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

// TestDeleteByIdSection - Tests para DELETE /sections/{id}
// Casos: 204 (eliminación exitosa), 404 (no encontrado), 417 (ID inválido)
func TestDeleteByIdSection(t *testing.T) {
	t.Run("Entity: Section, Method: Delete, Code: 200", func(t *testing.T) {
		id := "1"
		sectionId := 1

		expectedCode := 204

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("DeleteByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(nil)
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/sections/%s", id), nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.DeleteByID(response, request)

		assert.Equal(t, expectedCode, response.Code)
	})

	t.Run("Entity: Section, Method: Delete, Code: 404", func(t *testing.T) {
		expectedCode := 404
		id := "1000"
		sectionId := 1000
		expectedResponseBody := `{
			"message":"error: the requested resource was not found", 
			"status":"Not Found"
		}`

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("DeleteByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(error_message.ErrNotFound)
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/sections/%s", id), nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.DeleteByID(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: Delete, Code: 417", func(t *testing.T) {
		// Test ID inválido
		expectedCode := 417
		id := "1000a"
		sectionId := "1000a"
		expectedResponseBody := `{
			"message":"strconv.Atoi: parsing \"1000a\": invalid syntax", 
			"status":"Expectation Failed"
		}`

		sectionMock := tests.GetSectionServiceMock()
		vldMock := tests.GetSectionValidationMock()
		warehouseMock := &tests.WarehouseServiceMock{}
		sectionMock.On("DeleteByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(error_message.ErrNotFound)
		handler := handlers.GetSectionHandler(sectionMock, warehouseMock, vldMock)

		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/sections/%s", id), nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		handler.DeleteByID(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}
