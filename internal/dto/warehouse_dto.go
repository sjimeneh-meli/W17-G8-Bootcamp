package dto

// WarehouseResponse representa la respuesta de un warehouse
type WarehouseResponse struct {
	ID            int    `json:"id"`
	Address       string `json:"address"`
	Telephone     string `json:"telephone"`
	WareHouseCode string `json:"warehouse_code"`
}

// WarehouseRequest representa la solicitud para crear/actualizar un warehouse
type WarehouseRequest struct {
	Address       string `json:"address"`
	Telephone     string `json:"telephone"`
	WareHouseCode string `json:"warehouse_code"`
}
