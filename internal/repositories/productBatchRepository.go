package repositories

import "database/sql"

func GetProductBatchRepository(db *sql.DB) (ProductBatchRepositoryI, error) {
	return &productBatchRepository{
		database:  db,
		tablename: "product_batches",
	}, nil
}

type ProductBatchRepositoryI interface{}

type productBatchRepository struct {
	database  *sql.DB
	tablename string
}
