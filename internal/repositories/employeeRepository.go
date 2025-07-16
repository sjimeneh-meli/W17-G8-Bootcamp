package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

var employeeRepositoryInstance EmployeeRepositoryI

// GetNewEmployeeMySQLRepository - Creates and returns a new instance of MySqlEmployeeRepository using singleton pattern
// GetNewEmployeeMySQLRepository - Crea y retorna una nueva instancia de MySqlEmployeeRepository usando patrón singleton
func GetNewEmployeeMySQLRepository(db *sql.DB) EmployeeRepositoryI {
	if employeeRepositoryInstance != nil {
		return employeeRepositoryInstance
	}

	employeeRepositoryInstance = &MySqlEmployeeRepository{
		db: db,
	}
	return employeeRepositoryInstance
}

// EmployeeRepositoryI - Interface defining the contract for employee repository operations
// EmployeeRepositoryI - Interfaz que define el contrato para las operaciones del repositorio de empleados
type EmployeeRepositoryI interface {
	// GetAll - Retrieves all employees from the database and returns them as a map with employee ID as key
	// GetAll - Obtiene todos los empleados de la base de datos y los retorna como un mapa con el ID del empleado como clave
	GetAll(ctx context.Context) (map[int]models.Employee, error)

	// GetById - Retrieves a specific employee by their ID from the database
	// GetById - Obtiene un empleado específico por su ID de la base de datos
	GetById(ctx context.Context, id int) (models.Employee, error)

	// DeleteById - Removes an employee from the database by their ID
	// DeleteById - Elimina un empleado de la base de datos por su ID
	DeleteById(ctx context.Context, id int) error

	// Create - Inserts a new employee into the database and returns the created employee with its generated ID
	// Create - Inserta un nuevo empleado en la base de datos y retorna el empleado creado con su ID generado
	Create(ctx context.Context, employee models.Employee) (models.Employee, error)

	// Update - Modifies an existing employee in the database with partial updates support
	// Update - Modifica un empleado existente en la base de datos con soporte para actualizaciones parciales
	Update(ctx context.Context, employeeId int, employee models.Employee) (models.Employee, error)

	// GetCardNumberIds - Retrieves all card number IDs from the database for validation purposes
	// GetCardNumberIds - Obtiene todos los IDs de números de tarjeta de la base de datos para propósitos de validación
	GetCardNumberIds() ([]string, error)

	// ExistEmployeeById - Checks if an employee with the given ID exists in the database
	// ExistEmployeeById - Verifica si un empleado con el ID dado existe en la base de datos
	ExistEmployeeById(ctx context.Context, employeeId int) (bool, error)
}

// MySqlEmployeeRepository - MySQL implementation of the EmployeeRepositoryI interface
// MySqlEmployeeRepository - Implementación MySQL de la interfaz EmployeeRepositoryI
type MySqlEmployeeRepository struct {
	db *sql.DB // Database connection / Conexión a la base de datos
}

// GetAll - Retrieves all employees from the MySQL database and returns them as a map with employee ID as key
// GetAll - Obtiene todos los empleados de la base de datos MySQL y los retorna como un mapa con el ID del empleado como clave
func (r *MySqlEmployeeRepository) GetAll(ctx context.Context) (map[int]models.Employee, error) {
	employees := make(map[int]models.Employee)

	// SQL query to select all employee fields / Consulta SQL para seleccionar todos los campos del empleado
	query := "SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return employees, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	defer rows.Close()

	// Iterate through all rows and map each employee to the result map
	// Itera a través de todas las filas y mapea cada empleado al mapa de resultados
	for rows.Next() {
		employee := models.Employee{}
		err = rows.Scan(&employee.Id, &employee.CardNumberID, &employee.FirstName, &employee.LastName, &employee.WarehouseID)
		if err != nil {
			return employees, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
		}
		employees[employee.Id] = employee
	}

	return employees, nil
}

// GetById - Retrieves a specific employee by their ID from the MySQL database
// GetById - Obtiene un empleado específico por su ID de la base de datos MySQL
func (r *MySqlEmployeeRepository) GetById(ctx context.Context, id int) (models.Employee, error) {
	employee := models.Employee{}

	// SQL query to select employee by specific ID / Consulta SQL para seleccionar empleado por ID específico
	query := "SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&employee.Id, &employee.CardNumberID, &employee.FirstName, &employee.LastName, &employee.WarehouseID)
	if err != nil {
		// Handle case when no employee is found / Maneja el caso cuando no se encuentra ningún empleado
		if errors.Is(err, sql.ErrNoRows) {
			return models.Employee{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "employee with Id", id, "not exists.")
		}
		return models.Employee{}, err
	}

	return employee, nil
}

