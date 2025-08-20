// Package repositories_test - Locality Repository Unit Tests / Tests Unitarios del Repositorio Locality
// Database layer testing for locality CRUD operations with MySQL, country/province management and seller reports
// Testing de capa de datos para operaciones CRUD de localidades con MySQL, manejo de país/provincia y reportes de vendedores
package repositories_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/stretchr/testify/assert"
)

// TestSaveLocality - Tests for Save method with country/province management / Tests para método Save con manejo de país/provincia
// Cases: successful creation, country creation, province creation, locality conflict, various query errors
// Casos: creación exitosa, creación de país, creación de provincia, conflicto de localidad, varios errores de consulta
func TestSaveLocality(t *testing.T) {
	t.Run("create locality ok - existing country and province", func(t *testing.T) {
		// Test: Valid locality data with existing country and province creates locality successfully
		// Test: Datos válidos de localidad con país y provincia existentes crea localidad exitosamente
		inputLocality := models.Locality{
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		expectedLocality := models.Locality{
			Id:           1,
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Country exists
		mock.ExpectQuery("SELECT id FROM countries WHERE country_name = \\?").
			WithArgs("Test Country").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Province exists
		mock.ExpectQuery("SELECT id FROM provinces WHERE province_name = \\? AND id_country_fk = \\?").
			WithArgs("Test Province", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Locality doesn't exist - matching the exact multiline query format
		mock.ExpectQuery("SELECT EXISTS\\(\\s*SELECT 1 FROM localities WHERE locality_name = \\? AND province_id = \\?\\s*\\)").
			WithArgs("Test City", 1).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Insert locality
		mock.ExpectExec("INSERT INTO localities \\(locality_name, province_id\\) VALUES \\(\\?, \\?\\)").
			WithArgs("Test City", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		result, err := repository.Save(context.Background(), inputLocality)

		assert.NoError(t, err)
		assert.Equal(t, expectedLocality, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create locality ok - new country and province", func(t *testing.T) {
		// Test: Valid locality data with new country and province creates all entities successfully
		// Test: Datos válidos de localidad con país y provincia nuevos crea todas las entidades exitosamente
		inputLocality := models.Locality{
			LocalityName: "New City",
			ProvinceName: "New Province",
			CountryName:  "New Country",
		}

		expectedLocality := models.Locality{
			Id:           1,
			LocalityName: "New City",
			ProvinceName: "New Province",
			CountryName:  "New Country",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Country doesn't exist
		mock.ExpectQuery("SELECT id FROM countries WHERE country_name = \\?").
			WithArgs("New Country").
			WillReturnError(sql.ErrNoRows)

		// Insert new country
		mock.ExpectExec("INSERT INTO countries \\(country_name\\) VALUES \\(\\?\\)").
			WithArgs("New Country").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Province doesn't exist
		mock.ExpectQuery("SELECT id FROM provinces WHERE province_name = \\? AND id_country_fk = \\?").
			WithArgs("New Province", 1).
			WillReturnError(sql.ErrNoRows)

		// Insert new province
		mock.ExpectExec("INSERT INTO provinces \\(province_name, id_country_fk\\) VALUES \\(\\?, \\?\\)").
			WithArgs("New Province", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Locality doesn't exist - matching the exact multiline query format
		mock.ExpectQuery("SELECT EXISTS\\(\\s*SELECT 1 FROM localities WHERE locality_name = \\? AND province_id = \\?\\s*\\)").
			WithArgs("New City", 1).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		// Insert locality
		mock.ExpectExec("INSERT INTO localities \\(locality_name, province_id\\) VALUES \\(\\?, \\?\\)").
			WithArgs("New City", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		result, err := repository.Save(context.Background(), inputLocality)

		assert.NoError(t, err)
		assert.Equal(t, expectedLocality, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when locality already exists", func(t *testing.T) {
		// Test: Duplicate locality name and province prevents creation / Nombre de localidad y provincia duplicados previene creación
		inputLocality := models.Locality{
			LocalityName: "Existing City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Country exists
		mock.ExpectQuery("SELECT id FROM countries WHERE country_name = \\?").
			WithArgs("Test Country").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Province exists
		mock.ExpectQuery("SELECT id FROM provinces WHERE province_name = \\? AND id_country_fk = \\?").
			WithArgs("Test Province", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Locality already exists - matching the exact multiline query format
		mock.ExpectQuery("SELECT EXISTS\\(\\s*SELECT 1 FROM localities WHERE locality_name = \\? AND province_id = \\?\\s*\\)").
			WithArgs("Existing City", 1).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		result, err := repository.Save(context.Background(), inputLocality)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrAlreadyExists, err)
		assert.Equal(t, models.Locality{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when country query fails", func(t *testing.T) {
		// Test: Database error during country lookup / Error de base de datos durante búsqueda de país
		inputLocality := models.Locality{
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM countries WHERE country_name = \\?").
			WithArgs("Test Country").
			WillReturnError(error_message.ErrQuery)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		result, err := repository.Save(context.Background(), inputLocality)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Equal(t, models.Locality{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when country insert fails", func(t *testing.T) {
		// Test: Database error during new country creation / Error de base de datos durante creación de nuevo país
		inputLocality := models.Locality{
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "New Country",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM countries WHERE country_name = \\?").
			WithArgs("New Country").
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec("INSERT INTO countries \\(country_name\\) VALUES \\(\\?\\)").
			WithArgs("New Country").
			WillReturnError(error_message.ErrQuery)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		result, err := repository.Save(context.Background(), inputLocality)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Equal(t, models.Locality{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when province query fails", func(t *testing.T) {
		// Test: Database error during province lookup / Error de base de datos durante búsqueda de provincia
		inputLocality := models.Locality{
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM countries WHERE country_name = \\?").
			WithArgs("Test Country").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery("SELECT id FROM provinces WHERE province_name = \\? AND id_country_fk = \\?").
			WithArgs("Test Province", 1).
			WillReturnError(error_message.ErrQuery)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		result, err := repository.Save(context.Background(), inputLocality)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Equal(t, models.Locality{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when province insert fails", func(t *testing.T) {
		// Test: Database error during new province creation / Error de base de datos durante creación de nueva provincia
		inputLocality := models.Locality{
			LocalityName: "Test City",
			ProvinceName: "New Province",
			CountryName:  "Test Country",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM countries WHERE country_name = \\?").
			WithArgs("Test Country").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery("SELECT id FROM provinces WHERE province_name = \\? AND id_country_fk = \\?").
			WithArgs("New Province", 1).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec("INSERT INTO provinces \\(province_name, id_country_fk\\) VALUES \\(\\?, \\?\\)").
			WithArgs("New Province", 1).
			WillReturnError(error_message.ErrQuery)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		result, err := repository.Save(context.Background(), inputLocality)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Equal(t, models.Locality{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when locality existence check fails", func(t *testing.T) {
		// Test: Database error during locality existence check / Error de base de datos durante verificación de existencia de localidad
		inputLocality := models.Locality{
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM countries WHERE country_name = \\?").
			WithArgs("Test Country").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery("SELECT id FROM provinces WHERE province_name = \\? AND id_country_fk = \\?").
			WithArgs("Test Province", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery("SELECT EXISTS\\(\\s*SELECT 1 FROM localities WHERE locality_name = \\? AND province_id = \\?\\s*\\)").
			WithArgs("Test City", 1).
			WillReturnError(error_message.ErrQuery)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		result, err := repository.Save(context.Background(), inputLocality)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Equal(t, models.Locality{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when locality insert fails", func(t *testing.T) {
		// Test: Database error during locality insertion / Error de base de datos durante inserción de localidad
		inputLocality := models.Locality{
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM countries WHERE country_name = \\?").
			WithArgs("Test Country").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery("SELECT id FROM provinces WHERE province_name = \\? AND id_country_fk = \\?").
			WithArgs("Test Province", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery("SELECT EXISTS\\(\\s*SELECT 1 FROM localities WHERE locality_name = \\? AND province_id = \\?\\s*\\)").
			WithArgs("Test City", 1).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("INSERT INTO localities \\(locality_name, province_id\\) VALUES \\(\\?, \\?\\)").
			WithArgs("Test City", 1).
			WillReturnError(error_message.ErrQuery)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		result, err := repository.Save(context.Background(), inputLocality)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Equal(t, models.Locality{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when LastInsertId fails", func(t *testing.T) {
		// Test: Error retrieving generated ID after successful locality insert / Error al recuperar ID generado después de inserción exitosa de localidad
		inputLocality := models.Locality{
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT id FROM countries WHERE country_name = \\?").
			WithArgs("Test Country").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery("SELECT id FROM provinces WHERE province_name = \\? AND id_country_fk = \\?").
			WithArgs("Test Province", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery("SELECT EXISTS\\(\\s*SELECT 1 FROM localities WHERE locality_name = \\? AND province_id = \\?\\s*\\)").
			WithArgs("Test City", 1).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("INSERT INTO localities \\(locality_name, province_id\\) VALUES \\(\\?, \\?\\)").
			WithArgs("Test City", 1).
			WillReturnResult(sqlmock.NewErrorResult(error_message.ErrQuery))

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		result, err := repository.Save(context.Background(), inputLocality)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Equal(t, models.Locality{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestGetSellerReports - Tests for GetSellerReports method / Tests para método GetSellerReports
// Cases: all localities reports, specific locality report, locality not found, query errors
// Casos: reportes de todas las localidades, reporte de localidad específica, localidad no encontrada, errores de consulta
func TestGetSellerReports(t *testing.T) {
	t.Run("get all localities reports ok", func(t *testing.T) {
		// Test: Request for all localities returns complete seller reports / Solicitud de todas las localidades retorna reportes completos de vendedores
		expectedReports := []responses.LocalitySellerReport{
			{
				LocalityID:   1,
				LocalityName: "City One",
				SellerCount:  5,
			},
			{
				LocalityID:   2,
				LocalityName: "City Two",
				SellerCount:  3,
			},
			{
				LocalityID:   3,
				LocalityName: "City Three",
				SellerCount:  0, // Locality with no sellers / Localidad sin vendedores
			},
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "locality_name", "count"}).
			AddRow(1, "City One", 5).
			AddRow(2, "City Two", 3).
			AddRow(3, "City Three", 0)

		mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT\\(s.id\\) FROM localities l LEFT JOIN sellers s ON s.locality_id = l.id GROUP BY l.id, l.locality_name").
			WillReturnRows(rows)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		reports, err := repository.GetSellerReports(context.Background(), 0)

		assert.NoError(t, err)
		assert.Len(t, reports, 3)
		assert.Equal(t, expectedReports, reports)
		assert.Equal(t, 5, reports[0].SellerCount)
		assert.Equal(t, 0, reports[2].SellerCount) // Verify zero sellers case / Verificar caso de cero vendedores
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("get specific locality report ok", func(t *testing.T) {
		// Test: Request for specific locality returns its seller report / Solicitud de localidad específica retorna su reporte de vendedores
		localityID := 1
		expectedReports := []responses.LocalitySellerReport{
			{
				LocalityID:   1,
				LocalityName: "Specific City",
				SellerCount:  2,
			},
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Check locality exists
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(localityID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		// Get specific locality report
		rows := sqlmock.NewRows([]string{"id", "locality_name", "count"}).
			AddRow(1, "Specific City", 2)

		mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT\\(s.id\\) FROM localities l LEFT JOIN sellers s ON s.locality_id = l.id WHERE l.id = \\? GROUP BY l.id, l.locality_name").
			WithArgs(localityID).
			WillReturnRows(rows)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		reports, err := repository.GetSellerReports(context.Background(), localityID)

		assert.NoError(t, err)
		assert.Len(t, reports, 1)
		assert.Equal(t, expectedReports, reports)
		assert.Equal(t, localityID, reports[0].LocalityID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return empty slice when no localities exist", func(t *testing.T) {
		// Test: Database with no localities returns empty slice / Base de datos sin localidades retorna slice vacío
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "locality_name", "count"})

		mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT\\(s.id\\) FROM localities l LEFT JOIN sellers s ON s.locality_id = l.id GROUP BY l.id, l.locality_name").
			WillReturnRows(rows)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		reports, err := repository.GetSellerReports(context.Background(), 0)

		assert.NoError(t, err)
		assert.Len(t, reports, 0)
		assert.Empty(t, reports)
		// In Go, when you declare var slice []Type and never append to it, it remains nil
		// This is the actual behavior of the repository code
		assert.Nil(t, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when specific locality does not exist", func(t *testing.T) {
		// Test: Request for non-existent locality returns not found error / Solicitud de localidad inexistente retorna error not found
		nonExistentLocalityID := 999

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(nonExistentLocalityID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		reports, err := repository.GetSellerReports(context.Background(), nonExistentLocalityID)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)
		assert.Nil(t, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when locality existence check fails", func(t *testing.T) {
		// Test: Database error during locality existence verification / Error de base de datos durante verificación de existencia de localidad
		localityID := 1

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(localityID).
			WillReturnError(error_message.ErrFailedCheckingExistence)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		reports, err := repository.GetSellerReports(context.Background(), localityID)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrFailedCheckingExistence, err)
		assert.Nil(t, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when query fails for all localities", func(t *testing.T) {
		// Test: Database error during query for all localities / Error de base de datos durante consulta de todas las localidades
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT\\(s.id\\) FROM localities l LEFT JOIN sellers s ON s.locality_id = l.id GROUP BY l.id, l.locality_name").
			WillReturnError(error_message.ErrQueryingReport)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		reports, err := repository.GetSellerReports(context.Background(), 0)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQueryingReport, err)
		assert.Nil(t, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when query fails for specific locality", func(t *testing.T) {
		// Test: Database error during query for specific locality / Error de base de datos durante consulta de localidad específica
		localityID := 1

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(localityID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT\\(s.id\\) FROM localities l LEFT JOIN sellers s ON s.locality_id = l.id WHERE l.id = \\? GROUP BY l.id, l.locality_name").
			WithArgs(localityID).
			WillReturnError(error_message.ErrQueryingReport)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		reports, err := repository.GetSellerReports(context.Background(), localityID)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQueryingReport, err)
		assert.Nil(t, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when row scanning fails", func(t *testing.T) {
		// Test: Invalid data types cause row scanning errors / Tipos de datos inválidos causan errores de escaneo
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "locality_name", "count"}).
			AddRow("invalid_id", "Test City", 5) // Invalid ID type / Tipo de ID inválido

		mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT\\(s.id\\) FROM localities l LEFT JOIN sellers s ON s.locality_id = l.id GROUP BY l.id, l.locality_name").
			WillReturnRows(rows)

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		reports, err := repository.GetSellerReports(context.Background(), 0)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrFailedToScan, err)
		assert.Nil(t, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestExistById - Tests for ExistById method / Tests para método ExistById
// Cases: locality exists, locality doesn't exist, query error
// Casos: localidad existe, localidad no existe, error de consulta
func TestExistById(t *testing.T) {
	t.Run("locality exists", func(t *testing.T) {
		// Test: Existing locality ID returns true / ID de localidad existente retorna true
		localityID := 1

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(localityID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		exists, err := repository.ExistById(context.Background(), localityID)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("locality does not exist", func(t *testing.T) {
		// Test: Non-existent locality ID returns false / ID de localidad inexistente retorna false
		nonExistentLocalityID := 999

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(nonExistentLocalityID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		exists, err := repository.ExistById(context.Background(), nonExistentLocalityID)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when query fails", func(t *testing.T) {
		// Test: Database error during existence check / Error de base de datos durante verificación de existencia
		localityID := 1

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM localities WHERE id = \\?\\)").
			WithArgs(localityID).
			WillReturnError(errors.New("database connection error"))

		repositories.ResetLocalityRepositorySingleton()
		repository := repositories.NewSQLLocalityRepository(db)

		exists, err := repository.ExistById(context.Background(), localityID)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
