package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type productRecordRepository struct {
	DB *sql.DB
}

type IProductRecordRepository interface {
	Create(ctx context.Context, pr *models.ProductRecord) (*models.ProductRecord, error)
	GetReportByIdProduct(ctx context.Context, id int64) (*models.ProductRecordReport, error)
	ExistProductByID(ctx context.Context, id int64) (bool, error)
}

func NewProductRecordRepository(db *sql.DB) IProductRecordRepository {
	return &productRecordRepository{DB: db}
}

func (prr *productRecordRepository) Create(ctx context.Context, pr *models.ProductRecord) (*models.ProductRecord, error) {
	query := "INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES (?, ?, ?, ?)"
	result, err := prr.DB.ExecContext(ctx, query, pr.LastUpdateDate, pr.PurchasePrice, pr.SalePrice, pr.ProductID)
	if err != nil {
		return nil, fmt.Errorf("error to create product record: %w", err)
	}
	productRecordId, err := result.LastInsertId()

	if err != nil {
		return nil, fmt.Errorf("error to get last insert id: %w", err)
	}

	pr.ID = int(productRecordId)

	return pr, err
}

func (prr *productRecordRepository) GetReportByIdProduct(ctx context.Context, id int64) (*models.ProductRecordReport, error) {
	query := `
	SELECT p.id as product_id, p.description, count(*) as records_count
	FROM 
		products as p
	LEFT JOIN
		product_records as pr
	ON pr.product_id = p.id
	WHERE p.id = ?
	GROUP BY p.id, p.description`

	var productRecordReport models.ProductRecordReport

	err := prr.DB.QueryRowContext(ctx, query, id).Scan(&productRecordReport.ProductId, &productRecordReport.Description, &productRecordReport.RecordsCount)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error_message.ErrNotFound
		}
		return nil, fmt.Errorf("error to scan product record report %w", err)
	}

	return &productRecordReport, nil

}

func (prr *productRecordRepository) ExistProductByID(ctx context.Context, id int64) (bool, error) {
	query := "SELECT 1 FROM products WHERE id = ? LIMIT 1;"

	var exists int64
	err := prr.DB.QueryRowContext(ctx, query, id).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error al verificar la existencia del producto: %w", err)
	}

	return true, nil
}
