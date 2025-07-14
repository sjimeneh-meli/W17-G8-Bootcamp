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

// MySqlBuyerRepository implements BuyerRepositoryI for MySQL database operations
type MySqlBuyerRepository struct {
	db *sql.DB
}

// GetNewBuyerMySQLRepository creates and returns a new instance of MySqlBuyerRepository
func GetNewBuyerMySQLRepository(db *sql.DB) BuyerRepositoryI {
	return &MySqlBuyerRepository{
		db: db,
	}
}

// GetAll retrieves all buyers from the MySQL database
// Returns a map with buyer ID as key and buyer model as value
func (r *MySqlBuyerRepository) GetAll(ctx context.Context) (map[int]models.Buyer, error) {
	buyers := make(map[int]models.Buyer)

	query := "select id, id_card_number, first_name, last_name from buyers"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return buyers, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	defer rows.Close()

	tempBuyersMap := make(map[int]models.Buyer)
	for rows.Next() {
		buyer := models.Buyer{}
		err = rows.Scan(&buyer.Id, &buyer.CardNumberId, &buyer.FirstName, &buyer.LastName)
		if err != nil {
			return buyers, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
		}
		tempBuyersMap[buyer.Id] = buyer
	}

	buyers = tempBuyersMap
	return buyers, nil
}

// GetById retrieves a buyer by their ID from the MySQL database
// Returns an error if the buyer doesn't exist
func (r *MySqlBuyerRepository) GetById(ctx context.Context, id int) (models.Buyer, error) {
	buyer := models.Buyer{}

	query := "select id, id_card_number, first_name, last_name from buyers where id = ?"
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Err()
	if err != nil {
		return buyer, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	err = row.Scan(&buyer.Id, &buyer.CardNumberId, &buyer.FirstName, &buyer.LastName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Buyer{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", id, "doesn't exists.")
		}
		return models.Buyer{}, err
	}

	return buyer, nil
}

// DeleteById removes a buyer from the MySQL database by their ID
// Returns an error if the buyer doesn't exist
func (r *MySqlBuyerRepository) DeleteById(ctx context.Context, id int) error {
	query := "delete from buyers where id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", id, "doesn't exists.")
	}

	return nil
}

// Create inserts a new buyer into the MySQL database
// Returns the created buyer with its generated ID
func (r *MySqlBuyerRepository) Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error) {
	query := `insert into buyers (id_card_number, first_name, last_name) values (?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, buyer.CardNumberId, buyer.FirstName, buyer.LastName)

	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	buyer.Id = int(lastId)
	return buyer, nil
}

// Update modifies an existing buyer in the MySQL database
// Only updates the fields that are provided (non-empty values)
// Returns an error if the buyer doesn't exist
func (r *MySqlBuyerRepository) Update(ctx context.Context, buyerId int, buyer models.Buyer) (models.Buyer, error) {
	updates := []string{}
	values := []interface{}{}

	if buyer.FirstName != "" {
		updates = append(updates, "first_name = ?")
		values = append(values, buyer.FirstName)
	}
	if buyer.LastName != "" {
		updates = append(updates, "last_name = ?")
		values = append(values, buyer.LastName)
	}
	if buyer.CardNumberId != "" {
		updates = append(updates, "id_card_number = ?")
		values = append(values, buyer.CardNumberId)
	}

	query := "UPDATE buyers SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	values = append(values, buyerId)

	result, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err)
	}

	if rowsAffected == 0 {
		return models.Buyer{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", buyerId, "doesn't exists.")
	}

	updatedUser, err := r.GetById(ctx, buyerId)
	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	return updatedUser, nil
}

// GetCardNumberIds retrieves all card number IDs from the MySQL database
// Returns a slice of all existing card number IDs
func (r *MySqlBuyerRepository) GetCardNumberIds() ([]string, error) {
	cardNumberIds := []string{}

	query := "select id_card_number from buyers"
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

// ExistBuyerById checks if a buyer with the given ID exists in the MySQL database
// Returns true if the buyer exists, false otherwise
func (r *MySqlBuyerRepository) ExistBuyerById(ctx context.Context, buyerId int) (bool, error) {
	query := "SELECT 1 FROM buyers WHERE id = ? LIMIT 1;"

	var exists int64
	err := r.db.QueryRowContext(ctx, query, buyerId).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error al verificar la existencia del producto: %w", err)
	}
	return true, nil
}
