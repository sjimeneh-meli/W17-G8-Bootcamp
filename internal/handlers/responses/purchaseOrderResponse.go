package responses

import "time"

type PurchaseOrderResponse struct {
	OrderNumber     string    `json:"order_number"`
	OrderDate       time.Time `json:"order_date"`
	TrackingCode    string    `json:"tracking_code"`
	BuyerId         int       `json:"buyer_id"`
	ProductRecordId int       `json:"product_record_id"`
}
