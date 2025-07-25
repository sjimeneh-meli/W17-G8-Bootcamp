package repositories

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestProductRepository crea una nueva instancia del repository para tests evitando el singleton
// newTestProductRepository creates a new repository instance for tests avoiding the singleton
func newTestProductRepository(db *sql.DB) ProductRepository {
	return &service{
		db: db,
	}
}

// TestProductRepository_GetAll agrupa todos los tests para el método GetAll
// TestProductRepository_GetAll groups all tests for the GetAll method
func TestProductRepository_GetAll(t *testing.T) {
	// Test case: Successful retrieval of all products
	// Caso de prueba: Obtención exitosa de todos los productos
	t.Run("Success_ReturnsAllProducts", func(t *testing.T) {
		// Arrange - Configurar mock de base de datos y datos de prueba
		// Arrange - Set up database mock and test data
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		// Datos de productos esperados / Expected product data
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

		// Configurar expectativas del mock SQL / Set up SQL mock expectations
		rows := sqlmock.NewRows([]string{
			"id", "description", "expiration_rate", "freezing_rate", "height",
			"length", "net_weight", "product_code", "recommended_freezing_temperature",
			"width", "product_type_id", "seller_id",
		})

		for _, product := range expectedProducts {
			rows.AddRow(
				product.Id, product.Description, product.ExpirationRate, product.FreezingRate,
				product.Height, product.Length, product.NetWeight, product.ProductCode,
				product.RecommendedFreezingTemperature, product.Width, product.ProductTypeID, product.SellerID,
			)
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products")).
			WillReturnRows(rows)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		products, err := repo.GetAll(context.Background())

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Len(t, products, 2)
		assert.Equal(t, expectedProducts[0].ProductCode, products[0].ProductCode)
		assert.Equal(t, expectedProducts[1].ProductCode, products[1].ProductCode)

		// Verificar que todas las expectativas del mock se cumplieron
		// Verify that all mock expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Database query error
	// Caso de prueba: Error en la consulta de base de datos
	t.Run("Error_DatabaseQueryFailed", func(t *testing.T) {
		// Arrange - Configurar mock con error de consulta
		// Arrange - Set up mock with query error
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		expectedError := errors.New("database connection failed")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products")).
			WillReturnError(expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		products, err := repo.GetAll(context.Background())

		// Assert - Verificar que retorna error
		// Assert - Verify that it returns error
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "failed to query products")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Row scan error
	// Caso de prueba: Error al escanear filas
	t.Run("Error_RowScanFailed", func(t *testing.T) {
		// Arrange - Configurar mock con datos incompatibles para escaneo
		// Arrange - Set up mock with incompatible data for scanning
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		// Filas con tipos de datos incorrectos / Rows with incorrect data types
		rows := sqlmock.NewRows([]string{
			"id", "description", "expiration_rate", "freezing_rate", "height",
			"length", "net_weight", "product_code", "recommended_freezing_temperature",
			"width", "product_type_id", "seller_id",
		}).AddRow("invalid_id", "desc", 0.1, 0.8, 10.5, 15.7, 5.2, "PROD001", -18.0, 10.5, 1, nil)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products")).
			WillReturnRows(rows)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		products, err := repo.GetAll(context.Background())

		// Assert - Verificar que retorna error de escaneo
		// Assert - Verify that it returns scan error
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "failed to scan product row")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Empty result set
	// Caso de prueba: Conjunto de resultados vacío
	t.Run("Success_EmptyResultSet", func(t *testing.T) {
		// Arrange - Configurar mock sin filas
		// Arrange - Set up mock with no rows
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		rows := sqlmock.NewRows([]string{
			"id", "description", "expiration_rate", "freezing_rate", "height",
			"length", "net_weight", "product_code", "recommended_freezing_temperature",
			"width", "product_type_id", "seller_id",
		})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products")).
			WillReturnRows(rows)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		products, err := repo.GetAll(context.Background())

		// Assert - Verificar resultado vacío exitoso
		// Assert - Verify successful empty result
		assert.NoError(t, err)
		assert.Empty(t, products)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestProductRepository_GetByID agrupa todos los tests para el método GetByID
// TestProductRepository_GetByID groups all tests for the GetByID method
func TestProductRepository_GetByID(t *testing.T) {
	// Test case: Successful product retrieval by ID
	// Caso de prueba: Obtención exitosa de producto por ID
	t.Run("Success_ReturnsProductByID", func(t *testing.T) {
		// Arrange - Configurar mock de base de datos y datos de prueba
		// Arrange - Set up database mock and test data
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

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

		rows := sqlmock.NewRows([]string{
			"id", "description", "expiration_rate", "freezing_rate", "height",
			"length", "net_weight", "product_code", "recommended_freezing_temperature",
			"width", "product_type_id", "seller_id",
		}).AddRow(
			expectedProduct.Id, expectedProduct.Description, expectedProduct.ExpirationRate,
			expectedProduct.FreezingRate, expectedProduct.Height, expectedProduct.Length,
			expectedProduct.NetWeight, expectedProduct.ProductCode, expectedProduct.RecommendedFreezingTemperature,
			expectedProduct.Width, expectedProduct.ProductTypeID, expectedProduct.SellerID,
		)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products WHERE id = ?")).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		product, err := repo.GetByID(context.Background(), 1)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct.Id, product.Id)
		assert.Equal(t, expectedProduct.ProductCode, product.ProductCode)
		assert.Equal(t, expectedProduct.Description, product.Description)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Product not found
	// Caso de prueba: Producto no encontrado
	t.Run("Error_ProductNotFound", func(t *testing.T) {
		// Arrange - Configurar mock para retornar sql.ErrNoRows
		// Arrange - Set up mock to return sql.ErrNoRows
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products WHERE id = ?")).
			WithArgs(int64(999)).
			WillReturnError(sql.ErrNoRows)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		product, err := repo.GetByID(context.Background(), 999)

		// Assert - Verificar que retorna ErrNotFound
		// Assert - Verify that it returns ErrNotFound
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)
		assert.Equal(t, models.Product{}, product)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Database scan error
	// Caso de prueba: Error de escaneo de base de datos
	t.Run("Error_DatabaseScanError", func(t *testing.T) {
		// Arrange - Configurar mock con datos incompatibles
		// Arrange - Set up mock with incompatible data
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		rows := sqlmock.NewRows([]string{
			"id", "description", "expiration_rate", "freezing_rate", "height",
			"length", "net_weight", "product_code", "recommended_freezing_temperature",
			"width", "product_type_id", "seller_id",
		}).AddRow("invalid_id", "desc", 0.1, 0.8, 10.5, 15.7, 5.2, "PROD001", -18.0, 10.5, 1, nil)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products WHERE id = ?")).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		product, err := repo.GetByID(context.Background(), 1)

		// Assert - Verificar que retorna error de escaneo
		// Assert - Verify that it returns scan error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to scan product with id")
		assert.Equal(t, models.Product{}, product)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestProductRepository_Create agrupa todos los tests para el método Create
// TestProductRepository_Create groups all tests for the Create method
func TestProductRepository_Create(t *testing.T) {
	// Test case: Successful product creation
	// Caso de prueba: Creación exitosa de producto
	t.Run("Success_CreatesProduct", func(t *testing.T) {
		// Arrange - Configurar mock de base de datos y datos de prueba
		// Arrange - Set up database mock and test data
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		newProduct := models.Product{
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

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO products (description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")).
			WithArgs(newProduct.Description, newProduct.ExpirationRate, newProduct.FreezingRate, newProduct.Height,
				newProduct.Length, newProduct.NetWeight, newProduct.ProductCode, newProduct.RecommendedFreezingTemperature,
				newProduct.Width, newProduct.ProductTypeID, newProduct.SellerID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		createdProduct, err := repo.Create(context.Background(), newProduct)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Equal(t, int64(1), createdProduct.Id)
		assert.Equal(t, newProduct.ProductCode, createdProduct.ProductCode)
		assert.Equal(t, newProduct.Description, createdProduct.Description)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Database execution error
	// Caso de prueba: Error de ejecución de base de datos
	t.Run("Error_DatabaseExecutionFailed", func(t *testing.T) {
		// Arrange - Configurar mock con error de ejecución
		// Arrange - Set up mock with execution error
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		newProduct := models.Product{
			ProductCode: "PROD001",
			Description: "Test Product",
		}

		expectedError := errors.New("duplicate key constraint violation")
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO products (description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")).
			WithArgs(newProduct.Description, newProduct.ExpirationRate, newProduct.FreezingRate, newProduct.Height,
				newProduct.Length, newProduct.NetWeight, newProduct.ProductCode, newProduct.RecommendedFreezingTemperature,
				newProduct.Width, newProduct.ProductTypeID, newProduct.SellerID).
			WillReturnError(expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		createdProduct, err := repo.Create(context.Background(), newProduct)

		// Assert - Verificar que retorna error
		// Assert - Verify that it returns error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create product")
		assert.Equal(t, models.Product{}, createdProduct)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: LastInsertId error
	// Caso de prueba: Error en LastInsertId
	t.Run("Error_LastInsertIdFailed", func(t *testing.T) {
		// Arrange - Configurar mock con error en LastInsertId
		// Arrange - Set up mock with LastInsertId error
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		newProduct := models.Product{
			ProductCode: "PROD001",
			Description: "Test Product",
		}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO products (description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")).
			WithArgs(newProduct.Description, newProduct.ExpirationRate, newProduct.FreezingRate, newProduct.Height,
				newProduct.Length, newProduct.NetWeight, newProduct.ProductCode, newProduct.RecommendedFreezingTemperature,
				newProduct.Width, newProduct.ProductTypeID, newProduct.SellerID).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id not available")))

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		createdProduct, err := repo.Create(context.Background(), newProduct)

		// Assert - Verificar que retorna error
		// Assert - Verify that it returns error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get last insert id")
		assert.Equal(t, models.Product{}, createdProduct)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestProductRepository_CreateByBatch agrupa todos los tests para el método CreateByBatch
// TestProductRepository_CreateByBatch groups all tests for the CreateByBatch method
func TestProductRepository_CreateByBatch(t *testing.T) {
	// Test case: Successful batch creation
	// Caso de prueba: Creación exitosa en lote
	t.Run("Success_CreatesBatchProducts", func(t *testing.T) {
		// Arrange - Configurar mock de base de datos y datos de prueba
		// Arrange - Set up database mock and test data
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		products := []models.Product{
			{
				ProductCode: "PROD001",
				Description: "Test Product 1",
			},
			{
				ProductCode: "PROD002",
				Description: "Test Product 2",
			},
		}

		// Configurar expectativas de transacción / Set up transaction expectations
		mock.ExpectBegin()
		mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO products (description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"))

		// Expectativas para cada producto / Expectations for each product
		for i, product := range products {
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO products (description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")).
				WithArgs(product.Description, product.ExpirationRate, product.FreezingRate, product.Height,
					product.Length, product.NetWeight, product.ProductCode, product.RecommendedFreezingTemperature,
					product.Width, product.ProductTypeID, product.SellerID).
				WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
		}

		mock.ExpectCommit()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		createdProducts, err := repo.CreateByBatch(context.Background(), products)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Len(t, createdProducts, 2)
		assert.Equal(t, int64(1), createdProducts[0].Id)
		assert.Equal(t, int64(2), createdProducts[1].Id)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Transaction begin error
	// Caso de prueba: Error al iniciar transacción
	t.Run("Error_TransactionBeginFailed", func(t *testing.T) {
		// Arrange - Configurar mock con error de transacción
		// Arrange - Set up mock with transaction error
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		products := []models.Product{{ProductCode: "PROD001"}}

		expectedError := errors.New("failed to begin transaction")
		mock.ExpectBegin().WillReturnError(expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		createdProducts, err := repo.CreateByBatch(context.Background(), products)

		// Assert - Verificar que retorna error
		// Assert - Verify that it returns error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to begin transaction")
		assert.Nil(t, createdProducts)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Prepare statement error
	// Caso de prueba: Error al preparar declaración
	t.Run("Error_PrepareStatementFailed", func(t *testing.T) {
		// Arrange - Configurar mock con error de preparación
		// Arrange - Set up mock with prepare error
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		products := []models.Product{{ProductCode: "PROD001"}}

		mock.ExpectBegin()
		expectedError := errors.New("failed to prepare statement")
		mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO products (description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		createdProducts, err := repo.CreateByBatch(context.Background(), products)

		// Assert - Verificar que retorna error
		// Assert - Verify that it returns error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to prepare statement")
		assert.Nil(t, createdProducts)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Execution error during batch
	// Caso de prueba: Error de ejecución durante el lote
	t.Run("Error_ExecutionFailedDuringBatch", func(t *testing.T) {
		// Arrange - Configurar mock con error durante ejecución
		// Arrange - Set up mock with execution error
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		products := []models.Product{{ProductCode: "PROD001"}}

		mock.ExpectBegin()
		mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO products (description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"))

		expectedError := errors.New("constraint violation")
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO products (description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")).
			WithArgs(products[0].Description, products[0].ExpirationRate, products[0].FreezingRate, products[0].Height,
				products[0].Length, products[0].NetWeight, products[0].ProductCode, products[0].RecommendedFreezingTemperature,
				products[0].Width, products[0].ProductTypeID, products[0].SellerID).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		createdProducts, err := repo.CreateByBatch(context.Background(), products)

		// Assert - Verificar que retorna error
		// Assert - Verify that it returns error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute statement for product")
		assert.Nil(t, createdProducts)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestProductRepository_Update agrupa todos los tests para el método Update
// TestProductRepository_Update groups all tests for the Update method
func TestProductRepository_Update(t *testing.T) {
	// Test case: Successful product update
	// Caso de prueba: Actualización exitosa de producto
	t.Run("Success_UpdatesProduct", func(t *testing.T) {
		// Arrange - Configurar mock de base de datos y datos de prueba
		// Arrange - Set up database mock and test data
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		productToUpdate := models.Product{
			ProductCode: "PROD001_UPDATED",
			Description: "Updated Test Product",
		}

		mock.ExpectExec(regexp.QuoteMeta("UPDATE products SET description = ?, expiration_rate = ?, freezing_rate = ?, height = ?, length = ?, net_weight = ?, product_code = ?, recommended_freezing_temperature = ?, width = ?, product_type_id = ?, seller_id = ? WHERE id = ?")).
			WithArgs(productToUpdate.Description, productToUpdate.ExpirationRate, productToUpdate.FreezingRate,
				productToUpdate.Height, productToUpdate.Length, productToUpdate.NetWeight, productToUpdate.ProductCode,
				productToUpdate.RecommendedFreezingTemperature, productToUpdate.Width, productToUpdate.ProductTypeID,
				productToUpdate.SellerID, int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		updatedProduct, err := repo.Update(context.Background(), 1, productToUpdate)

		// Assert - Verificar resultados
		// Assert - Verify results
		assert.NoError(t, err)
		assert.Equal(t, int64(1), updatedProduct.Id)
		assert.Equal(t, productToUpdate.ProductCode, updatedProduct.ProductCode)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Product not found for update
	// Caso de prueba: Producto no encontrado para actualizar
	t.Run("Error_ProductNotFoundForUpdate", func(t *testing.T) {
		// Arrange - Configurar mock sin filas afectadas
		// Arrange - Set up mock with no rows affected
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		productToUpdate := models.Product{ProductCode: "PROD001"}

		mock.ExpectExec(regexp.QuoteMeta("UPDATE products SET description = ?, expiration_rate = ?, freezing_rate = ?, height = ?, length = ?, net_weight = ?, product_code = ?, recommended_freezing_temperature = ?, width = ?, product_type_id = ?, seller_id = ? WHERE id = ?")).
			WithArgs(productToUpdate.Description, productToUpdate.ExpirationRate, productToUpdate.FreezingRate,
				productToUpdate.Height, productToUpdate.Length, productToUpdate.NetWeight, productToUpdate.ProductCode,
				productToUpdate.RecommendedFreezingTemperature, productToUpdate.Width, productToUpdate.ProductTypeID,
				productToUpdate.SellerID, int64(999)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		updatedProduct, err := repo.Update(context.Background(), 999, productToUpdate)

		// Assert - Verificar que retorna ErrNotFound
		// Assert - Verify that it returns ErrNotFound
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)
		assert.Equal(t, models.Product{}, updatedProduct)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Database execution error
	// Caso de prueba: Error de ejecución de base de datos
	t.Run("Error_DatabaseExecutionFailed", func(t *testing.T) {
		// Arrange - Configurar mock con error de ejecución
		// Arrange - Set up mock with execution error
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		productToUpdate := models.Product{ProductCode: "PROD001"}

		expectedError := errors.New("database constraint violation")
		mock.ExpectExec(regexp.QuoteMeta("UPDATE products SET description = ?, expiration_rate = ?, freezing_rate = ?, height = ?, length = ?, net_weight = ?, product_code = ?, recommended_freezing_temperature = ?, width = ?, product_type_id = ?, seller_id = ? WHERE id = ?")).
			WithArgs(productToUpdate.Description, productToUpdate.ExpirationRate, productToUpdate.FreezingRate,
				productToUpdate.Height, productToUpdate.Length, productToUpdate.NetWeight, productToUpdate.ProductCode,
				productToUpdate.RecommendedFreezingTemperature, productToUpdate.Width, productToUpdate.ProductTypeID,
				productToUpdate.SellerID, int64(1)).
			WillReturnError(expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		updatedProduct, err := repo.Update(context.Background(), 1, productToUpdate)

		// Assert - Verificar que retorna error
		// Assert - Verify that it returns error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update product")
		assert.Equal(t, models.Product{}, updatedProduct)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestProductRepository_Delete agrupa todos los tests para el método Delete
// TestProductRepository_Delete groups all tests for the Delete method
func TestProductRepository_Delete(t *testing.T) {
	// Test case: Successful product deletion
	// Caso de prueba: Eliminación exitosa de producto
	t.Run("Success_DeletesProduct", func(t *testing.T) {
		// Arrange - Configurar mock de base de datos
		// Arrange - Set up database mock
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = ?")).
			WithArgs(int64(1)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		err = repo.Delete(context.Background(), 1)

		// Assert - Verificar que no hay error
		// Assert - Verify that there is no error
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Product not found for deletion
	// Caso de prueba: Producto no encontrado para eliminar
	t.Run("Error_ProductNotFoundForDeletion", func(t *testing.T) {
		// Arrange - Configurar mock sin filas afectadas
		// Arrange - Set up mock with no rows affected
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = ?")).
			WithArgs(int64(999)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		err = repo.Delete(context.Background(), 999)

		// Assert - Verificar que retorna ErrNotFound
		// Assert - Verify that it returns ErrNotFound
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Database execution error
	// Caso de prueba: Error de ejecución de base de datos
	t.Run("Error_DatabaseExecutionFailed", func(t *testing.T) {
		// Arrange - Configurar mock con error de ejecución
		// Arrange - Set up mock with execution error
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		expectedError := errors.New("foreign key constraint violation")
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = ?")).
			WithArgs(int64(1)).
			WillReturnError(expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		err = repo.Delete(context.Background(), 1)

		// Assert - Verificar que retorna error
		// Assert - Verify that it returns error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete product")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestProductRepository_Exists agrupa todos los tests para el método Exists
// TestProductRepository_Exists groups all tests for the Exists method
func TestProductRepository_Exists(t *testing.T) {
	// Test case: Product exists
	// Caso de prueba: El producto existe
	t.Run("Success_ProductExists", func(t *testing.T) {
		// Arrange - Configurar mock que retorna producto existente
		// Arrange - Set up mock that returns existing product
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM products WHERE id = ? LIMIT 1")).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		exists, err := repo.Exists(context.Background(), 1)

		// Assert - Verificar que el producto existe
		// Assert - Verify that the product exists
		assert.NoError(t, err)
		assert.True(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Product does not exist
	// Caso de prueba: El producto no existe
	t.Run("Success_ProductDoesNotExist", func(t *testing.T) {
		// Arrange - Configurar mock que retorna sql.ErrNoRows
		// Arrange - Set up mock that returns sql.ErrNoRows
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM products WHERE id = ? LIMIT 1")).
			WithArgs(int64(999)).
			WillReturnError(sql.ErrNoRows)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		exists, err := repo.Exists(context.Background(), 999)

		// Assert - Verificar que el producto no existe
		// Assert - Verify that the product does not exist
		assert.NoError(t, err)
		assert.False(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Database query error
	// Caso de prueba: Error de consulta de base de datos
	t.Run("Error_DatabaseQueryFailed", func(t *testing.T) {
		// Arrange - Configurar mock con error de consulta
		// Arrange - Set up mock with query error
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		expectedError := errors.New("database connection error")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM products WHERE id = ? LIMIT 1")).
			WithArgs(int64(1)).
			WillReturnError(expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		exists, err := repo.Exists(context.Background(), 1)

		// Assert - Verificar que retorna error
		// Assert - Verify that it returns error
		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), "error checking existence of product")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestProductRepository_ExistsByProductCode agrupa todos los tests para el método ExistsByProductCode
// TestProductRepository_ExistsByProductCode groups all tests for the ExistsByProductCode method
func TestProductRepository_ExistsByProductCode(t *testing.T) {
	// Test case: Product code exists
	// Caso de prueba: El código de producto existe
	t.Run("Success_ProductCodeExists", func(t *testing.T) {
		// Arrange - Configurar mock que retorna código existente
		// Arrange - Set up mock that returns existing code
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM products WHERE product_code = ? LIMIT 1")).
			WithArgs("PROD001").
			WillReturnRows(rows)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		exists, err := repo.ExistsByProductCode(context.Background(), "PROD001")

		// Assert - Verificar que el código existe
		// Assert - Verify that the code exists
		assert.NoError(t, err)
		assert.True(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Product code does not exist
	// Caso de prueba: El código de producto no existe
	t.Run("Success_ProductCodeDoesNotExist", func(t *testing.T) {
		// Arrange - Configurar mock que retorna sql.ErrNoRows
		// Arrange - Set up mock that returns sql.ErrNoRows
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM products WHERE product_code = ? LIMIT 1")).
			WithArgs("NONEXISTENT").
			WillReturnError(sql.ErrNoRows)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		exists, err := repo.ExistsByProductCode(context.Background(), "NONEXISTENT")

		// Assert - Verificar que el código no existe
		// Assert - Verify that the code does not exist
		assert.NoError(t, err)
		assert.False(t, exists)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test case: Database query error
	// Caso de prueba: Error de consulta de base de datos
	t.Run("Error_DatabaseQueryFailed", func(t *testing.T) {
		// Arrange - Configurar mock con error de consulta
		// Arrange - Set up mock with query error
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := newTestProductRepository(db)

		expectedError := errors.New("database timeout")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM products WHERE product_code = ? LIMIT 1")).
			WithArgs("PROD001").
			WillReturnError(expectedError)

		// Act - Ejecutar el método bajo prueba
		// Act - Execute the method under test
		exists, err := repo.ExistsByProductCode(context.Background(), "PROD001")

		// Assert - Verificar que retorna error
		// Assert - Verify that it returns error
		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), "error checking existence of product code")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestNewProductRepository prueba el constructor del ProductRepository
// TestNewProductRepository tests the ProductRepository constructor
func TestNewProductRepository(t *testing.T) {
	// Test case: Successful repository creation
	// Caso de prueba: Creación exitosa del repositorio
	t.Run("Success_CreatesRepository", func(t *testing.T) {
		// Arrange - Crear mock de base de datos
		// Arrange - Create database mock
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		// Act - Crear nuevo repositorio
		// Act - Create new repository
		repo := NewProductRepository(db)

		// Assert - Verificar que el repositorio se creó correctamente
		// Assert - Verify that the repository was created correctly
		assert.NotNil(t, repo)
	})
}
