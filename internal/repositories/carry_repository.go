package repositories

import (
	"database/sql"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type CarryRepository interface {
	Create(carry models.Carry) (models.Carry, error)
	ExistsByCid(cid string) (bool, error)
}

type CarryRepositoryImpl struct {
	db *sql.DB
}

func NewCarryRepository(db *sql.DB) *CarryRepositoryImpl {
	return &CarryRepositoryImpl{db: db}
}

func (r *CarryRepositoryImpl) Create(carry models.Carry) (models.Carry, error) {
	result, err := r.db.Exec(
		"INSERT INTO `carriers`(`cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES (?,?,?,?,?)",
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

func (r *CarryRepositoryImpl) ExistsByCid(cid string) (bool, error) {
	row := r.db.QueryRow("SELECT COUNT(*) FROM `carriers` WHERE `cid` = ?", cid)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%w: %v", error_message.ErrInternalServerError, err)
	}

	return count > 0, nil
}
