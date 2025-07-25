package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestProductHandler_GetAll agrupa todos los tests para el método GetAll
// TestProductHandler_GetAll groups all tests for the GetAll method
func TestProductHandler_GetAll(t *testing.T) {
	// Test case: Successful retrieval of all products
	// Caso de prueba: Obtención exitosa de todos los productos
	t.Run("Success_ReturnsAllProducts", func(t *testing.T) {
		// Arrange - Configurar mocks y datos de prueba
		// Arrange - Set up mocks and test data
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		// Datos de productos de prueba / Test product data
		expectedProducts := []models.Product{
			{
				Id:                             1,
				ProductCode:                    "PROD001",
				Description:                    "Test Product 1",
				Width:                          10.5,
				Height:                         20.3,
				Length:                         15.7,
				NetWeight:                      5.2,
				ExpirationRate:                 0.1,
				RecommendedFreezingTemperature: -18.0,
				FreezingRate:                   0.8,
				ProductTypeID:                  1,
				SellerID:                       nil,
			},
			{
				Id:                             2,
				ProductCode:                    "PROD002",
				Description:                    "Test Product 2",
				Width:                          8.5,
				Height:                         15.3,
				Length:                         12.7,
				NetWeight:                      3.2,
				ExpirationRate:                 0.2,
				RecommendedFreezingTemperature: -20.0,
				FreezingRate:                   0.9,
				ProductTypeID:                  2,
				SellerID:                       nil,
			},
		}

		// Configurar expectativas del mock / Set up mock expectations
		mockService.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return(expectedProducts, nil)

		// Crear request HTTP de prueba / Create test HTTP request
		req := httptest.NewRequest("GET", "/products", nil)
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.GetAll(rec, req)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.Equal(t, http.StatusOK, rec.Code)

		var response responses.DataResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verificar que la respuesta contiene los productos esperados
		// Verify that the response contains the expected products
		responseData, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Len(t, responseData, 2)

		mockService.AssertExpectations(t)
	})

	// Test case: Service returns timeout error
	// Caso de prueba: El servicio retorna error de timeout
	t.Run("Error_TimeoutExceeded", func(t *testing.T) {
		// Arrange - Configurar mocks con error de timeout
		// Arrange - Set up mocks with timeout error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		// Configurar mock para retornar error de timeout / Set up mock to return timeout error
		mockService.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return([]models.Product{}, context.DeadlineExceeded)

		req := httptest.NewRequest("GET", "/products", nil)
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.GetAll(rec, req)

		// Assert - Verificar que retorna error de timeout
		// Assert - Verify that it returns timeout error
		assert.Equal(t, http.StatusGatewayTimeout, rec.Code)
		assert.Contains(t, rec.Body.String(), "the request took too long to process")

		mockService.AssertExpectations(t)
	})

	// Test case: Service returns internal server error
	// Caso de prueba: El servicio retorna error interno del servidor
	t.Run("Error_InternalServerError", func(t *testing.T) {
		// Arrange - Configurar mocks con error interno
		// Arrange - Set up mocks with internal error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		expectedError := errors.New("database connection failed")
		mockService.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return([]models.Product{}, expectedError)

		req := httptest.NewRequest("GET", "/products", nil)
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.GetAll(rec, req)

		// Assert - Verificar que retorna error interno del servidor
		// Assert - Verify that it returns internal server error
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), expectedError.Error())

		mockService.AssertExpectations(t)
	})
}

