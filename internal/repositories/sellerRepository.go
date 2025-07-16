package repositories

import (
	"database/sql"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

var sellerRepositoryInstance SellerRepository

func NewSQLSellerRepository(db *sql.DB) SellerRepository {
	if sellerRepositoryInstance != nil {
		return sellerRepositoryInstance
	}

	sellerRepositoryInstance = &SQLSellerRepository{db: db}
	return sellerRepositoryInstance
}

type SellerRepository interface {
	GetAll() ([]models.Seller, error)
	Save(seller models.Seller) ([]models.Seller, error)
	Update(id int, seller models.Seller) ([]models.Seller, error)
	Delete(id int) error
}

type SQLSellerRepository struct {
	db *sql.DB
}

func (r *SQLSellerRepository) GetAll() ([]models.Seller, error) {
	rows, err := r.db.Query("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sellers []models.Seller
	for rows.Next() {
		var s models.Seller
		if err := rows.Scan(&s.Id, &s.CID, &s.CompanyName, &s.Address, &s.Telephone, &s.LocalityID); err != nil {
			return nil, err
		}
		sellers = append(sellers, s)
	}
	return sellers, nil
}

func (r *SQLSellerRepository) Save(seller models.Seller) ([]models.Seller, error) {
	// Validar que no exista otro con el mismo CID
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM sellers WHERE cid = ?)", seller.CID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, error_message.ErrAlreadyExists
	}

	res, err := r.db.Exec("INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)",
		seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.LocalityID)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	seller.Id = int(lastID)
	return []models.Seller{seller}, nil
}

func (r *SQLSellerRepository) Update(id int, seller models.Seller) ([]models.Seller, error) {
	var existing models.Seller
	err := r.db.QueryRow("SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = ?", id).
		Scan(&existing.Id, &existing.CID, &existing.CompanyName, &existing.Address, &existing.Telephone, &existing.LocalityID)

	if err == sql.ErrNoRows {
		return nil, error_message.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Solo actualiza los campos no vacíos
	if seller.CID != "" {
		existing.CID = seller.CID
	}
	if seller.CompanyName != "" {
		existing.CompanyName = seller.CompanyName
	}
	if seller.Address != "" {
		existing.Address = seller.Address
	}
	if seller.Telephone != "" {
		existing.Telephone = seller.Telephone
	}

	if seller.LocalityID != 0 {
		existing.LocalityID = seller.LocalityID
	}

	_, err = r.db.Exec("UPDATE sellers SET cid = ?, company_name = ?, address = ?, telephone = ?, locality_id = ? WHERE id = ?",
		existing.CID, existing.CompanyName, existing.Address, existing.Telephone, existing.LocalityID, id)

	if err != nil {
		return nil, err
	}

	return []models.Seller{existing}, nil
}

func (r *SQLSellerRepository) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM sellers WHERE id = ?", id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return error_message.ErrNotFound
	}
	return nil
}
