package repositories_test

import (
	"context"
	"fmt"
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

	})

	t.Run("Fails because of an internal server error scanning database results", func(t *testing.T) {

	})

}
