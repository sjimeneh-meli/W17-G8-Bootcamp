package requests

import "time"

type ProductRecordRequest struct {
	ID             int       `json:"id"`
	LastUpdateDate time.Time `json:"last_update_date"`
	PurchasePrice  float64   `json:"purchase_price"`
	SalePrice      float64   `json:"sale_price"`
	ProductID      int       `json:"product_id"`
}
