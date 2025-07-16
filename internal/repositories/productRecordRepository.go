package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

var productRecordInstance IProductRecordRepository

// NewProductRecordRepository - Constructor function that creates a new repository instance
// NewProductRecordRepository - Función constructora que crea una nueva instancia del repositorio
func NewProductRecordRepository(db *sql.DB) IProductRecordRepository {
	if productRecordInstance != nil {
		return productRecordInstance
	}

	productRecordInstance = &productRecordRepository{DB: db}
	return productRecordInstance
}

// productRecordRepository - Repository implementation for product records operations
// productRecordRepository - Implementación del repositorio para operaciones de registros de productos
type productRecordRepository struct {
	DB *sql.DB // Database connection / Conexión a la base de datos
}

// IProductRecordRepository - Interface defining the contract for product records repository
// IProductRecordRepository - Interfaz que define el contrato para el repositorio de registros de productos
type IProductRecordRepository interface {
	// Create - Creates a new product record in the database
	// Create - Crea un nuevo registro de producto en la base de datos
	Create(ctx context.Context, pr *models.ProductRecord) (*models.ProductRecord, error)

	// GetReportByIdProduct - Generates a report for a specific product with records count
	// GetReportByIdProduct - Genera un reporte para un producto específico con el conteo de registros
	GetReportByIdProduct(ctx context.Context, id int64) (*models.ProductRecordReport, error)

	// GetReport - Generates a report showing all products information and the count of their records
	// GetReport - Genera un reporte que muestra la información de todos los productos y el conteo de sus registros
	GetReport(ctx context.Context) ([]*models.ProductRecordReport, error)

	// ExistProductRecordByID - Checks if a product record exists in the database
	// ExistProductRecordByID - Verifica si un registro de producto existe en la base de datos
	ExistProductRecordByID(ctx context.Context, id int64) bool
}

// Create - Inserts a new product record into the database and returns the created record with its ID
// Create - Inserta un nuevo registro de producto en la base de datos y retorna el registro creado con su ID
func (prr *productRecordRepository) Create(ctx context.Context, pr *models.ProductRecord) (*models.ProductRecord, error) {
	// SQL query to insert a new product record / Consulta SQL para insertar un nuevo registro de producto
	query := "INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES (?, ?, ?, ?)"

	// Execute the insert query with the provided context and product record data
	// Ejecuta la consulta de inserción con el contexto proporcionado y los datos del registro de producto
	result, err := prr.DB.ExecContext(ctx, query, pr.LastUpdateDate, pr.PurchasePrice, pr.SalePrice, pr.ProductID)
	if err != nil {
		return nil, fmt.Errorf("error to create product record: %w", err)
	}

	// Get the auto-generated ID from the database / Obtiene el ID autogenerado de la base de datos
	productRecordId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error to get last insert id: %w", err)
	}

	// Set the generated ID to the product record / Asigna el ID generado al registro de producto
	pr.ID = int(productRecordId)

	return pr, err
}

// GetReportByIdProduct - Generates a report showing product information and the count of its records
// GetReportByIdProduct - Genera un reporte que muestra la información del producto y el conteo de sus registros
func (prr *productRecordRepository) GetReportByIdProduct(ctx context.Context, id int64) (*models.ProductRecordReport, error) {
	// Complex SQL query using LEFT JOIN to get product info and count records
	// Uses count(*) since we're filtering by WHERE clause, ensuring at least one product exists
	// Consulta SQL compleja usando LEFT JOIN para obtener información del producto y contar registros
	// Usa count(*) ya que estamos filtrando por WHERE clause, asegurando que al menos un producto existe
	query := `
	SELECT p.id as product_id, p.description, count(*) as records_count
	FROM 
		products as p
	LEFT JOIN
		product_records as pr
	ON pr.product_id = p.id
	WHERE p.id = ?
	GROUP BY p.id, p.description`

	var productRecordReport models.ProductRecordReport

	// Execute query and scan results into the report struct
	// Ejecuta la consulta y escanea los resultados en la estructura del reporte
	err := prr.DB.QueryRowContext(ctx, query, id).Scan(
		&productRecordReport.ProductId,
		&productRecordReport.Description,
		&productRecordReport.RecordsCount,
	)

	if err != nil {
		// Handle case when no product is found / Maneja el caso cuando no se encuentra ningún producto
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error_message.ErrDependencyNotFound
		}
		return nil, fmt.Errorf("error to scan product record report %w", err)
	}

	return &productRecordReport, nil
}

