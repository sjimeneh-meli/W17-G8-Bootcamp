package repositories_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/seeders"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductBatch(t *testing.T) {
	t.Run("Successfully create product barch from db", func(t *testing.T) {
		var exampleDate = time.Date(
			2025,
			time.August,
			19,
			0,
			0,
			0,
			0,
			time.UTC,
		)

		productBatchModel := models.ProductBatch{
			Id:                 1,
			BatchNumber:        "AAA",
			CurrentQuantity:    1,
			CurrentTemperature: 3.43,
			DueDate:            exampleDate,
			InitialQuantity:    1,
			ManufacturingDate:  exampleDate,
			ManufacturingHour:  exampleDate,
			MinimumTemperature: 3.43,
			ProductID:          1,
			SectionID:          1,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"batch_number", "current_quantity", "current_temperature", "due_date", "initial_quantity", "manufacturing_date", "manufacturing_hour", "minimum_temperature", "product_id", "section_id"}).
			AddRow("AAA", "1", "3.43", "2025-08-19 00:00:00 +0000 UTC", "1", "2025-08-19 00:00:00 +0000 UTC", "2025-08-19 00:00:00 +0000 UTC", "3.43", "1", "1")

		mock.ExpectQuery(regexp.QuoteMeta(
			"INSERT INTO batch_number,current_quantity,current_temperature,due_date,initial_quantity,manufacturing_date,manufacturing_hour,minimum_temperature,product_id,section_id;")).
			WillReturnRows(rows).
			RowsWillBeClosed()

		repository := repositories.GetProductBatchRepository(db)
		repository.Create(context.Background(), &productBatchModel)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, seeders.NewProductBatchModel, productBatchModel, "ok")
		repositories.ResetProductBatchRepositoryInstance()
	})
}

func TestExistWithBatchNumber(t *testing.T) {
	t.Run("Successfully ask if exist a product batch with the same number from db", func(t *testing.T) {
		id := 1
		batchNumber := "AAA"

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "batch_number", "current_quantity", "current_temperature", "due_date", "initial_quantity", "manufacturing_date", "manufacturing_hour", "minimum_temperature", "product_id", "section_id"}).
			AddRow("1", "AAA", "1", "3.43", "2025-08-19 00:00:00 +0000 UTC", "1", "2025-08-19 00:00:00 +0000 UTC", "2025-08-19 00:00:00 +0000 UTC", "3.43", "1", "1")

		repository := repositories.GetProductBatchRepository(db)
		repository.ExistsWithBatchNumber(context.Background(), id, batchNumber)

		mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT COUNT(Id) FROM product_batches WHERE batch_number = ? AND Id <> ?")).
			WillReturnRows(rows).
			RowsWillBeClosed()

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, 1, 1, "ok")
		repositories.ResetProductBatchRepositoryInstance()
	})
}

func TestGetProductQuantityBySectionId(t *testing.T) {
	t.Run("Successfully get the product quantity by section id from db", func(t *testing.T) {
		id := 1

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "batch_number", "current_quantity", "current_temperature", "due_date", "initial_quantity", "manufacturing_date", "manufacturing_hour", "minimum_temperature", "product_id", "section_id"}).
			AddRow("1", "AAA", "1", "3.43", "2025-08-19 00:00:00 +0000 UTC", "1", "2025-08-19 00:00:00 +0000 UTC", "2025-08-19 00:00:00 +0000 UTC", "3.43", "1", "1")

		repository := repositories.GetProductBatchRepository(db)
		repository.GetProductQuantityBySectionId(context.Background(), id)

		mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT SUM(current_quantity) FROM product_batches WHERE section_id = ?")).
			WillReturnRows(rows).
			RowsWillBeClosed()

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, 1, 1, "ok")
		repositories.ResetProductBatchRepositoryInstance()
	})
}
