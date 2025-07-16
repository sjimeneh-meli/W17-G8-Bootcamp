package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

// Carry table and field constants / Constantes de tabla y campos de transportista
const (
	carryTable = "carriers"

	// Field groups for better maintainability / Grupos de campos para mejor mantenibilidad
	carryFields       = "`id`, `cid`, `company_name`, `address`, `telephone`, `locality_id`"
	carryInsertFields = "`cid`, `company_name`, `address`, `telephone`, `locality_id`"
)

// Carry query strings - organized by operation type / Cadenas de consulta de transportista - organizadas por tipo de operación
var (
	// INSERT queries / Consultas INSERT
	queryCreateCarry = fmt.Sprintf("INSERT INTO `%s`(%s) VALUES (?,?,?,?,?)", carryTable, carryInsertFields)

	// SELECT queries / Consultas SELECT
	queryExistsByCid = fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE `cid` = ?", carryTable)

	// Report queries / Consultas de reportes
	queryGetCarryReportsByLocality = "SELECT l.id, l.locality_name, COUNT(c.id) AS carriers_count FROM localities l LEFT JOIN carriers c ON l.id = c.locality_id WHERE l.id = ? GROUP BY l.id"
	queryGetAllCarryReports        = "SELECT l.id, l.locality_name, COUNT(c.id) AS carriers_count FROM localities l LEFT JOIN carriers c ON l.id = c.locality_id GROUP BY l.id"
)

var carryRepositoryInstance CarryRepository

// NewCarryRepository - Creates and returns a new instance of CarryRepositoryImpl using singleton pattern
// NewCarryRepository - Crea y retorna una nueva instancia de CarryRepositoryImpl usando patrón singleton
func NewCarryRepository(db *sql.DB) CarryRepository {
	if carryRepositoryInstance != nil {
		return carryRepositoryInstance
	}

	carryRepositoryInstance = &CarryRepositoryImpl{db: db}
	return carryRepositoryInstance
}

// CarryRepository - Interface defining the contract for carry repository operations
// CarryRepository - Interfaz que define el contrato para las operaciones del repositorio de transportistas
type CarryRepository interface {
	// Create - Inserts a new carry (carrier) into the database and returns the created carry with its generated ID
	// Create - Inserta un nuevo transportista en la base de datos y retorna el transportista creado con su ID generado
	Create(ctx context.Context, carry models.Carry) (models.Carry, error)

	// ExistsByCid - Checks if a carry with the given CID already exists in the database
	// ExistsByCid - Verifica si un transportista con el CID dado ya existe en la base de datos
	ExistsByCid(ctx context.Context, cid string) (bool, error)

	// GetCarryReportsByLocality - Retrieves carry reports for all localities or a specific locality by ID, showing carrier counts
	// GetCarryReportsByLocality - Obtiene reportes de transportistas para todas las localidades o una localidad específica por ID, mostrando conteos de transportistas
	GetCarryReportsByLocality(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error)
}

// CarryRepositoryImpl - Implementation of the CarryRepository interface
// CarryRepositoryImpl - Implementación de la interfaz CarryRepository
type CarryRepositoryImpl struct {
	db *sql.DB // Database connection / Conexión a la base de datos
}

// Create - Inserts a new carry (carrier) into the database and returns the created carry with its generated ID
// Create - Inserta un nuevo transportista en la base de datos y retorna el transportista creado con su ID generado
func (r *CarryRepositoryImpl) Create(ctx context.Context, carry models.Carry) (models.Carry, error) {
	// Execute insert statement with carry data / Ejecutar declaración de inserción con datos del transportista
	result, err := r.db.ExecContext(ctx, queryCreateCarry,
		carry.Cid, carry.CompanyName, carry.Address, carry.Telephone, carry.LocalityId,
	)
	if err != nil {
		return models.Carry{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	// Get the auto-generated ID and assign it to the carry / Obtener el ID autogenerado y asignarlo al transportista
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return models.Carry{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	carry.Id = int(lastInsertId)

	return carry, nil
}

// ExistsByCid - Checks if a carry with the given CID already exists in the database
// ExistsByCid - Verifica si un transportista con el CID dado ya existe en la base de datos
func (r *CarryRepositoryImpl) ExistsByCid(ctx context.Context, cid string) (bool, error) {
	// Execute query to count carries with the given CID / Ejecutar consulta para contar transportistas con el CID dado
	row := r.db.QueryRowContext(ctx, queryExistsByCid, cid)

	var count int
	err := row.Scan(&count)
	if err != nil {
		// If no rows found, CID doesn't exist (not an error) / Si no se encuentran filas, el CID no existe (no es un error)
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return count > 0, nil
}

// GetCarryReportsByLocality - Retrieves carry reports showing locality information and carrier counts
// GetCarryReportsByLocality - Obtiene reportes de transportistas mostrando información de localidad y conteos de transportistas
func (r *CarryRepositoryImpl) GetCarryReportsByLocality(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
	var rows *sql.Rows
	var err error

	// Execute different query based on whether specific locality ID is requested
	// Ejecutar consulta diferente basada en si se solicita un ID de localidad específico
	if localityID != 0 {
		// Query for specific locality using LEFT JOIN to include localities with zero carriers
		// Consulta para localidad específica usando LEFT JOIN para incluir localidades sin transportistas
		rows, err = r.db.QueryContext(ctx, queryGetCarryReportsByLocality, localityID)
	} else {
		// Query for all localities using LEFT JOIN to include localities with zero carriers
		// Consulta para todas las localidades usando LEFT JOIN para incluir localidades sin transportistas
		rows, err = r.db.QueryContext(ctx, queryGetAllCarryReports)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	defer rows.Close()

	// Iterate through all rows and scan each report into the results slice
	// Itera a través de todas las filas y escanea cada reporte en el slice de resultados
	var reports []responses.LocalityCarryReport
	for rows.Next() {
		var report responses.LocalityCarryReport
		err := rows.Scan(&report.LocalityId, &report.LocalityName, &report.CarriersCount)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
		}
		reports = append(reports, report)
	}

	// Check for any errors that occurred during iteration / Verificar si ocurrieron errores durante la iteración
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return reports, nil
}
