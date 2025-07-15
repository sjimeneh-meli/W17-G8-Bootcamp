package responses

type EmployeeResponse struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"id_card_number"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}
type EmployeeResponsePost struct {
	CardNumberID string `json:"id_card_number"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}
