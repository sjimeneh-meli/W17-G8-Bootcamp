package repositories_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestMySqlBuyerRepositoryGetAll(t *testing.T) {
	t.Run("Successfully returns filled buyers map when there is data on db response", func(t *testing.T) {
		expectedBuyers := map[int]models.Buyer{
			1: {
				Id:           1,
				CardNumberId: "CARD-001",
				FirstName:    "Ignacio",
				LastName:     "Garcia",
			},
			2: {
				Id:           2,
				CardNumberId: "CARD-002",
				FirstName:    "Jesus",
				LastName:     "Ortega",
			},
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name"}).
			AddRow("1", "CARD-001", "Ignacio", "Garcia").
			AddRow("2", "CARD-002", "Jesus", "Ortega")

		mock.ExpectQuery("select id, id_card_number, first_name, last_name from buyers").WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		dbBuyers, err := repository.GetAll(context.Background())

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedBuyers, dbBuyers, "dbBuyers should be a map with two buyers")
	})

	t.Run("Successfully returns empty buyers map when there isn't data on db response", func(t *testing.T) {
		expectedBuyers := map[int]models.Buyer{}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name"})

		mock.ExpectQuery("select id, id_card_number, first_name, last_name from buyers").WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		dbBuyers, err := repository.GetAll(context.Background())

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedBuyers, dbBuyers, "dbBuyers should be an empty map of buyers")
	})

	t.Run("Fails because of an internal server error querying the database", func(t *testing.T) {
		expectedError := error_message.ErrInternalServerError
		expectedBuyers := map[int]models.Buyer{}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery("select id, id_card_number, first_name, last_name from buyers").WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		dbBuyers, err := repository.GetAll(context.Background())

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedBuyers, dbBuyers, "dbBuyers should be an empty map of buyers")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})

	t.Run("Fails because of an internal server error scanning database results", func(t *testing.T) {
		expectedBuyers := map[int]models.Buyer{}
		expectedError := error_message.ErrInternalServerError
		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name"}).
			AddRow("1", "CARD-001", "Ignacio", "Garcia").
			AddRow("STRING", "CARD-002", "Jesus", "Ortega")

		mock.ExpectQuery("select id, id_card_number, first_name, last_name from buyers").WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		dbBuyers, err := repository.GetAll(context.Background())

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedBuyers, dbBuyers, "dbBuyers should be an empty map of buyers")
		assert.ErrorIs(t, err, expectedError)
	})
}

func TestGetById(t *testing.T) {
	t.Run("Successfully return searched buyer from db", func(t *testing.T) {
		expectedBuyer := models.Buyer{
			Id:           10,
			CardNumberId: "CARD-001",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		row := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name"}).
			AddRow("10", "CARD-001", "Ignacio", "Garcia")
		mock.ExpectQuery(
			"select id, id_card_number, first_name, last_name from buyers where id = ?").
			WithArgs(expectedBuyer.Id).
			WillReturnRows(row).
			RowsWillBeClosed()

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}
		dbBuyer, err := repository.GetById(context.Background(), expectedBuyer.Id)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedBuyer, dbBuyer)

	})

	t.Run("Fails because of an internal server error on the row returned from the database", func(t *testing.T) {
		expectedBuyer := models.Buyer{}
		expectedError := error_message.ErrInternalServerError

		searchId := 1
		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery("select id, id_card_number, first_name, last_name from buyers where id = ?").
			WithArgs(searchId).
			WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		buyerDb, err := repository.GetById(context.Background(), searchId)

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, expectedError)
		assert.Equal(t, expectedBuyer, buyerDb)
	})

	t.Run("Fails because of an internal server error scanning database results", func(t *testing.T) {
		expectedBuyer := models.Buyer{}
		expectedError := error_message.ErrInternalServerError
		searchId := 1
		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name"}).
			AddRow("STRING", "CARD-001", "Ignacio", "Garcia")

		mock.ExpectQuery("select id, id_card_number, first_name, last_name from buyers where id = ?").
			WithArgs(searchId).
			WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		dbBuyers, err := repository.GetById(context.Background(), searchId)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedBuyer, dbBuyers, "dbBuyers should be an empty map of buyers")
		assert.ErrorIs(t, err, expectedError)
	})

	t.Run("Fails because there is no row returned from query", func(t *testing.T) {
		expectedBuyer := models.Buyer{}
		expectedError := error_message.ErrNotFound

		searchId := 10

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()
		row := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name"})

		mock.ExpectQuery("select id, id_card_number, first_name, last_name from buyers where id = ?").
			WithArgs(searchId).
			WillReturnRows(row).RowsWillBeClosed()

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		buyerDb, err := repository.GetById(context.Background(), searchId)

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, expectedError)
		assert.Equal(t, expectedBuyer, buyerDb)
	})

}

