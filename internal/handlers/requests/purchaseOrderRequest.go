package requests

import "time"

type PurchaseOrderRequest struct {
	Data PurchaseOrderAttributes `json:"data"`
}

type PurchaseOrderAttributes struct {
	OrderNumber     string    `json:"order_number"`
	OrderDate       time.Time `json:"order_date"`
	TrackingCode    string    `json:"tracking_code"`
	BuyerId         int       `json:"buyer_id"`
	ProductRecordId int       `json:"product_record_id"`
}
