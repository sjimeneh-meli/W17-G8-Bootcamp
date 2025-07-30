// Package repositories - Warehouse Repository Unit Tests / Tests Unitarios del Repositorio Warehouse
// Database layer testing for warehouse CRUD operations with MySQL and warehouse code uniqueness
// Testing de capa de datos para operaciones CRUD de almacenes con MySQL y unicidad de código de almacén
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

// setupWarehouseMock - Helper function to setup warehouse test dependencies
// setupWarehouseMock - Función helper para configurar dependencias de test de almacén
// Creates mock database, sqlmock instance and expected warehouse model
// Crea base de datos mock, instancia sqlmock y modelo de almacén esperado
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

// TestWarehouseRepository_GetById - Tests for GetById method / Tests para método GetById
// Cases: successful retrieval, warehouse not found
// Casos: recuperación exitosa, almacén no encontrado
func TestWarehouseRepository_GetById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Test: Valid warehouse ID returns complete warehouse data / ID de almacén válido retorna datos completos del almacén

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
		// Test: Non-existent warehouse ID returns not found error / ID de almacén inexistente retorna error not found
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

// TestWarehouseRepository_GetAll - Tests for GetAll method / Tests para método GetAll
// Cases: successful retrieval, internal server error, row scanning error
// Casos: recuperación exitosa, error interno del servidor, error de escaneo de filas
func TestWarehouseRepository_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Test: Database with warehouses returns complete list / Base de datos con almacenes retorna lista completa
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
		// Test: Database connection error during query / Error de conexión de base de datos durante consulta
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
		// Test: Invalid data types cause row scanning errors / Tipos de datos inválidos causan errores de escaneo
		db, mock, _ := setupWarehouseMock(t)
		defer db.Close()
		// Create row with invalid data to cause scan failure / Crear una fila con datos inválidos para que falle el scan
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

// TestWarehouseRepository_ExistsByCode - Tests for ExistsByCode method / Tests para método ExistsByCode
// Cases: code exists, code does not exist, internal server error
// Casos: código existe, código no existe, error interno del servidor
func TestWarehouseRepository_ExistsByCode(t *testing.T) {
	t.Run("success - code exists", func(t *testing.T) {
		// Test: Existing warehouse code returns true / Código de almacén existente retorna true
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
		// Test: Non-existent warehouse code returns false / Código de almacén inexistente retorna false
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
		// Test: Database connection error during existence check / Error de conexión durante verificación de existencia
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

// TestWarehouseRepository_Create - Tests for Create method / Tests para método Create
// Cases: successful creation, warehouse code conflict, internal server error, LastInsertId error
// Casos: creación exitosa, conflicto de código de almacén, error interno del servidor, error de LastInsertId
func TestWarehouseRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Test: Valid warehouse data creates new warehouse with generated ID / Datos válidos de almacén crean nuevo almacén con ID generado
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()

		// Configure mock to return generated ID / Configurar el mock para que devuelva el ID generado
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
		// Test: Duplicate warehouse code prevents creation / Código de almacén duplicado previene creación
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
		// Test: Database connection error during insert / Error de conexión de base de datos durante inserción
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
		// Test: Error retrieving generated ID after successful insert / Error al recuperar ID generado después de inserción exitosa
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()

		// Mock that returns result but LastInsertId fails / Mock que devuelve un resultado pero LastInsertId falla
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

// TestWarehouseRepository_Update - Tests for Update method / Tests para método Update
// Cases: successful update, warehouse code conflict, not found, RowsAffected error
// Casos: actualización exitosa, conflicto de código de almacén, no encontrado, error de RowsAffected
func TestWarehouseRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Test: Valid warehouse data updates existing warehouse / Datos válidos de almacén actualizan almacén existente
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
		// Test: Duplicate warehouse code prevents update / Código de almacén duplicado previene actualización
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
		// Test: Update non-existent warehouse returns error / Actualizar almacén inexistente retorna error
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()

		// Mock that returns 0 affected rows / Mock que devuelve 0 filas afectadas
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
		// Test: Error checking affected rows after update / Error al verificar filas afectadas después de actualización
		db, mock, expectedWarehouse := setupWarehouseMock(t)
		defer db.Close()

		// Mock that returns result but RowsAffected fails / Mock que devuelve un resultado pero RowsAffected falla
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

// TestWarehouseRepository_Delete - Tests for Delete method / Tests para método Delete
// Cases: successful deletion, not found, internal server error, no rows affected (success case)
// Casos: eliminación exitosa, no encontrado, error interno del servidor, sin filas afectadas (caso exitoso)
func TestWarehouseRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Test: Valid warehouse ID deletes warehouse successfully / ID de almacén válido elimina almacén exitosamente
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
		// Test: Database error during delete operation / Error de base de datos durante operación de eliminación
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
		// Test: Database connection error during delete / Error de conexión de base de datos durante eliminación
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
		// Test: Delete operation with no affected rows still succeeds / Operación de eliminación sin filas afectadas aún es exitosa
		db, mock, _ := setupWarehouseMock(t)
		defer db.Close()

		// Mock that returns 0 affected rows (not an error in Delete) / Mock que devuelve 0 filas afectadas (no es un error en Delete)
		mock.ExpectExec("DELETE FROM `warehouse` WHERE `id` = \\?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repo := NewWarehouseRepository(db)
		err := repo.Delete(context.Background(), 1)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
