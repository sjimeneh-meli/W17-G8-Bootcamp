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

func GetNewEmployeeMySQLRepository(db *sql.DB) EmployeeRepositoryI {
	if employeeRepositoryInstance != nil {
		return employeeRepositoryInstance
	}

	employeeRepositoryInstance = &MySqlEmployeeRepository{
		db: db,
	}
	return employeeRepositoryInstance
}

type EmployeeRepositoryI interface {
	GetAll(ctx context.Context) (map[int]models.Employee, error)
	GetById(ctx context.Context, id int) (models.Employee, error)
	DeleteById(ctx context.Context, id int) error
	Create(ctx context.Context, employee models.Employee) (models.Employee, error)
	Update(ctx context.Context, employeeId int, employee models.Employee) (models.Employee, error)
	GetCardNumberIds() ([]string, error)
	ExistEmployeeById(ctx context.Context, employeeId int) (bool, error)
}

type MySqlEmployeeRepository struct {
	db *sql.DB
}

func (r *MySqlEmployeeRepository) GetAll(ctx context.Context) (map[int]models.Employee, error) {
	employees := make(map[int]models.Employee)

	query := "SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return employees, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	defer rows.Close()

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

func (r *MySqlEmployeeRepository) GetById(ctx context.Context, id int) (models.Employee, error) {
	employee := models.Employee{}

	query := "SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&employee.Id, &employee.CardNumberID, &employee.FirstName, &employee.LastName, &employee.WarehouseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Employee{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "employee with Id", id, "not exists.")
		}
		return models.Employee{}, err
	}

	return employee, nil
}

func (r *MySqlEmployeeRepository) DeleteById(ctx context.Context, id int) error {
	query := "DELETE FROM employees WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Employee with Id", id, "doesn't exist.")
	}

	return nil
}

func (r *MySqlEmployeeRepository) Create(ctx context.Context, employee models.Employee) (models.Employee, error) {
	query := `INSERT INTO employees (id_card_number, first_name, last_name, warehouse_id) VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, employee.CardNumberID, employee.FirstName, employee.LastName, employee.WarehouseID)
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	employee.Id = int(lastId)
	return employee, nil
}

func (r *MySqlEmployeeRepository) Update(ctx context.Context, employeeId int, employee models.Employee) (models.Employee, error) {
	updates := []string{}
	values := []interface{}{}

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

	query := "UPDATE employees SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	values = append(values, employeeId)

	result, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err)
	}

	if rowsAffected == 0 {
		return models.Employee{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "employee with Id", employeeId, "not exists.")
	}

	updatedEmployee, err := r.GetById(ctx, employeeId)
	if err != nil {
		return models.Employee{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	return updatedEmployee, nil
}

func (r *MySqlEmployeeRepository) GetCardNumberIds() ([]string, error) {
	cardNumberIds := []string{}

	query := "SELECT id_card_number FROM employees"
	rows, err := r.db.Query(query)
	if err != nil {
		return []string{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	defer rows.Close()

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

func (r *MySqlEmployeeRepository) ExistEmployeeById(ctx context.Context, employeeId int) (bool, error) {
	query := "SELECT 1 FROM employees WHERE id = ? LIMIT 1"

	var exists int64
	err := r.db.QueryRowContext(ctx, query, employeeId).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error verifying employee existence: %w", err)
	}
	return true, nil
}
