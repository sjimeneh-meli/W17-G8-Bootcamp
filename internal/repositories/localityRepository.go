package repositories

import (
	"context"
	"database/sql"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type LocalityRepository interface {
	Save(ctx context.Context, locality models.Locality) (models.Locality, error)
	GetSellerReports(ctx context.Context, localityID int) ([]responses.LocalitySellerReport, error)
	ExistById(ctx context.Context, localityID int) (bool, error)
}
type SQLLocalityRepository struct {
	db *sql.DB
}

func NewSQLLocalityRepository(db *sql.DB) *SQLLocalityRepository {
	return &SQLLocalityRepository{db: db}
}

func (r *SQLLocalityRepository) Save(ctx context.Context, locality models.Locality) (models.Locality, error) {
	// 1. Buscar o insertar el país
	var countryID int
	err := r.db.QueryRowContext(ctx, "SELECT id FROM countries WHERE country_name = ?", locality.CountryName).Scan(&countryID)
	if err == sql.ErrNoRows {
		res, err := r.db.ExecContext(ctx, "INSERT INTO countries (country_name) VALUES (?)", locality.CountryName)
		if err != nil {
			return models.Locality{}, error_message.ErrQuery
		}
		lastID, _ := res.LastInsertId()
		countryID = int(lastID)
	} else if err != nil {
		return models.Locality{}, error_message.ErrQuery
	}

	// 2. Buscar o insertar la provincia
	var provinceID int
	err = r.db.QueryRowContext(ctx, "SELECT id FROM provinces WHERE province_name = ? AND id_country_fk = ?", locality.ProvinceName, countryID).Scan(&provinceID)
	if err == sql.ErrNoRows {
		res, err := r.db.ExecContext(ctx, "INSERT INTO provinces (province_name, id_country_fk) VALUES (?, ?)", locality.ProvinceName, countryID)
		if err != nil {
			return models.Locality{}, error_message.ErrQuery
		}
		lastID, _ := res.LastInsertId()
		provinceID = int(lastID)
	} else if err != nil {
		return models.Locality{}, error_message.ErrQuery
	}

	// ✅ 3. Verificar si ya existe una localidad con ese nombre y esa provincia
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

	// 4. Insertar nueva localidad (el ID será auto-generado)
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO localities (locality_name, province_id) VALUES (?, ?)",
		locality.LocalityName, provinceID,
	)

	if err != nil {
		return models.Locality{}, error_message.ErrQuery
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return models.Locality{}, error_message.ErrQuery
	}
	locality.Id = int(lastID)

	return locality, nil
}

func (r *SQLLocalityRepository) GetSellerReports(ctx context.Context, localityID int) ([]responses.LocalitySellerReport, error) {
	// Si se especifica un ID, verificar que exista
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

	// Si se pide un solo locality
	if localityID != 0 {
		rows, err = r.db.QueryContext(ctx, `
			SELECT l.id, l.locality_name, COUNT(s.id)
			FROM localities l
			LEFT JOIN sellers s ON s.locality_id = l.id
			WHERE l.id = ?
			GROUP BY l.id, l.locality_name
		`, localityID)
	} else {
		// Si no se envía ID, traer todos
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

func (r *SQLLocalityRepository) ExistById(ctx context.Context, localityID int) (bool, error) {
	row := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM localities WHERE id = ?)", localityID)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, error_message.ErrQuery
	}
	return exists, nil
}
