// Package repositories_test - InboundOrder Repository Unit Tests / Tests Unitarios del Repositorio InboundOrder
// Database layer testing for inbound order CRUD operations with MySQL and order number uniqueness
// Testing de capa de datos para operaciones CRUD de órdenes de entrada con MySQL y unicidad de número de orden
package repositories_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/stretchr/testify/assert"
)

// TestMySqlInboundOrderRepository_GetAllInboundOrdersReports - Tests for GetAllInboundOrdersReports method
// TestMySqlInboundOrderRepository_GetAllInboundOrdersReports - Tests para método GetAllInboundOrdersReports
// Cases: successful retrieval with data, empty result, database query error, row scanning error
// Casos: recuperación exitosa con datos, resultado vacío, error de consulta de base de datos, error de escaneo de filas
func TestMySqlInboundOrderRepository_GetAllInboundOrdersReports(t *testing.T) {
	t.Run("Successfully returns filled reports slice when there is data on db response", func(t *testing.T) {
		// Test: Database with inbound orders returns complete reports slice
		// Test: Base de datos con órdenes de entrada retorna slice completo de reportes
		repositories.ResetInboundOrderRepositoryInstance()

		expectedReports := []models.InboundOrderReport{
			{Id: 1, IdCardNumber: "CARD-001", FirstName: "Carlos", LastName: "Ruiz", InboundOrderCount: 5},
			{Id: 2, IdCardNumber: "CARD-002", FirstName: "Maria", LastName: "Lopez", InboundOrderCount: 3},
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "inbound_orders_count"}).
			AddRow(1, "CARD-001", "Carlos", "Ruiz", 5).
			AddRow(2, "CARD-002", "Maria", "Lopez", 3)

		mock.ExpectQuery(`SELECT e.id, e.id_card_number, e.first_name, e.last_name, COUNT\(io.id\) AS inbound_orders_count`).
			WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.GetAllInboundOrdersReports(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedReports, result)
		mock.ExpectationsWereMet()
	})

	t.Run("Successfully returns empty reports slice when there isn't data on db response", func(t *testing.T) {
		// Test: Empty database returns empty slice
		// Test: Base de datos vacía retorna slice vacío
		repositories.ResetInboundOrderRepositoryInstance()

		expectedReports := []models.InboundOrderReport{}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "inbound_orders_count"})

		mock.ExpectQuery(`SELECT e.id, e.id_card_number, e.first_name, e.last_name, COUNT\(io.id\) AS inbound_orders_count`).
			WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.GetAllInboundOrdersReports(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedReports, result)
		mock.ExpectationsWereMet()
	})

	t.Run("Fails because of an internal server error querying the database", func(t *testing.T) {
		// Test: Database connection error during query
		// Test: Error de conexión de base de datos durante consulta
		repositories.ResetInboundOrderRepositoryInstance()

		expectedError := error_message.ErrInternalServerError

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery(`SELECT e.id, e.id_card_number, e.first_name, e.last_name, COUNT\(io.id\) AS inbound_orders_count`).
			WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.GetAllInboundOrdersReports(context.Background())

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.ErrorIs(t, err, expectedError)
		mock.ExpectationsWereMet()
	})

	t.Run("Fails because of row scanning error", func(t *testing.T) {
		// Test: Row scanning error during result processing
		// Test: Error de escaneo de filas durante procesamiento de resultados
		repositories.ResetInboundOrderRepositoryInstance()

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		// Create rows with wrong data types to cause scanning error
		// Crear filas con tipos de datos incorrectos para causar error de escaneo
		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "inbound_orders_count"}).
			AddRow("invalid_id", "CARD-001", "Carlos", "Ruiz", 5) // Invalid ID type

		mock.ExpectQuery(`SELECT e.id, e.id_card_number, e.first_name, e.last_name, COUNT\(io.id\) AS inbound_orders_count`).
			WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.GetAllInboundOrdersReports(context.Background())

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.ErrorIs(t, err, error_message.ErrInternalServerError)
		mock.ExpectationsWereMet()
	})
}

