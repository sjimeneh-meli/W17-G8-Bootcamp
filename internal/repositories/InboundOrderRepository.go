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

func GetNewInboundOrderMySQLRepository(db *sql.DB) InboundOrderRepositoryI {
	if inboundOrderRepositoryInstance != nil {
		return inboundOrderRepositoryInstance
	}
	inboundOrderRepositoryInstance = &MySqlInboundOrderRepository{
		db: db,
	}
	return inboundOrderRepositoryInstance
}

type InboundOrderRepositoryI interface {
	GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error)
	GetInboundOrdersReportByEmployeeId(ctx context.Context, employeeId int) (models.InboundOrderReport, error)
	Create(ctx context.Context, inbound models.InboundOrder) (models.InboundOrder, error)
	ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error)
}

type MySqlInboundOrderRepository struct {
	db *sql.DB
}

func (r *MySqlInboundOrderRepository) GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error) {
	reports := []models.InboundOrderReport{}

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

func (r *MySqlInboundOrderRepository) GetInboundOrdersReportByEmployeeId(ctx context.Context, employeeId int) (models.InboundOrderReport, error) {
	report := models.InboundOrderReport{}

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
		if errors.Is(err, sql.ErrNoRows) {
			return models.InboundOrderReport{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "employee with Id", employeeId, "doesn't exist.")
		}
		return models.InboundOrderReport{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	return report, nil
}
func (r *MySqlInboundOrderRepository) Create(ctx context.Context, inbound models.InboundOrder) (models.InboundOrder, error) {
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

	id, err := result.LastInsertId()
	if err != nil {
		return models.InboundOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	inbound.Id = int(id)
	return inbound, nil
}

func (r *MySqlInboundOrderRepository) ExistsByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM inbound_orders WHERE order_number = ?)"
	err := r.db.QueryRowContext(ctx, query, orderNumber).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	return exists, nil
}
