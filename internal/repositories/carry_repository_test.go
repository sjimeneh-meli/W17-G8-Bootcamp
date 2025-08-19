// Package repositories - Carry Repository Unit Tests / Tests Unitarios del Repositorio Carry
// Database layer testing for carry CRUD operations with MySQL and CID uniqueness
// Testing de capa de datos para operaciones CRUD de transportistas con MySQL y unicidad de CID
package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupCarryRepositoryTest(t *testing.T) (*CarryRepositoryImpl, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	repo := &CarryRepositoryImpl{db: db}

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestCarryRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, mock, cleanup := setupCarryRepositoryTest(t)
		defer cleanup()

		carry := models.Carry{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		mock.ExpectExec("INSERT INTO `carriers`").
			WithArgs(carry.Cid, carry.CompanyName, carry.Address, carry.Telephone, carry.LocalityId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := repo.Create(context.Background(), carry)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, carry.Cid, result.Cid)
		assert.Equal(t, carry.CompanyName, result.CompanyName)
		assert.Equal(t, carry.Address, result.Address)
		assert.Equal(t, carry.Telephone, result.Telephone)
		assert.Equal(t, carry.LocalityId, result.LocalityId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		repo, mock, cleanup := setupCarryRepositoryTest(t)
		defer cleanup()

		carry := models.Carry{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		mock.ExpectExec("INSERT INTO `carriers`").
			WithArgs(carry.Cid, carry.CompanyName, carry.Address, carry.Telephone, carry.LocalityId).
			WillReturnError(errors.New("database error"))

		result, err := repo.Create(context.Background(), carry)

		assert.Error(t, err)
		assert.Equal(t, models.Carry{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("last_insert_id_error", func(t *testing.T) {
		repo, mock, cleanup := setupCarryRepositoryTest(t)
		defer cleanup()

		carry := models.Carry{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		mock.ExpectExec("INSERT INTO `carriers`").
			WithArgs(carry.Cid, carry.CompanyName, carry.Address, carry.Telephone, carry.LocalityId).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))

		result, err := repo.Create(context.Background(), carry)

		assert.Error(t, err)
		assert.Equal(t, models.Carry{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCarryRepository_ExistsByCid(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		repo, mock, cleanup := setupCarryRepositoryTest(t)
		defer cleanup()

		cid := "CID123"
		expectedCount := 1

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `carriers`").
			WithArgs(cid).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		exists, err := repo.ExistsByCid(context.Background(), cid)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_exists", func(t *testing.T) {
		repo, mock, cleanup := setupCarryRepositoryTest(t)
		defer cleanup()

		cid := "CID123"
		expectedCount := 0

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `carriers`").
			WithArgs(cid).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		exists, err := repo.ExistsByCid(context.Background(), cid)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		repo, mock, cleanup := setupCarryRepositoryTest(t)
		defer cleanup()

		cid := "CID123"

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `carriers`").
			WithArgs(cid).
			WillReturnError(errors.New("database error"))

		exists, err := repo.ExistsByCid(context.Background(), cid)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCarryRepository_GetCarryReportsByLocality(t *testing.T) {
	t.Run("specific_locality_with_carriers", func(t *testing.T) {
		repo, mock, cleanup := setupCarryRepositoryTest(t)
		defer cleanup()

		localityID := 1
		expectedReports := []responses.LocalityCarryReport{
			{
				LocalityId:    1,
				LocalityName:  "Test Locality",
				CarriersCount: 2,
			},
		}

		rows := sqlmock.NewRows([]string{"id", "locality_name", "carriers_count"}).
			AddRow(expectedReports[0].LocalityId, expectedReports[0].LocalityName, expectedReports[0].CarriersCount)

		mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT").
			WithArgs(localityID).
			WillReturnRows(rows)

		reports, err := repo.GetCarryReportsByLocality(context.Background(), localityID)

		assert.NoError(t, err)
		assert.Equal(t, expectedReports, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("specific_locality_without_carriers", func(t *testing.T) {
		repo, mock, cleanup := setupCarryRepositoryTest(t)
		defer cleanup()

		localityID := 1
		expectedReports := []responses.LocalityCarryReport{
			{
				LocalityId:    1,
				LocalityName:  "Test Locality",
				CarriersCount: 0,
			},
		}

		rows := sqlmock.NewRows([]string{"id", "locality_name", "carriers_count"}).
			AddRow(expectedReports[0].LocalityId, expectedReports[0].LocalityName, expectedReports[0].CarriersCount)

		mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT").
			WithArgs(localityID).
			WillReturnRows(rows)

		reports, err := repo.GetCarryReportsByLocality(context.Background(), localityID)

		assert.NoError(t, err)
		assert.Equal(t, expectedReports, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("all_localities", func(t *testing.T) {
		repo, mock, cleanup := setupCarryRepositoryTest(t)
		defer cleanup()

		expectedReports := []responses.LocalityCarryReport{
			{
				LocalityId:    1,
				LocalityName:  "Test Locality 1",
				CarriersCount: 2,
			},
			{
				LocalityId:    2,
				LocalityName:  "Test Locality 2",
				CarriersCount: 0,
			},
		}

		rows := sqlmock.NewRows([]string{"id", "locality_name", "carriers_count"}).
			AddRow(expectedReports[0].LocalityId, expectedReports[0].LocalityName, expectedReports[0].CarriersCount).
			AddRow(expectedReports[1].LocalityId, expectedReports[1].LocalityName, expectedReports[1].CarriersCount)

		mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT").
			WillReturnRows(rows)

		reports, err := repo.GetCarryReportsByLocality(context.Background(), 0)

		assert.NoError(t, err)
		assert.Equal(t, expectedReports, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		repo, mock, cleanup := setupCarryRepositoryTest(t)
		defer cleanup()

		localityID := 1

		mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT").
			WithArgs(localityID).
			WillReturnError(errors.New("database error"))

		reports, err := repo.GetCarryReportsByLocality(context.Background(), localityID)

		assert.Error(t, err)
		assert.Nil(t, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan_error", func(t *testing.T) {
		repo, mock, cleanup := setupCarryRepositoryTest(t)
		defer cleanup()

		localityID := 1

		// Create rows with wrong data type to cause scan error
		rows := sqlmock.NewRows([]string{"id", "locality_name", "carriers_count"}).
			AddRow("invalid_id", "Test Locality", "invalid_count")

		mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT").
			WithArgs(localityID).
			WillReturnRows(rows)

		reports, err := repo.GetCarryReportsByLocality(context.Background(), localityID)

		assert.Error(t, err)
		assert.Nil(t, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows_error", func(t *testing.T) {
		repo, mock, cleanup := setupCarryRepositoryTest(t)
		defer cleanup()

		localityID := 1

		rows := sqlmock.NewRows([]string{"id", "locality_name", "carriers_count"}).
			AddRow(1, "Test Locality", 2).
			RowError(0, errors.New("row error"))

		mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT").
			WithArgs(localityID).
			WillReturnRows(rows)

		reports, err := repo.GetCarryReportsByLocality(context.Background(), localityID)

		assert.Error(t, err)
		assert.Nil(t, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
