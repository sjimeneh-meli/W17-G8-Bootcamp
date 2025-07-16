package repositories

import (
	"context"
	"database/sql"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/database"
)

var productBatchRepositoryInstance ProductBatchRepositoryI

// GetProductBatchRepository - Creates and returns a new instance of productBatchRepository using singleton pattern
// GetProductBatchRepository - Crea y retorna una nueva instancia de productBatchRepository usando patrón singleton
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

// ProductBatchRepositoryI - Interface defining the contract for product batch repository operations
// ProductBatchRepositoryI - Interfaz que define el contrato para las operaciones del repositorio de lotes de productos
type ProductBatchRepositoryI interface {
	// Create - Inserts a new product batch into the database and assigns the generated ID to the model
	// Create - Inserta un nuevo lote de producto en la base de datos y asigna el ID generado al modelo
	Create(ctx context.Context, model *models.ProductBatch) error

	// GetProductQuantityBySectionId - Retrieves the total quantity of products in a specific section
	// GetProductQuantityBySectionId - Obtiene la cantidad total de productos en una sección específica
	GetProductQuantityBySectionId(ctx context.Context, id int) int

	// ExistsWithBatchNumber - Checks if a product batch exists with the given batch number, excluding a specific ID
	// ExistsWithBatchNumber - Verifica si existe un lote de producto con el número de lote dado, excluyendo un ID específico
	ExistsWithBatchNumber(ctx context.Context, id int, batchNumber string) bool
}

// productBatchRepository - Implementation of ProductBatchRepositoryI using a generic database helper
// productBatchRepository - Implementación de ProductBatchRepositoryI usando un helper genérico de base de datos
type productBatchRepository struct {
	database  *sql.DB // Database connection / Conexión a la base de datos
	tablename string  // Table name for product batches / Nombre de tabla para lotes de productos
}

// Create - Inserts a new product batch into the database with all required fields and sets the generated ID
// Create - Inserta un nuevo lote de producto en la base de datos con todos los campos requeridos y establece el ID generado
func (r *productBatchRepository) Create(ctx context.Context, model *models.ProductBatch) error {
	// Prepare data map with all product batch fields / Preparar mapa de datos con todos los campos del lote de producto
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

	// Execute insert operation using generic database helper / Ejecutar operación de inserción usando helper genérico de base de datos
	result, err := database.Insert(ctx, r.database, r.tablename, data)

	if err != nil {
		return err
	}

	// Get the auto-generated ID and assign it to the model / Obtener el ID autogenerado y asignarlo al modelo
	newID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	model.Id = int(newID)
	return nil
}

// ExistsWithBatchNumber - Checks if a product batch exists with the given batch number, excluding a specific ID for update scenarios
// ExistsWithBatchNumber - Verifica si existe un lote de producto con el número de lote dado, excluyendo un ID específico para escenarios de actualización
func (r *productBatchRepository) ExistsWithBatchNumber(ctx context.Context, id int, batchNumber string) bool {
	// Query to count batches with same batch number but different ID / Consulta para contar lotes con el mismo número pero diferente ID
	row := database.SelectOne(ctx, r.database, r.tablename, []string{"COUNT(Id)"}, "batch_number = ? AND Id <> ?", batchNumber, id)
	var count string
	if err := row.Scan(&count); err != nil {
		// Return true on error to be safe for validation / Retornar true en caso de error para ser seguro en la validación
		return true
	}

	return count != "0"
}

// GetProductQuantityBySectionId - Calculates and returns the total current quantity of all product batches in a specific section
// GetProductQuantityBySectionId - Calcula y retorna la cantidad total actual de todos los lotes de productos en una sección específica
func (r *productBatchRepository) GetProductQuantityBySectionId(ctx context.Context, id int) int {
	// Query to sum current quantities for all batches in the section / Consulta para sumar cantidades actuales de todos los lotes en la sección
	row := database.SelectOne(ctx, r.database, r.tablename, []string{"SUM(current_quantity)"}, "section_id = ?", id)
	var quantity int
	if err := row.Scan(&quantity); err != nil {
		// Return existing quantity on scan error / Retornar cantidad existente en caso de error de escaneo
		return quantity
	}

	return quantity
}
