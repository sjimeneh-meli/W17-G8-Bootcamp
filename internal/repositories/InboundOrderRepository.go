package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

var inboundOrderRepositoryInstance InboundOrderRepositoryI

// GetNewInboundOrderMySQLRepository - Creates and returns a new instance of MySqlInboundOrderRepository using singleton pattern
// GetNewInboundOrderMySQLRepository - Crea y retorna una nueva instancia de MySqlInboundOrderRepository usando patrón singleton
func GetNewInboundOrderMySQLRepository(db *sql.DB) InboundOrderRepositoryI {
	if inboundOrderRepositoryInstance != nil {
		return inboundOrderRepositoryInstance
	}
	inboundOrderRepositoryInstance = &MySqlInboundOrderRepository{
		db: db,
	}
	return inboundOrderRepositoryInstance
}

// InboundOrderRepositoryI - Interface defining the contract for inbound order repository operations
// InboundOrderRepositoryI - Interfaz que define el contrato para las operaciones del repositorio de órdenes de entrada
type InboundOrderRepositoryI interface {
	// GetAllInboundOrdersReports - Retrieves inbound order reports for all employees with their order counts
	// GetAllInboundOrdersReports - Obtiene reportes de órdenes de entrada para todos los empleados con sus conteos de órdenes
	GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error)

	// GetInboundOrdersReportByEmployeeId - Retrieves an inbound order report for a specific employee ID
	// GetInboundOrdersReportByEmployeeId - Obtiene un reporte de órdenes de entrada para un ID de empleado específico
	GetInboundOrdersReportByEmployeeId(ctx context.Context, employeeId int) (models.InboundOrderReport, error)

	// Create - Inserts a new inbound order into the database and returns the created order with its generated ID
	// Create - Inserta una nueva orden de entrada en la base de datos y retorna la orden creada con su ID generado
	Create(ctx context.Context, inbound models.InboundOrder) (models.InboundOrder, error)

	// ExistsByOrderNumber - Checks if an inbound order with the given order number already exists in the database
	// ExistsByOrderNumber - Verifica si una orden de entrada con el número de orden dado ya existe en la base de datos
	ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error)
}

// MySqlInboundOrderRepository - MySQL implementation of the InboundOrderRepositoryI interface
// MySqlInboundOrderRepository - Implementación MySQL de la interfaz InboundOrderRepositoryI
type MySqlInboundOrderRepository struct {
	db *sql.DB // Database connection / Conexión a la base de datos
}

// GetAllInboundOrdersReports - Retrieves inbound order reports for all employees showing employee info and their order counts
// GetAllInboundOrdersReports - Obtiene reportes de órdenes de entrada para todos los empleados mostrando información del empleado y sus conteos de órdenes
func (r *MySqlInboundOrderRepository) GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error) {
	reports := []models.InboundOrderReport{}

	// Complex SQL query using INNER JOIN to get employee info and count their inbound orders
	// Consulta SQL compleja usando INNER JOIN para obtener información del empleado y contar sus órdenes de entrada
	query := `
	SELECT e.id, e.id_card_number, e.first_name, e.last_name, COUNT(io.id) AS inbound_orders_count
	FROM employees e
	INNER JOIN inbound_orders io ON io.employee_id = e.id
	GROUP BY e.id
	ORDER BY e.id;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return []models.InboundOrderReport{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	defer rows.Close()

	// Iterate through all rows and scan each report into the results slice
	// Itera a través de todas las filas y escanea cada reporte en el slice de resultados
	for rows.Next() {
		var report models.InboundOrderReport
		err := rows.Scan(&report.Id, &report.IdCardNumber, &report.FirstName, &report.LastName, &report.InboundOrderCount)
		if err != nil {
			return []models.InboundOrderReport{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
		}
		reports = append(reports, report)
	}

	return reports, nil
}

// GetInboundOrdersReportByEmployeeId - Retrieves an inbound order report for a specific employee showing their info and order count
// GetInboundOrdersReportByEmployeeId - Obtiene un reporte de órdenes de entrada para un empleado específico mostrando su información y conteo de órdenes
func (r *MySqlInboundOrderRepository) GetInboundOrdersReportByEmployeeId(ctx context.Context, employeeId int) (models.InboundOrderReport, error) {
	report := models.InboundOrderReport{}

	// Complex SQL query using INNER JOIN to get specific employee info and count their inbound orders
	// Consulta SQL compleja usando INNER JOIN para obtener información específica del empleado y contar sus órdenes de entrada
	query := `
	SELECT e.id, e.id_card_number, e.first_name, e.last_name, COUNT(io.id) AS inbound_orders_count
	FROM employees e
	INNER JOIN inbound_orders io ON io.employee_id = e.id
	WHERE e.id = ?
	GROUP BY e.id;
	`

	row := r.db.QueryRowContext(ctx, query, employeeId)
	err := row.Scan(&report.Id, &report.IdCardNumber, &report.FirstName, &report.LastName, &report.InboundOrderCount)

	if err != nil {
		// Handle case when no employee is found / Maneja el caso cuando no se encuentra ningún empleado
		if errors.Is(err, sql.ErrNoRows) {
			return models.InboundOrderReport{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "employee with Id", employeeId, "doesn't exist.")
		}
		return models.InboundOrderReport{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	return report, nil
}

// Create - Inserts a new inbound order into the MySQL database and returns the created order with its generated ID
// Create - Inserta una nueva orden de entrada en la base de datos MySQL y retorna la orden creada con su ID generado
func (r *MySqlInboundOrderRepository) Create(ctx context.Context, inbound models.InboundOrder) (models.InboundOrder, error) {
	// SQL query to insert new inbound order with all required fields / Consulta SQL para insertar nueva orden de entrada con todos los campos requeridos
	query := `
		INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		inbound.OrderDate,
		inbound.OrderNumber,
		inbound.EmployeeId,
		inbound.ProductBatchId,
		inbound.WarehouseId,
	)
	if err != nil {
		return models.InboundOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	// Get the auto-generated ID from the database / Obtiene el ID autogenerado de la base de datos
	id, err := result.LastInsertId()
	if err != nil {
		return models.InboundOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	inbound.Id = int(id)
	return inbound, nil
}

// ExistsByOrderNumber - Checks if an inbound order with the given order number already exists in the MySQL database
// ExistsByOrderNumber - Verifica si una orden de entrada con el número de orden dado ya existe en la base de datos MySQL
func (r *MySqlInboundOrderRepository) ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	var exists bool
	// Simple query using EXISTS for efficient existence check / Consulta simple usando EXISTS para verificación eficiente de existencia
	query := "SELECT EXISTS(SELECT 1 FROM inbound_orders WHERE order_number = ?)"
	err := r.db.QueryRowContext(ctx, query, orderNumber).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	return exists, nil
}
