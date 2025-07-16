package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/database"
)

var sectionRepositoryInstance SectionRepositoryI

// GetSectionRepository - Creates and returns a new instance of sectionRepository using singleton pattern
// GetSectionRepository - Crea y retorna una nueva instancia de sectionRepository usando patrón singleton
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

// SectionRepositoryI - Interface defining the contract for section repository operations
// SectionRepositoryI - Interfaz que define el contrato para las operaciones del repositorio de secciones
type SectionRepositoryI interface {
	// GetAll - Retrieves all sections from the database
	// GetAll - Obtiene todas las secciones de la base de datos
	GetAll(ctx context.Context) ([]*models.Section, error)

	// GetByID - Retrieves a specific section by its ID from the database
	// GetByID - Obtiene una sección específica por su ID de la base de datos
	GetByID(ctx context.Context, id int) (*models.Section, error)

	// Create - Inserts a new section into the database and assigns the generated ID to the model
	// Create - Inserta una nueva sección en la base de datos y asigna el ID generado al modelo
	Create(ctx context.Context, model *models.Section) error

	// Update - Modifies an existing section in the database with the provided data
	// Update - Modifica una sección existente en la base de datos con los datos proporcionados
	Update(ctx context.Context, model *models.Section) error

	// ExistWithID - Checks if a section with the given ID exists in the database
	// ExistWithID - Verifica si una sección con el ID dado existe en la base de datos
	ExistWithID(ctx context.Context, id int) bool

	// ExistsWithSectionNumber - Checks if a section exists with the given section number, excluding a specific ID
	// ExistsWithSectionNumber - Verifica si existe una sección con el número de sección dado, excluyendo un ID específico
	ExistsWithSectionNumber(ctx context.Context, id int, sectionNumber string) bool

	// DeleteByID - Removes a section from the database by its ID
	// DeleteByID - Elimina una sección de la base de datos por su ID
	DeleteByID(ctx context.Context, id int) error
}

// sectionRepository - Implementation of SectionRepositoryI using a generic database helper
// sectionRepository - Implementación de SectionRepositoryI usando un helper genérico de base de datos
type sectionRepository struct {
	database  *sql.DB // Database connection / Conexión a la base de datos
	tablename string  // Table name for sections / Nombre de tabla para secciones
}

// GetAll - Retrieves all sections from the database and returns them as a slice of pointers
// GetAll - Obtiene todas las secciones de la base de datos y las retorna como un slice de punteros
func (r *sectionRepository) GetAll(ctx context.Context) ([]*models.Section, error) {
	// Define columns to select from the sections table / Definir columnas a seleccionar de la tabla de secciones
	columns := []string{"Id", "section_number", "current_capacity", "current_temperature", "maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id"}

	// Execute select query using generic database helper / Ejecutar consulta select usando helper genérico de base de datos
	rows, err := database.Select(ctx, r.database, r.tablename, columns, "")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sections []*models.Section

	// Iterate through all rows and scan each section into the results slice
	// Itera a través de todas las filas y escanea cada sección en el slice de resultados
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

// GetByID - Retrieves a specific section by its ID from the database
// GetByID - Obtiene una sección específica por su ID de la base de datos
func (r *sectionRepository) GetByID(ctx context.Context, id int) (*models.Section, error) {
	// Define columns to select from the sections table / Definir columnas a seleccionar de la tabla de secciones
	columns := []string{"Id", "section_number", "current_capacity", "current_temperature", "maximum_capacity", "minimum_capacity", "minimum_temperature", "product_type_id", "warehouse_id"}
	// Execute select query for specific ID using generic database helper / Ejecutar consulta select para ID específico usando helper genérico de base de datos
	row := database.SelectOne(ctx, r.database, r.tablename, columns, "Id = ?", id)

	var section models.Section

	// Scan the row into the section model / Escanear la fila en el modelo de sección
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

// Create - Inserts a new section into the database with all required fields and sets the generated ID
// Create - Inserta una nueva sección en la base de datos con todos los campos requeridos y establece el ID generado
func (r *sectionRepository) Create(ctx context.Context, model *models.Section) error {
	// Prepare data map with all section fields / Preparar mapa de datos con todos los campos de la sección
	data := make(map[any]any)
	data["section_number"] = model.SectionNumber
	data["current_capacity"] = model.CurrentCapacity
	data["current_temperature"] = model.CurrentTemperature
	data["maximum_capacity"] = model.MaximumCapacity
	data["minimum_capacity"] = model.MinimumCapacity
	data["minimum_temperature"] = model.MinimumTemperature
	data["product_type_id"] = model.ProductTypeID
	data["warehouse_id"] = model.WarehouseID

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

// Update - Modifies an existing section in the database with the provided data using a direct SQL statement
// Update - Modifica una sección existente en la base de datos con los datos proporcionados usando una declaración SQL directa
func (r *sectionRepository) Update(ctx context.Context, model *models.Section) error {
	// Build SQL update statement with all section fields / Construir declaración SQL de actualización con todos los campos de sección
	sqlStatement := fmt.Sprintf("UPDATE %s SET `section_number`=?, `current_capacity`=?, `current_temperature`=?, `maximum_capacity`=?, `minimum_capacity`=?, `minimum_temperature`=?, `product_type_id`=?, `warehouse_id`=? WHERE `Id`=?", r.tablename)
	// Execute update statement with section data / Ejecutar declaración de actualización con datos de sección
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

// ExistWithID - Checks if a section with the given ID exists in the database
// ExistWithID - Verifica si una sección con el ID dado existe en la base de datos
func (r *sectionRepository) ExistWithID(ctx context.Context, id int) bool {
	// Query to count sections with the given ID / Consulta para contar secciones con el ID dado
	row := database.SelectOne(ctx, r.database, r.tablename, []string{"COUNT(Id)"}, "Id = ?", id)
	var count int
	if err := row.Scan(&count); err != nil {
		// Return true on error to be safe for validation / Retornar true en caso de error para ser seguro en la validación
		return true
	}

	return count != 0
}

// ExistsWithSectionNumber - Checks if a section exists with the given section number, excluding a specific ID for update scenarios
// ExistsWithSectionNumber - Verifica si existe una sección con el número de sección dado, excluyendo un ID específico para escenarios de actualización
func (r *sectionRepository) ExistsWithSectionNumber(ctx context.Context, id int, sectionNumber string) bool {
	// Query to count sections with same section number but different ID / Consulta para contar secciones con el mismo número pero diferente ID
	row := database.SelectOne(ctx, r.database, r.tablename, []string{"COUNT(Id)"}, "section_number = ? AND Id <> ?", sectionNumber, id)
	var count string
	if err := row.Scan(&count); err != nil {
		// Return true on error to be safe for validation / Retornar true en caso de error para ser seguro en la validación
		return true
	}

	return count != "0"
}

// DeleteByID - Removes a section from the database by its ID using the generic database helper
// DeleteByID - Elimina una sección de la base de datos por su ID usando el helper genérico de base de datos
func (r *sectionRepository) DeleteByID(ctx context.Context, id int) error {
	// Execute delete operation using generic database helper / Ejecutar operación de eliminación usando helper genérico de base de datos
	_, err := database.Delete(ctx, r.database, r.tablename, "Id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
