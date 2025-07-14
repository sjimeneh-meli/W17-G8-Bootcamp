package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type PurchaseOrderRepositoryI interface {
	GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error)
	Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error)
	ExistPurchaseOrderByOrderNumber(ctx context.Context, orderNumber string) (bool, error)
	GetPurchaseOrdersReportByBuyerId(ctx context.Context, buyerId int) (models.PurchaseOrderReport, error)
	GetAllPurchaseOrdersReports(ctx context.Context) ([]models.PurchaseOrderReport, error)
}

type MySqlPurchaseOrderRepository struct {
	db *sql.DB
}

func GetNewPurchaseOrderMySQLRepository(db *sql.DB) PurchaseOrderRepositoryI {
	return &MySqlPurchaseOrderRepository{
		db: db,
	}
}

func (r *MySqlPurchaseOrderRepository) GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error) {
	orders := make(map[int]models.PurchaseOrder)

	query := "select id, order_number, order_date, tracking_code, buyer_id, product_record_id from purchase_orders"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return orders, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}
	defer rows.Close()

	tempOrdersMap := make(map[int]models.PurchaseOrder)

	for rows.Next() {
		order := models.PurchaseOrder{}
		err := rows.Scan(&order.Id, &order.OrderNumber, &order.OrderDate, &order.TrackingCode, &order.BuyerId, &order.ProductRecordId)
		if err != nil {
			return orders, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
		}

		tempOrdersMap[order.Id] = order
	}

	orders = tempOrdersMap
	return orders, nil
}

func (r *MySqlPurchaseOrderRepository) Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error) {
	query := `insert into purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id)
	values (?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, order.OrderNumber, order.OrderDate, order.TrackingCode, order.BuyerId, order.ProductRecordId)
	if err != nil {
		return models.PurchaseOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return models.PurchaseOrder{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	order.Id = int(lastId)
	return order, nil
}

func (r *MySqlPurchaseOrderRepository) GetPurchaseOrdersReportByBuyerId(ctx context.Context, buyerId int) (models.PurchaseOrderReport, error) {
	report := models.PurchaseOrderReport{}

	query := `select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as "purchase_orders_count"
from productos_frescos.buyers b
inner join productos_frescos.purchase_orders po on po.buyer_id = b.id
where b.id = ?
group by b.id;`
	row := r.db.QueryRowContext(ctx, query, buyerId)

	err := row.Scan(&report.Id, &report.IdCardNumber, &report.FirstName, &report.LastName, &report.PurchaseOrderCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PurchaseOrderReport{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", buyerId, "doesn't exists.")
		}
		return models.PurchaseOrderReport{}, fmt.Errorf("%w - %s", error_message.ErrInternalServerError, err.Error())
	}

	return report, nil
}

func (r *MySqlPurchaseOrderRepository) GetAllPurchaseOrdersReports(ctx context.Context) ([]models.PurchaseOrderReport, error) {
	reports := []models.PurchaseOrderReport{}

	query := `select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as "purchase_orders_count"
from productos_frescos.buyers b
inner join productos_frescos.purchase_orders po on po.buyer_id = b.id
group by b.id
order by b.id;`

	rows, err := r.db.QueryContext(ctx, query)
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

func (r *MySqlPurchaseOrderRepository) ExistPurchaseOrderByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	query := "SELECT 1 FROM purchase_orders WHERE order_number = ? LIMIT 1;"

	var exists int64
	err := r.db.QueryRowContext(ctx, query, orderNumber).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error al verificar la existencia del order number: %w", err)
	}
	return true, nil
}
