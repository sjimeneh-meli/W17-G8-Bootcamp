package models

// Employee representa el modelo interno de un empleado
type Employee struct {
	ID           int
	CardNumberID string
	FirstName    string
	LastName     string
	WarehouseID  int
	// warehouse tabla relacionada req2
}

//type EmployeeDoc struct {
//	ID           int    `json:"id"`
//	CardNumberID string `json:"card_number_id"`
//	FirstName    string `json:"first_name"`
//	LastName     string `json:"last_name"`
//	WarehouseID  int    `json:"warehouse_id"`
//}
