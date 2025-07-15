package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/database"
)

var sectionRepositoryInstance SectionRepositoryI

func GetSectionRepository(db *sql.DB) SectionRepositoryI {
	if sectionRepositoryInstance != nil {
		return sectionRepositoryInstance
	}

	sectionRepositoryInstance = &sectionRepository{
		database:  db,
		tablename: "sections",
	}

	return sectionRepositoryInstance
}

type SectionRepositoryI interface {
	GetAll(ctx context.Context) ([]*models.Section, error)
	GetByID(ctx context.Context, id int) (*models.Section, error)
	Create(ctx context.Context, model *models.Section) error
	Update(ctx context.Context, model *models.Section) error
	ExistWithID(ctx context.Context, id int) bool
	ExistsWithSectionNumber(ctx context.Context, id int, sectionNumber string) bool
	DeleteByID(ctx context.Context, id int) error
}

type sectionRepository struct {
	database  *sql.DB
	tablename string
}

func (r *sectionRepository) GetAll(ctx context.Context) ([]*models.Section, error) {
	columns := []string{"Id", "section_number", "current_capacity", "current_temperature", "maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id"}

	rows, err := database.Select(ctx, r.database, r.tablename, columns, "")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sections []*models.Section

	for rows.Next() {
		var section models.Section
		if err := rows.Scan(
			&section.Id,
			&section.SectionNumber,
			&section.CurrentCapacity,
			&section.CurrentTemperature,
			&section.MaximumCapacity,
			&section.MinimumCapacity,
			&section.MinimumTemperature,
			&section.ProductTypeID,
			&section.WarehouseID,
		); err != nil {
			return sections, err
		}

		sections = append(sections, &section)
	}

	return sections, nil

}

func (r *sectionRepository) GetByID(ctx context.Context, id int) (*models.Section, error) {
	columns := []string{"Id", "section_number", "current_capacity", "current_temperature", "maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id"}
	row := database.SelectOne(ctx, r.database, r.tablename, columns, "Id = ?", id)

	var section models.Section

	if err := row.Scan(
		&section.Id,
		&section.SectionNumber,
		&section.CurrentCapacity,
		&section.CurrentTemperature,
		&section.MaximumCapacity,
		&section.MinimumCapacity,
		&section.MinimumTemperature,
		&section.ProductTypeID,
		&section.WarehouseID,
	); err != nil {
		return nil, err
	}

	return &section, nil
}

func (r *sectionRepository) Create(ctx context.Context, model *models.Section) error {
	data := make(map[any]any)
	data["section_number"] = model.SectionNumber
	data["current_capacity"] = model.CurrentCapacity
	data["current_temperature"] = model.CurrentTemperature
	data["maximum_capacity"] = model.MaximumCapacity
	data["minimum_capacity"] = model.MinimumCapacity
	data["minimum_temperature"] = model.MinimumTemperature
	data["product_type_id"] = model.ProductTypeID
	data["warehouse_id"] = model.WarehouseID

	result, err := database.Insert(ctx, r.database, r.tablename, data)

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

func (r *sectionRepository) Update(ctx context.Context, model *models.Section) error {
	sqlStatement := fmt.Sprintf("UPDATE %s SET `section_number`=?, `current_capacity`=?, `current_temperature`=?, `maximum_capacity`=?, `minimum_capacity`=?, `minimum_temperature`=?, `product_type_id`=?, `warehouse_id`=? WHERE `Id`=?", r.tablename)
	_, err := r.database.Exec(sqlStatement,
		model.SectionNumber,
		model.CurrentCapacity,
		model.CurrentTemperature,
		model.MaximumCapacity,
		model.MinimumCapacity,
		model.MinimumTemperature,
		model.ProductTypeID,
		model.WarehouseID,
		model.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *sectionRepository) ExistWithID(ctx context.Context, id int) bool {
	row := database.SelectOne(ctx, r.database, r.tablename, []string{"COUNT(Id)"}, "Id = ?", id)
	var count int
	if err := row.Scan(&count); err != nil {
		return true
	}

	return count != 0
}

func (r *sectionRepository) ExistsWithSectionNumber(ctx context.Context, id int, sectionNumber string) bool {
	row := database.SelectOne(ctx, r.database, r.tablename, []string{"COUNT(Id)"}, "section_number = ? AND Id <> ?", sectionNumber, id)
	var count string
	if err := row.Scan(&count); err != nil {
		return true
	}

	return count != "0"
}

func (r *sectionRepository) DeleteByID(ctx context.Context, id int) error {
	_, err := database.Delete(ctx, r.database, r.tablename, "Id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