// TestProductHandler_Create agrupa todos los tests para el método Create
// TestProductHandler_Create groups all tests for the Create method
func TestProductHandler_Create(t *testing.T) {
	// Test case: Successful product creation
	// Caso de prueba: Creación exitosa de producto
	t.Run("Success_CreatesProduct", func(t *testing.T) {
		// Arrange - Configurar mocks y datos de prueba
		// Arrange - Set up mocks and test data
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		// Request de producto válido / Valid product request
		productRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
			SellerID:                       nil,
		}

		// Producto esperado que retorna el servicio / Expected product returned by service
		expectedProduct := models.Product{
			Id:                             1,
			ProductCode:                    productRequest.ProductCode,
			Description:                    productRequest.Description,
			Width:                          productRequest.Width,
			Height:                         productRequest.Height,
			Length:                         productRequest.Length,
			NetWeight:                      productRequest.NetWeight,
			ExpirationRate:                 productRequest.ExpirationRate,
			RecommendedFreezingTemperature: productRequest.RecommendedFreezingTemperature,
			FreezingRate:                   productRequest.FreezingRate,
			ProductTypeID:                  productRequest.ProductTypeID,
			SellerID:                       productRequest.SellerID,
		}

		// Configurar expectativas del mock / Set up mock expectations
		mockService.On("Create", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("models.Product")).Return(expectedProduct, nil)

		// Crear request HTTP con JSON válido / Create HTTP request with valid JSON
		reqBody, _ := json.Marshal(productRequest)
		req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Create(rec, req)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response responses.DataResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		mockService.AssertExpectations(t)
	})

	// Test case: Invalid JSON payload
	// Caso de prueba: Payload JSON inválido
	t.Run("Error_InvalidJSON", func(t *testing.T) {
		// Arrange - Configurar handler sin expectativas del mock
		// Arrange - Set up handler without mock expectations
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		// JSON inválido / Invalid JSON
		req := httptest.NewRequest("POST", "/products", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Create(rec, req)

		// Assert - Verificar que retorna error de bad request
		// Assert - Verify that it returns bad request error
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid request payload")

		mockService.AssertExpectations(t)
	})

	// Test case: Validation error
	// Caso de prueba: Error de validación
	t.Run("Error_ValidationFailed", func(t *testing.T) {
		// Arrange - Configurar request con campos faltantes
		// Arrange - Set up request with missing fields
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		// Request con campos requeridos faltantes / Request with missing required fields
		invalidRequest := requests.ProductRequest{
			ProductCode: "", // Campo requerido vacío / Required field empty
			Description: "Test Product",
		}

		reqBody, _ := json.Marshal(invalidRequest)
		req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Create(rec, req)

		// Assert - Verificar que retorna error de entidad no procesable
		// Assert - Verify that it returns unprocessable entity error
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		mockService.AssertExpectations(t)
	})

	// Test case: Service returns timeout error
	// Caso de prueba: El servicio retorna error de timeout
	t.Run("Error_TimeoutExceeded", func(t *testing.T) {
		// Arrange - Configurar mocks con error de timeout
		// Arrange - Set up mocks with timeout error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		productRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
		}

		mockService.On("Create", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("models.Product")).Return(models.Product{}, context.DeadlineExceeded)

		reqBody, _ := json.Marshal(productRequest)
		req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Create(rec, req)

		// Assert - Verificar que retorna error de timeout
		// Assert - Verify that it returns timeout error
		assert.Equal(t, http.StatusGatewayTimeout, rec.Code)
		assert.Contains(t, rec.Body.String(), "the request took too long to process")

		mockService.AssertExpectations(t)
	})

	// Test case: Product already exists
	// Caso de prueba: El producto ya existe
	t.Run("Error_ProductAlreadyExists", func(t *testing.T) {
		// Arrange - Configurar mocks con error de conflicto
		// Arrange - Set up mocks with conflict error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		productRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
		}

		mockService.On("Create", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("models.Product")).Return(models.Product{}, error_message.ErrAlreadyExists)

		reqBody, _ := json.Marshal(productRequest)
		req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Create(rec, req)

		// Assert - Verificar que retorna error de conflicto
		// Assert - Verify that it returns conflict error
		assert.Equal(t, http.StatusConflict, rec.Code)

		mockService.AssertExpectations(t)
	})

	// Test case: Internal server error
	// Caso de prueba: Error interno del servidor
	t.Run("Error_InternalServerError", func(t *testing.T) {
		// Arrange - Configurar mocks con error interno
		// Arrange - Set up mocks with internal error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		productRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
		}

		expectedError := errors.New("database error")
		mockService.On("Create", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("models.Product")).Return(models.Product{}, expectedError)

		reqBody, _ := json.Marshal(productRequest)
		req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Create(rec, req)

		// Assert - Verificar que retorna error interno del servidor
		// Assert - Verify that it returns internal server error
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		mockService.AssertExpectations(t)
	})
}

