// Package handlers_test - Tests de integración para Carry Handler
// Tests HTTP comprehensivos para operaciones CRUD de la entidad Carry
package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
)

// newCarryHandlerMock - Factory para crear handler con servicio mock
func newCarryHandlerMock() (*tests.CarryServiceMock, *handlers.CarryHandler) {
	mock := &tests.CarryServiceMock{}
	h := handlers.NewCarryHandler(mock)
	return mock, h
}

// TestCarryHandler_Create - Tests para POST /carries
// Casos: 201 (creación exitosa), 400 (validación), 409 (CID duplicado), 400 (JSON inválido), 500 (error interno)
func TestCarryHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		mock.CreateCarryFunc = func(ctx context.Context, carry models.Carry) (models.Carry, error) {
			carry.Id = 1
			return carry, nil
		}

		body := `{"cid":"CID123","company_name":"Test Company","address":"Test Address","telephone":"123456789","locality_id":1}`
		req := httptest.NewRequest("POST", "/carries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Contains(t, rr.Body.String(), "CID123")
	})

	t.Run("fail_with_missing_required_fields", func(t *testing.T) {
		_, handler := newCarryHandlerMock()

		body := `{"cid":"","company_name":"","address":"","telephone":"","locality_id":0}`
		req := httptest.NewRequest("POST", "/carries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("fail_with_invalid_field_lengths", func(t *testing.T) {
		_, handler := newCarryHandlerMock()

		// CID too long
		body := `{"cid":"` + strings.Repeat("A", 256) + `","company_name":"Test","address":"Test","telephone":"123","locality_id":1}`
		req := httptest.NewRequest("POST", "/carries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("cid_already_exists", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		mock.CreateCarryFunc = func(ctx context.Context, carry models.Carry) (models.Carry, error) {
			return models.Carry{}, error_message.ErrAlreadyExists
		}

		body := `{"cid":"CID123","company_name":"Test Company","address":"Test Address","telephone":"123456789","locality_id":1}`
		req := httptest.NewRequest("POST", "/carries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("locality_not_found", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		mock.CreateCarryFunc = func(ctx context.Context, carry models.Carry) (models.Carry, error) {
			return models.Carry{}, error_message.ErrNotFound
		}

		body := `{"cid":"CID123","company_name":"Test Company","address":"Test Address","telephone":"123456789","locality_id":1}`
		req := httptest.NewRequest("POST", "/carries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("invalid_JSON_format", func(t *testing.T) {
		_, handler := newCarryHandlerMock()

		body := `{"cid":"CID123","company_name":"Test Company","address":"Test Address","telephone":"123456789","locality_id":1`
		req := httptest.NewRequest("POST", "/carries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("create_carry_internal_server_error", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		mock.CreateCarryFunc = func(ctx context.Context, carry models.Carry) (models.Carry, error) {
			return models.Carry{}, error_message.ErrInternalServerError
		}

		body := `{"cid":"CID123","company_name":"Test Company","address":"Test Address","telephone":"123456789","locality_id":1}`
		req := httptest.NewRequest("POST", "/carries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("service_returns_unexpected_error", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		mock.CreateCarryFunc = func(ctx context.Context, carry models.Carry) (models.Carry, error) {
			return models.Carry{}, errors.New("unexpected database error")
		}

		body := `{"cid":"CID123","company_name":"Test Company","address":"Test Address","telephone":"123456789","locality_id":1}`
		req := httptest.NewRequest("POST", "/carries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("request_timeout_simulation", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		mock.CreateCarryFunc = func(ctx context.Context, carry models.Carry) (models.Carry, error) {
			// Simulate a timeout by checking if context is cancelled
			select {
			case <-ctx.Done():
				return models.Carry{}, ctx.Err()
			default:
				return models.Carry{}, errors.New("unexpected database error")
			}
		}

		body := `{"cid":"CID123","company_name":"Test Company","address":"Test Address","telephone":"123456789","locality_id":1}`
		req := httptest.NewRequest("POST", "/carries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		// Create a context that's already cancelled to simulate timeout
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusRequestTimeout, rr.Code)
	})
}

// TestCarryHandler_GetCarryReportByLocality - Tests para GET /carries/reports
// Casos: 200 (exitoso), 400 (ID inválido), 404 (localidad no encontrada), 500 (error interno)
func TestCarryHandler_GetCarryReportByLocality(t *testing.T) {
	t.Run("get_all_localities_report_success", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		expectedReports := []responses.LocalityCarryReport{
			{LocalityId: 1, LocalityName: "Locality 1", CarriersCount: 5},
			{LocalityId: 2, LocalityName: "Locality 2", CarriersCount: 3},
		}
		mock.GetCarryReportByLocalityFunc = func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
			return expectedReports, nil
		}

		req := httptest.NewRequest("GET", "/carries/reports", nil)
		rr := httptest.NewRecorder()

		handler.GetCarryReportByLocality(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Locality 1")
		assert.Contains(t, rr.Body.String(), "Locality 2")
	})

	t.Run("get_specific_locality_report_success", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		expectedReports := []responses.LocalityCarryReport{
			{LocalityId: 1, LocalityName: "Locality 1", CarriersCount: 5},
		}
		mock.GetCarryReportByLocalityFunc = func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
			assert.Equal(t, 1, localityID)
			return expectedReports, nil
		}

		req := httptest.NewRequest("GET", "/carries/reports?id=1", nil)
		rr := httptest.NewRecorder()

		handler.GetCarryReportByLocality(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Locality 1")
	})

	t.Run("get_specific_locality_report_not_found", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		mock.GetCarryReportByLocalityFunc = func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
			return nil, error_message.ErrNotFound
		}

		req := httptest.NewRequest("GET", "/carries/reports?id=999", nil)
		rr := httptest.NewRecorder()

		handler.GetCarryReportByLocality(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("invalid_locality_ID_format", func(t *testing.T) {
		_, handler := newCarryHandlerMock()

		req := httptest.NewRequest("GET", "/carries/reports?id=invalid", nil)
		rr := httptest.NewRecorder()

		handler.GetCarryReportByLocality(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("negative_locality_ID", func(t *testing.T) {
		_, handler := newCarryHandlerMock()

		req := httptest.NewRequest("GET", "/carries/reports?id=-1", nil)
		rr := httptest.NewRecorder()

		handler.GetCarryReportByLocality(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("zero_locality_ID", func(t *testing.T) {
		_, handler := newCarryHandlerMock()

		req := httptest.NewRequest("GET", "/carries/reports?id=0", nil)
		rr := httptest.NewRecorder()

		handler.GetCarryReportByLocality(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("internal_server_error", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		mock.GetCarryReportByLocalityFunc = func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
			return nil, error_message.ErrInternalServerError
		}

		req := httptest.NewRequest("GET", "/carries/reports", nil)
		rr := httptest.NewRecorder()

		handler.GetCarryReportByLocality(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("service_returns_unexpected_error", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		mock.GetCarryReportByLocalityFunc = func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
			return nil, errors.New("unexpected database error")
		}

		req := httptest.NewRequest("GET", "/carries/reports", nil)
		rr := httptest.NewRecorder()

		handler.GetCarryReportByLocality(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("request_timeout_simulation", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		mock.GetCarryReportByLocalityFunc = func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
			// Simulate a timeout by checking if context is cancelled
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return nil, errors.New("unexpected database error")
			}
		}

		req := httptest.NewRequest("GET", "/carries/reports", nil)

		// Create a context that's already cancelled to simulate timeout
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.GetCarryReportByLocality(rr, req)

		assert.Equal(t, http.StatusRequestTimeout, rr.Code)
	})

	t.Run("empty_query_parameter", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		expectedReports := []responses.LocalityCarryReport{
			{LocalityId: 1, LocalityName: "Locality 1", CarriersCount: 5},
		}
		mock.GetCarryReportByLocalityFunc = func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
			assert.Equal(t, 0, localityID) // Should be 0 for all localities
			return expectedReports, nil
		}

		req := httptest.NewRequest("GET", "/carries/reports", nil)
		rr := httptest.NewRecorder()

		handler.GetCarryReportByLocality(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Locality 1")
	})

	t.Run("very_large_locality_ID", func(t *testing.T) {
		mock, handler := newCarryHandlerMock()
		expectedReports := []responses.LocalityCarryReport{
			{LocalityId: 999999999, LocalityName: "Large Locality", CarriersCount: 1},
		}
		mock.GetCarryReportByLocalityFunc = func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
			assert.Equal(t, 999999999, localityID)
			return expectedReports, nil
		}

		req := httptest.NewRequest("GET", "/carries/reports?id=999999999", nil)
		rr := httptest.NewRecorder()

		handler.GetCarryReportByLocality(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code) // Should be valid
		assert.Contains(t, rr.Body.String(), "Large Locality")
	})
}