// TestMySqlInboundOrderRepository_GetInboundOrdersReportByEmployeeId - Tests for GetInboundOrdersReportByEmployeeId method
// TestMySqlInboundOrderRepository_GetInboundOrdersReportByEmployeeId - Tests para método GetInboundOrdersReportByEmployeeId
// Cases: existing employee, non-existent employee, database query error
// Casos: empleado existente, empleado inexistente, error de consulta de base de datos
func TestMySqlInboundOrderRepository_GetInboundOrdersReportByEmployeeId(t *testing.T) {
	t.Run("Successfully returns report when employee exists and has inbound orders", func(t *testing.T) {
		// Test: Employee exists and has orders, returns report
		// Test: Empleado existe y tiene órdenes, retorna reporte
		repositories.ResetInboundOrderRepositoryInstance()

		expectedReport := models.InboundOrderReport{
			Id:                1,
			IdCardNumber:      "CARD-001",
			FirstName:         "Carlos",
			LastName:          "Ruiz",
			InboundOrderCount: 5,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "inbound_orders_count"}).
			AddRow(1, "CARD-001", "Carlos", "Ruiz", 5)

		mock.ExpectQuery(`SELECT e.id, e.id_card_number, e.first_name, e.last_name, COUNT\(io.id\) AS inbound_orders_count`).
			WithArgs(1).WillReturnRows(rows)

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.GetInboundOrdersReportByEmployeeId(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedReport, result)
		mock.ExpectationsWereMet()
	})

	t.Run("Fails when employee doesn't exist", func(t *testing.T) {
		// Test: Employee doesn't exist, returns not found error
		// Test: Empleado no existe, retorna error not found
		repositories.ResetInboundOrderRepositoryInstance()

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery(`SELECT e.id, e.id_card_number, e.first_name, e.last_name, COUNT\(io.id\) AS inbound_orders_count`).
			WithArgs(999).WillReturnError(sql.ErrNoRows)

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.GetInboundOrdersReportByEmployeeId(context.Background(), 999)

		assert.Error(t, err)
		assert.Equal(t, models.InboundOrderReport{}, result)
		assert.ErrorIs(t, err, error_message.ErrNotFound)
		mock.ExpectationsWereMet()
	})

	t.Run("Fails because of database query error", func(t *testing.T) {
		// Test: Database connection error during query
		// Test: Error de conexión de base de datos durante consulta
		repositories.ResetInboundOrderRepositoryInstance()

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery(`SELECT e.id, e.id_card_number, e.first_name, e.last_name, COUNT\(io.id\) AS inbound_orders_count`).
			WithArgs(1).WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.GetInboundOrdersReportByEmployeeId(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, models.InboundOrderReport{}, result)
		assert.ErrorIs(t, err, error_message.ErrInternalServerError)
		mock.ExpectationsWereMet()
	})
}

