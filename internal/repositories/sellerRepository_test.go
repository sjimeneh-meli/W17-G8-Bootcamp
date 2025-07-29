package repositories_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	t.Run("create ok", func(t *testing.T) {
		expectedSeller := models.Seller{
			Id:          1,
			CID:         "SEL-001",
			CompanyName: "Test Company",
			Address:     "Test Address 123",
			Telephone:   "555-1234",
			LocalityID:  1,
		}

		inputSeller := models.Seller{
			CID:         "SEL-001",
			CompanyName: "Test Company",
			Address:     "Test Address 123",
			Telephone:   "555-1234",
			LocalityID:  1,
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM sellers WHERE cid = \\?\\)").
			WithArgs("SEL-001").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectExec("INSERT INTO sellers \\(cid, company_name, address, telephone, locality_id\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("SEL-001", "Test Company", "Test Address 123", "555-1234", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		sellerCreated, err := repository.Save(inputSeller)
		assert.NoError(t, err)
		assert.Len(t, sellerCreated, 1)
		assert.Equal(t, expectedSeller, sellerCreated[0])

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create conflict", func(t *testing.T) {
		// Arrange - Input data
		// 1. Crear conexión mock
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// 2. Crear instancia del repositorio con la base de datos mock
		repositories.ResetSellerRepositorySingleton()
		repo := repositories.NewSQLSellerRepository(db)

		// 3. Seller de prueba que intenta guardarse
		inputSeller := models.Seller{
			CID:         "SEL-001",
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "555-1234",
			LocalityID:  1,
		}

		// 4. Mockear la consulta EXISTS para el CID (ya existe)
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM sellers WHERE cid = \\?\\)").
			WithArgs(inputSeller.CID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))

		// 5. Ejecutar el método Save
		result, err := repo.Save(inputSeller)

		// 6. Validaciones
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, error_message.ErrAlreadyExists, err)

		// 7. Verificar que todas las expectativas se cumplieron
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("should return error when CID checking query fails", func(t *testing.T) {
		// Arrange - Input data
		inputSeller := models.Seller{
			CID:         "SEL-001",
			CompanyName: "Test Company",
			Address:     "Test Address 123",
			Telephone:   "555-1234",
			LocalityID:  1,
		}

		// Arrange - Mock database setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange - First query fails (CID checking)
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM sellers WHERE cid = \\?\\)").
			WithArgs("SEL-001").
			WillReturnError(error_message.ErrFailedCheckingExistence)

		// Arrange - Create repository instance
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute Save method
		sellerCreated, err := repository.Save(inputSeller)

		// Assert - Verify error handling
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrFailedCheckingExistence, err)
		assert.Nil(t, sellerCreated)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when locality checking query fails", func(t *testing.T) {
		// Arrange - Input data
		inputSeller := models.Seller{
			CID:         "SEL-001",
			CompanyName: "Test Company",
			Address:     "Test Address 123",
			Telephone:   "555-1234",
			LocalityID:  1,
		}

		// Arrange - Mock database setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange - First query succeeds (CID doesn't exist)
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM sellers WHERE cid = \\?\\)").
			WithArgs("SEL-001").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Arrange - Second query fails (locality checking)
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(1).
			WillReturnError(error_message.ErrFailedCheckingExistence)

		// Arrange - Create repository instance
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute Save method
		sellerCreated, err := repository.Save(inputSeller)

		// Assert - Verify error handling
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrFailedCheckingExistence, err)
		assert.Nil(t, sellerCreated)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when locality does not exist", func(t *testing.T) {
		// Arrange - Input data
		inputSeller := models.Seller{
			CID:         "SEL-001",
			CompanyName: "Test Company",
			Address:     "Test Address 123",
			Telephone:   "555-1234",
			LocalityID:  999, // Non-existent locality
		}

		// Arrange - Mock database setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange - First query succeeds (CID doesn't exist)
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM sellers WHERE cid = \\?\\)").
			WithArgs("SEL-001").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Arrange - Second query succeeds but locality doesn't exist
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(999).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Arrange - Create repository instance
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute Save method
		sellerCreated, err := repository.Save(inputSeller)

		// Assert - Verify dependency error
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrDependencyNotFound, err)
		assert.Nil(t, sellerCreated)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when INSERT query fails", func(t *testing.T) {
		// Arrange - Input data
		inputSeller := models.Seller{
			CID:         "SEL-001",
			CompanyName: "Test Company",
			Address:     "Test Address 123",
			Telephone:   "555-1234",
			LocalityID:  1,
		}

		// Arrange - Mock database setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange - First query succeeds (CID doesn't exist)
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM sellers WHERE cid = \\?\\)").
			WithArgs("SEL-001").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Arrange - Second query succeeds (locality exists)
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		// Arrange - INSERT fails
		mock.ExpectExec("INSERT INTO sellers \\(cid, company_name, address, telephone, locality_id\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("SEL-001", "Test Company", "Test Address 123", "555-1234", 1).
			WillReturnError(error_message.ErrQuery)

		// Arrange - Create repository instance
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute Save method
		sellerCreated, err := repository.Save(inputSeller)

		// Assert - Verify insert error
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Nil(t, sellerCreated)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when LastInsertId fails", func(t *testing.T) {
		// Arrange - Input data
		inputSeller := models.Seller{
			CID:         "SEL-001",
			CompanyName: "Test Company",
			Address:     "Test Address 123",
			Telephone:   "555-1234",
			LocalityID:  1,
		}

		// Arrange - Mock database setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange - First query succeeds (CID doesn't exist)
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM sellers WHERE cid = \\?\\)").
			WithArgs("SEL-001").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Arrange - Second query succeeds (locality exists)
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		// Arrange - INSERT succeeds but LastInsertId fails
		mock.ExpectExec("INSERT INTO sellers \\(cid, company_name, address, telephone, locality_id\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
			WithArgs("SEL-001", "Test Company", "Test Address 123", "555-1234", 1).
			WillReturnResult(sqlmock.NewErrorResult(error_message.ErrQuery)) // Result that fails on LastInsertId

		// Arrange - Create repository instance
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute Save method
		sellerCreated, err := repository.Save(inputSeller)

		// Assert - Verify LastInsertId error
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err) // Original error from LastInsertId
		assert.Nil(t, sellerCreated)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestGet(t *testing.T) {
	t.Run("find all", func(t *testing.T) {
		expectedSellers := []models.Seller{
			{
				Id:          1,
				CID:         "SEL-001",
				CompanyName: "Company One",
				Address:     "Address One 123",
				Telephone:   "555-0001",
				LocalityID:  1,
			},
			{
				Id:          2,
				CID:         "SEL-002",
				CompanyName: "Company Two",
				Address:     "Address Two 456",
				Telephone:   "555-0002",
				LocalityID:  2,
			},
			{
				Id:          3,
				CID:         "SEL-003",
				CompanyName: "Company Three",
				Address:     "Address Three 789",
				Telephone:   "555-0003",
				LocalityID:  1,
			},
		}

		// Arrange - Mock database setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
			AddRow(1, "SEL-001", "Company One", "Address One 123", "555-0001", 1).
			AddRow(2, "SEL-002", "Company Two", "Address Two 456", "555-0002", 2).
			AddRow(3, "SEL-003", "Company Three", "Address Three 789", "555-0003", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers").
			WillReturnRows(rows)

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		sellersResult, err := repository.GetAll()

		assert.NoError(t, err)
		assert.Len(t, sellersResult, 3)
		assert.Equal(t, expectedSellers, sellersResult)

		assert.Equal(t, "SEL-001", sellersResult[0].CID)
		assert.Equal(t, "Company Two", sellersResult[1].CompanyName)
		assert.Equal(t, 1, sellersResult[2].LocalityID)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return empty slice when no sellers exist", func(t *testing.T) {
		// Arrange - Expected empty result
		expectedSellers := []models.Seller{}

		// Arrange - Mock database setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange - Configure mock expectation: Return empty rows
		rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})
		// No AddRow() calls → Empty result set

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers").
			WillReturnRows(rows)

		// Arrange - Create repository instance
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute GetAll method
		sellersResult, err := repository.GetAll()

		// Assert - Verify empty results
		assert.NoError(t, err)
		assert.Len(t, sellersResult, 0) // Verify zero elements
		assert.Equal(t, expectedSellers, sellersResult)
		assert.Empty(t, sellersResult)  // Alternative way to check empty slice
		assert.NotNil(t, sellersResult) // Slice exists but is empty (not nil)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when database query fails", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers").
			WillReturnError(error_message.ErrQuery)

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		sellersResult, err := repository.GetAll()

		assert.Error(t, err)
		assert.Nil(t, sellersResult)
		assert.Equal(t, error_message.ErrQuery, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when row scanning fails", func(t *testing.T) {
		// Arrange - Mock database setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange - Configure mock expectation: Return rows with wrong data types to cause scan error
		rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
			AddRow("invalid_id", "SEL-001", "Company One", "Address One", "555-0001", 1) // String instead of int for ID

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers").
			WillReturnRows(rows)

		// Arrange - Create repository instance
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute GetAll method
		sellersResult, err := repository.GetAll()

		// Assert - Verify error handling during scan
		assert.Error(t, err)            // Must have error (scan/conversion error)
		assert.Nil(t, sellersResult)    // Should return nil slice on scan error
		assert.NotEmpty(t, err.Error()) // Error message should not be empty

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestUpdate(t *testing.T) {
	t.Run("update ok", func(t *testing.T) {

		sellerId := 1

		// Existing seller data in database (before update)
		existingSeller := models.Seller{
			Id:          sellerId,
			CID:         "SEL-001-OLD",
			CompanyName: "Old Company Name",
			Address:     "Old Address 123",
			Telephone:   "555-0000",
			LocalityID:  1,
		}

		partialUpdateSeller := models.Seller{
			CID:         "SEL-001-UPDATED",
			CompanyName: "Updated Company Name",
			// Address and Telephone not to update (partial update)
			LocalityID: 2,
		}

		expectedSeller := models.Seller{
			Id:          existingSeller.Id,
			CID:         partialUpdateSeller.CID,         //  Updated
			CompanyName: partialUpdateSeller.CompanyName, //  Updated
			Address:     existingSeller.Address,          //  No updated
			Telephone:   existingSeller.Telephone,        //  No updated
			LocalityID:  partialUpdateSeller.LocalityID,  //  Updated
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
			AddRow(existingSeller.Id, existingSeller.CID, existingSeller.CompanyName, existingSeller.Address, existingSeller.Telephone, existingSeller.LocalityID)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(sellerId).
			WillReturnRows(rows)

		mock.ExpectExec("UPDATE sellers SET cid = \\?, company_name = \\?, address = \\?, telephone = \\?, locality_id = \\? WHERE id = \\?").
			WithArgs(expectedSeller.CID, expectedSeller.CompanyName, expectedSeller.Address, expectedSeller.Telephone, expectedSeller.LocalityID, sellerId).
			WillReturnResult(sqlmock.NewResult(0, 1)) // 0 = no LastInsertId needed, 1 = rows affected

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		updatedSellers, err := repository.Update(sellerId, partialUpdateSeller)

		assert.NoError(t, err)
		assert.Len(t, updatedSellers, 1)
		assert.Equal(t, expectedSeller, updatedSellers[0])

		// Assert - Verify specific merged fields
		assert.Equal(t, partialUpdateSeller.CID, updatedSellers[0].CID)                 // Updated field
		assert.Equal(t, partialUpdateSeller.CompanyName, updatedSellers[0].CompanyName) // Updated field
		assert.Equal(t, existingSeller.Address, updatedSellers[0].Address)              // Preserved field
		assert.Equal(t, existingSeller.Telephone, updatedSellers[0].Telephone)          // Preserved field
		assert.Equal(t, partialUpdateSeller.LocalityID, updatedSellers[0].LocalityID)   // Updated field

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when seller does not exist", func(t *testing.T) {
		sellerId := 999 // Non-existent seller ID
		partialUpdateSeller := models.Seller{
			CID:         "SEL-001-UPDATED",
			CompanyName: "Updated Company Name",
			LocalityID:  2,
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Mock SELECT: Return sql.ErrNoRows (seller not found)
		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(sellerId).
			WillReturnError(sql.ErrNoRows)

		// Note: No ExpectExec for UPDATE because it should not reach that point
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute Update method
		updatedSellers, err := repository.Update(sellerId, partialUpdateSeller)

		// Assert - Verify error and nil result
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)
		assert.Nil(t, updatedSellers)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when SELECT query fails", func(t *testing.T) {
		sellerId := 1
		partialUpdateSeller := models.Seller{
			CID:         "SEL-001-UPDATED",
			CompanyName: "Updated Company Name",
			LocalityID:  2,
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Mock SELECT: Return database error (not sql.ErrNoRows)
		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(sellerId).
			WillReturnError(error_message.ErrQuery)

		// Note: No ExpectExec for UPDATE because it should not reach that point
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute Update method
		updatedSellers, err := repository.Update(sellerId, partialUpdateSeller)

		// Assert - Verify error and nil result
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Nil(t, updatedSellers)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when UPDATE query fails", func(t *testing.T) {
		sellerId := 1

		// Existing seller data in database (before update)
		existingSeller := models.Seller{
			Id:          sellerId,
			CID:         "SEL-001-OLD",
			CompanyName: "Old Company Name",
			Address:     "Old Address 123",
			Telephone:   "555-0000",
			LocalityID:  1,
		}

		partialUpdateSeller := models.Seller{
			CID:         "SEL-001-UPDATED",
			CompanyName: "Updated Company Name",
			LocalityID:  2,
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Mock SELECT: Return existing seller successfully
		rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
			AddRow(existingSeller.Id, existingSeller.CID, existingSeller.CompanyName, existingSeller.Address, existingSeller.Telephone, existingSeller.LocalityID)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(sellerId).
			WillReturnRows(rows)

		// Mock UPDATE: Return database error
		mock.ExpectExec("UPDATE sellers SET cid = \\?, company_name = \\?, address = \\?, telephone = \\?, locality_id = \\? WHERE id = \\?").
			WithArgs(partialUpdateSeller.CID, partialUpdateSeller.CompanyName, existingSeller.Address, existingSeller.Telephone, partialUpdateSeller.LocalityID, sellerId).
			WillReturnError(error_message.ErrQuery)

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute Update method
		updatedSellers, err := repository.Update(sellerId, partialUpdateSeller)

		// Assert - Verify error and nil result
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Nil(t, updatedSellers)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestDelete(t *testing.T) {

	t.Run("delete ok", func(t *testing.T) {
		// Arrange - Define seller ID to delete
		sellerIdToDelete := 1

		// Arrange - Mock database setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange - Configure mock expectation: DELETE affects 1 row (successful deletion)
		mock.ExpectExec("DELETE FROM sellers WHERE id = \\?").
			WithArgs(sellerIdToDelete).
			WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

		// Arrange - Create repository instance
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute Delete method
		err = repository.Delete(sellerIdToDelete)

		// Assert - Verify successful deletion (no error)
		assert.NoError(t, err)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when seller does not exist", func(t *testing.T) {
		// Arrange - Define seller ID that doesn't exist
		nonExistentSellerId := 999

		// Arrange - Mock database setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange - Configure mock expectation: DELETE affects 0 rows (seller doesn't exist)
		mock.ExpectExec("DELETE FROM sellers WHERE id = \\?").
			WithArgs(nonExistentSellerId).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

		// Arrange - Create repository instance
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute Delete method
		err = repository.Delete(nonExistentSellerId)

		// Assert - Verify ErrNotFound error
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when DELETE query fails", func(t *testing.T) {
		// Arrange - Define seller ID
		sellerId := 1

		// Arrange - Mock database setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange - Configure mock expectation: DELETE query fails
		mock.ExpectExec("DELETE FROM sellers WHERE id = \\?").
			WithArgs(sellerId).
			WillReturnError(errors.New("database connection error"))

		// Arrange - Create repository instance
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute Delete method
		err = repository.Delete(sellerId)

		// Assert - Verify ErrQuery error
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when RowsAffected fails", func(t *testing.T) {
		// Arrange - Define seller ID
		sellerId := 1

		// Arrange - Mock database setup
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange - Configure mock expectation: DELETE succeeds but RowsAffected fails
		mock.ExpectExec("DELETE FROM sellers WHERE id = \\?").
			WithArgs(sellerId).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

		// Arrange - Create repository instance
		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		// Act - Execute Delete method
		err = repository.Delete(sellerId)

		// Assert - Verify ErrQuery error
		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)

		// Assert - Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}
