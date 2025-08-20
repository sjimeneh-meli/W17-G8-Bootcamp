// Package repositories_test - Employee Repository Unit Tests / Tests Unitarios del Repositorio Employee
// Database layer testing for employee CRUD operations with MySQL and card number ID uniqueness
// Testing de capa de datos para operaciones CRUD de empleados con MySQL y unicidad de ID de tarjeta
package repositories_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/stretchr/testify/assert"
)

// TestMySqlEmployeeRepositoryGetAll - Tests for GetAll method / Tests para método GetAll
// Cases: successful retrieval with data, empty result, database query error, row scanning error
// Casos: recuperación exitosa con datos, resultado vacío, error de consulta de base de datos, error de escaneo de filas
func TestMySqlEmployeeRepositoryGetAll(t *testing.T) {
	t.Run("Successfully returns filled employees map when there is data on db response", func(t *testing.T) {
		// Test: Database with employees returns complete map / Base de datos con empleados retorna mapa completo
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedEmployees := map[int]models.Employee{
			1: {Id: 1, CardNumberID: "CARD-001", FirstName: "Pedro", LastName: "Martinez", WarehouseID: 1},
			2: {Id: 2, CardNumberID: "CARD-002", FirstName: "Sofia", LastName: "Rodriguez", WarehouseID: 2},
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
			AddRow(1, "CARD-001", "Pedro", "Martinez", 1).
			AddRow(2, "CARD-002", "Sofia", "Rodriguez", 2)

		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees").WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployees, err := repository.GetAll(context.Background())

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedEmployees, dbEmployees, "dbEmployees should be a map with two employees")
	})

	t.Run("Successfully returns empty employees map when there isn't data on db response", func(t *testing.T) {
		// Test: Empty database returns empty map / Base de datos vacía retorna mapa vacío
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedEmployees := map[int]models.Employee{}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"})

		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees").WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployees, err := repository.GetAll(context.Background())

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedEmployees, dbEmployees, "dbEmployees should be an empty map of employees")
	})

	t.Run("Fails because of an internal server error querying the database", func(t *testing.T) {
		// Test: Database connection error during query / Error de conexión de base de datos durante consulta
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedError := error_message.ErrInternalServerError
		expectedEmployees := map[int]models.Employee{}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees").
			WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployees, err := repository.GetAll(context.Background())

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedEmployees, dbEmployees, "dbEmployees should be an empty map of employees")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})

	t.Run("Fails because of an internal server error scanning database results", func(t *testing.T) {
		// Test: Invalid data types cause row scanning errors / Tipos de datos inválidos causan errores de escaneo
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedEmployees := map[int]models.Employee{
			1: {Id: 1, CardNumberID: "CARD-001", FirstName: "Pedro", LastName: "Martinez", WarehouseID: 1},
		}
		expectedError := error_message.ErrInternalServerError
		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
			AddRow(1, "CARD-001", "Pedro", "Martinez", 1).
			AddRow("STRING", "CARD-002", "Sofia", "Rodriguez", 2)

		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees").WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployees, err := repository.GetAll(context.Background())

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedEmployees, dbEmployees, "dbEmployees should be an empty map of employees")
		assert.ErrorIs(t, err, expectedError)
	})
}

// TestMySqlEmployeeRepositoryGetById - Tests for GetById method / Tests para método GetById
// Cases: successful retrieval, employee not found, internal server error
// Casos: recuperación exitosa, empleado no encontrado, error interno del servidor
func TestMySqlEmployeeRepositoryGetById(t *testing.T) {
	t.Run("Successfully returns employee when found in database", func(t *testing.T) {
		// Test: Valid employee ID returns complete employee data / ID de empleado válido retorna datos completos del empleado
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedEmployee := models.Employee{
			Id:           1,
			CardNumberID: "CARD-001",
			FirstName:    "Pedro",
			LastName:     "Martinez",
			WarehouseID:  1,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
			AddRow(1, "CARD-001", "Pedro", "Martinez", 1)

		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(1).WillReturnRows(rows)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.GetById(context.Background(), 1)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should match expected employee")
	})

	t.Run("Fails because employee doesn't exist in database", func(t *testing.T) {
		// Test: Non-existent employee ID returns not found error / ID de empleado inexistente retorna error not found
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedEmployee := models.Employee{}
		expectedError := error_message.ErrNotFound

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(1).WillReturnError(error_message.ErrNotFound)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.GetById(context.Background(), 1)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should be empty")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrNotFound")
	})

	t.Run("Fails because of an internal server error querying the database", func(t *testing.T) {
		// Test: Database connection error during retrieval / Error de conexión de base de datos durante recuperación
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedEmployee := models.Employee{}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(1).WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.GetById(context.Background(), 1)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should be empty")
	})
}

