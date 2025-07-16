package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type InboundOrderRepositoryI interface {
	GetAllInboundOrdersReports(ctx context.Context) ([]models.InboundOrderReport, error)
	GetInboundOrdersReportByEmployeeId(ctx context.Context, employeeId int) (models.InboundOrderReport, error)
}

type MySqlInboundOrderRepository struct {
	db *sql.DB
}

func GetNewInboundOrderMySQLRepository(db *sql.DB) InboundOrderRepositoryI {
	return &MySqlInboundOrderRepository{
		db: db,
	}
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
