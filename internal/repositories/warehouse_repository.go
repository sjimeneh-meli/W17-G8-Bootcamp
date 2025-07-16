package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

// Warehouse table and field constants / Constantes de tabla y campos de almacén
const (
	warehouseTable = "warehouse"

	// Field groups for better maintainability / Grupos de campos para mejor mantenibilidad
	warehouseFields       = "`id`, `address`, `telephone`, `warehouse_code`, `minimum_capacity`, `minimum_temperature`"
	warehouseInsertFields = "`address`, `telephone`, `warehouse_code`, `minimum_capacity`, `minimum_temperature`, `locality_id`"
	warehouseUpdateFields = "`address` = ?, `telephone` = ?, `warehouse_code` = ?, `minimum_capacity` = ?, `minimum_temperature` = ?"
)

// Warehouse query strings - organized by operation type / Cadenas de consulta de almacén - organizadas por tipo de operación
var (
	// SELECT queries / Consultas SELECT
	queryGetAllWarehouses = fmt.Sprintf("SELECT %s FROM `%s`", warehouseFields, warehouseTable)
	queryGetWarehouseById = fmt.Sprintf("SELECT %s FROM `%s` WHERE `id` = ?", warehouseFields, warehouseTable)
	queryExistsByCode     = fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE `warehouse_code` = ?", warehouseTable)

	// INSERT queries / Consultas INSERT
	queryCreateWarehouse = fmt.Sprintf("INSERT INTO `%s`(%s) VALUES (?,?,?,?,?,?)", warehouseTable, warehouseInsertFields)

	// UPDATE queries / Consultas UPDATE
	queryUpdateWarehouse = fmt.Sprintf("UPDATE `%s` SET %s WHERE `id` = ?", warehouseTable, warehouseUpdateFields)

	// DELETE queries / Consultas DELETE
	queryDeleteWarehouse = fmt.Sprintf("DELETE FROM `%s` WHERE `id` = ?", warehouseTable)

	warehouseRepositoryInstance WarehouseRepository
)

// NewWarehouseRepository - Creates and returns a new instance of WarehouseRepositoryImpl using singleton pattern
// NewWarehouseRepository - Crea y retorna una nueva instancia de WarehouseRepositoryImpl usando patrón singleton
func NewWarehouseRepository(db *sql.DB) WarehouseRepository {
	if warehouseRepositoryInstance != nil {
		return warehouseRepositoryInstance
	}
	warehouseRepositoryInstance = &WarehouseRepositoryImpl{db: db}
	return warehouseRepositoryInstance
}

// WarehouseRepository - Interface defining the contract for warehouse repository operations
// WarehouseRepository - Interfaz que define el contrato para las operaciones del repositorio de almacenes
type WarehouseRepository interface {
	// GetAll - Retrieves all warehouses from the database
	// GetAll - Obtiene todos los almacenes de la base de datos
	GetAll(ctx context.Context) ([]models.Warehouse, error)

	// Create - Inserts a new warehouse into the database and returns the created warehouse with its generated ID
	// Create - Inserta un nuevo almacén en la base de datos y retorna el almacén creado con su ID generado
	Create(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error)

	// ExistsByCode - Checks if a warehouse with the given warehouse code already exists in the database
	// ExistsByCode - Verifica si un almacén con el código de almacén dado ya existe en la base de datos
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// GetById - Retrieves a specific warehouse by its ID from the database
	// GetById - Obtiene un almacén específico por su ID de la base de datos
	GetById(ctx context.Context, id int) (models.Warehouse, error)

	// Delete - Removes a warehouse from the database by its ID
	// Delete - Elimina un almacén de la base de datos por su ID
	Delete(ctx context.Context, id int) error

	// Update - Modifies an existing warehouse in the database and returns the updated warehouse
	// Update - Modifica un almacén existente en la base de datos y retorna el almacén actualizado
	Update(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error)
}

// WarehouseRepositoryImpl - Implementation of the WarehouseRepository interface
// WarehouseRepositoryImpl - Implementación de la interfaz WarehouseRepository
type WarehouseRepositoryImpl struct {
	db *sql.DB // Database connection / Conexión a la base de datos
}