// TestMySqlInboundOrderRepository_Create - Tests for Create method
// TestMySqlInboundOrderRepository_Create - Tests para método Create
// Cases: successful creation, database execution error, last insert ID error
// Casos: creación exitosa, error de ejecución de base de datos, error de último ID insertado
func TestMySqlInboundOrderRepository_Create(t *testing.T) {
	t.Run("Successfully creates inbound order and returns it with generated ID", func(t *testing.T) {
		// Test: Valid inbound order data creates record successfully
		// Test: Datos válidos de orden de entrada crean registro exitosamente
		repositories.ResetInboundOrderRepositoryInstance()

		orderDate := time.Now()
		inputOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		expectedOrder := models.InboundOrder{
			Id:             1,
			OrderDate:      orderDate,
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec(`INSERT INTO inbound_orders`).
			WithArgs(orderDate, "ORD-001", 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.Create(context.Background(), inputOrder)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrder, result)
		mock.ExpectationsWereMet()
	})

	t.Run("Fails because of database execution error", func(t *testing.T) {
		// Test: Database execution error during insert
		// Test: Error de ejecución de base de datos durante inserción
		repositories.ResetInboundOrderRepositoryInstance()

		orderDate := time.Now()
		inputOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec(`INSERT INTO inbound_orders`).
			WithArgs(orderDate, "ORD-001", 1, 1, 1).
			WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.Create(context.Background(), inputOrder)

		assert.Error(t, err)
		assert.Equal(t, models.InboundOrder{}, result)
		assert.ErrorIs(t, err, error_message.ErrInternalServerError)
		mock.ExpectationsWereMet()
	})

	t.Run("Fails because of last insert ID error", func(t *testing.T) {
		// Test: Error getting last insert ID after successful execution
		// Test: Error obteniendo último ID insertado después de ejecución exitosa
		repositories.ResetInboundOrderRepositoryInstance()

		orderDate := time.Now()
		inputOrder := models.InboundOrder{
			OrderDate:      orderDate,
			OrderNumber:    "ORD-001",
			EmployeeId:     1,
			ProductBatchId: 1,
			WarehouseId:    1,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec(`INSERT INTO inbound_orders`).
			WithArgs(orderDate, "ORD-001", 1, 1, 1).
			WillReturnResult(sqlmock.NewErrorResult(error_message.ErrInternalServerError))

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.Create(context.Background(), inputOrder)

		assert.Error(t, err)
		assert.Equal(t, models.InboundOrder{}, result)
		assert.ErrorIs(t, err, error_message.ErrInternalServerError)
		mock.ExpectationsWereMet()
	})
}

// TestMySqlInboundOrderRepository_ExistsByOrderNumber - Tests for ExistsByOrderNumber method
// TestMySqlInboundOrderRepository_ExistsByOrderNumber - Tests para método ExistsByOrderNumber
// Cases: order number exists, order number doesn't exist, database query error
// Casos: número de orden existe, número de orden no existe, error de consulta de base de datos
func TestMySqlInboundOrderRepository_ExistsByOrderNumber(t *testing.T) {
	t.Run("Returns true when order number exists", func(t *testing.T) {
		// Test: Order number exists in database, returns true
		// Test: Número de orden existe en base de datos, retorna true
		repositories.ResetInboundOrderRepositoryInstance()

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"exists"}).AddRow(true)

		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM inbound_orders WHERE order_number = \?\)`).
			WithArgs("ORD-001").WillReturnRows(rows)

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.ExistsByOrderNumber(context.Background(), "ORD-001")

		assert.NoError(t, err)
		assert.True(t, result)
		mock.ExpectationsWereMet()
	})

	t.Run("Returns false when order number doesn't exist", func(t *testing.T) {
		// Test: Order number doesn't exist in database, returns false
		// Test: Número de orden no existe en base de datos, retorna false
		repositories.ResetInboundOrderRepositoryInstance()

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"exists"}).AddRow(false)

		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM inbound_orders WHERE order_number = \?\)`).
			WithArgs("ORD-999").WillReturnRows(rows)

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.ExistsByOrderNumber(context.Background(), "ORD-999")

		assert.NoError(t, err)
		assert.False(t, result)
		mock.ExpectationsWereMet()
	})

	t.Run("Fails because of database query error", func(t *testing.T) {
		// Test: Database connection error during query
		// Test: Error de conexión de base de datos durante consulta
		repositories.ResetInboundOrderRepositoryInstance()

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM inbound_orders WHERE order_number = \?\)`).
			WithArgs("ORD-001").WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.GetNewInboundOrderMySQLRepository(db)

		result, err := repository.ExistsByOrderNumber(context.Background(), "ORD-001")

		assert.Error(t, err)
		assert.False(t, result)
		assert.ErrorIs(t, err, error_message.ErrInternalServerError)
		mock.ExpectationsWereMet()
	})
}

// TestGetNewInboundOrderMySQLRepository_Singleton - Test for singleton pattern to improve coverage
// TestGetNewInboundOrderMySQLRepository_Singleton - Test del patrón singleton para mejorar coverage
// Validates that the repository follows singleton pattern correctly
// Valida que el repositorio sigue el patrón singleton correctamente
func TestGetNewInboundOrderMySQLRepository_Singleton(t *testing.T) {
	t.Run("GetNewInboundOrderMySQLRepository_returns_same_instance_when_already_exists", func(t *testing.T) {
		// Test: Singleton pattern ensures same instance is returned
		// Test: Patrón singleton asegura que se retorna la misma instancia

		// Reset singleton instance to ensure clean test
		// Resetea instancia singleton para asegurar test limpio
		inboundOrderRepositoryInstance := &repositories.MySqlInboundOrderRepository{}

		db1, _, err1 := sqlmock.New()
		if err1 != nil {
			fmt.Println("failed to open sqlmock database:", err1)
		}
		defer db1.Close()

		db2, _, err2 := sqlmock.New()
		if err2 != nil {
			fmt.Println("failed to open sqlmock database:", err2)
		}
		defer db2.Close()

		// First call - creates new instance / Primera llamada - crea nueva instancia
		repository1 := repositories.GetNewInboundOrderMySQLRepository(db1)

		// Second call - should return same instance (ignores new database)
		// Segunda llamada - debe retornar la misma instancia (ignora nueva base de datos)
		repository2 := repositories.GetNewInboundOrderMySQLRepository(db2)

		// Verify they are the same instance / Verifica que son la misma instancia
		assert.Equal(t, repository1, repository2)

		// Clean up singleton for next tests / Limpia singleton para próximos tests
		_ = inboundOrderRepositoryInstance
	})
}
