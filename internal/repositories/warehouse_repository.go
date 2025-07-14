package repositories

import (
	"database/sql"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type WarehouseRepository interface {
	GetAll() ([]models.Warehouse, error)
	Create(warehouse models.Warehouse) (models.Warehouse, error)
	ExistsByCode(code string) (bool, error)
	GetById(id int) (models.Warehouse, error)
	Delete(id int) error
	Update(id int, warehouse models.Warehouse) (models.Warehouse, error)
}

type WarehouseRepositoryImpl struct {
	db *sql.DB
}

func NewWarehouseRepository(db *sql.DB) *WarehouseRepositoryImpl {
	return &WarehouseRepositoryImpl{db: db}
}

func (r *WarehouseRepositoryImpl) GetAll() ([]models.Warehouse, error) {
	//1. Correr la query
	rows, err := r.db.Query("SELECT `id` , `address` , `telephone` , `warehouse_code` , `minimum_capacity` , `minimum_temperature` FROM `warehouse`")
	if err != nil {
		fmt.Printf("warehouse: %v\n", err.Error())
		return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	//2. Mapear los resultados
	warehouses := []models.Warehouse{}
	for rows.Next() {
		var w models.Warehouse
		err := rows.Scan(&w.Id, &w.Address, &w.Telephone, &w.WareHouseCode, &w.MinimumCapacity, &w.MinimumTemperature)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
		}
		warehouses = append(warehouses, w)

	}
	return warehouses, nil

}

func (r *WarehouseRepositoryImpl) Create(warehouse models.Warehouse) (models.Warehouse, error) {
	result, err := r.db.Exec("INSERT INTO `warehouse`(`address`, `telephone`, `warehouse_code` , `minimum_capacity` , `minimum_temperature` , `locality_id` ) VALUES (?,?,?,?,?,?)",
		warehouse.Address, warehouse.Telephone, warehouse.WareHouseCode, warehouse.MinimumCapacity, warehouse.MinimumTemperature, warehouse.LocalityId)
	if err != nil {
		fmt.Printf("warehouse: %v\n", err.Error())
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	warehouse.Id = int(lastInsertId)

	return warehouse, nil
}

func (r *WarehouseRepositoryImpl) ExistsByCode(code string) (bool, error) {
	row := r.db.QueryRow("SELECT COUNT(*) FROM `warehouse` WHERE `warehouse_code` = ?", code)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return count > 0, nil
}

func (r *WarehouseRepositoryImpl) GetById(id int) (models.Warehouse, error) {
	row := r.db.QueryRow(
		"SELECT `id`,`address`, `telephone`, `warehouse_code` , `minimum_capacity` , `minimum_temperature` FROM `warehouse` WHERE `id` = ?",
		id,
	)

	var warehouse models.Warehouse
	err := row.Scan(&warehouse.Id, &warehouse.Address, &warehouse.Telephone, &warehouse.WareHouseCode, &warehouse.MinimumCapacity, &warehouse.MinimumTemperature)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Warehouse{}, fmt.Errorf("%w: warehouse with id %d", error_message.ErrNotFound, id)
		}
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return warehouse, nil
}

func (r *WarehouseRepositoryImpl) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM `warehouse` WHERE `id` = ? ", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%w: warehouse with id %d", error_message.ErrNotFound, id)
		}
		fmt.Printf("warehouse: %v\n", err.Error())
		return fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return nil
}

func (r *WarehouseRepositoryImpl) Update(id int, warehouse models.Warehouse) (models.Warehouse, error) {
	result, err := r.db.Exec("UPDATE `warehouse` SET `address` = ?, `telephone` = ?, `warehouse_code` = ?, `minimum_capacity` = ?, `minimum_temperature` = ? WHERE `id` = ?",
		warehouse.Address, warehouse.Telephone, warehouse.WareHouseCode, warehouse.MinimumCapacity, warehouse.MinimumTemperature, id)
	if err != nil {
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	if rowsAffected == 0 {
		return models.Warehouse{}, fmt.Errorf("%w: warehouse with id %d", error_message.ErrNotFound, id)
	}

	return warehouse, nil
}
