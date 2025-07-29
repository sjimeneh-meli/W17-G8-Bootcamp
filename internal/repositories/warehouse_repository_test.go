package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
)

var (
	warehouseTestFields = []string{
		"`id`", "`address`", "`telephone`", "`warehouse_code`",
		"`minimum_capacity`", "`minimum_temperature`", "`locality_id`",
	}
)

func setupWarehouseMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, models.Warehouse) {
	expectedWarehouse := models.Warehouse{
		Id:                 1,
		Address:            "Dirección 1",
		Telephone:          "123456789",
		WareHouseCode:      "WH001",
		MinimumCapacity:    100,
		MinimumTemperature: 5.0,
		LocalityId:         1,
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "Failed to create mock database")

	return db, mock, expectedWarehouse
}

func TestWarehouseRepository_GetById(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()

		rows := mock.NewRows(warehouseTestFields).AddRow(
			expectedWarehouse.Id, expectedWarehouse.Address, expectedWarehouse.Telephone,
			expectedWarehouse.WareHouseCode, expectedWarehouse.MinimumCapacity,
			expectedWarehouse.MinimumTemperature, expectedWarehouse.LocalityId,
		)

		mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `warehouse` WHERE `id` = ?", strings.Join(warehouseTestFields, ", "))).
			WithArgs(expectedWarehouse.Id).
			WillReturnRows(rows)

		repo := NewWarehouseRepository(db)

		result, err := repo.GetById(context.Background(), expectedWarehouse.Id)

		assert.NoError(t, err)
		assert.Equal(t, expectedWarehouse, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail not found", func(t *testing.T) {
		db, mock, _ := setupWarehouseMock(t)
		defer db.Close()
		mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `warehouse` WHERE `id` = ?", strings.Join(warehouseTestFields, ", "))).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(warehouseTestFields))
		repo := NewWarehouseRepository(db)
		result, err := repo.GetById(context.Background(), 1)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
func TestWarehouseRepository_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()
		rows := mock.NewRows(warehouseTestFields).AddRow(
			expectedWarehouse.Id, expectedWarehouse.Address, expectedWarehouse.Telephone,
			expectedWarehouse.WareHouseCode, expectedWarehouse.MinimumCapacity,
			expectedWarehouse.MinimumTemperature, expectedWarehouse.LocalityId,
		)
		mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `warehouse`", strings.Join(warehouseTestFields, ", "))).
			WillReturnRows(rows)
		repo := NewWarehouseRepository(db)
		result, err := repo.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []models.Warehouse{expectedWarehouse}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail internal server error", func(t *testing.T) {
		db, mock, _ := setupWarehouseMock(t)
		defer db.Close()
		mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `warehouse`", strings.Join(warehouseTestFields, ", "))).
			WillReturnError(sql.ErrConnDone)
		repo := NewWarehouseRepository(db)
		result, err := repo.GetAll(context.Background())
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail rows scan", func(t *testing.T) {
		db, mock, _ := setupWarehouseMock(t)
		defer db.Close()
		// Crear una fila con datos inválidos para que falle el scan
		rows := mock.NewRows(warehouseTestFields).AddRow(
			"invalid_id", "Dirección 1", "123456789",
			"WH001", "invalid_capacity", "invalid_temperature", "invalid_locality",
		)
		mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `warehouse`", strings.Join(warehouseTestFields, ", "))).
			WillReturnRows(rows)
		repo := NewWarehouseRepository(db)
		result, err := repo.GetAll(context.Background())
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestWarehouseRepository_ExistsByCode(t *testing.T) {
	t.Run("success - code exists", func(t *testing.T) {
		db, mock, _ := setupWarehouseMock(t)
		defer db.Close()

		rows := mock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `warehouse` WHERE `warehouse_code` = \\?").
			WithArgs("WH001").
			WillReturnRows(rows)

		repo := NewWarehouseRepository(db)
		exists, err := repo.ExistsByCode(context.Background(), "WH001")

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - code does not exist", func(t *testing.T) {
		db, mock, _ := setupWarehouseMock(t)
		defer db.Close()

		rows := mock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `warehouse` WHERE `warehouse_code` = \\?").
			WithArgs("WH999").
			WillReturnRows(rows)

		repo := NewWarehouseRepository(db)
		exists, err := repo.ExistsByCode(context.Background(), "WH999")

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Fail internal server error", func(t *testing.T) {
		db, mock, _ := setupWarehouseMock(t)
		defer db.Close()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `warehouse` WHERE `warehouse_code` = \\?").
			WithArgs("WH001").
			WillReturnError(sql.ErrConnDone)

		repo := NewWarehouseRepository(db)
		exists, err := repo.ExistsByCode(context.Background(), "WH001")

		assert.Error(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()

		// Configurar el mock para que devuelva el ID generado
		mock.ExpectExec("INSERT INTO `warehouse`\\(`address`, `telephone`, `warehouse_code`, `minimum_capacity`, `minimum_temperature`, `locality_id`\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)").
			WithArgs(expectedWarehouse.Address, expectedWarehouse.Telephone, expectedWarehouse.WareHouseCode, expectedWarehouse.MinimumCapacity, expectedWarehouse.MinimumTemperature, expectedWarehouse.LocalityId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewWarehouseRepository(db)
		result, err := repo.Create(context.Background(), expectedWarehouse)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, expectedWarehouse.Address, result.Address)
		assert.Equal(t, expectedWarehouse.Telephone, result.Telephone)
		assert.Equal(t, expectedWarehouse.WareHouseCode, result.WareHouseCode)
		assert.Equal(t, expectedWarehouse.MinimumCapacity, result.MinimumCapacity)
		assert.Equal(t, expectedWarehouse.MinimumTemperature, result.MinimumTemperature)
		assert.Equal(t, expectedWarehouse.LocalityId, result.LocalityId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail conflict warehouse code", func(t *testing.T) {
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()
		mock.ExpectExec("INSERT INTO `warehouse`\\(`address`, `telephone`, `warehouse_code`, `minimum_capacity`, `minimum_temperature`, `locality_id`\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)").
			WithArgs(expectedWarehouse.Address, expectedWarehouse.Telephone, expectedWarehouse.WareHouseCode, expectedWarehouse.MinimumCapacity, expectedWarehouse.MinimumTemperature, expectedWarehouse.LocalityId).
			WillReturnError(errors.New("warehouse code already exists"))
		repo := NewWarehouseRepository(db)
		result, err := repo.Create(context.Background(), expectedWarehouse)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail internal server error", func(t *testing.T) {
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()
		mock.ExpectExec("INSERT INTO `warehouse`\\(`address`, `telephone`, `warehouse_code`, `minimum_capacity`, `minimum_temperature`, `locality_id`\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)").
			WithArgs(expectedWarehouse.Address, expectedWarehouse.Telephone, expectedWarehouse.WareHouseCode, expectedWarehouse.MinimumCapacity, expectedWarehouse.MinimumTemperature, expectedWarehouse.LocalityId).
			WillReturnError(sql.ErrConnDone)
		repo := NewWarehouseRepository(db)
		result, err := repo.Create(context.Background(), expectedWarehouse)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail LastInsertId error", func(t *testing.T) {
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()

		// Mock que devuelve un resultado pero LastInsertId falla
		mock.ExpectExec("INSERT INTO `warehouse`\\(`address`, `telephone`, `warehouse_code`, `minimum_capacity`, `minimum_temperature`, `locality_id`\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)").
			WithArgs(expectedWarehouse.Address, expectedWarehouse.Telephone, expectedWarehouse.WareHouseCode, expectedWarehouse.MinimumCapacity, expectedWarehouse.MinimumTemperature, expectedWarehouse.LocalityId).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("LastInsertId failed")))

		repo := NewWarehouseRepository(db)
		result, err := repo.Create(context.Background(), expectedWarehouse)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
func TestWarehouseRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()
		mock.ExpectExec("UPDATE `warehouse` SET `address` = \\?, `telephone` = \\?, `warehouse_code` = \\?, `minimum_capacity` = \\?, `minimum_temperature` = \\? WHERE `id` = \\?").
			WithArgs(expectedWarehouse.Address, expectedWarehouse.Telephone, expectedWarehouse.WareHouseCode, expectedWarehouse.MinimumCapacity, expectedWarehouse.MinimumTemperature, expectedWarehouse.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		repo := NewWarehouseRepository(db)
		result, err := repo.Update(context.Background(), expectedWarehouse.Id, expectedWarehouse)
		assert.NoError(t, err)
		assert.Equal(t, expectedWarehouse, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail conflict warehouse code", func(t *testing.T) {
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()
		mock.ExpectExec("UPDATE `warehouse` SET `address` = \\?, `telephone` = \\?, `warehouse_code` = \\?, `minimum_capacity` = \\?, `minimum_temperature` = \\? WHERE `id` = \\?").
			WithArgs(expectedWarehouse.Address, expectedWarehouse.Telephone, expectedWarehouse.WareHouseCode, expectedWarehouse.MinimumCapacity, expectedWarehouse.MinimumTemperature, expectedWarehouse.Id).
			WillReturnError(errors.New("warehouse code already exists"))
		repo := NewWarehouseRepository(db)
		result, err := repo.Update(context.Background(), expectedWarehouse.Id, expectedWarehouse)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail not found - no rows affected", func(t *testing.T) {
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()

		// Mock que devuelve 0 filas afectadas
		mock.ExpectExec("UPDATE `warehouse` SET `address` = \\?, `telephone` = \\?, `warehouse_code` = \\?, `minimum_capacity` = \\?, `minimum_temperature` = \\? WHERE `id` = \\?").
			WithArgs(expectedWarehouse.Address, expectedWarehouse.Telephone, expectedWarehouse.WareHouseCode, expectedWarehouse.MinimumCapacity, expectedWarehouse.MinimumTemperature, expectedWarehouse.Id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repo := NewWarehouseRepository(db)
		result, err := repo.Update(context.Background(), expectedWarehouse.Id, expectedWarehouse)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail RowsAffected error", func(t *testing.T) {
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()

		// Mock que devuelve un resultado pero RowsAffected falla
		mock.ExpectExec("UPDATE `warehouse` SET `address` = \\?, `telephone` = \\?, `warehouse_code` = \\?, `minimum_capacity` = \\?, `minimum_temperature` = \\? WHERE `id` = \\?").
			WithArgs(expectedWarehouse.Address, expectedWarehouse.Telephone, expectedWarehouse.WareHouseCode, expectedWarehouse.MinimumCapacity, expectedWarehouse.MinimumTemperature, expectedWarehouse.Id).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("RowsAffected failed")))

		repo := NewWarehouseRepository(db)
		result, err := repo.Update(context.Background(), expectedWarehouse.Id, expectedWarehouse)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()
		mock.ExpectExec("DELETE FROM `warehouse` WHERE `id` = ?").
			WithArgs(expectedWarehouse.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		repo := NewWarehouseRepository(db)
		err := repo.Delete(context.Background(), expectedWarehouse.Id)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail not found", func(t *testing.T) {
		db, mock, _ := setupWarehouseMock(t)
		defer db.Close()
		mock.ExpectExec("DELETE FROM `warehouse` WHERE `id` = ?").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)
		repo := NewWarehouseRepository(db)
		err := repo.Delete(context.Background(), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail internal server error", func(t *testing.T) {
		db, mock, _ := setupWarehouseMock(t)
		defer db.Close()
		mock.ExpectExec("DELETE FROM `warehouse` WHERE `id` = ?").
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)
		repo := NewWarehouseRepository(db)
		err := repo.Delete(context.Background(), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Success - no rows affected", func(t *testing.T) {
		db, mock, _ := setupWarehouseMock(t)
		defer db.Close()

		// Mock que devuelve 0 filas afectadas (no es un error en Delete)
		mock.ExpectExec("DELETE FROM `warehouse` WHERE `id` = \\?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repo := NewWarehouseRepository(db)
		err := repo.Delete(context.Background(), 1)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
