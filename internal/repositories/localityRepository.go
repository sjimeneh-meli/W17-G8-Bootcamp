package repositories

import (
	"context"
	"database/sql"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

var localityRepositoryInstance LocalityRepository

// NewSQLLocalityRepository - Creates and returns a new instance of SQLLocalityRepository using singleton pattern
// NewSQLLocalityRepository - Crea y retorna una nueva instancia de SQLLocalityRepository usando patrón singleton
func NewSQLLocalityRepository(db *sql.DB) LocalityRepository {
	if localityRepositoryInstance != nil {
		return localityRepositoryInstance
	}

	localityRepositoryInstance = &SQLLocalityRepository{db: db}
	return localityRepositoryInstance
}

// LocalityRepository - Interface defining the contract for locality repository operations
// LocalityRepository - Interfaz que define el contrato para las operaciones del repositorio de localidades
type LocalityRepository interface {
	// Save - Creates a new locality in the database, handling country and province relationships
	// Save - Crea una nueva localidad en la base de datos, manejando las relaciones de país y provincia
	Save(ctx context.Context, locality models.Locality) (models.Locality, error)

	// GetSellerReports - Retrieves seller reports for all localities or a specific locality by ID
	// GetSellerReports - Obtiene reportes de vendedores para todas las localidades o una localidad específica por ID
	GetSellerReports(ctx context.Context, localityID int) ([]responses.LocalitySellerReport, error)

	// ExistById - Checks if a locality with the given ID exists in the database
	// ExistById - Verifica si una localidad con el ID dado existe en la base de datos
	ExistById(ctx context.Context, localityID int) (bool, error)
}

// SQLLocalityRepository - SQL implementation of the LocalityRepository interface
// SQLLocalityRepository - Implementación SQL de la interfaz LocalityRepository
type SQLLocalityRepository struct {
	db *sql.DB // Database connection / Conexión a la base de datos
}

// Save - Creates a new locality in the database with automatic country and province management
// Save - Crea una nueva localidad en la base de datos con manejo automático de país y provincia
func (r *SQLLocalityRepository) Save(ctx context.Context, locality models.Locality) (models.Locality, error) {
	// 1. Find or insert the country / 1. Buscar o insertar el país
	var countryID int
	err := r.db.QueryRowContext(ctx, "SELECT id FROM countries WHERE country_name = ?", locality.CountryName).Scan(&countryID)
	if err == sql.ErrNoRows {
		// Create new country if it doesn't exist / Crear nuevo país si no existe
		res, err := r.db.ExecContext(ctx, "INSERT INTO countries (country_name) VALUES (?)", locality.CountryName)
		if err != nil {
			return models.Locality{}, error_message.ErrQuery
		}
		lastID, _ := res.LastInsertId()
		countryID = int(lastID)
	} else if err != nil {
		return models.Locality{}, error_message.ErrQuery
	}

	// 2. Find or insert the province / 2. Buscar o insertar la provincia
	var provinceID int
	err = r.db.QueryRowContext(ctx, "SELECT id FROM provinces WHERE province_name = ? AND id_country_fk = ?", locality.ProvinceName, countryID).Scan(&provinceID)
	if err == sql.ErrNoRows {
		// Create new province if it doesn't exist / Crear nueva provincia si no existe
		res, err := r.db.ExecContext(ctx, "INSERT INTO provinces (province_name, id_country_fk) VALUES (?, ?)", locality.ProvinceName, countryID)
		if err != nil {
			return models.Locality{}, error_message.ErrQuery
		}
		lastID, _ := res.LastInsertId()
		provinceID = int(lastID)
	} else if err != nil {
		return models.Locality{}, error_message.ErrQuery
	}

	// 3. Check if locality already exists with same name and province / 3. Verificar si ya existe una localidad con ese nombre y esa provincia
	var exists bool
	err = r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM localities WHERE locality_name = ? AND province_id = ?
		)
	`, locality.LocalityName, provinceID).Scan(&exists)
	if err != nil {
		return models.Locality{}, error_message.ErrQuery
	}
	if exists {
		return models.Locality{}, error_message.ErrAlreadyExists
	}

	// 4. Insert new locality with auto-generated ID / 4. Insertar nueva localidad (el ID será auto-generado)
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO localities (locality_name, province_id) VALUES (?, ?)",
		locality.LocalityName, provinceID,
	)

	if err != nil {
		return models.Locality{}, error_message.ErrQuery
	}

	// Get the auto-generated ID and assign it to the locality / Obtener el ID autogenerado y asignarlo a la localidad
	lastID, err := res.LastInsertId()
	if err != nil {
		return models.Locality{}, error_message.ErrQuery
	}
	locality.Id = int(lastID)

	return locality, nil
}

// GetSellerReports - Retrieves seller reports showing locality information and seller counts
// GetSellerReports - Obtiene reportes de vendedores mostrando información de localidad y conteos de vendedores
func (r *SQLLocalityRepository) GetSellerReports(ctx context.Context, localityID int) ([]responses.LocalitySellerReport, error) {
	// If specific ID provided, verify locality exists / Si se especifica un ID, verificar que exista
	if localityID != 0 {
		var exists bool
		err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM localities WHERE id = ?)", localityID).Scan(&exists)
		if err != nil {
			return nil, error_message.ErrFailedCheckingExistence
		}
		if !exists {
			return nil, error_message.ErrNotFound
		}
	}

	var (
		rows *sql.Rows
		err  error
	)

	// Execute different query based on whether specific locality ID is requested
	// Ejecutar consulta diferente basada en si se solicita un ID de localidad específico
	if localityID != 0 {
		// Query for specific locality using LEFT JOIN to include localities with zero sellers
		// Consulta para localidad específica usando LEFT JOIN para incluir localidades sin vendedores
		rows, err = r.db.QueryContext(ctx, `
			SELECT l.id, l.locality_name, COUNT(s.id)
			FROM localities l
			LEFT JOIN sellers s ON s.locality_id = l.id
			WHERE l.id = ?
			GROUP BY l.id, l.locality_name
		`, localityID)
	} else {
		// Query for all localities using LEFT JOIN to include localities with zero sellers
		// Consulta para todas las localidades usando LEFT JOIN para incluir localidades sin vendedores
		rows, err = r.db.QueryContext(ctx, `
			SELECT l.id, l.locality_name, COUNT(s.id)
			FROM localities l
			LEFT JOIN sellers s ON s.locality_id = l.id
			GROUP BY l.id, l.locality_name
		`)
	}

	if err != nil {
		return nil, error_message.ErrQueryingReport
	}
	defer rows.Close()

	// Iterate through all rows and scan each report into the results slice
	// Itera a través de todas las filas y escanea cada reporte en el slice de resultados
	var reports []responses.LocalitySellerReport
	for rows.Next() {
		var r responses.LocalitySellerReport
		if err := rows.Scan(&r.LocalityID, &r.LocalityName, &r.SellerCount); err != nil {
			return nil, error_message.ErrFailedToScan
		}
		reports = append(reports, r)
	}

	return reports, nil
}

// ExistById - Checks if a locality with the given ID exists in the database
// ExistById - Verifica si una localidad con el ID dado existe en la base de datos
func (r *SQLLocalityRepository) ExistById(ctx context.Context, localityID int) (bool, error) {
	// Simple query using EXISTS for efficient existence check / Consulta simple usando EXISTS para verificación eficiente de existencia
	row := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM localities WHERE id = ?)", localityID)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, error_message.ErrQuery
	}
	return exists, nil
}
