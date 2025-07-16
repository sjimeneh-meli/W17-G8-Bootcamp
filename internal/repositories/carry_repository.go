package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

// Carry table and field constants
const (
	carryTable = "carriers"

	// Field groups for better maintainability
	carryFields       = "`id`, `cid`, `company_name`, `address`, `telephone`, `locality_id`"
	carryInsertFields = "`cid`, `company_name`, `address`, `telephone`, `locality_id`"
)

// Carry query strings - organized by operation type
var (
	// INSERT queries
	queryCreateCarry = fmt.Sprintf("INSERT INTO `%s`(%s) VALUES (?,?,?,?,?)", carryTable, carryInsertFields)

	// SELECT queries
	queryExistsByCid = fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE `cid` = ?", carryTable)

	// Report queries
	queryGetCarryReportsByLocality = "SELECT l.id, l.locality_name, COUNT(c.id) AS carriers_count FROM localities l LEFT JOIN carriers c ON l.id = c.locality_id WHERE l.id = ? GROUP BY l.id"
	queryGetAllCarryReports        = "SELECT l.id, l.locality_name, COUNT(c.id) AS carriers_count FROM localities l LEFT JOIN carriers c ON l.id = c.locality_id GROUP BY l.id"
)

type CarryRepository interface {
	Create(ctx context.Context, carry models.Carry) (models.Carry, error)
	ExistsByCid(ctx context.Context, cid string) (bool, error)
	GetCarryReportsByLocality(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error)
}

type CarryRepositoryImpl struct {
	db *sql.DB
}

func NewCarryRepository(db *sql.DB) *CarryRepositoryImpl {
	return &CarryRepositoryImpl{db: db}
}

func (r *CarryRepositoryImpl) Create(ctx context.Context, carry models.Carry) (models.Carry, error) {
	result, err := r.db.ExecContext(ctx, queryCreateCarry,
		carry.Cid, carry.CompanyName, carry.Address, carry.Telephone, carry.LocalityId,
	)
	if err != nil {
		return models.Carry{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return models.Carry{}, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}
	carry.Id = int(lastInsertId)

	return carry, nil
}

func (r *CarryRepositoryImpl) ExistsByCid(ctx context.Context, cid string) (bool, error) {
	row := r.db.QueryRowContext(ctx, queryExistsByCid, cid)

	var count int
	err := row.Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			// Si no hay filas, significa que no existe el CID
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return count > 0, nil
}

func (r *CarryRepositoryImpl) GetCarryReportsByLocality(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
	var rows *sql.Rows
	var err error

	if localityID != 0 {
		rows, err = r.db.QueryContext(ctx, queryGetCarryReportsByLocality, localityID)
	} else {
		rows, err = r.db.QueryContext(ctx, queryGetAllCarryReports)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	defer rows.Close()

	var reports []responses.LocalityCarryReport
	for rows.Next() {
		var report responses.LocalityCarryReport
		err := rows.Scan(&report.LocalityId, &report.LocalityName, &report.CarriersCount)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
		}
		reports = append(reports, report)
	}

	// Verificar si hubo error durante la iteración
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return reports, nil
}
