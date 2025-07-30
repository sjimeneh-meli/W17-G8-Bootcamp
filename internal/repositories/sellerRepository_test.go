// Package repositories_test - Seller Repository Unit Tests / Tests Unitarios del Repositorio Seller
// Database layer testing for seller CRUD operations with MySQL, CID uniqueness and locality validation
// Testing de capa de datos para operaciones CRUD de vendedores con MySQL, unicidad de CID y validación de localidad
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

// TestPost - Tests for Save method with validation / Tests para método Save con validación
// Cases: successful creation, CID conflict, CID check error, locality check error, non-existent locality, insert error, LastInsertId error
// Casos: creación exitosa, conflicto de CID, error de verificación CID, error de verificación localidad, localidad inexistente, error de inserción, error de LastInsertId
func TestPost(t *testing.T) {
	t.Run("create ok", func(t *testing.T) {
		// Test: Valid seller data with unique CID and existing locality creates seller successfully
		// Test: Datos válidos de vendedor con CID único y localidad existente crea vendedor exitosamente
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
		// Test: Duplicate CID prevents seller creation / CID duplicado previene creación de vendedor

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repositories.ResetSellerRepositorySingleton()
		repo := repositories.NewSQLSellerRepository(db)

		inputSeller := models.Seller{
			CID:         "SEL-001",
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "555-1234",
			LocalityID:  1,
		}

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM sellers WHERE cid = \\?\\)").
			WithArgs(inputSeller.CID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))

		result, err := repo.Save(inputSeller)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, error_message.ErrAlreadyExists, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("should return error when CID checking query fails", func(t *testing.T) {
		// Test: Database error during CID existence check / Error de base de datos durante verificación de existencia de CID
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
			WillReturnError(error_message.ErrFailedCheckingExistence)

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		sellerCreated, err := repository.Save(inputSeller)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrFailedCheckingExistence, err)
		assert.Nil(t, sellerCreated)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when locality checking query fails", func(t *testing.T) {
		// Test: Database error during locality existence check / Error de base de datos durante verificación de existencia de localidad
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
			WillReturnError(error_message.ErrFailedCheckingExistence)

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		sellerCreated, err := repository.Save(inputSeller)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrFailedCheckingExistence, err)
		assert.Nil(t, sellerCreated)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when locality does not exist", func(t *testing.T) {
		// Test: Non-existent locality prevents seller creation / Localidad inexistente previene creación de vendedor
		inputSeller := models.Seller{
			CID:         "SEL-001",
			CompanyName: "Test Company",
			Address:     "Test Address 123",
			Telephone:   "555-1234",
			LocalityID:  999, // Non-existent locality / Localidad inexistente
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM sellers WHERE cid = \\?\\)").
			WithArgs("SEL-001").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(999).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		sellerCreated, err := repository.Save(inputSeller)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrDependencyNotFound, err)
		assert.Nil(t, sellerCreated)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when INSERT query fails", func(t *testing.T) {
		// Test: Database error during seller insertion / Error de base de datos durante inserción de vendedor
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
			WillReturnError(error_message.ErrQuery)

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		sellerCreated, err := repository.Save(inputSeller)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Nil(t, sellerCreated)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when LastInsertId fails", func(t *testing.T) {
		// Test: Error retrieving generated ID after successful insert / Error al recuperar ID generado después de inserción exitosa
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
			WillReturnResult(sqlmock.NewErrorResult(error_message.ErrQuery))

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		sellerCreated, err := repository.Save(inputSeller)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Nil(t, sellerCreated)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