// TestMySqlEmployeeRepositoryDeleteById - Tests for DeleteById method / Tests para método DeleteById
// Cases: successful deletion, employee not found, internal server error
// Casos: eliminación exitosa, empleado no encontrado, error interno del servidor
func TestMySqlEmployeeRepositoryDeleteById(t *testing.T) {
	t.Run("Successfully deletes employee when exists in database", func(t *testing.T) {
		// Test: Valid employee ID deletes employee successfully / ID de empleado válido elimina empleado exitosamente
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("DELETE FROM employees WHERE id = ?").
			WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		err = repository.DeleteById(context.Background(), 1)

		assert.Nil(t, err, "err should be nil")
	})

	t.Run("Fails because employee doesn't exist in database", func(t *testing.T) {
		// Test: Delete non-existent employee returns not found error / Eliminar empleado inexistente retorna error not found
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedError := error_message.ErrNotFound

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("DELETE FROM employees WHERE id = ?").
			WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		err = repository.DeleteById(context.Background(), 1)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrNotFound")
	})

	t.Run("Fails because of an internal server error executing delete", func(t *testing.T) {
		// Test: Database connection error during delete operation / Error de conexión de base de datos durante operación de eliminación
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedError := error_message.ErrInternalServerError

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("DELETE FROM employees WHERE id = ?").
			WithArgs(1).WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		err = repository.DeleteById(context.Background(), 1)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})
}

