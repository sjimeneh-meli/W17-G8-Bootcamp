package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

// Warehouse table and field constants
const (
	warehouseTable = "warehouse"

	// Field groups for better maintainability
	warehouseFields       = "`id`, `address`, `telephone`, `warehouse_code`, `minimum_capacity`, `minimum_temperature`"
	warehouseInsertFields = "`address`, `telephone`, `warehouse_code`, `minimum_capacity`, `minimum_temperature`, `locality_id`"
	warehouseUpdateFields = "`address` = ?, `telephone` = ?, `warehouse_code` = ?, `minimum_capacity` = ?, `minimum_temperature` = ?"
)

// Warehouse query strings - organized by operation type
var (
	// SELECT queries
	queryGetAllWarehouses = fmt.Sprintf("SELECT %s FROM `%s`", warehouseFields, warehouseTable)
	queryGetWarehouseById = fmt.Sprintf("SELECT %s FROM `%s` WHERE `id` = ?", warehouseFields, warehouseTable)
	queryExistsByCode     = fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE `warehouse_code` = ?", warehouseTable)

	// INSERT queries
	queryCreateWarehouse = fmt.Sprintf("INSERT INTO `%s`(%s) VALUES (?,?,?,?,?,?)", warehouseTable, warehouseInsertFields)

	// UPDATE queries
	queryUpdateWarehouse = fmt.Sprintf("UPDATE `%s` SET %s WHERE `id` = ?", warehouseTable, warehouseUpdateFields)

	// DELETE queries
	queryDeleteWarehouse = fmt.Sprintf("DELETE FROM `%s` WHERE `id` = ?", warehouseTable)

	warehouseRepositoryInstance WarehouseRepository
)

func NewWarehouseRepository(db *sql.DB) WarehouseRepository {
	if warehouseRepositoryInstance != nil {
		return warehouseRepositoryInstance
	}
	warehouseRepositoryInstance = &WarehouseRepositoryImpl{db: db}
	return warehouseRepositoryInstance
}

type WarehouseRepository interface {
	GetAll(ctx context.Context) ([]models.Warehouse, error)
	Create(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	GetById(ctx context.Context, id int) (models.Warehouse, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error)
}

type WarehouseRepositoryImpl struct {
	db *sql.DB
}

func (r *WarehouseRepositoryImpl) GetAll(ctx context.Context) ([]models.Warehouse, error) {
	//1. Correr la query
	rows, err := r.db.QueryContext(ctx, queryGetAllWarehouses)
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

func (r *WarehouseRepositoryImpl) Create(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error) {
	result, err := r.db.ExecContext(ctx, queryCreateWarehouse,
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

func (r *WarehouseRepositoryImpl) ExistsByCode(ctx context.Context, code string) (bool, error) {
	row := r.db.QueryRowContext(ctx, queryExistsByCode, code)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return count > 0, nil
}

func (r *WarehouseRepositoryImpl) GetById(ctx context.Context, id int) (models.Warehouse, error) {
	row := r.db.QueryRowContext(ctx, queryGetWarehouseById, id)

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

func (r *WarehouseRepositoryImpl) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, queryDeleteWarehouse, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%w: warehouse with id %d", error_message.ErrNotFound, id)
		}
		fmt.Printf("warehouse: %v\n", err.Error())
		return fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return nil
}

func (r *WarehouseRepositoryImpl) Update(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error) {
	result, err := r.db.ExecContext(ctx, queryUpdateWarehouse,
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