// TestGet - Tests for GetAll method / Tests para método GetAll
// Cases: successful retrieval, empty result, database query error, row scanning error
// Casos: recuperación exitosa, resultado vacío, error de consulta de base de datos, error de escaneo de filas
func TestGet(t *testing.T) {
	t.Run("find all", func(t *testing.T) {
		// Test: Database with sellers returns complete list / Base de datos con vendedores retorna lista completa
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
		// Test: Empty database returns empty slice / Base de datos vacía retorna slice vacío
		expectedSellers := []models.Seller{}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"})

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers").
			WillReturnRows(rows)

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		sellersResult, err := repository.GetAll()

		assert.NoError(t, err)
		assert.Len(t, sellersResult, 0)
		assert.Equal(t, expectedSellers, sellersResult)
		assert.Empty(t, sellersResult)
		assert.NotNil(t, sellersResult)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when database query fails", func(t *testing.T) {
		// Test: Database connection error during query / Error de conexión de base de datos durante consulta

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
		// Test: Invalid data types cause row scanning errors / Tipos de datos inválidos causan errores de escaneo
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
			AddRow("invalid_id", "SEL-001", "Company One", "Address One", "555-0001", 1)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers").
			WillReturnRows(rows)

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		sellersResult, err := repository.GetAll()

		assert.Error(t, err)
		assert.Nil(t, sellersResult)
		assert.NotEmpty(t, err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

// TestUpdateSeller - Tests for Update method / Tests para método Update
// Cases: successful update, seller not found, SELECT query error, UPDATE query error
// Casos: actualización exitosa, vendedor no encontrado, error de consulta SELECT, error de consulta UPDATE
func TestUpdateSeller(t *testing.T) {
	t.Run("update ok", func(t *testing.T) {
		// Test: Valid partial update modifies existing seller successfully / Actualización parcial válida modifica vendedor existente exitosamente

		sellerId := 1

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
			// Address and Telephone not to update (partial update) / Address y Telephone no se actualizan (actualización parcial)
			LocalityID: 2,
		}

		expectedSeller := models.Seller{
			Id:          existingSeller.Id,
			CID:         partialUpdateSeller.CID,         //  Updated / Actualizado
			CompanyName: partialUpdateSeller.CompanyName, //  Updated / Actualizado
			Address:     existingSeller.Address,          //  No updated / No actualizado
			Telephone:   existingSeller.Telephone,        //  No updated / No actualizado
			LocalityID:  partialUpdateSeller.LocalityID,  //  Updated / Actualizado
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
			WillReturnResult(sqlmock.NewResult(0, 1)) // 0 = no LastInsertId needed, 1 = rows affected / 0 = no se necesita LastInsertId, 1 = filas afectadas

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		updatedSellers, err := repository.Update(sellerId, partialUpdateSeller)

		assert.NoError(t, err)
		assert.Len(t, updatedSellers, 1)
		assert.Equal(t, expectedSeller, updatedSellers[0])

		assert.Equal(t, partialUpdateSeller.CID, updatedSellers[0].CID)                 // Updated field / Campo actualizado
		assert.Equal(t, partialUpdateSeller.CompanyName, updatedSellers[0].CompanyName) // Updated field / Campo actualizado
		assert.Equal(t, existingSeller.Address, updatedSellers[0].Address)              // Preserved field / Campo preservado
		assert.Equal(t, existingSeller.Telephone, updatedSellers[0].Telephone)          // Preserved field / Campo preservado
		assert.Equal(t, partialUpdateSeller.LocalityID, updatedSellers[0].LocalityID)   // Updated field / Campo actualizado

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when seller does not exist", func(t *testing.T) {
		// Test: Update non-existent seller returns not found error / Actualizar vendedor inexistente retorna error not found
		sellerId := 999
		partialUpdateSeller := models.Seller{
			CID:         "SEL-001-UPDATED",
			CompanyName: "Updated Company Name",
			LocalityID:  2,
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(sellerId).
			WillReturnError(sql.ErrNoRows)

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		updatedSellers, err := repository.Update(sellerId, partialUpdateSeller)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)
		assert.Nil(t, updatedSellers)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when SELECT query fails", func(t *testing.T) {
		// Test: Database error during seller retrieval for update / Error de base de datos durante recuperación de vendedor para actualización
		sellerId := 1
		partialUpdateSeller := models.Seller{
			CID:         "SEL-001-UPDATED",
			CompanyName: "Updated Company Name",
			LocalityID:  2,
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(sellerId).
			WillReturnError(error_message.ErrQuery)

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		updatedSellers, err := repository.Update(sellerId, partialUpdateSeller)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Nil(t, updatedSellers)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when UPDATE query fails", func(t *testing.T) {
		// Test: Database error during seller update execution / Error de base de datos durante ejecución de actualización de vendedor
		sellerId := 1

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

		rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
			AddRow(existingSeller.Id, existingSeller.CID, existingSeller.CompanyName, existingSeller.Address, existingSeller.Telephone, existingSeller.LocalityID)

		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = \\?").
			WithArgs(sellerId).
			WillReturnRows(rows)

		mock.ExpectExec("UPDATE sellers SET cid = \\?, company_name = \\?, address = \\?, telephone = \\?, locality_id = \\? WHERE id = \\?").
			WithArgs(partialUpdateSeller.CID, partialUpdateSeller.CompanyName, existingSeller.Address, existingSeller.Telephone, partialUpdateSeller.LocalityID, sellerId).
			WillReturnError(error_message.ErrQuery)

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		updatedSellers, err := repository.Update(sellerId, partialUpdateSeller)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Nil(t, updatedSellers)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

// TestDelete - Tests for Delete method / Tests para método Delete
// Cases: successful deletion, seller not found, DELETE query error, RowsAffected error
// Casos: eliminación exitosa, vendedor no encontrado, error de consulta DELETE, error de RowsAffected
func TestDelete(t *testing.T) {

	t.Run("delete ok", func(t *testing.T) {
		// Test: Valid seller ID deletes seller successfully / ID de vendedor válido elimina vendedor exitosamente
		sellerIdToDelete := 1

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("DELETE FROM sellers WHERE id = \\?").
			WithArgs(sellerIdToDelete).
			WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected / 1 fila afectada

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		err = repository.Delete(sellerIdToDelete)

		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when seller does not exist", func(t *testing.T) {
		// Test: Delete non-existent seller returns not found error / Eliminar vendedor inexistente retorna error not found
		nonExistentSellerId := 999

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("DELETE FROM sellers WHERE id = \\?").
			WithArgs(nonExistentSellerId).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected / 0 filas afectadas

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		err = repository.Delete(nonExistentSellerId)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when DELETE query fails", func(t *testing.T) {
		// Test: Database connection error during delete operation / Error de conexión de base de datos durante operación de eliminación
		sellerId := 1

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("DELETE FROM sellers WHERE id = \\?").
			WithArgs(sellerId).
			WillReturnError(errors.New("database connection error"))

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		err = repository.Delete(sellerId)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when RowsAffected fails", func(t *testing.T) {
		// Test: Error checking affected rows after delete / Error al verificar filas afectadas después de eliminación
		sellerId := 1

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectExec("DELETE FROM sellers WHERE id = \\?").
			WithArgs(sellerId).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

		repositories.ResetSellerRepositorySingleton()
		repository := repositories.NewSQLSellerRepository(db)

		err = repository.Delete(sellerId)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

}
