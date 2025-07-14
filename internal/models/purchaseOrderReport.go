package models

type PurchaseOrderReport struct {
	Id                 int    `json:"id"`
	IdCardNumber       int    `json:"id_card_number"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	PurchaseOrderCount int    `json:"purchase_orders_count"`
}
