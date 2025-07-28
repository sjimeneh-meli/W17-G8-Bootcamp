package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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
