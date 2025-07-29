package repositories_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

func TestGetByIdSection(t *testing.T) {
	t.Run("Successfully return searched section from db", func(t *testing.T) {
		section := &models.Section{
			Id:                 1,
			SectionNumber:      "A-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		row := mock.NewRows([]string{"Id", "section_number", "current_capacity", "current_temperature", "maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id"}).
			AddRow("1", "A-01", "2", "3.43", "2", "2", "2", "2", "1")

		mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT Id,section_number,current_capacity,current_temperature,maximum_capacity,minimum_capacity,minimum_temperature,product_type_id,warehouse_id FROM sections WHERE Id = ?")).
			WithArgs(section.Id).
			WillReturnRows(row).
			RowsWillBeClosed()

		repository := repositories.GetSectionRepository(db)
		sectionDB, err := repository.GetByID(context.Background(), section.Id)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, section, sectionDB)
		repositories.ResetSectionRepositoryInstance()
	})
}

func TestGetAllSections(t *testing.T) {
	t.Run("Successfully return all sections from db", func(t *testing.T) {
		sections := []*models.Section{
			{
				Id:                 1,
				SectionNumber:      "A-01",
				CurrentCapacity:    2,
				CurrentTemperature: 3.43,
				MaximumCapacity:    2,
				MinimumCapacity:    2,
				MinimumTemperature: 2,
				ProductTypeID:      2,
				WarehouseID:        1,
			},
			{
				Id:                 2,
				SectionNumber:      "B-01",
				CurrentCapacity:    2,
				CurrentTemperature: 3.43,
				MaximumCapacity:    2,
				MinimumCapacity:    2,
				MinimumTemperature: 2,
				ProductTypeID:      2,
				WarehouseID:        1,
			},
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"Id", "section_number", "current_capacity", "current_temperature", "maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id"}).
			AddRow("1", "A-01", "2", "3.43", "2", "2", "2", "2", "1").
			AddRow("2", "B-01", "2", "3.43", "2", "2", "2", "2", "1")

		mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT Id,section_number,current_capacity,current_temperature,maximum_capacity,minimum_capacity,minimum_temperature,product_type_id,warehouse_id FROM sections")).
			WillReturnRows(rows).
			RowsWillBeClosed()

		repository := repositories.GetSectionRepository(db)
		sectionDB, err := repository.GetAll(context.Background())

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, sections, sectionDB, "ok")
		repositories.ResetSectionRepositoryInstance()
	})
}

func TestDeleteByIdSections(t *testing.T) {
	t.Run("Successfully delete section from db", func(t *testing.T) {
		sectionID := 10

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectExec("DELETE FROM sections WHERE Id = ?").
			WithArgs(sectionID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		repository := repositories.GetSectionRepository(db)

		err = repository.DeleteByID(context.Background(), sectionID)

		assert.Nil(t, err)
		repositories.ResetSectionRepositoryInstance()
	})
}
