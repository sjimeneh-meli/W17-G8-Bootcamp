package tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/config"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPost(t *testing.T) {
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
		expectedResponseCode := 409
		expectedResponseBody := `{
    		"status": "Conflict",
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
}

func TestGetAll(t *testing.T) {
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

}

func TestGetById(t *testing.T) {

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
}

func TestDeleteById(t *testing.T) {
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
}

func TestUpdateOkSection(t *testing.T) {
	db := database.InitDB(&config.Config{
		Database: config.Database{
			DBUser:     "root",
			DBPassword: "test",
			DBHost:     "localhost",
			DBPort:     "3306",
			DBName:     "productos_frescos",
		},
	})
	defer db.Close()

	rp := repositories.GetSectionRepository(db)
	warehouseRp := repositories.NewWarehouseRepository(db)
	srv := services.GetSectionService(rp)
	warehouseSrv := services.NewWarehouseService(warehouseRp)
	vld := validations.GetSectionValidation()
	hdCreateFunc := handlers.GetSectionHandler(srv, warehouseSrv, vld).Update

	request := httptest.NewRequest(http.MethodPut, "/sections/2", strings.NewReader(`
		{
			"section_number": "AA-01",
			"current_capacity": 50,
			"current_temperature": -18,
			"maximum_capacity": 100,
			"minimum_capacity": 10,
			"minimum_temperature": -22,
			"product_type_id": 2,
			"warehouse_id": 1
		}
	`))
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "2")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
	response := httptest.NewRecorder()
	hdCreateFunc(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.JSONEq(t, `{
		"data": {
			"id": 2,
			"section_number": "AA-01",
			"current_capacity": 50,
			"current_temperature": -18,
			"maximum_capacity": 100,
			"minimum_capacity": 10,
			"minimum_temperature": -22,
			"product_type_id": 2,
			"warehouse_id": 1
		}
	}`, response.Body.String())
}

func TestUpdateNotFoundSection(t *testing.T) {
	db := database.InitDB(&config.Config{
		Database: config.Database{
			DBUser:     "root",
			DBPassword: "test",
			DBHost:     "localhost",
			DBPort:     "3306",
			DBName:     "productos_frescos",
		},
	})
	defer db.Close()

	rp := repositories.GetSectionRepository(db)
	warehouseRp := repositories.NewWarehouseRepository(db)
	srv := services.GetSectionService(rp)
	warehouseSrv := services.NewWarehouseService(warehouseRp)
	vld := validations.GetSectionValidation()
	hdCreateFunc := handlers.GetSectionHandler(srv, warehouseSrv, vld).Update

	request := httptest.NewRequest(http.MethodPut, "/sections/10000", strings.NewReader(`
		{
			"section_number": "Z-01",
			"current_capacity": 2,
			"current_temperature": 3.43,
			"maximum_capacity": 2,
			"minimum_capacity": 2,
			"minimum_temperature": 2,
			"product_type_id": 2,
			"warehouse_id": 2
		}
	`))
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "10000")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
	response := httptest.NewRecorder()
	hdCreateFunc(response, request)

	require.Equal(t, http.StatusNotFound, response.Code)
	require.JSONEq(t, `{
    	"status": "Not Found",
    	"message": "error: the requested resource was not found"
	}`, response.Body.String())
}
