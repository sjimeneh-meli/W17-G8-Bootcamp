package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type LocalityRepository interface {
	Save(locality models.Locality) (models.Locality, error)
	GetSellerReports(localityID int) ([]responses.LocalitySellerReport, error)
}
type SQLLocalityRepository struct {
	db *sql.DB
}

func NewSQLLocalityRepository(db *sql.DB) *SQLLocalityRepository {
	return &SQLLocalityRepository{db: db}
}

func (r *SQLLocalityRepository) Save(locality models.Locality) (models.Locality, error) {
	// 1. Buscar o insertar el país
	var countryID int
	err := r.db.QueryRow("SELECT id FROM countries WHERE country_name = ?", locality.CountryName).Scan(&countryID)
	if err == sql.ErrNoRows {
		res, err := r.db.Exec("INSERT INTO countries (country_name) VALUES (?)", locality.CountryName)
		if err != nil {
			return models.Locality{}, err
		}
		lastID, _ := res.LastInsertId()
		countryID = int(lastID)
	} else if err != nil {
		return models.Locality{}, err
	}

	// 2. Buscar o insertar la provincia
	var provinceID int
	err = r.db.QueryRow("SELECT id FROM provinces WHERE province_name = ? AND id_country_fk = ?", locality.ProvinceName, countryID).Scan(&provinceID)
	if err == sql.ErrNoRows {
		res, err := r.db.Exec("INSERT INTO provinces (province_name, id_country_fk) VALUES (?, ?)", locality.ProvinceName, countryID)
		if err != nil {
			return models.Locality{}, err
		}
		lastID, _ := res.LastInsertId()
		provinceID = int(lastID)
	} else if err != nil {
		return models.Locality{}, err
	}

	// ✅ 3. Verificar si ya existe una localidad con ese nombre y esa provincia
	var exists bool
	err = r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM localities WHERE locality_name = ? AND province_id = ?
		)
	`, locality.LocalityName, provinceID).Scan(&exists)
	if err != nil {
		return models.Locality{}, err
	}
	if exists {
		return models.Locality{}, errors.New("locality already exists") // 409 Conflict
	}

	// 4. Insertar nueva localidad (el ID será auto-generado)
	res, err := r.db.Exec(
		"INSERT INTO localities (locality_name, province_id) VALUES (?, ?)",
		locality.LocalityName, provinceID,
	)
	fmt.Println(res)
	if err != nil {
		return models.Locality{}, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return models.Locality{}, err
	}
	locality.Id = int(lastID)

	return locality, nil
}

func (r *SQLLocalityRepository) GetSellerReports(localityID int) ([]responses.LocalitySellerReport, error) {
	// Si se especifica un ID, verificar que exista
	if localityID != 0 {
		var exists bool
		err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM localities WHERE id = ?)", localityID).Scan(&exists)
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
		rows, err = r.db.Query(`
			SELECT l.id, l.locality_name, COUNT(s.id)
			FROM localities l
			LEFT JOIN sellers s ON s.locality_id = l.id
			WHERE l.id = ?
			GROUP BY l.id, l.locality_name
		`, localityID)
	} else {
		// Si no se envía ID, traer todos
		rows, err = r.db.Query(`
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
