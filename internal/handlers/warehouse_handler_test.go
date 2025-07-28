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
	t.Run("invalid JSON format", func(t *testing.T) {
		_, handler := newWarehouseHandlerMock()
		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX1","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2`
		req := httptest.NewRequest("POST", "/warehouses", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.Create(rr, req)
		assert.Equalf(t, http.StatusBadRequest, rr.Code,
			"Código inesperado. Body: %s", rr.Body.String())
	})
	t.Run("validate code uniqueness internal server error", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.ValidateCodeUniquenessFunc = func(ctx context.Context, code string) error {
			return error_message.ErrInternalServerError
		}
		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX1","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("POST", "/warehouses", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.Create(rr, req)
		assert.Equalf(t, http.StatusInternalServerError, rr.Code,
			"Código inesperado. Body: %s", rr.Body.String())
	})
	t.Run("create warehouse internal server error", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.CreateFunc = func(ctx context.Context, w models.Warehouse) (models.Warehouse, error) {
			return models.Warehouse{}, error_message.ErrInternalServerError
		}
		mock.ValidateCodeUniquenessFunc = func(ctx context.Context, code string) error { return nil }
		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX1","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("POST", "/warehouses", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler.Create(rr, req)
		assert.Equalf(t, http.StatusInternalServerError, rr.Code,
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

	t.Run("Find All Empty List", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.GetAllFunc = func(ctx context.Context) ([]models.Warehouse, error) {
			return []models.Warehouse{}, nil
		}
		request := httptest.NewRequest("GET", "/warehouse", nil)
		response := httptest.NewRecorder()
		handler.GetAll(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code, "Código inesperado. Body: %s", response.Body.String())
	})

	t.Run("Find All Internal Server Error", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.GetAllFunc = func(ctx context.Context) ([]models.Warehouse, error) {
			return []models.Warehouse{}, error_message.ErrInternalServerError
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

	t.Run("Find By Id Internal Server Error", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{}, error_message.ErrInternalServerError
		}
		request := httptest.NewRequest("GET", "/warehouses/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
		response := httptest.NewRecorder()
		handler.GetById(response, request)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("Find By Id Invalid ID", func(t *testing.T) {
		_, handler := newWarehouseHandlerMock()
		request := httptest.NewRequest("GET", "/warehouses/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "invalid")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
		response := httptest.NewRecorder()
		handler.GetById(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Find By Id Missing ID", func(t *testing.T) {
		_, handler := newWarehouseHandlerMock()
		request := httptest.NewRequest("GET", "/warehouses/{id}", nil)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
		response := httptest.NewRecorder()
		handler.GetById(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestWarehouseHandler_Update(t *testing.T) {
	t.Run("Update Success", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.UpdateFunc = func(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error) {
			return warehouse, nil
		}
		mock.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{
				Id:                 id,
				Address:            "Dirección Original",
				Telephone:          "123456789",
				WareHouseCode:      "WH001",
				MinimumCapacity:    100,
				MinimumTemperature: 5.0,
				LocalityId:         1,
			}, nil
		}
		mock.ValidateCodeUniquenessFunc = func(ctx context.Context, code string) error {
			return nil
		}
		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX1","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("PATCH", "/warehouses/{id}", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		rr := httptest.NewRecorder()
		handler.Update(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
	t.Run("Update Fail not found", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.UpdateFunc = func(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error) {
			return models.Warehouse{}, error_message.ErrNotFound
		}
		mock.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{}, error_message.ErrNotFound
		}
		request := httptest.NewRequest("PATCH", "/warehouses/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
		response := httptest.NewRecorder()
		handler.Update(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})
	t.Run("Update Invalid JSON", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{
				Id:                 1,
				Address:            "Dirección Original",
				Telephone:          "123456789",
				WareHouseCode:      "WH001",
				MinimumCapacity:    100,
				MinimumTemperature: 5.0,
				LocalityId:         1,
			}, nil
		}
		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX1","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2`
		req := httptest.NewRequest("PATCH", "/warehouses/{id}", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		rr := httptest.NewRecorder()
		handler.Update(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("Update Invalid ID", func(t *testing.T) {
		_, handler := newWarehouseHandlerMock()
		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX1","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("PATCH", "/warehouses/{id}", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		rr := httptest.NewRecorder()
		handler.Update(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("Update Missing ID", func(t *testing.T) {
		_, handler := newWarehouseHandlerMock()
		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX1","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("PATCH", "/warehouses/{id}", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
		rr := httptest.NewRecorder()
		handler.Update(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("Update GetById Internal Server Error", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{}, error_message.ErrInternalServerError
		}
		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX1","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("PATCH", "/warehouses/{id}", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		rr := httptest.NewRecorder()
		handler.Update(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
	t.Run("Update Code Uniqueness Conflict", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{
				Id:                 1,
				Address:            "Dirección Original",
				Telephone:          "123456789",
				WareHouseCode:      "WH001",
				MinimumCapacity:    100,
				MinimumTemperature: 5.0,
				LocalityId:         1,
			}, nil
		}
		mock.ValidateCodeUniquenessFunc = func(ctx context.Context, code string) error {
			return error_message.ErrAlreadyExists
		}
		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX2","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("PATCH", "/warehouses/{id}", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		rr := httptest.NewRecorder()
		handler.Update(rr, req)
		assert.Equal(t, http.StatusConflict, rr.Code)
	})
	t.Run("Update Internal Server Error", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.UpdateFunc = func(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error) {
			return models.Warehouse{}, error_message.ErrInternalServerError
		}
		mock.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{
				Id:                 1,
				Address:            "Dirección Original",
				Telephone:          "123456789",
				WareHouseCode:      "WH001",
				MinimumCapacity:    100,
				MinimumTemperature: 5.0,
				LocalityId:         1,
			}, nil
		}
		mock.ValidateCodeUniquenessFunc = func(ctx context.Context, code string) error {
			return nil
		}
		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX1","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("PATCH", "/warehouses/{id}", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		rr := httptest.NewRecorder()
		handler.Update(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
	t.Run("Update Invalid Structure", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{
				Id:                 1,
				Address:            "Dirección Original",
				Telephone:          "123456789",
				WareHouseCode:      "WH001",
				MinimumCapacity:    100,
				MinimumTemperature: 5.0,
				LocalityId:         1,
			}, nil
		}
		body := `{"address":"","telephone":"1234","warehouse_code":"XX1","minimum_capacity":-10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("PATCH", "/warehouses/{id}", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		rr := httptest.NewRecorder()
		handler.Update(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("Update Code Uniqueness Internal Server Error", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{
				Id:                 1,
				Address:            "Dirección Original",
				Telephone:          "123456789",
				WareHouseCode:      "WH001",
				MinimumCapacity:    100,
				MinimumTemperature: 5.0,
				LocalityId:         1,
			}, nil
		}
		mock.ValidateCodeUniquenessFunc = func(ctx context.Context, code string) error {
			return error_message.ErrInternalServerError
		}
		body := `{"address":"Dirección","telephone":"1234","warehouse_code":"XX2","minimum_capacity":10,"minimum_temperature":5.0,"locality_id":2}`
		req := httptest.NewRequest("PATCH", "/warehouses/{id}", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		rr := httptest.NewRecorder()
		handler.Update(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestWarehouseHandler_Delete(t *testing.T) {
	t.Run("Delete Success", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.DeleteFunc = func(ctx context.Context, id int) error {
			return nil
		}
		request := httptest.NewRequest("DELETE", "/warehouses/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
		response := httptest.NewRecorder()
		handler.Delete(response, request)
		assert.Equal(t, http.StatusNoContent, response.Code)
	})
	t.Run("Delete Fail not found", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.DeleteFunc = func(ctx context.Context, id int) error {
			return error_message.ErrNotFound
		}
		request := httptest.NewRequest("DELETE", "/warehouses/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
		response := httptest.NewRecorder()
		handler.Delete(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})
	t.Run("Delete Invalid ID", func(t *testing.T) {
		_, handler := newWarehouseHandlerMock()
		request := httptest.NewRequest("DELETE", "/warehouses/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "invalid")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
		response := httptest.NewRecorder()
		handler.Delete(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	t.Run("Delete Missing ID", func(t *testing.T) {
		_, handler := newWarehouseHandlerMock()
		request := httptest.NewRequest("DELETE", "/warehouses/{id}", nil)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chi.NewRouteContext()))
		response := httptest.NewRecorder()
		handler.Delete(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	t.Run("Delete Internal Server Error", func(t *testing.T) {
		mock, handler := newWarehouseHandlerMock()
		mock.DeleteFunc = func(ctx context.Context, id int) error {
			return error_message.ErrInternalServerError
		}
		request := httptest.NewRequest("DELETE", "/warehouses/{id}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
		response := httptest.NewRecorder()
		handler.Delete(response, request)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}
