package repositories

import (
	"database/sql"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/database"
)

var productBatchRepositoryInstance ProductBatchRepositoryI

func GetProductBatchRepository(db *sql.DB) ProductBatchRepositoryI {
	if productBatchRepositoryInstance != nil {
		return productBatchRepositoryInstance
	}

	productBatchRepositoryInstance = &productBatchRepository{
		database:  db,
		tablename: "product_batches",
	}
	return productBatchRepositoryInstance
}

type ProductBatchRepositoryI interface {
	Create(model *models.ProductBatch) error
	GetProductQuantityBySectionId(id int) int
	ExistsWithBatchNumber(id int, batchNumber string) bool
}

type productBatchRepository struct {
	database  *sql.DB
	tablename string
}

func (r *productBatchRepository) Create(model *models.ProductBatch) error {
	data := make(map[any]any)
	data["batch_number"] = model.BatchNumber
	data["current_quantity"] = model.CurrentQuantity
	data["current_temperature"] = model.CurrentTemperature
	data["due_date"] = model.DueDate
	data["initial_quantity"] = model.InitialQuantity
	data["manufacturing_date"] = model.ManufacturingDate
	data["manufacturing_hour"] = model.ManufacturingHour
	data["minimum_temperature"] = model.MinimumTemperature
	data["product_id"] = model.ProductID
	data["section_id"] = model.SectionID

	result, err := database.Insert(r.database, r.tablename, data)

	if err != nil {
		return err
	}

	newID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	model.Id = int(newID)
	return nil
}

func (r *productBatchRepository) ExistsWithBatchNumber(id int, batchNumber string) bool {
	row := database.SelectOne(r.database, r.tablename, []string{"COUNT(Id)"}, "batch_number = ? AND Id <> ?", batchNumber, id)
	var count string
	if err := row.Scan(&count); err != nil {
		return true
	}

	return count != "0"
}

func (r *productBatchRepository) GetProductQuantityBySectionId(id int) int {
	row := database.SelectOne(r.database, r.tablename, []string{"SUM(current_quantity)"}, "section_id = ?", id)
	var quantity int
	if err := row.Scan(&quantity); err != nil {
		return quantity
	}

	return quantity
}
