package models

import "time"

type PurchaseOrder struct {
	Id              int       `json:"id"`
	OrderNumber     string    `json:"order_number"`
	OrderDate       time.Time `json:"order_date"`
	TrackingCode    string    `json:"tracking_code"`
	BuyerId         int       `json:"buyer_id"`
	ProductRecordId int       `json:"product_record_id"`
}

type PurchaseOrderReport struct {
	Id                 int    `json:"id"`
	IdCardNumber       string `json:"id_card_number"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	PurchaseOrderCount int    `json:"purchase_orders_count"`
}