// GetAll - Retrieves all warehouses from the database and returns them as a slice
// GetAll - Obtiene todos los almacenes de la base de datos y los retorna como un slice
func (r *WarehouseRepositoryImpl) GetAll(ctx context.Context) ([]models.Warehouse, error) {
	// Execute query to select all warehouse fields / Ejecutar consulta para seleccionar todos los campos del almacén
	rows, err := r.db.QueryContext(ctx, queryGetAllWarehouses)
	if err != nil {
		fmt.Printf("warehouse: %v\n", err.Error())
		return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	// Iterate through all rows and scan each warehouse into the results slice
	// Itera a través de todas las filas y escanea cada almacén en el slice de resultados
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

// Create - Inserts a new warehouse into the database and returns the created warehouse with its generated ID
// Create - Inserta un nuevo almacén en la base de datos y retorna el almacén creado con su ID generado
func (r *WarehouseRepositoryImpl) Create(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error) {
	// Execute insert statement with warehouse data / Ejecutar declaración de inserción con datos del almacén
	result, err := r.db.ExecContext(ctx, queryCreateWarehouse,
		warehouse.Address, warehouse.Telephone, warehouse.WareHouseCode, warehouse.MinimumCapacity, warehouse.MinimumTemperature, warehouse.LocalityId)
	if err != nil {
		fmt.Printf("warehouse: %v\n", err.Error())
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	// Get the auto-generated ID and assign it to the warehouse / Obtener el ID autogenerado y asignarlo al almacén
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	warehouse.Id = int(lastInsertId)

	return warehouse, nil
}

// ExistsByCode - Checks if a warehouse with the given warehouse code already exists in the database
// ExistsByCode - Verifica si un almacén con el código de almacén dado ya existe en la base de datos
func (r *WarehouseRepositoryImpl) ExistsByCode(ctx context.Context, code string) (bool, error) {
	// Execute query to count warehouses with the given code / Ejecutar consulta para contar almacenes con el código dado
	row := r.db.QueryRowContext(ctx, queryExistsByCode, code)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return count > 0, nil
}

// GetById - Retrieves a specific warehouse by its ID from the database
// GetById - Obtiene un almacén específico por su ID de la base de datos
func (r *WarehouseRepositoryImpl) GetById(ctx context.Context, id int) (models.Warehouse, error) {
	// Execute query to select warehouse by specific ID / Ejecutar consulta para seleccionar almacén por ID específico
	row := r.db.QueryRowContext(ctx, queryGetWarehouseById, id)

	var warehouse models.Warehouse
	err := row.Scan(&warehouse.Id, &warehouse.Address, &warehouse.Telephone, &warehouse.WareHouseCode, &warehouse.MinimumCapacity, &warehouse.MinimumTemperature)
	if err != nil {
		// Handle case when no warehouse is found / Manejar el caso cuando no se encuentra ningún almacén
		if err == sql.ErrNoRows {
			return models.Warehouse{}, fmt.Errorf("%w: warehouse with id %d", error_message.ErrNotFound, id)
		}
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return warehouse, nil
}

// Delete - Removes a warehouse from the database by its ID
// Delete - Elimina un almacén de la base de datos por su ID
func (r *WarehouseRepositoryImpl) Delete(ctx context.Context, id int) error {
	// Execute delete statement for the specified warehouse ID / Ejecutar declaración de eliminación para el ID del almacén especificado
	_, err := r.db.ExecContext(ctx, queryDeleteWarehouse, id)
	if err != nil {
		// Handle case when no warehouse is found / Manejar el caso cuando no se encuentra ningún almacén
		if err == sql.ErrNoRows {
			return fmt.Errorf("%w: warehouse with id %d", error_message.ErrNotFound, id)
		}
		fmt.Printf("warehouse: %v\n", err.Error())
		return fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return nil
}

// Update - Modifies an existing warehouse in the database and returns the updated warehouse
// Update - Modifica un almacén existente en la base de datos y retorna el almacén actualizado
func (r *WarehouseRepositoryImpl) Update(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error) {
	// Execute update statement with warehouse data / Ejecutar declaración de actualización con datos del almacén
	result, err := r.db.ExecContext(ctx, queryUpdateWarehouse,
		warehouse.Address, warehouse.Telephone, warehouse.WareHouseCode, warehouse.MinimumCapacity, warehouse.MinimumTemperature, id)
	if err != nil {
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	// Check if any rows were affected to confirm update / Verificar si alguna fila fue afectada para confirmar la actualización
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Warehouse{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	// If no rows affected, warehouse doesn't exist / Si ninguna fila fue afectada, el almacén no existe
	if rowsAffected == 0 {
		return models.Warehouse{}, fmt.Errorf("%w: warehouse with id %d", error_message.ErrNotFound, id)
	}

	return warehouse, nil
}
