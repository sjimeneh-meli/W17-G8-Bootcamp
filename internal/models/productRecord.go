package models

import "time"

type ProductRecord struct {
	ID             int       `json:"id"`
	LastUpdateDate time.Time `json:"last_update_date"`
	PurchasePrice  float64   `json:"purchase_price"`
	SalePrice      float64   `json:"sale_price"`
	ProductID      int64     `json:"product_id"`
}

type ProductRecordReport struct {
	ProductId    int64  `json:"product_id"`
	Description  string `json:"description"`
	RecordsCount int64  `json:"records_count"`
}
