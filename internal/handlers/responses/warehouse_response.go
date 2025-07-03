package responses

type WarehouseResponse struct {
	ID                 int     `json:"id"`
	Address            string  `json:"address"`
	Telephone          string  `json:"telephone"`
	WareHouseCode      string  `json:"warehouse_code"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MinimumTemperature float64 `json:"minimum_temperature"`
}
