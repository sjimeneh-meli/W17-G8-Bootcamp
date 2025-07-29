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
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPostSection(t *testing.T) {
	t.Run("Entity: Section, Method: POST, Code: 201", func(t *testing.T) {
		expectedResponseCode := 201
		expectedResponseBody := `{
			"data": {
				"id": 29,  
				"section_number": "K-01",
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
			"section_number": "K-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		sectionRequest := &models.Section{
			Id:                 0,
			SectionNumber:      "K-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		request := httptest.NewRequest(http.MethodPost, "/sections", actualRequest)
		response := httptest.NewRecorder()
		sectionMock := tests.GetSectionServiceMock()

		sectionMock.On("Create", mock.AnythingOfType("*context.timerCtx"), sectionRequest).Return(nil)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)
		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: Section, Method: POST, Code: 422", func(t *testing.T) {
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

		sectionRequest := &models.Section{
			Id:                 0,
			SectionNumber:      "K-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		request := httptest.NewRequest(http.MethodPost, "/sections", actualRequest)
		response := httptest.NewRecorder()
		sectionMock := tests.GetSectionServiceMock()

		sectionMock.On("Create", mock.AnythingOfType("*context.timerCtx"), sectionRequest).Return(nil)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)
		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())

	})

	t.Run("Entity: Section, Method: POST, Code: 409", func(t *testing.T) {
		expectedResponseCode := 417
		expectedResponseBody := `{
    		"status": "Expectation Failed",
    		"message": "already exist a section with the same number"
		}`

		actualRequest := strings.NewReader(`{
			"section_number": "K-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		sectionRequest := &models.Section{
			Id:                 0,
			SectionNumber:      "K-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		request := httptest.NewRequest(http.MethodPost, "/sections", actualRequest)
		response := httptest.NewRecorder()
		sectionMock := tests.GetSectionServiceMock()

		sectionMock.On("Create", mock.AnythingOfType("*context.timerCtx"), sectionRequest).Return(errors.New("already exist a section with the same number"))
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)
		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())

	})

	t.Run("Entity: Section, Method: POST, Code: 417", func(t *testing.T) {
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

		sectionRequest := &models.Section{
			Id:                 0,
			SectionNumber:      "K-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		request := httptest.NewRequest(http.MethodPost, "/sections", actualRequest)
		response := httptest.NewRecorder()
		sectionMock := tests.GetSectionServiceMock()

		sectionMock.On("Create", mock.AnythingOfType("*context.timerCtx"), sectionRequest).Return(nil)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)
		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())

	})
}

func TestGetAllSections(t *testing.T) {
	t.Run("Entity: Section, Method: Get, Code: 200", func(t *testing.T) {
		expectedCode := 200

		sections := []*models.Section{
			{
				Id:                 1,
				SectionNumber:      "K-01",
				CurrentCapacity:    2,
				CurrentTemperature: 3.43,
				MaximumCapacity:    2,
				MinimumCapacity:    2,
				MinimumTemperature: 2,
				ProductTypeID:      2,
				WarehouseID:        1,
			},
		}
		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(sections, nil)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

		request := httptest.NewRequest(http.MethodGet, "/sections", nil)
		response := httptest.NewRecorder()

		handler.GetAll(response, request)

		require.Equal(t, expectedCode, response.Code)
	})

	t.Run("Entity: Section, Method: Get, Code: 404", func(t *testing.T) {
		expectedCode := 404

		sections := []*models.Section{}
		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(sections, error_message.ErrNotFound)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

		request := httptest.NewRequest(http.MethodGet, "/sections", nil)
		response := httptest.NewRecorder()

		handler.GetAll(response, request)

		require.Equal(t, expectedCode, response.Code)
	})

}