func TestDeleteById(t *testing.T) {
	t.Run("Successfully delete buyer from db", func(t *testing.T) {
		searchId := 10

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("delete from buyers where id = ?").
			WithArgs(searchId).
			WillReturnResult(sqlmock.NewResult(0, 1))

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		err = repository.DeleteById(context.Background(), searchId)

		assert.Nil(t, err)
	})

	t.Run("Fails because of an internal server error querying the database", func(t *testing.T) {
		expectedError := error_message.ErrInternalServerError
		searchId := 10

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("delete from buyers where id = ?").WithArgs(searchId).
			WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		err = repository.DeleteById(context.Background(), searchId)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})

	t.Run("Fails because of an internal server error getting affected rows", func(t *testing.T) {
		expectedError := error_message.ErrInternalServerError
		searchId := 10

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("delete from buyers where id = ?").WithArgs(searchId).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("error at RowsAffected")))

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		err = repository.DeleteById(context.Background(), searchId)

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, expectedError)

	})

	t.Run("Fails because the buyer id intended to delete isn't found", func(t *testing.T) {
		expectedError := error_message.ErrNotFound
		searchId := 10

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("delete from buyers where id = ?").WithArgs(searchId).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		err = repository.DeleteById(context.Background(), searchId)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, expectedError)
	})
}

func TestCreate(t *testing.T) {
	t.Run("Successfully create a new buyer record on db", func(t *testing.T) {
		expectedBuyer := models.Buyer{
			Id:           17,
			CardNumberId: "CARD-0017",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		buyerToInsert := models.Buyer{
			Id:           0,
			CardNumberId: "CARD-0017",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta("insert into buyers (id_card_number, first_name, last_name) values (?, ?, ?)")).
			WithArgs(buyerToInsert.CardNumberId, buyerToInsert.FirstName, buyerToInsert.LastName).
			WillReturnResult(sqlmock.NewResult(17, 1))

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		buyerDb, err := repository.Create(context.Background(), buyerToInsert)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedBuyer, buyerDb)
	})

	t.Run("Fails because of an internal server error excecuting insert", func(t *testing.T) {
		expectedBuyer := models.Buyer{}
		expectedError := error_message.ErrInternalServerError

		buyerToInsert := models.Buyer{
			Id:           0,
			CardNumberId: "CARD-0017",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta("insert into buyers (id_card_number, first_name, last_name) values (?, ?, ?)")).
			WithArgs(buyerToInsert.CardNumberId, buyerToInsert.FirstName, buyerToInsert.LastName).
			WillReturnError(errors.New("error excecuting query"))

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		buyerDb, err := repository.Create(context.Background(), buyerToInsert)

		assert.NotNil(t, err, "err should not be nil")
		assert.ErrorIs(t, err, expectedError)
		assert.Equal(t, expectedBuyer, buyerDb)
	})

	t.Run("Fails because of an internal server error getting the inserted record Id", func(t *testing.T) {
		expectedBuyer := models.Buyer{}
		expectedError := error_message.ErrInternalServerError

		buyerToInsert := models.Buyer{
			Id:           0,
			CardNumberId: "CARD-0017",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta("insert into buyers (id_card_number, first_name, last_name) values (?, ?, ?)")).
			WithArgs(buyerToInsert.CardNumberId, buyerToInsert.FirstName, buyerToInsert.LastName).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("error getting LastInsertId")))

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		buyerDb, err := repository.Create(context.Background(), buyerToInsert)

		assert.NotNil(t, err, "err should not be nil")
		assert.ErrorIs(t, err, expectedError)
		assert.Equal(t, expectedBuyer, buyerDb)
	})

}
