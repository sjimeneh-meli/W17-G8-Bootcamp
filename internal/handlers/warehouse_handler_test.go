package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
)

// Helper para nuevo mock+handler por test
func newWarehouseHandlerMock() (*tests.WarehouseServiceMock, *handlers.WarehouseHandler) {
	mock := &tests.WarehouseServiceMock{}
	h := handlers.NewWarehouseHandler(mock)
	return mock, h
}

func TestWarehouseHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.CreateFunc = func(ctx context.Context, w models.Warehouse) (models.Warehouse, error) {
			return w, nil
		}
		mock.ValidateCodeUniquenessFunc = func(ctx context.Context, code string) error { return nil }

		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX1","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("POST", "/warehouses", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.Create(rr, req)
		assert.Equalf(t, http.StatusCreated, rr.Code,
			"Código inesperado. Body: %s", rr.Body.String())
	})
	t.Run("fail with no contains necessary fields", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.CreateFunc = func(ctx context.Context, w models.Warehouse) (models.Warehouse, error) {
			return w, nil
		}
		mock.ValidateCodeUniquenessFunc = func(ctx context.Context, code string) error { return nil }
		body := `{"address":"Dirección","telephon":"1234","warehouse_code":"XX1","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("POST", "/warehouses", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.Create(rr, req)
		assert.Equalf(t, http.StatusUnprocessableEntity, rr.Code,
			"Código inesperado. Body: %s", rr.Body.String())
	})
	t.Run("warehouse_code already exists", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.CreateFunc = func(ctx context.Context, w models.Warehouse) (models.Warehouse, error) {
			return w, nil
		}
		mock.ValidateCodeUniquenessFunc = func(ctx context.Context, code string) error {
			return error_message.ErrAlreadyExists
		}
		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX1","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("POST", "/warehouses", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.Create(rr, req)
		assert.Equalf(t, http.StatusConflict, rr.Code,
			"Código inesperado. Body: %s", rr.Body.String())
	})

}

func TestWarehouseHandler_Read(t *testing.T) {
	t.Run("Find By Id Success", func(t *testing.T) {
		warehouseService := &tests.WarehouseServiceMock{}
		warehouseService.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{
				Id:                 id,
				Address:            "Dirección de Prueba",
				Telephone:          "123456789",
				WareHouseCode:      "WH001",
				MinimumCapacity:    100,
				MinimumTemperature: 5.0,
				LocalityId:         1,
			}, nil
		}
		handler := handlers.NewWarehouseHandler(warehouseService)
		request := httptest.NewRequest("GET", "/warehouses/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
		response := httptest.NewRecorder()
		handler.GetById(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		expectedID := 1
		assert.Equal(t, expectedID, 1)
	})

	t.Run("Find All Succes", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		warehouses := []models.Warehouse{
			{Id: 1, Address: "Dirección 1", Telephone: "123456789", WareHouseCode: "WH001", MinimumCapacity: 100, MinimumTemperature: 5.0, LocalityId: 1},
			{Id: 2, Address: "Dirección 2", Telephone: "123456789", WareHouseCode: "WH002", MinimumCapacity: 100, MinimumTemperature: 5.0, LocalityId: 2},
		}
		mock.GetAllFunc = func(ctx context.Context) ([]models.Warehouse, error) {
			return warehouses, nil
		}
		request := httptest.NewRequest("GET", "/warehouse", nil)
		response := httptest.NewRecorder()
		handler.GetAll(response, request)
		assert.Equal(t, http.StatusOK, response.Code, "Código inesperado. Body: %s", response.Body.String())
	})

	t.Run("Find All Fail", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.GetAllFunc = func(ctx context.Context) ([]models.Warehouse, error) {
			return []models.Warehouse{}, error_message.ErrNotFound
		}
		request := httptest.NewRequest("GET", "/warehouse", nil)
		response := httptest.NewRecorder()
		handler.GetAll(response, request)
		assert.Equal(t, http.StatusInternalServerError, response.Code, "Código inesperado. Body: %s", response.Body.String())
	})

	t.Run("Find By Id Fail", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{}, error_message.ErrNotFound
		}
		request := httptest.NewRequest("GET", "/warehouses/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
		response := httptest.NewRecorder()
		handler.GetById(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}
