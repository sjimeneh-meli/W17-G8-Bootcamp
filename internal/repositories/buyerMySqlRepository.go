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

type MySqlBuyerRepository struct {
	db *sql.DB
}

func GetNewBuyerMySQLRepository(db *sql.DB) BuyerRepositoryI {
	return &MySqlBuyerRepository{
		db: db,
	}
}

func (r *MySqlBuyerRepository) GetAll(ctx context.Context) (map[int]models.Buyer, error) {
	buyers := make(map[int]models.Buyer)
	rows, err := r.db.QueryContext(ctx, "select id, id_card_number, first_name, last_name from buyers")
	if err != nil {
		return buyers, err
	}
	defer rows.Close()

	tempBuyersMap := make(map[int]models.Buyer)
	for rows.Next() {
		buyer := models.Buyer{}
		err = rows.Scan(&buyer.Id, &buyer.CardNumberId, &buyer.FirstName, &buyer.LastName)
		if err != nil {
			return buyers, err
		}
		tempBuyersMap[buyer.Id] = buyer
	}

	buyers = tempBuyersMap
	return buyers, nil
}

func (r *MySqlBuyerRepository) GetById(ctx context.Context, id int) (models.Buyer, error) {
	buyer := models.Buyer{}

	row := r.db.QueryRowContext(ctx, "select id, id_card_number, first_name, last_name from buyers where id = ?", id)
	err := row.Err()
	if err != nil {
		return buyer, nil
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

func (r *MySqlBuyerRepository) DeleteById(ctx context.Context, id int) error {

	result, err := r.db.ExecContext(ctx, "delete from buyers where id = ?", id)
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

func (r *MySqlBuyerRepository) Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error) {
	result, err := r.db.ExecContext(ctx, `insert into buyers (id, id_card_number, first_name, last_name) values (?, ?, ?, ?)`, buyer.Id, buyer.CardNumberId, buyer.FirstName, buyer.LastName)

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

func (r *MySqlBuyerRepository) GetCardNumberIds() ([]string, error) {
	cardNumberIds := []string{}

	rows, err := r.db.Query("select id_card_number from buyers")
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

func (r *MySqlBuyerRepository) GetPurchaseOrdersReportByBuyerId(ctx context.Context, buyerId int) (models.PurchaseOrderReport, error) {

	report := models.PurchaseOrderReport{}
	row := r.db.QueryRowContext(ctx, `select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as "purchase_orders_count"
from productos_frescos.buyers b
inner join productos_frescos.purchase_orders po on po.buyer_id = b.id
where b.id = ?
group by b.id;`, buyerId)

	err := row.Scan(&report.Id, &report.IdCardNumber, &report.FirstName, &report.LastName, &report.PurchaseOrderCount)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return models.PurchaseOrderReport{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", buyerId, "doesn't exists.")
		}

		return models.PurchaseOrderReport{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	return report, nil
}

func (r *MySqlBuyerRepository) GetPurchaseOrdersReport(ctx context.Context) ([]models.PurchaseOrderReport, error) {

	reports := []models.PurchaseOrderReport{}
	rows, err := r.db.QueryContext(ctx, `select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as "purchase_orders_count"
from productos_frescos.buyers b
inner join productos_frescos.purchase_orders po on po.buyer_id = b.id
group by b.id;`)

	if err != nil {
		return []models.PurchaseOrderReport{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		report := models.PurchaseOrderReport{}
		err := rows.Scan(&report.Id, &report.IdCardNumber, &report.FirstName, &report.LastName, &report.PurchaseOrderCount)
		if err != nil {
			return []models.PurchaseOrderReport{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
		}

		reports = append(reports, report)
	}
	return reports, nil
}
