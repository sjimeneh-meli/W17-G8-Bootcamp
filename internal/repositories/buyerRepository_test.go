package repositories_test

import (
	"context"
	"database/sql"
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

		mock.ExpectQuery("select id, id_card_number, first_name, last_name from buyers").
			WillReturnError(error_message.ErrInternalServerError)

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

func TestUpdate(t *testing.T) {
	t.Run("Successfully updates a buyer", func(t *testing.T) {
		expectedBuyer := models.Buyer{
			Id:           18,
			CardNumberId: "CARD-0018",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		searchId := 18
		buyerToUpdate := models.Buyer{
			FirstName: "Nacho",
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE buyers SET first_name = ? WHERE id = ?")).
			WithArgs(buyerToUpdate.FirstName, searchId).
			WillReturnResult(sqlmock.NewResult(0, 1))

		row := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name"}).
			AddRow("18", "CARD-0018", "Ignacio", "Garcia")

		mock.ExpectQuery("select id, id_card_number, first_name, last_name from buyers where id = ?").
			WithArgs(searchId).
			WillReturnRows(row)

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		buyerDb, err := repository.Update(context.Background(), searchId, buyerToUpdate)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedBuyer, buyerDb)

	})

	t.Run("Fails because of an internal error excecuting the update", func(t *testing.T) {
		expectedBuyer := models.Buyer{}
		expectedError := error_message.ErrInternalServerError
		searchId := 18

		buyerToUpdate := models.Buyer{
			LastName: "Garcia",
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE buyers SET last_name = ? WHERE id = ?")).
			WithArgs(buyerToUpdate.LastName, searchId).
			WillReturnError(errors.New("error excecuting the query"))

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		buyerDb, err := repository.Update(context.Background(), searchId, buyerToUpdate)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, expectedError)
		assert.Equal(t, expectedBuyer, buyerDb)
	})

	t.Run("Fails because of an internal error getting the affected rows", func(t *testing.T) {
		expectedError := error_message.ErrInternalServerError
		expectedBuyer := models.Buyer{}

		buyerToUpdate := models.Buyer{
			CardNumberId: "CARD-00019",
		}
		searchId := 10

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE buyers SET id_card_number = ? WHERE id = ?")).
			WithArgs(buyerToUpdate.CardNumberId, searchId).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("error at RowsAffected")))

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		buyerDb, err := repository.Update(context.Background(), searchId, buyerToUpdate)

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, expectedError)
		assert.Equal(t, expectedBuyer, buyerDb)
	})

	t.Run("Fails because the buyer Id doesn't exists", func(t *testing.T) {
		expectedBuyer := models.Buyer{}
		expectedError := error_message.ErrNotFound

		searchId := 18
		buyerToUpdate := models.Buyer{
			FirstName: "Nacho",
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE buyers SET first_name = ? WHERE id = ?")).
			WithArgs(buyerToUpdate.FirstName, searchId).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		buyerDb, err := repository.Update(context.Background(), searchId, buyerToUpdate)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, expectedError, "error should be of type ErrNotFound")
		assert.Equal(t, expectedBuyer, buyerDb)
	})

	t.Run("Fails obtaining the updated buyer", func(t *testing.T) {
		expectedBuyer := models.Buyer{}
		expectedError := error_message.ErrInternalServerError

		searchId := 18
		buyerToUpdate := models.Buyer{
			FirstName: "Nacho",
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE buyers SET first_name = ? WHERE id = ?")).
			WithArgs(buyerToUpdate.FirstName, searchId).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("select id, id_card_number, first_name, last_name from buyers where id = ?").
			WithArgs(searchId).
			WillReturnError(errors.New("error on GetById"))

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		buyerDb, err := repository.Update(context.Background(), searchId, buyerToUpdate)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, expectedError, "err should be of type internal server error")
		assert.Equal(t, expectedBuyer, buyerDb)
	})
}

func TestGetCardNumberIds(t *testing.T) {
	t.Run("Fails because of an internal server error on the query", func(t *testing.T) {
		expectedCardNumberIdList := []string{}
		expectedError := error_message.ErrInternalServerError

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery("select id_card_number from buyers").
			WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		dbCardNumberIdsList, err := repository.GetCardNumberIds()

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, expectedError)
		assert.Equal(t, expectedCardNumberIdList, dbCardNumberIdsList)
	})

	t.Run("Successfully returns a cardNumberId list", func(t *testing.T) {
		expectedCardNumberIdList := []string{
			"CARD-001", "CARD-002",
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id_card_number"}).
			AddRow("CARD-001").AddRow("CARD-002")
		mock.ExpectQuery("select id_card_number from buyers").
			WillReturnRows(rows)

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		dbCardNumberIdsList, err := repository.GetCardNumberIds()

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedCardNumberIdList, dbCardNumberIdsList)
	})
}

func TestExistBuyerById(t *testing.T) {
	t.Run("Successfully returns true when buyer exists", func(t *testing.T) {
		expectedExists := true
		searchId := 1

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		row := mock.NewRows([]string{"1"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM buyers WHERE id = ? LIMIT 1")).
			WithArgs(searchId).
			WillReturnRows(row)

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		exists, err := repository.ExistBuyerById(context.Background(), searchId)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedExists, exists)
	})

	t.Run("Successfully returns false when buyer doesn't exist", func(t *testing.T) {
		expectedExists := false
		searchId := 999

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM buyers WHERE id = ? LIMIT 1")).
			WithArgs(searchId).
			WillReturnError(sql.ErrNoRows)

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		exists, err := repository.ExistBuyerById(context.Background(), searchId)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedExists, exists)
	})

	t.Run("Fails because of an internal server error querying the database", func(t *testing.T) {
		expectedExists := false
		searchId := 1

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM buyers WHERE id = ? LIMIT 1")).
			WithArgs(searchId).
			WillReturnError(errors.New("database connection error"))

		repository := repositories.MySqlBuyerRepository{
			Db: db,
		}

		exists, err := repository.ExistBuyerById(context.Background(), searchId)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, error_message.ErrInternalServerError)
		assert.Equal(t, expectedExists, exists)
	})
}

func TestGetNewBuyerMySQLRepository(t *testing.T) {
	t.Run("Successfully returns a new buyer repository", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		repository := repositories.GetNewBuyerMySQLRepository(db)

		assert.NotNil(t, repository, "repository shouldn't be nil")
	})

	t.Run("Multiple calls to GetNewBuyerMySQLRepository should return the same instance", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		repository := repositories.GetNewBuyerMySQLRepository(db)

		assert.NotNil(t, repository, "repository shouldn't be nil")

		repository2 := repositories.GetNewBuyerMySQLRepository(nil)

		assert.Equal(t, repository, repository2, "repository and repository2 should be the same")
	})
}
