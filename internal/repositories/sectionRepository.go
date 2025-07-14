package repositories

import (
	"database/sql"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	tools "github.com/sajimenezher_meli/meli-frescos-8/pkg"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

func GetSectionRepository(loader *loader.StorageJSON[models.Section], db *sql.DB) (SectionRepositoryI, error) {
	storage, err := loader.ReadAll()
	if err != nil {
		return nil, err
	}

	return &sectionRepository{
		storage:   storage,
		loader:    loader,
		database:  db,
		tablename: "sections",
	}, nil
}

type SectionRepositoryI interface {
	GetAll() ([]*models.Section, error)
	GetByID(id int) (*models.Section, error)
	Create(model *models.Section) error
	Update(model *models.Section) error
	ExistsWithSectionNumber(id int, sectionNumber string) bool
	DeleteByID(id int) error
}

type sectionRepository struct {
	storage   map[int]models.Section
	loader    *loader.StorageJSON[models.Section]
	database  *sql.DB
	tablename string
}

func (r *sectionRepository) GetAll() ([]*models.Section, error) {
	columns := []string{"Id", "section_number", "current_capacity", "current_temperature", "maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id"}

	rows, err := r.Select(r.tablename, columns, "")

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

func (r *sectionRepository) GetByID(id int) (*models.Section, error) {
	columns := []string{"Id", "section_number", "current_capacity", "current_temperature", "maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id"}
	row := r.SelectOne(r.tablename, columns, "Id = ?", id)

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

func (r *sectionRepository) Create(model *models.Section) error {
	data := make(map[any]any)
	data["section_number"] = model.SectionNumber
	data["current_capacity"] = model.CurrentCapacity
	data["current_temperature"] = model.CurrentTemperature
	data["maximum_capacity"] = model.MaximumCapacity
	data["minimum_capacity"] = model.MinimumCapacity
	data["minimum_temperature"] = model.MinimumTemperature
	data["product_type_id"] = model.ProductTypeID
	data["warehouse_id"] = model.WarehouseID

	result, err := r.Insert(r.tablename, data)

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

func (r *sectionRepository) Update(model *models.Section) error {
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

func (r *sectionRepository) ExistsWithSectionNumber(id int, sectionNumber string) bool {
	row := r.SelectOne(r.tablename, []string{"COUNT(Id)"}, "section_number = ? AND Id <> ?", sectionNumber, id)
	var count string
	if err := row.Scan(&count); err != nil {
		return true
	}

	return count != "0"
}

func (r *sectionRepository) DeleteByID(id int) error {
	_, err := r.Delete(r.tablename, "Id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *sectionRepository) SelectOne(tablename string, fields []string, condition string, values ...any) *sql.Row {
	columns := tools.SliceToString(fields, ",")
	sqlStatement := fmt.Sprintf("SELECT %s FROM %s", columns, tablename)
	if condition != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, condition)
	}

	return r.database.QueryRow(sqlStatement, values...)
}

func (r *sectionRepository) Select(tablename string, fields []string, condition string, values ...any) (*sql.Rows, error) {
	columns := tools.SliceToString(fields, ",")
	sqlStatement := fmt.Sprintf("SELECT %s FROM %s", columns, tablename)
	if condition != "" {
		sqlStatement = fmt.Sprintf("%s WHERE %s", sqlStatement, condition)
	}

	return r.database.Query(sqlStatement, values...)
}

func (r *sectionRepository) Insert(tablename string, data map[any]any) (sql.Result, error) {
	keys, values := tools.GetSlicesOfKeyAndValuesFromMap(data)
	columns := tools.SliceToString(keys, ",")
	placeholders := tools.SliceToString(tools.FillNewSlice(len(data), "?"), ",")

	sqlStatement := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s);", tablename, columns, placeholders)

	return r.database.Exec(sqlStatement, values...)
}

func (r *sectionRepository) Delete(tablename string, condition string, values ...any) (sql.Result, error) {
	sqlStatement := fmt.Sprintf("DELETE FROM %s WHERE %s;", tablename, condition)

	return r.database.Exec(sqlStatement, values...)
}
