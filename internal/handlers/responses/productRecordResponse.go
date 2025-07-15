package responses

import "time"

type ProductRecordResponse struct {
	LastUpdateDate time.Time `json:"last_update_date"`
	PurchasePrice  float64   `json:"purchase_price"`
	SalePrice      float64   `json:"sale_price"`
	ProductID      int64     `json:"product_id"`
}