// DeleteById - Removes an employee from the MySQL database by their ID
// DeleteById - Elimina un empleado de la base de datos MySQL por su ID
func (r *MySqlEmployeeRepository) DeleteById(ctx context.Context, id int) error {
	// SQL query to delete employee by ID / Consulta SQL para eliminar empleado por ID
	query := "DELETE FROM employees WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err)
	}

	// Check if any rows were affected to confirm deletion / Verifica si alguna fila fue afectada para confirmar la eliminación
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err)
	}

	// If no rows affected, employee doesn't exist / Si ninguna fila fue afectada, el empleado no existe
	if rowsAffected == 0 {
		return fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Employee with Id", id, "doesn't exist.")
	}

	return nil
}

// Create - Inserts a new employee into the MySQL database and returns the created employee with its generated ID
// Create - Inserta un nuevo empleado en la base de datos MySQL y retorna el empleado creado con su ID generado
func (r *MySqlEmployeeRepository) Create(ctx context.Context, employee models.Employee) (models.Employee, error) {
	// SQL query to insert new employee / Consulta SQL para insertar nuevo empleado
	query := `INSERT INTO employees (id_card_number, first_name, last_name, warehouse_id) VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, employee.CardNumberID, employee.FirstName, employee.LastName, employee.WarehouseID)
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	// Get the auto-generated ID from the database / Obtiene el ID autogenerado de la base de datos
	lastId, err := result.LastInsertId()
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	employee.Id = int(lastId)
	return employee, nil
}

// Update - Modifies an existing employee in the MySQL database with support for partial updates
// Update - Modifica un empleado existente en la base de datos MySQL con soporte para actualizaciones parciales
func (r *MySqlEmployeeRepository) Update(ctx context.Context, employeeId int, employee models.Employee) (models.Employee, error) {
	updates := []string{}
	values := []interface{}{}

	// Build dynamic UPDATE query based on provided fields / Construye consulta UPDATE dinámica basada en campos proporcionados
	if employee.FirstName != "" {
		updates = append(updates, "first_name = ?")
		values = append(values, employee.FirstName)
	}
	if employee.LastName != "" {
		updates = append(updates, "last_name = ?")
		values = append(values, employee.LastName)
	}
	if employee.CardNumberID != "" {
		updates = append(updates, "id_card_number = ?")
		values = append(values, employee.CardNumberID)
	}
	if employee.WarehouseID != 0 {
		updates = append(updates, "warehouse_id = ?")
		values = append(values, employee.WarehouseID)
	}

	// Execute dynamic UPDATE query / Ejecuta consulta UPDATE dinámica
	query := "UPDATE employees SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	values = append(values, employeeId)

	result, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	// Check if any rows were affected to confirm update / Verifica si alguna fila fue afectada para confirmar la actualización
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err)
	}

	// If no rows affected, employee doesn't exist / Si ninguna fila fue afectada, el empleado no existe
	if rowsAffected == 0 {
		return models.Employee{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "employee with Id", employeeId, "not exists.")
	}

	// Retrieve and return the updated employee / Obtiene y retorna el empleado actualizado
	updatedEmployee, err := r.GetById(ctx, employeeId)
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	return updatedEmployee, nil
}

// GetCardNumberIds - Retrieves all card number IDs from the MySQL database for validation purposes
// GetCardNumberIds - Obtiene todos los IDs de números de tarjeta de la base de datos MySQL para propósitos de validación
func (r *MySqlEmployeeRepository) GetCardNumberIds() ([]string, error) {
	cardNumberIds := []string{}

	// SQL query to select all card number IDs / Consulta SQL para seleccionar todos los IDs de números de tarjeta
	query := "SELECT id_card_number FROM employees"
	rows, err := r.db.Query(query)
	if err != nil {
		return []string{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	defer rows.Close()

	// Iterate through all rows and collect card number IDs / Itera a través de todas las filas y recolecta los IDs de números de tarjeta
	for rows.Next() {
		cardNumberId := ""
		err = rows.Scan(&cardNumberId)
		if err != nil {
			return []string{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
		}

		cardNumberIds = append(cardNumberIds, cardNumberId)
	}

	return cardNumberIds, nil
}

// ExistEmployeeById - Checks if an employee with the given ID exists in the MySQL database
// ExistEmployeeById - Verifica si un empleado con el ID dado existe en la base de datos MySQL
func (r *MySqlEmployeeRepository) ExistEmployeeById(ctx context.Context, employeeId int) (bool, error) {
	// Simple query to check employee existence using LIMIT 1 for efficiency / Consulta simple para verificar existencia del empleado usando LIMIT 1 por eficiencia
	query := "SELECT 1 FROM employees WHERE id = ? LIMIT 1"

	var exists int64
	err := r.db.QueryRowContext(ctx, query, employeeId).Scan(&exists)

	if err != nil {
		// If no rows found, employee doesn't exist (not an error) / Si no se encuentran filas, el empleado no existe (no es un error)
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error verifying employee existence: %w", err)
	}
	return true, nil
}
