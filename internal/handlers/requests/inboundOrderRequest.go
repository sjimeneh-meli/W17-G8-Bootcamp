package requests

import "time"

type InboundOrderRequest struct {
	Data InboundOrderAttributes `json:"data"`
}

type InboundOrderAttributes struct {
	OrderDate      time.Time `json:"order_date"`
	OrderNumber    string    `json:"order_number"`
	EmployeeId     int       `json:"employee_id"`
	ProductBatchId int       `json:"product_batch_id"`
	WarehouseId    int       `json:"warehouse_id"`
}
