package models

import "time"

type InboundOrder struct {
	Id             int       `json:"id"`
	OrderDate      time.Time `json:"order_date"`
	OrderNumber    string    `json:"order_number"`
	EmployeeId     int       `json:"employee_id"`
	ProductBatchId int       `json:"product_batch_id"`
	WarehouseId    int       `json:"warehouse_id"`
}

type InboundOrderReport struct {
	Id                int    `json:"id"`
	IdCardNumber      string `json:"id_card_number"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	InboundOrderCount int    `json:"inbound_orders_count"`
}
