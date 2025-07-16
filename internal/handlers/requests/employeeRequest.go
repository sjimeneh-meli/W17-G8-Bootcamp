package requests

type EmployeeRequest struct {
	CardNumberID string `json:"id_card_number"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}