// GetReport - Generates a report showing products information and the count of its records
// GetReport - Genera un reporte que muestra la información de los productos y el conteo de sus registros
func (prr *productRecordRepository) GetReport(ctx context.Context) (reports []*models.ProductRecordReport, err error) {

	// Complex SQL query using LEFT JOIN to get product info and count records
	// Uses const for immutability and count(pr.product_id) to count only non-null values (actual records)
	// Consulta SQL compleja usando LEFT JOIN para obtener información del producto y contar registros
	// Usa const para inmutabilidad y count(pr.product_id) para contar solo valores no nulos (registros reales)
	const query = `
        SELECT
            p.id as product_id,
            p.description,
            count(pr.product_id) as records_count
        FROM
            products p
        LEFT JOIN
            product_records pr ON pr.product_id = p.id
        GROUP BY
            p.id, p.description`

	// Prepare statement for better performance and security against SQL injection
	// Prepara la consulta para mejor rendimiento y seguridad contra inyección SQL
	stmt, err := prr.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare product report query: %w", err)
	}
	// Ensure statement is closed to prevent resource leaks
	// Asegura que la consulta se cierre para prevenir fugas de recursos
	defer stmt.Close()

	// Execute query and scan results into the report struct
	// Ejecuta la consulta y escanea los resultados en la estructura del reporte
	rows, err := stmt.QueryContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to execute product report query: %w", err)
	}
	// Ensure rows are closed to prevent resource leaks
	// Asegura que las filas se cierren para prevenir fugas de recursos
	defer rows.Close()

	// Initialize empty slice to hold the reports
	// Inicializa un slice vacío para almacenar los reportes
	reports = []*models.ProductRecordReport{}

	// Iterate through all rows and scan each one into a ProductRecordReport struct
	// Itera a través de todas las filas y escanea cada una en una estructura ProductRecordReport
	for rows.Next() {
		var pr models.ProductRecordReport

		// Scan current row values into the struct fields
		// Escanea los valores de la fila actual en los campos de la estructura
		if err := rows.Scan(&pr.ProductId, &pr.Description, &pr.RecordsCount); err != nil {
			return nil, fmt.Errorf("failed to scan product record row: %w", err)
		}
		// Append the scanned report to the results slice
		// Agrega el reporte escaneado al slice de resultados
		reports = append(reports, &pr)
	}

	// Check for any errors that occurred during iteration
	// Verifica si ocurrieron errores durante la iteración
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product record rows: %w", err)
	}

	return reports, nil
}

// ExistProductRecordByID - Validates if a product record exists in the database by its ID
// ExistProductRecordByID - Valida si un registro de producto existe en la base de datos por su ID
func (prr *productRecordRepository) ExistProductRecordByID(ctx context.Context, id int64) bool {
	// Simple query to check product record existence using LIMIT 1 for efficiency
	// Consulta simple para verificar la existencia del registro de producto usando LIMIT 1 por eficiencia
	query := "SELECT 1 FROM product_records WHERE id = ? LIMIT 1;"

	// Use int64 to match the database return type for consistency
	// Usa int64 para coincidir con el tipo de retorno de la base de datos por consistencia
	var exists int64
	err := prr.DB.QueryRowContext(ctx, query, id).Scan(&exists)

	if err != nil {
		// If no rows found, product record doesn't exist (not an error)
		// Si no se encuentran filas, el registro del producto no existe (no es un error)
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
		// Any other error is a real database error / Cualquier otro error es un error real de base de datos
		return false
	}

	// If we reach here, product record exists / Si llegamos aquí, el registro de producto existe
	return true
}