// TestProductHandler_Get agrupa todos los tests para el método Get
// TestProductHandler_Get groups all tests for the Get method
func TestProductHandler_Get(t *testing.T) {
	// Test case: Successful product retrieval by ID
	// Caso de prueba: Obtención exitosa de producto por ID
	t.Run("Success_ReturnsProductByID", func(t *testing.T) {
		// Arrange - Configurar mocks y datos de prueba
		// Arrange - Set up mocks and test data
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		expectedProduct := models.Product{
			Id:                             1,
			ProductCode:                    "PROD001",
			Description:                    "Test Product",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
			SellerID:                       nil,
		}

		mockService.On("GetByID", mock.AnythingOfType("*context.timerCtx"), int64(1)).Return(expectedProduct, nil)

		// Crear request con parámetro ID / Create request with ID parameter
		req := httptest.NewRequest("GET", "/products/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"1"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Get(rec, req)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.Equal(t, http.StatusOK, rec.Code)

		var response responses.DataResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		mockService.AssertExpectations(t)
	})

	// Test case: Invalid ID parameter (empty)
	// Caso de prueba: Parámetro ID inválido (vacío)
	t.Run("Error_EmptyIDParameter", func(t *testing.T) {
		// Arrange - Configurar request sin parámetro ID
		// Arrange - Set up request without ID parameter
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		req := httptest.NewRequest("GET", "/products/", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{""},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Get(rec, req)

		// Assert - Verificar que retorna error de bad request
		// Assert - Verify that it returns bad request error
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "id parameter is required")

		mockService.AssertExpectations(t)
	})

	// Test case: Invalid ID parameter (non-numeric)
	// Caso de prueba: Parámetro ID inválido (no numérico)
	t.Run("Error_NonNumericIDParameter", func(t *testing.T) {
		// Arrange - Configurar request con ID no numérico
		// Arrange - Set up request with non-numeric ID
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		req := httptest.NewRequest("GET", "/products/abc", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"abc"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Get(rec, req)

		// Assert - Verificar que retorna error de bad request
		// Assert - Verify that it returns bad request error
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid id parameter")

		mockService.AssertExpectations(t)
	})

	// Test case: Service returns timeout error
	// Caso de prueba: El servicio retorna error de timeout
	t.Run("Error_TimeoutExceeded", func(t *testing.T) {
		// Arrange - Configurar mocks con error de timeout
		// Arrange - Set up mocks with timeout error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		mockService.On("GetByID", mock.AnythingOfType("*context.timerCtx"), int64(1)).Return(models.Product{}, context.DeadlineExceeded)

		req := httptest.NewRequest("GET", "/products/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"1"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Get(rec, req)

		// Assert - Verificar que retorna error de timeout
		// Assert - Verify that it returns timeout error
		assert.Equal(t, http.StatusGatewayTimeout, rec.Code)
		assert.Contains(t, rec.Body.String(), "the request took too long to process")

		mockService.AssertExpectations(t)
	})

	// Test case: Product not found
	// Caso de prueba: Producto no encontrado
	t.Run("Error_ProductNotFound", func(t *testing.T) {
		// Arrange - Configurar mocks con error de no encontrado
		// Arrange - Set up mocks with not found error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		mockService.On("GetByID", mock.AnythingOfType("*context.timerCtx"), int64(999)).Return(models.Product{}, error_message.ErrNotFound)

		req := httptest.NewRequest("GET", "/products/999", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"999"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Get(rec, req)

		// Assert - Verificar que retorna error de no encontrado
		// Assert - Verify that it returns not found error
		assert.Equal(t, http.StatusNotFound, rec.Code)

		mockService.AssertExpectations(t)
	})

	// Test case: Internal server error
	// Caso de prueba: Error interno del servidor
	t.Run("Error_InternalServerError", func(t *testing.T) {
		// Arrange - Configurar mocks con error interno
		// Arrange - Set up mocks with internal error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		expectedError := errors.New("database error")
		mockService.On("GetByID", mock.AnythingOfType("*context.timerCtx"), int64(1)).Return(models.Product{}, expectedError)

		req := httptest.NewRequest("GET", "/products/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"1"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Get(rec, req)

		// Assert - Verificar que retorna error interno del servidor
		// Assert - Verify that it returns internal server error
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		mockService.AssertExpectations(t)
	})
}

// TestProductHandler_Update agrupa todos los tests para el método Update
// TestProductHandler_Update groups all tests for the Update method
func TestProductHandler_Update(t *testing.T) {
	// Test case: Successful product update
	// Caso de prueba: Actualización exitosa de producto
	t.Run("Success_UpdatesProduct", func(t *testing.T) {
		// Arrange - Configurar mocks y datos de prueba
		// Arrange - Set up mocks and test data
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		productRequest := requests.ProductRequest{
			ProductCode:                    "PROD001_UPDATED",
			Description:                    "Updated Test Product",
			Width:                          12.5,
			Height:                         22.3,
			Length:                         17.7,
			NetWeight:                      6.2,
			ExpirationRate:                 0.15,
			RecommendedFreezingTemperature: -19.0,
			FreezingRate:                   0.85,
			ProductTypeID:                  2,
			SellerID:                       nil,
		}

		expectedProduct := models.Product{
			Id:                             1,
			ProductCode:                    productRequest.ProductCode,
			Description:                    productRequest.Description,
			Width:                          productRequest.Width,
			Height:                         productRequest.Height,
			Length:                         productRequest.Length,
			NetWeight:                      productRequest.NetWeight,
			ExpirationRate:                 productRequest.ExpirationRate,
			RecommendedFreezingTemperature: productRequest.RecommendedFreezingTemperature,
			FreezingRate:                   productRequest.FreezingRate,
			ProductTypeID:                  productRequest.ProductTypeID,
			SellerID:                       productRequest.SellerID,
		}

		mockService.On("Update", mock.AnythingOfType("*context.timerCtx"), int64(1), mock.AnythingOfType("models.Product")).Return(expectedProduct, nil)

		reqBody, _ := json.Marshal(productRequest)
		req := httptest.NewRequest("PUT", "/products/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"1"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Update(rec, req)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.Equal(t, http.StatusOK, rec.Code)

		var response responses.DataResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		mockService.AssertExpectations(t)
	})

	// Test case: Invalid ID parameter
	// Caso de prueba: Parámetro ID inválido
	t.Run("Error_InvalidIDParameter", func(t *testing.T) {
		// Arrange - Configurar request con ID inválido
		// Arrange - Set up request with invalid ID
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		productRequest := requests.ProductRequest{
			ProductCode: "PROD001",
			Description: "Test Product",
		}

		reqBody, _ := json.Marshal(productRequest)
		req := httptest.NewRequest("PUT", "/products/invalid", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"invalid"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Update(rec, req)

		// Assert - Verificar que retorna error de bad request
		// Assert - Verify that it returns bad request error
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		mockService.AssertExpectations(t)
	})

	// Test case: Invalid JSON payload
	// Caso de prueba: Payload JSON inválido
	t.Run("Error_InvalidJSON", func(t *testing.T) {
		// Arrange - Configurar request con JSON inválido
		// Arrange - Set up request with invalid JSON
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		req := httptest.NewRequest("PUT", "/products/1", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"1"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Update(rec, req)

		// Assert - Verificar que retorna error de bad request
		// Assert - Verify that it returns bad request error
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid request payload")

		mockService.AssertExpectations(t)
	})

	// Test case: Service returns timeout error
	// Caso de prueba: El servicio retorna error de timeout
	t.Run("Error_TimeoutExceeded", func(t *testing.T) {
		// Arrange - Configurar mocks con error de timeout
		// Arrange - Set up mocks with timeout error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		productRequest := requests.ProductRequest{
			ProductCode: "PROD001",
			Description: "Test Product",
		}

		mockService.On("Update", mock.AnythingOfType("*context.timerCtx"), int64(1), mock.AnythingOfType("models.Product")).Return(models.Product{}, context.DeadlineExceeded)

		reqBody, _ := json.Marshal(productRequest)
		req := httptest.NewRequest("PUT", "/products/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"1"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Update(rec, req)

		// Assert - Verificar que retorna error de timeout
		// Assert - Verify that it returns timeout error
		assert.Equal(t, http.StatusGatewayTimeout, rec.Code)
		assert.Contains(t, rec.Body.String(), "the request took too long to process")

		mockService.AssertExpectations(t)
	})

	// Test case: Product not found
	// Caso de prueba: Producto no encontrado
	t.Run("Error_ProductNotFound", func(t *testing.T) {
		// Arrange - Configurar mocks con error de no encontrado
		// Arrange - Set up mocks with not found error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		productRequest := requests.ProductRequest{
			ProductCode: "PROD001",
			Description: "Test Product",
		}

		mockService.On("Update", mock.AnythingOfType("*context.timerCtx"), int64(999), mock.AnythingOfType("models.Product")).Return(models.Product{}, error_message.ErrNotFound)

		reqBody, _ := json.Marshal(productRequest)
		req := httptest.NewRequest("PUT", "/products/999", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"999"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Update(rec, req)

		// Assert - Verificar que retorna error de no encontrado
		// Assert - Verify that it returns not found error
		assert.Equal(t, http.StatusNotFound, rec.Code)

		mockService.AssertExpectations(t)
	})

	// Test case: Internal server error
	// Caso de prueba: Error interno del servidor
	t.Run("Error_InternalServerError", func(t *testing.T) {
		// Arrange - Configurar mocks con error interno
		// Arrange - Set up mocks with internal error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		productRequest := requests.ProductRequest{
			ProductCode: "PROD001",
			Description: "Test Product",
		}

		expectedError := errors.New("database error")
		mockService.On("Update", mock.AnythingOfType("*context.timerCtx"), int64(1), mock.AnythingOfType("models.Product")).Return(models.Product{}, expectedError)

		reqBody, _ := json.Marshal(productRequest)
		req := httptest.NewRequest("PUT", "/products/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"1"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Update(rec, req)

		// Assert - Verificar que retorna error interno del servidor
		// Assert - Verify that it returns internal server error
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		mockService.AssertExpectations(t)
	})
}

// TestProductHandler_Delete agrupa todos los tests para el método Delete
// TestProductHandler_Delete groups all tests for the Delete method
func TestProductHandler_Delete(t *testing.T) {
	// Test case: Successful product deletion
	// Caso de prueba: Eliminación exitosa de producto
	t.Run("Success_DeletesProduct", func(t *testing.T) {
		// Arrange - Configurar mocks y datos de prueba
		// Arrange - Set up mocks and test data
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		mockService.On("Delete", mock.AnythingOfType("*context.timerCtx"), int64(1)).Return(nil)

		req := httptest.NewRequest("DELETE", "/products/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"1"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Delete(rec, req)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())

		mockService.AssertExpectations(t)
	})

	// Test case: Invalid ID parameter
	// Caso de prueba: Parámetro ID inválido
	t.Run("Error_InvalidIDParameter", func(t *testing.T) {
		// Arrange - Configurar request con ID inválido
		// Arrange - Set up request with invalid ID
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		req := httptest.NewRequest("DELETE", "/products/invalid", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"invalid"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Delete(rec, req)

		// Assert - Verificar que retorna error de bad request
		// Assert - Verify that it returns bad request error
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		mockService.AssertExpectations(t)
	})

	// Test case: Service returns timeout error
	// Caso de prueba: El servicio retorna error de timeout
	t.Run("Error_TimeoutExceeded", func(t *testing.T) {
		// Arrange - Configurar mocks con error de timeout
		// Arrange - Set up mocks with timeout error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		mockService.On("Delete", mock.AnythingOfType("*context.timerCtx"), int64(1)).Return(context.DeadlineExceeded)

		req := httptest.NewRequest("DELETE", "/products/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"1"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Delete(rec, req)

		// Assert - Verificar que retorna error de timeout
		// Assert - Verify that it returns timeout error
		assert.Equal(t, http.StatusGatewayTimeout, rec.Code)
		assert.Contains(t, rec.Body.String(), "the request took too long to process")

		mockService.AssertExpectations(t)
	})

	// Test case: Product not found
	// Caso de prueba: Producto no encontrado
	t.Run("Error_ProductNotFound", func(t *testing.T) {
		// Arrange - Configurar mocks con error de no encontrado
		// Arrange - Set up mocks with not found error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		mockService.On("Delete", mock.AnythingOfType("*context.timerCtx"), int64(999)).Return(error_message.ErrNotFound)

		req := httptest.NewRequest("DELETE", "/products/999", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"999"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Delete(rec, req)

		// Assert - Verificar que retorna error de no encontrado
		// Assert - Verify that it returns not found error
		assert.Equal(t, http.StatusNotFound, rec.Code)

		mockService.AssertExpectations(t)
	})

	// Test case: Internal server error
	// Caso de prueba: Error interno del servidor
	t.Run("Error_InternalServerError", func(t *testing.T) {
		// Arrange - Configurar mocks con error interno
		// Arrange - Set up mocks with internal error
		mockService := new(tests.ProductServiceMock)
		handler := NewProductHandler(mockService)

		expectedError := errors.New("database error")
		mockService.On("Delete", mock.AnythingOfType("*context.timerCtx"), int64(1)).Return(expectedError)

		req := httptest.NewRequest("DELETE", "/products/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"1"},
			},
		}))
		rec := httptest.NewRecorder()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		handler.Delete(rec, req)

		// Assert - Verificar que retorna error interno del servidor
		// Assert - Verify that it returns internal server error
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		mockService.AssertExpectations(t)
	})
}

// TestParseID agrupa todos los tests para la función parseID
// TestParseID groups all tests for the parseID function
func TestParseID(t *testing.T) {
	// Test case: Successful ID parsing
	// Caso de prueba: Parsing exitoso de ID
	t.Run("Success_ParsesIDCorrectly", func(t *testing.T) {
		// Arrange - Configurar request con ID válido
		// Arrange - Set up request with valid ID
		req := httptest.NewRequest("GET", "/products/123", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"123"},
			},
		}))

		// Act - Ejecutar parseID
		// Act - Execute parseID
		id, err := parseID(req)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Equal(t, int64(123), id)
	})

	// Test case: Empty ID parameter
	// Caso de prueba: Parámetro ID vacío
	t.Run("Error_EmptyIDParameter", func(t *testing.T) {
		// Arrange - Configurar request sin ID
		// Arrange - Set up request without ID
		req := httptest.NewRequest("GET", "/products/", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{""},
			},
		}))

		// Act - Ejecutar parseID
		// Act - Execute parseID
		id, err := parseID(req)

		// Assert - Verificar que retorna error
		// Assert - Verify that it returns error
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.Contains(t, err.Error(), "id parameter is required")
	})

	// Test case: Invalid ID parameter (non-numeric)
	// Caso de prueba: Parámetro ID inválido (no numérico)
	t.Run("Error_NonNumericIDParameter", func(t *testing.T) {
		// Arrange - Configurar request con ID no numérico
		// Arrange - Set up request with non-numeric ID
		req := httptest.NewRequest("GET", "/products/abc", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"abc"},
			},
		}))

		// Act - Ejecutar parseID
		// Act - Execute parseID
		id, err := parseID(req)

		// Assert - Verificar que retorna error
		// Assert - Verify that it returns error
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.Contains(t, err.Error(), "invalid id parameter")
		assert.Contains(t, err.Error(), "abc")
	})

	// Test case: Negative ID parameter
	// Caso de prueba: Parámetro ID negativo
	t.Run("Success_ParsesNegativeID", func(t *testing.T) {
		// Arrange - Configurar request con ID negativo
		// Arrange - Set up request with negative ID
		req := httptest.NewRequest("GET", "/products/-1", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{"-1"},
			},
		}))

		// Act - Ejecutar parseID
		// Act - Execute parseID
		id, err := parseID(req)

		// Assert - Verificar que parsea correctamente números negativos
		// Assert - Verify that it correctly parses negative numbers
		assert.NoError(t, err)
		assert.Equal(t, int64(-1), id)
	})

	// Test case: Large ID parameter
	// Caso de prueba: Parámetro ID grande
	t.Run("Success_ParsesLargeID", func(t *testing.T) {
		// Arrange - Configurar request con ID grande
		// Arrange - Set up request with large ID
		largeID := "9223372036854775807" // Max int64
		req := httptest.NewRequest("GET", "/products/"+largeID, nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
			URLParams: chi.RouteParams{
				Keys:   []string{"id"},
				Values: []string{largeID},
			},
		}))

		// Act - Ejecutar parseID
		// Act - Execute parseID
		id, err := parseID(req)

		// Assert - Verificar que parsea correctamente números grandes
		// Assert - Verify that it correctly parses large numbers
		assert.NoError(t, err)
		expectedID, _ := strconv.ParseInt(largeID, 10, 64)
		assert.Equal(t, expectedID, id)
	})
}

// TestNewProductHandler prueba el constructor del ProductHandler
// TestNewProductHandler tests the ProductHandler constructor
func TestNewProductHandler(t *testing.T) {
	// Test case: Successful handler creation
	// Caso de prueba: Creación exitosa del handler
	t.Run("Success_CreatesHandler", func(t *testing.T) {
		// Arrange - Crear mock service
		// Arrange - Create mock service
		mockService := new(tests.ProductServiceMock)

		// Act - Crear nuevo handler
		// Act - Create new handler
		handler := NewProductHandler(mockService)

		// Assert - Verificar que el handler se creó correctamente
		// Assert - Verify that the handler was created correctly
		assert.NotNil(t, handler)
		assert.Equal(t, mockService, handler.service)
	})
}