// TestMySqlEmployeeRepositoryCreate - Tests for Create method / Tests para método Create
// Cases: successful creation, internal server error executing insert, internal server error getting LastInsertId
// Casos: creación exitosa, error interno del servidor ejecutando inserción, error interno del servidor obteniendo LastInsertId
func TestMySqlEmployeeRepositoryCreate(t *testing.T) {
	t.Run("Successfully creates employee in database", func(t *testing.T) {
		// Test: Valid employee data creates new employee with generated ID / Datos válidos de empleado crean nuevo empleado con ID generado
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		inputEmployee := models.Employee{
			CardNumberID: "CARD-001",
			FirstName:    "Pedro",
			LastName:     "Martinez",
			WarehouseID:  1,
		}
		expectedEmployee := models.Employee{
			Id:           1,
			CardNumberID: "CARD-001",
			FirstName:    "Pedro",
			LastName:     "Martinez",
			WarehouseID:  1,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("INSERT INTO employees \\(id_card_number, first_name, last_name, warehouse_id\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
			WithArgs("CARD-001", "Pedro", "Martinez", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.Create(context.Background(), inputEmployee)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should match expected employee")
	})

	t.Run("Fails because of an internal server error executing insert", func(t *testing.T) {
		// Test: Database connection error during employee insertion / Error de conexión de base de datos durante inserción de empleado
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		inputEmployee := models.Employee{
			CardNumberID: "CARD-001",
			FirstName:    "Pedro",
			LastName:     "Martinez",
			WarehouseID:  1,
		}
		expectedEmployee := models.Employee{}
		expectedError := error_message.ErrInternalServerError

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("INSERT INTO employees \\(id_card_number, first_name, last_name, warehouse_id\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
			WithArgs("CARD-001", "Pedro", "Martinez", 1).
			WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.Create(context.Background(), inputEmployee)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should be empty")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})

	t.Run("Fails because of an internal server error getting last insert id", func(t *testing.T) {
		// Test: Error retrieving generated ID after successful insert / Error al recuperar ID generado después de inserción exitosa
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		inputEmployee := models.Employee{
			CardNumberID: "CARD-001",
			FirstName:    "Pedro",
			LastName:     "Martinez",
			WarehouseID:  1,
		}
		expectedEmployee := models.Employee{}
		expectedError := error_message.ErrInternalServerError

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("INSERT INTO employees \\(id_card_number, first_name, last_name, warehouse_id\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
			WithArgs("CARD-001", "Pedro", "Martinez", 1).
			WillReturnResult(sqlmock.NewErrorResult(error_message.ErrInternalServerError))

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.Create(context.Background(), inputEmployee)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should be empty")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})
}

// TestMySqlEmployeeRepositoryUpdate - Tests for Update method / Tests para método Update
// Cases: successful update, employee not found, internal server error executing update
// Casos: actualización exitosa, empleado no encontrado, error interno del servidor ejecutando actualización
func TestMySqlEmployeeRepositoryUpdate(t *testing.T) {
	t.Run("Successfully updates employee when exists in database", func(t *testing.T) {
		// Test: Valid employee data updates existing employee / Datos válidos de empleado actualizan empleado existente
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		inputEmployee := models.Employee{
			FirstName: "UpdatedName",
			LastName:  "UpdatedLastName",
		}
		expectedEmployee := models.Employee{
			Id:           1,
			CardNumberID: "CARD-001",
			FirstName:    "UpdatedName",
			LastName:     "UpdatedLastName",
			WarehouseID:  1,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		// Mock the UPDATE query / Mock de la consulta UPDATE
		mock.ExpectExec("UPDATE employees SET first_name = \\?, last_name = \\? WHERE id = \\?").
			WithArgs("UpdatedName", "UpdatedLastName", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Mock the GetById query after update / Mock de la consulta GetById después de actualización
		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
			AddRow(1, "CARD-001", "UpdatedName", "UpdatedLastName", 1)
		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(1).WillReturnRows(rows)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.Update(context.Background(), 1, inputEmployee)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should match expected employee")
	})

	t.Run("Fails because employee doesn't exist in database", func(t *testing.T) {
		// Test: Update non-existent employee returns not found error / Actualizar empleado inexistente retorna error not found
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		inputEmployee := models.Employee{
			FirstName: "UpdatedName",
		}
		expectedEmployee := models.Employee{}
		expectedError := error_message.ErrNotFound

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("UPDATE employees SET first_name = \\? WHERE id = \\?").
			WithArgs("UpdatedName", 1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.Update(context.Background(), 1, inputEmployee)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should be empty")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrNotFound")
	})

	t.Run("Fails because of an internal server error executing update", func(t *testing.T) {
		// Test: Database connection error during employee update / Error de conexión de base de datos durante actualización de empleado
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		inputEmployee := models.Employee{
			FirstName: "UpdatedName",
		}
		expectedEmployee := models.Employee{}
		expectedError := error_message.ErrInternalServerError

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("UPDATE employees SET first_name = \\? WHERE id = \\?").
			WithArgs("UpdatedName", 1).
			WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.Update(context.Background(), 1, inputEmployee)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should be empty")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})
}

// TestMySqlEmployeeRepositoryGetCardNumberIds - Tests for GetCardNumberIds method / Tests para método GetCardNumberIds
// Cases: successful retrieval with data, empty database, internal server error
// Casos: recuperación exitosa con datos, base de datos vacía, error interno del servidor
func TestMySqlEmployeeRepositoryGetCardNumberIds(t *testing.T) {
	t.Run("Successfully returns all card number ids when there is data in database", func(t *testing.T) {
		// Test: Database with card numbers returns complete list / Base de datos con números de tarjeta retorna lista completa
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedCardNumberIds := []string{"CARD-001", "CARD-002", "CARD-003"}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id_card_number"}).
			AddRow("CARD-001").
			AddRow("CARD-002").
			AddRow("CARD-003")

		mock.ExpectQuery("SELECT id_card_number FROM employees").WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		cardNumberIds, err := repository.GetCardNumberIds()

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedCardNumberIds, cardNumberIds, "cardNumberIds should match expected")
	})

	t.Run("Successfully returns empty slice when there is no data in database", func(t *testing.T) {
		// Test: Empty database returns empty slice / Base de datos vacía retorna slice vacío
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedCardNumberIds := []string{}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id_card_number"})

		mock.ExpectQuery("SELECT id_card_number FROM employees").WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		cardNumberIds, err := repository.GetCardNumberIds()

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedCardNumberIds, cardNumberIds, "cardNumberIds should be empty slice")
	})

	t.Run("Fails because of an internal server error querying the database", func(t *testing.T) {
		// Test: Database connection error during card number retrieval / Error de conexión de base de datos durante recuperación de números de tarjeta
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedCardNumberIds := []string{}
		expectedError := error_message.ErrInternalServerError

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT id_card_number FROM employees").
			WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		cardNumberIds, err := repository.GetCardNumberIds()

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedCardNumberIds, cardNumberIds, "cardNumberIds should be empty slice")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})
}

// TestMySqlEmployeeRepositoryExistEmployeeById - Tests for ExistEmployeeById method / Tests para método ExistEmployeeById
// Cases: employee exists, employee doesn't exist, internal server error
// Casos: empleado existe, empleado no existe, error interno del servidor
func TestMySqlEmployeeRepositoryExistEmployeeById(t *testing.T) {
	t.Run("Successfully returns true when employee exists in database", func(t *testing.T) {
		// Test: Existing employee ID returns true / ID de empleado existente retorna true
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"1"}).AddRow(1)

		mock.ExpectQuery("SELECT 1 FROM employees WHERE id = \\? LIMIT 1").
			WithArgs(1).WillReturnRows(rows)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		exists, err := repository.ExistEmployeeById(context.Background(), 1)

		assert.Nil(t, err, "err should be nil")
		assert.True(t, exists, "exists should be true")
	})

	t.Run("Successfully returns false when employee doesn't exist in database", func(t *testing.T) {
		// Test: Non-existent employee ID returns false / ID de empleado inexistente retorna false
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT 1 FROM employees WHERE id = \\? LIMIT 1").
			WithArgs(1).WillReturnError(sql.ErrNoRows)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		exists, err := repository.ExistEmployeeById(context.Background(), 1)

		assert.Nil(t, err, "err should be nil")
		assert.False(t, exists, "exists should be false")
	})

	t.Run("Fails because of an internal server error querying the database", func(t *testing.T) {
		// Test: Database connection error during existence check / Error de conexión de base de datos durante verificación de existencia
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT 1 FROM employees WHERE id = \\? LIMIT 1").
			WithArgs(1).WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		exists, err := repository.ExistEmployeeById(context.Background(), 1)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.False(t, exists, "exists should be false")
	})
}

// TestGetNewEmployeeMySQLRepository - Tests for singleton pattern in GetNewEmployeeMySQLRepository
// TestGetNewEmployeeMySQLRepository - Tests para patrón singleton en GetNewEmployeeMySQLRepository
func TestGetNewEmployeeMySQLRepository(t *testing.T) {
	t.Run("Successfully creates new employee repository when instance doesn't exist", func(t *testing.T) {
		// Test: First call creates new instance / Primera llamada crea nueva instancia
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		assert.NotNil(t, repository, "repository should not be nil")
		mock.ExpectationsWereMet()
	})

	t.Run("Returns same instance when repository already exists", func(t *testing.T) {
		// Test: Subsequent calls return same instance / Llamadas subsecuentes retornan la misma instancia
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()

		db1, mock1, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db1.Close()

		db2, mock2, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db2.Close()

		// First call creates instance
		repository1 := repositories.GetNewEmployeeMySQLRepository(db1)
		// Second call should return same instance (singleton pattern)
		repository2 := repositories.GetNewEmployeeMySQLRepository(db2)

		assert.NotNil(t, repository1, "repository1 should not be nil")
		assert.NotNil(t, repository2, "repository2 should not be nil")
		assert.Equal(t, repository1, repository2, "Both repositories should be the same instance")
		mock1.ExpectationsWereMet()
		mock2.ExpectationsWereMet()
	})
}

// Additional tests to cover missing paths in existing methods
// Tests adicionales para cubrir paths faltantes en métodos existentes

func TestMySqlEmployeeRepositoryGetById_AdditionalCoverage(t *testing.T) {
	t.Run("Fails because of scan error with invalid data type", func(t *testing.T) {
		// Test: Database returns invalid data types causing scan error / Base de datos retorna tipos de datos inválidos causando error de scan
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedEmployee := models.Employee{}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		// Return invalid data type for ID field
		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
			AddRow("INVALID_ID", "CARD-001", "Pedro", "Martinez", 1)

		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(1).WillReturnRows(rows)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.GetById(context.Background(), 1)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should be empty")
	})
}

func TestMySqlEmployeeRepositoryDeleteById_AdditionalCoverage(t *testing.T) {
	t.Run("Fails because of rows affected error", func(t *testing.T) {
		// Test: Error when getting rows affected after successful exec / Error al obtener filas afectadas después de exec exitoso
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		expectedError := error_message.ErrInternalServerError

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("DELETE FROM employees WHERE id = ?").
			WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(error_message.ErrInternalServerError))

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		err = repository.DeleteById(context.Background(), 1)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})
}

func TestMySqlEmployeeRepositoryUpdate_AdditionalCoverage(t *testing.T) {
	t.Run("Successfully updates employee with CardNumberID field", func(t *testing.T) {
		// Test: Update only CardNumberID field / Actualizar solo el campo CardNumberID
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		inputEmployee := models.Employee{
			CardNumberID: "UPDATED-CARD",
		}
		expectedEmployee := models.Employee{
			Id:           1,
			CardNumberID: "UPDATED-CARD",
			FirstName:    "Pedro",
			LastName:     "Martinez",
			WarehouseID:  1,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		// Mock the UPDATE query / Mock de la consulta UPDATE
		mock.ExpectExec("UPDATE employees SET id_card_number = \\? WHERE id = \\?").
			WithArgs("UPDATED-CARD", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Mock the GetById query after update / Mock de la consulta GetById después de actualización
		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
			AddRow(1, "UPDATED-CARD", "Pedro", "Martinez", 1)
		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(1).WillReturnRows(rows)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.Update(context.Background(), 1, inputEmployee)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should match expected employee")
	})

	t.Run("Successfully updates employee with WarehouseID field", func(t *testing.T) {
		// Test: Update only WarehouseID field / Actualizar solo el campo WarehouseID
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		inputEmployee := models.Employee{
			WarehouseID: 99,
		}
		expectedEmployee := models.Employee{
			Id:           1,
			CardNumberID: "CARD-001",
			FirstName:    "Pedro",
			LastName:     "Martinez",
			WarehouseID:  99,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		// Mock the UPDATE query / Mock de la consulta UPDATE
		mock.ExpectExec("UPDATE employees SET warehouse_id = \\? WHERE id = \\?").
			WithArgs(99, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Mock the GetById query after update / Mock de la consulta GetById después de actualización
		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
			AddRow(1, "CARD-001", "Pedro", "Martinez", 99)
		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(1).WillReturnRows(rows)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.Update(context.Background(), 1, inputEmployee)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should match expected employee")
	})

	t.Run("Successfully updates employee with all fields", func(t *testing.T) {
		// Test: Update all fields / Actualizar todos los campos
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		inputEmployee := models.Employee{
			CardNumberID: "ALL-UPDATED",
			FirstName:    "UpdatedFirst",
			LastName:     "UpdatedLast",
			WarehouseID:  999,
		}
		expectedEmployee := models.Employee{
			Id:           1,
			CardNumberID: "ALL-UPDATED",
			FirstName:    "UpdatedFirst",
			LastName:     "UpdatedLast",
			WarehouseID:  999,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		// Mock the UPDATE query with all fields / Mock de la consulta UPDATE con todos los campos
		mock.ExpectExec("UPDATE employees SET first_name = \\?, last_name = \\?, id_card_number = \\?, warehouse_id = \\? WHERE id = \\?").
			WithArgs("UpdatedFirst", "UpdatedLast", "ALL-UPDATED", 999, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Mock the GetById query after update / Mock de la consulta GetById después de actualización
		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
			AddRow(1, "ALL-UPDATED", "UpdatedFirst", "UpdatedLast", 999)
		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(1).WillReturnRows(rows)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.Update(context.Background(), 1, inputEmployee)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should match expected employee")
	})

	t.Run("Fails because of rows affected error", func(t *testing.T) {
		// Test: Error when getting rows affected after update / Error al obtener filas afectadas después de actualización
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		inputEmployee := models.Employee{
			FirstName: "UpdatedName",
		}
		expectedEmployee := models.Employee{}
		expectedError := error_message.ErrInternalServerError

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("UPDATE employees SET first_name = \\? WHERE id = \\?").
			WithArgs("UpdatedName", 1).
			WillReturnResult(sqlmock.NewErrorResult(error_message.ErrInternalServerError))

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.Update(context.Background(), 1, inputEmployee)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should be empty")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})

	t.Run("Fails because GetById fails after successful update", func(t *testing.T) {
		// Test: Update succeeds but GetById fails afterward / Actualización exitosa pero GetById falla después
		// Reset singleton for testing / Resetear singleton para testing
		repositories.ResetEmployeeRepositoryInstance()
		inputEmployee := models.Employee{
			FirstName: "UpdatedName",
		}
		expectedEmployee := models.Employee{}
		expectedError := error_message.ErrInternalServerError

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		// Mock successful UPDATE query / Mock de consulta UPDATE exitosa
		mock.ExpectExec("UPDATE employees SET first_name = \\? WHERE id = \\?").
			WithArgs("UpdatedName", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Mock failing GetById query / Mock de consulta GetById que falla
		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(1).WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.GetNewEmployeeMySQLRepository(db)

		dbEmployee, err := repository.Update(context.Background(), 1, inputEmployee)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedEmployee, dbEmployee, "dbEmployee should be empty")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})
}