func TestGetByIdSection(t *testing.T) {

	t.Run("Entity: Section, Method: Get, Code: 200", func(t *testing.T) {
		expectedCode := 200
		expectedResponseBody := `{
			"data": {
				"id": 1, 
				"section_number": "K-01",
				"current_capacity": 2,
				"current_temperature": 3.43,
				"maximum_capacity": 2,
				"minimum_capacity": 2,
				"minimum_temperature": 2,
				"product_type_id": 2,
				"warehouse_id": 1
			}
		}`

		id := "1"
		sectionId := 1
		section := &models.Section{
			Id:                 1,
			SectionNumber:      "K-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(section, nil)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

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
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(&models.Section{}, error_message.ErrNotFound)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

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
		expectedCode := 417
		id := "100a"
		sectionId := "100s"
		expectedResponseBody := `{
			"message":"strconv.Atoi: parsing \"100a\": invalid syntax", 
			"status":"Expectation Failed"
		}`

		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(&models.Section{}, error_message.ErrNotFound)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

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

func TestUpdate(t *testing.T) {
	t.Run("Entity: Section, Method: Update, Code: 200", func(t *testing.T) {
		id := "1"
		sectionId := 1

		expectedCode := 200
		expectedResponseBody := `{
			"data": {
				"id": 1, 
				"section_number": "K-01",
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
			"section_number": "K-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		section := &models.Section{
			Id:                 1,
			SectionNumber:      "K-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(section, nil)
		sectionMock.On("Update", mock.AnythingOfType("*context.timerCtx"), section).Return(nil)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

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

		section := &models.Section{
			Id:                 1,
			SectionNumber:      "K-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(&models.Section{}, error_message.ErrNotFound)
		sectionMock.On("Update", mock.AnythingOfType("*context.timerCtx"), sectionId, section).Return(error_message.ErrNotFound)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

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

		section := &models.Section{
			Id:                 1,
			SectionNumber:      "K-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(&models.Section{}, error_message.ErrNotFound)
		sectionMock.On("Update", mock.AnythingOfType("*context.timerCtx"), sectionId, section).Return(error_message.ErrNotFound)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

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
		id := "1"
		sectionId := 1

		expectedCode := 422
		expectedResponseBody := `{
			"message":"current_capacity: cannot be blank.", 
			"status":"Unprocessable Entity"
		}`

		sectionRequest := strings.NewReader(`{
			"section_number": "K-01",
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 1
		}`)

		section := &models.Section{
			Id:                 1,
			SectionNumber:      "K-01",
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(section, nil)
		sectionMock.On("Update", mock.AnythingOfType("*context.timerCtx"), section).Return(nil)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

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
		id := "1"
		sectionId := 1

		expectedCode := 417
		expectedResponseBody := `{
			"message":"already exist a section with the same number", 
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

		section := &models.Section{
			Id: 1,
		}

		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(section, nil)
		sectionMock.On("Update", mock.AnythingOfType("*context.timerCtx"), section).Return(errors.New("already exist a section with the same number"))
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

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
		sectionMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(section, nil)
		sectionMock.On("Update", mock.AnythingOfType("*context.timerCtx"), section).Return(errors.New("json: cannot unmarshal number into Go struct field SectionRequest.section_number of type string"))
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

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

func TestDeleteByIdSection(t *testing.T) {
	t.Run("Entity: Section, Method: Delete, Code: 200", func(t *testing.T) {
		id := "1"
		sectionId := 1

		expectedCode := 204

		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("DeleteByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(nil)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

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
		sectionMock.On("DeleteByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(error_message.ErrNotFound)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

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
		expectedCode := 417
		id := "1000a"
		sectionId := "1000a"
		expectedResponseBody := `{
			"message":"strconv.Atoi: parsing \"1000a\": invalid syntax", 
			"status":"Expectation Failed"
		}`

		sectionMock := tests.GetSectionServiceMock()
		sectionMock.On("DeleteByID", mock.AnythingOfType("*context.timerCtx"), sectionId).Return(error_message.ErrNotFound)
		vld := validations.GetSectionValidation()
		handler := handlers.GetSectionHandler(sectionMock, nil, vld)

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
