package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type PurchaseOrderRepositoryI interface {
	GetAll(ctx context.Context) (map[int]models.PurchaseOrder, error)
	Create(ctx context.Context, order models.PurchaseOrder) (models.PurchaseOrder, error)
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
	return models.PurchaseOrder{}, nil
	//Implementar luego
}
