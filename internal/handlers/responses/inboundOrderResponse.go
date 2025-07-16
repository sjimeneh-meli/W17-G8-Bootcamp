package responses

import "time"

type InboundOrderResponse struct {
	Id             int       `json:"id"`
	OrderNumber    string    `json:"order_number"`
	OrderDate      time.Time `json:"order_date"`
	EmployeeId     int       `json:"employee_id"`
	ProductBatchId int       `json:"product_batch_id"`
	WarehouseId    int       `json:"warehouse_id"`
}
