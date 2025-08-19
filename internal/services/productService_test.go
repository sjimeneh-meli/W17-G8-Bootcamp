package services

import (
	"context"
	"errors"
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestProductService_GetAll agrupa todos los tests para el método GetAll
// TestProductService_GetAll groups all tests for the GetAll method
func TestProductService_GetAll(t *testing.T) {
	// Test case: Successful retrieval of all products
	// Caso de prueba: Obtención exitosa de todos los productos
	t.Run("Success_ReturnsAllProducts", func(t *testing.T) {
		// Arrange - Reset singleton and set up mocks
		// Arrange - Resetear singleton y configurar mocks
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

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

		mockRepo.On("GetAll", mock.Anything).Return(expectedProducts, nil)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.GetAll(context.Background())

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, result)
		assert.Len(t, result, 2)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Repository returns empty slice
	// Caso de prueba: El repositorio retorna slice vacío
	t.Run("Success_ReturnsEmptySlice", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock to return empty slice
		// Arrange - Resetear singleton y configurar mock para retornar slice vacío
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		expectedProducts := []models.Product{}
		mockRepo.On("GetAll", mock.Anything).Return(expectedProducts, nil)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.GetAll(context.Background())

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, result)
		assert.Len(t, result, 0)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Repository returns timeout error
	// Caso de prueba: El repositorio retorna error de timeout
	t.Run("Error_TimeoutExceeded", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with timeout error
		// Arrange - Resetear singleton y configurar mock con error de timeout
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		mockRepo.On("GetAll", mock.Anything).Return([]models.Product{}, context.DeadlineExceeded)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.GetAll(context.Background())

		// Assert - Verificar que retorna error de timeout
		// Assert - Verify that it returns timeout error
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Repository returns database error
	// Caso de prueba: El repositorio retorna error de base de datos
	t.Run("Error_DatabaseError", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with database error
		// Arrange - Resetear singleton y configurar mock con error de base de datos
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		expectedError := errors.New("database connection failed")
		mockRepo.On("GetAll", mock.Anything).Return([]models.Product{}, expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.GetAll(context.Background())

		// Assert - Verificar que retorna error de base de datos
		// Assert - Verify that it returns database error
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})
}

// TestProductService_GetByID agrupa todos los tests para el método GetByID
// TestProductService_GetByID groups all tests for the GetByID method
func TestProductService_GetByID(t *testing.T) {
	// Test case: Successful retrieval of product by ID
	// Caso de prueba: Obtención exitosa de producto por ID
	t.Run("Success_ReturnsProductByID", func(t *testing.T) {
		// Arrange - Reset singleton and set up mocks and test data
		// Arrange - Resetear singleton y configurar mocks y datos de prueba
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

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

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedProduct, nil)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.GetByID(context.Background(), 1)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, result)
		assert.Equal(t, int64(1), result.Id)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Product not found
	// Caso de prueba: Producto no encontrado
	t.Run("Error_ProductNotFound", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with not found error
		// Arrange - Resetear singleton y configurar mock con error de no encontrado
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(999)).Return(models.Product{}, error_message.ErrNotFound)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.GetByID(context.Background(), 999)

		// Assert - Verificar que retorna error de no encontrado
		// Assert - Verify that it returns not found error
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Repository returns timeout error
	// Caso de prueba: El repositorio retorna error de timeout
	t.Run("Error_TimeoutExceeded", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with timeout error
		// Arrange - Resetear singleton y configurar mock con error de timeout
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(models.Product{}, context.DeadlineExceeded)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.GetByID(context.Background(), 1)

		// Assert - Verificar que retorna error de timeout
		// Assert - Verify that it returns timeout error
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Database error
	// Caso de prueba: Error de base de datos
	t.Run("Error_DatabaseError", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with database error
		// Arrange - Resetear singleton y configurar mock con error de base de datos
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		expectedError := errors.New("database connection failed")
		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(models.Product{}, expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.GetByID(context.Background(), 1)

		// Assert - Verificar que retorna error de base de datos
		// Assert - Verify that it returns database error
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})
}

// TestProductService_Create agrupa todos los tests para el método Create
// TestProductService_Create groups all tests for the Create method
func TestProductService_Create(t *testing.T) {
	// Test case: Successful product creation
	// Caso de prueba: Creación exitosa de producto
	t.Run("Success_CreatesProduct", func(t *testing.T) {
		// Arrange - Reset singleton and set up mocks and test data
		// Arrange - Resetear singleton y configurar mocks y datos de prueba
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		inputProduct := models.Product{
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

		expectedProduct := inputProduct
		expectedProduct.Id = 1

		mockRepo.On("Create", mock.Anything, inputProduct).Return(expectedProduct, nil)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.Create(context.Background(), inputProduct)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, result)
		assert.Equal(t, int64(1), result.Id)
		assert.Equal(t, inputProduct.ProductCode, result.ProductCode)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Product already exists
	// Caso de prueba: El producto ya existe
	t.Run("Error_ProductAlreadyExists", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with already exists error
		// Arrange - Resetear singleton y configurar mock con error de ya existe
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		inputProduct := models.Product{
			ProductCode: "PROD001",
			Description: "Test Product",
		}

		mockRepo.On("Create", mock.Anything, inputProduct).Return(models.Product{}, error_message.ErrAlreadyExists)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.Create(context.Background(), inputProduct)

		// Assert - Verificar que retorna error de ya existe
		// Assert - Verify that it returns already exists error
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrAlreadyExists, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Repository returns timeout error
	// Caso de prueba: El repositorio retorna error de timeout
	t.Run("Error_TimeoutExceeded", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with timeout error
		// Arrange - Resetear singleton y configurar mock con error de timeout
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		inputProduct := models.Product{
			ProductCode: "PROD001",
			Description: "Test Product",
		}

		mockRepo.On("Create", mock.Anything, inputProduct).Return(models.Product{}, context.DeadlineExceeded)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.Create(context.Background(), inputProduct)

		// Assert - Verificar que retorna error de timeout
		// Assert - Verify that it returns timeout error
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Database error
	// Caso de prueba: Error de base de datos
	t.Run("Error_DatabaseError", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with database error
		// Arrange - Resetear singleton y configurar mock con error de base de datos
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		inputProduct := models.Product{
			ProductCode: "PROD001",
			Description: "Test Product",
		}

		expectedError := errors.New("database connection failed")
		mockRepo.On("Create", mock.Anything, inputProduct).Return(models.Product{}, expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.Create(context.Background(), inputProduct)

		// Assert - Verificar que retorna error de base de datos
		// Assert - Verify that it returns database error
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})
}

// TestProductService_CreateByBatch agrupa todos los tests para el método CreateByBatch
// TestProductService_CreateByBatch groups all tests for the CreateByBatch method
func TestProductService_CreateByBatch(t *testing.T) {
	// Test case: Successful batch creation
	// Caso de prueba: Creación exitosa en lote
	t.Run("Success_CreatesBatchProducts", func(t *testing.T) {
		// Arrange - Reset singleton and set up mocks and test data
		// Arrange - Resetear singleton y configurar mocks y datos de prueba
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		inputProducts := []models.Product{
			{
				ProductCode:   "PROD001",
				Description:   "Test Product 1",
				Width:         10.5,
				ProductTypeID: 1,
			},
			{
				ProductCode:   "PROD002",
				Description:   "Test Product 2",
				Width:         8.5,
				ProductTypeID: 2,
			},
		}

		expectedProducts := []models.Product{
			{
				Id:            1,
				ProductCode:   "PROD001",
				Description:   "Test Product 1",
				Width:         10.5,
				ProductTypeID: 1,
			},
			{
				Id:            2,
				ProductCode:   "PROD002",
				Description:   "Test Product 2",
				Width:         8.5,
				ProductTypeID: 2,
			},
		}

		mockRepo.On("CreateByBatch", mock.Anything, inputProducts).Return(expectedProducts, nil)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.CreateByBatch(context.Background(), inputProducts)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, result)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(1), result[0].Id)
		assert.Equal(t, int64(2), result[1].Id)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Empty batch
	// Caso de prueba: Lote vacío
	t.Run("Success_EmptyBatch", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock for empty batch
		// Arrange - Resetear singleton y configurar mock para lote vacío
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		inputProducts := []models.Product{}
		expectedProducts := []models.Product{}

		mockRepo.On("CreateByBatch", mock.Anything, inputProducts).Return(expectedProducts, nil)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.CreateByBatch(context.Background(), inputProducts)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, result)
		assert.Len(t, result, 0)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Repository returns timeout error
	// Caso de prueba: El repositorio retorna error de timeout
	t.Run("Error_TimeoutExceeded", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with timeout error
		// Arrange - Resetear singleton y configurar mock con error de timeout
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		inputProducts := []models.Product{
			{ProductCode: "PROD001", Description: "Test Product 1"},
		}

		mockRepo.On("CreateByBatch", mock.Anything, inputProducts).Return([]models.Product{}, context.DeadlineExceeded)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.CreateByBatch(context.Background(), inputProducts)

		// Assert - Verificar que retorna error de timeout
		// Assert - Verify that it returns timeout error
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Database error
	// Caso de prueba: Error de base de datos
	t.Run("Error_DatabaseError", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with database error
		// Arrange - Resetear singleton y configurar mock con error de base de datos
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		inputProducts := []models.Product{
			{ProductCode: "PROD001", Description: "Test Product 1"},
		}

		expectedError := errors.New("database connection failed")
		mockRepo.On("CreateByBatch", mock.Anything, inputProducts).Return([]models.Product{}, expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.CreateByBatch(context.Background(), inputProducts)

		// Assert - Verificar que retorna error de base de datos
		// Assert - Verify that it returns database error
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})
}

// TestProductService_Update agrupa todos los tests para el método Update
// TestProductService_Update groups all tests for the Update method
func TestProductService_Update(t *testing.T) {
	// Test case: Successful product update
	// Caso de prueba: Actualización exitosa de producto
	t.Run("Success_UpdatesProduct", func(t *testing.T) {
		// Arrange - Reset singleton and set up mocks and test data
		// Arrange - Resetear singleton y configurar mocks y datos de prueba
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		updateProduct := models.Product{
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

		expectedProduct := updateProduct
		expectedProduct.Id = 1

		mockRepo.On("Update", mock.Anything, int64(1), updateProduct).Return(expectedProduct, nil)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.Update(context.Background(), 1, updateProduct)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, result)
		assert.Equal(t, int64(1), result.Id)
		assert.Equal(t, "PROD001_UPDATED", result.ProductCode)
		assert.Equal(t, "Updated Test Product", result.Description)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Product not found
	// Caso de prueba: Producto no encontrado
	t.Run("Error_ProductNotFound", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with not found error
		// Arrange - Resetear singleton y configurar mock con error de no encontrado
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		updateProduct := models.Product{
			ProductCode: "PROD999",
			Description: "Non-existent Product",
		}

		mockRepo.On("Update", mock.Anything, int64(999), updateProduct).Return(models.Product{}, error_message.ErrNotFound)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.Update(context.Background(), 999, updateProduct)

		// Assert - Verificar que retorna error de no encontrado
		// Assert - Verify that it returns not found error
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Repository returns timeout error
	// Caso de prueba: El repositorio retorna error de timeout
	t.Run("Error_TimeoutExceeded", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with timeout error
		// Arrange - Resetear singleton y configurar mock con error de timeout
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		updateProduct := models.Product{
			ProductCode: "PROD001",
			Description: "Test Product",
		}

		mockRepo.On("Update", mock.Anything, int64(1), updateProduct).Return(models.Product{}, context.DeadlineExceeded)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.Update(context.Background(), 1, updateProduct)

		// Assert - Verificar que retorna error de timeout
		// Assert - Verify that it returns timeout error
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Database error
	// Caso de prueba: Error de base de datos
	t.Run("Error_DatabaseError", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with database error
		// Arrange - Resetear singleton y configurar mock con error de base de datos
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		updateProduct := models.Product{
			ProductCode: "PROD001",
			Description: "Test Product",
		}

		expectedError := errors.New("database connection failed")
		mockRepo.On("Update", mock.Anything, int64(1), updateProduct).Return(models.Product{}, expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		result, err := service.Update(context.Background(), 1, updateProduct)

		// Assert - Verificar que retorna error de base de datos
		// Assert - Verify that it returns database error
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})
}

// TestProductService_Delete agrupa todos los tests para el método Delete
// TestProductService_Delete groups all tests for the Delete method
func TestProductService_Delete(t *testing.T) {
	// Test case: Successful product deletion
	// Caso de prueba: Eliminación exitosa de producto
	t.Run("Success_DeletesProduct", func(t *testing.T) {
		// Arrange - Reset singleton and set up mocks and test data
		// Arrange - Resetear singleton y configurar mocks y datos de prueba
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		err := service.Delete(context.Background(), 1)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Product not found
	// Caso de prueba: Producto no encontrado
	t.Run("Error_ProductNotFound", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with not found error
		// Arrange - Resetear singleton y configurar mock con error de no encontrado
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(999)).Return(error_message.ErrNotFound)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		err := service.Delete(context.Background(), 999)

		// Assert - Verificar que retorna error de no encontrado
		// Assert - Verify that it returns not found error
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Repository returns timeout error
	// Caso de prueba: El repositorio retorna error de timeout
	t.Run("Error_TimeoutExceeded", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with timeout error
		// Arrange - Resetear singleton y configurar mock con error de timeout
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		mockRepo.On("Delete", mock.Anything, int64(1)).Return(context.DeadlineExceeded)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		err := service.Delete(context.Background(), 1)

		// Assert - Verificar que retorna error de timeout
		// Assert - Verify that it returns timeout error
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Database error
	// Caso de prueba: Error de base de datos
	t.Run("Error_DatabaseError", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with database error
		// Arrange - Resetear singleton y configurar mock con error de base de datos
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		expectedError := errors.New("database connection failed")
		mockRepo.On("Delete", mock.Anything, int64(1)).Return(expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		err := service.Delete(context.Background(), 1)

		// Assert - Verificar que retorna error de base de datos
		// Assert - Verify that it returns database error
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockRepo.AssertExpectations(t)
	})
}

// TestProductService_ExistById agrupa todos los tests para el método ExistById
// TestProductService_ExistById groups all tests for the ExistById method
func TestProductService_ExistById(t *testing.T) {
	// Test case: Product exists
	// Caso de prueba: El producto existe
	t.Run("Success_ProductExists", func(t *testing.T) {
		// Arrange - Reset singleton and set up mocks and test data
		// Arrange - Resetear singleton y configurar mocks y datos de prueba
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		mockRepo.On("Exists", mock.Anything, int64(1)).Return(true, nil)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		exists, err := service.ExistById(context.Background(), 1)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.True(t, exists)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Product does not exist
	// Caso de prueba: El producto no existe
	t.Run("Success_ProductDoesNotExist", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock for product that does not exist
		// Arrange - Resetear singleton y configurar mock para producto que no existe
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		mockRepo.On("Exists", mock.Anything, int64(999)).Return(false, nil)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		exists, err := service.ExistById(context.Background(), 999)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Repository returns timeout error
	// Caso de prueba: El repositorio retorna error de timeout
	t.Run("Error_TimeoutExceeded", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with timeout error
		// Arrange - Resetear singleton y configurar mock con error de timeout
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		mockRepo.On("Exists", mock.Anything, int64(1)).Return(false, context.DeadlineExceeded)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		exists, err := service.ExistById(context.Background(), 1)

		// Assert - Verificar que retorna error de timeout
		// Assert - Verify that it returns timeout error
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})

	// Test case: Database error
	// Caso de prueba: Error de base de datos
	t.Run("Error_DatabaseError", func(t *testing.T) {
		// Arrange - Reset singleton and set up mock with database error
		// Arrange - Resetear singleton y configurar mock con error de base de datos
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		expectedError := errors.New("database connection failed")
		mockRepo.On("Exists", mock.Anything, int64(1)).Return(false, expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		exists, err := service.ExistById(context.Background(), 1)

		// Assert - Verificar que retorna error de base de datos
		// Assert - Verify that it returns database error
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})
}

// TestNewProductService agrupa todos los tests para el constructor NewProductService
// TestNewProductService groups all tests for the NewProductService constructor
func TestNewProductService(t *testing.T) {
	// Test case: Successful service creation
	// Caso de prueba: Creación exitosa del servicio
	t.Run("Success_CreatesService", func(t *testing.T) {
		// Arrange - Reset singleton and create mock repository
		// Arrange - Resetear singleton y crear mock repository
		productServiceInstance = nil
		mockRepo := new(tests.ProductRepositoryMock)

		// Act - Crear nuevo servicio
		// Act - Create new service
		service := NewProductService(mockRepo)

		// Assert - Verificar que el servicio se creó correctamente
		// Assert - Verify that the service was created correctly
		assert.NotNil(t, service)

		// Verificar que la instancia implementa la interfaz
		// Verify that the instance implements the interface
		var _ ProductService = service
	})

	// Test case: Singleton pattern behavior
	// Caso de prueba: Comportamiento de patrón singleton
	t.Run("Success_SingletonBehavior", func(t *testing.T) {
		// Arrange - Reset singleton instance for test isolation
		// Arrange - Resetear instancia singleton para aislamiento de test
		productServiceInstance = nil

		mockRepo1 := new(tests.ProductRepositoryMock)
		mockRepo2 := new(tests.ProductRepositoryMock)

		// Act - Crear dos instancias del servicio
		// Act - Create two service instances
		service1 := NewProductService(mockRepo1)
		service2 := NewProductService(mockRepo2)

		// Assert - Verificar que retorna la misma instancia (singleton)
		// Assert - Verify that it returns the same instance (singleton)
		assert.Equal(t, service1, service2)

		// Reset singleton instance after test
		// Resetear instancia singleton después del test
		productServiceInstance = nil
	})

	// Test case: Nil repository handling
	// Caso de prueba: Manejo de repositorio nulo
	t.Run("Success_NilRepository", func(t *testing.T) {
		// Arrange - Reset singleton for test isolation
		// Arrange - Resetear singleton para aislamiento de test
		productServiceInstance = nil

		// Act - Crear servicio con repositorio nulo
		// Act - Create service with nil repository
		service := NewProductService(nil)

		// Assert - Verificar que el servicio se crea incluso con repositorio nulo
		// Assert - Verify that the service is created even with nil repository
		assert.NotNil(t, service)

		// Reset singleton instance after test
		// Resetear instancia singleton después del test
		productServiceInstance = nil
	})
}
