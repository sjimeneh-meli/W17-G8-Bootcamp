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

func TestMySqlEmployeeRepositoryGetAll(t *testing.T) {
	t.Run("Successfully returns filled employees map when there is data on db response", func(t *testing.T) {
		// Reset singleton for testing
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
		// Reset singleton for testing
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
		// Reset singleton for testing
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
		// Reset singleton for testing
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

func TestMySqlEmployeeRepositoryGetById(t *testing.T) {
	t.Run("Successfully returns employee when found in database", func(t *testing.T) {
		// Reset singleton for testing
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
		// Reset singleton for testing
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
		// Reset singleton for testing
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

func TestMySqlEmployeeRepositoryDeleteById(t *testing.T) {
	t.Run("Successfully deletes employee when exists in database", func(t *testing.T) {
		// Reset singleton for testing
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
		// Reset singleton for testing
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
		// Reset singleton for testing
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

func TestMySqlEmployeeRepositoryCreate(t *testing.T) {
	t.Run("Successfully creates employee in database", func(t *testing.T) {
		// Reset singleton for testing
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
		// Reset singleton for testing
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
		// Reset singleton for testing
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

func TestMySqlEmployeeRepositoryUpdate(t *testing.T) {
	t.Run("Successfully updates employee when exists in database", func(t *testing.T) {
		// Reset singleton for testing
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

		// Mock the UPDATE query
		mock.ExpectExec("UPDATE employees SET first_name = \\?, last_name = \\? WHERE id = \\?").
			WithArgs("UpdatedName", "UpdatedLastName", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Mock the GetById query after update
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
		// Reset singleton for testing
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
		// Reset singleton for testing
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

func TestMySqlEmployeeRepositoryGetCardNumberIds(t *testing.T) {
	t.Run("Successfully returns all card number ids when there is data in database", func(t *testing.T) {
		// Reset singleton for testing
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
		// Reset singleton for testing
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
		// Reset singleton for testing
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

func TestMySqlEmployeeRepositoryExistEmployeeById(t *testing.T) {
	t.Run("Successfully returns true when employee exists in database", func(t *testing.T) {
		// Reset singleton for testing
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
		// Reset singleton for testing
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
		// Reset singleton for testing
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
