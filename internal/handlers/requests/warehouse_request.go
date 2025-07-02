package requests

type WarehouseRequest struct {
	Address            string  `json:"address"`
	Telephone          string  `json:"telephone"`
	WareHouseCode      string  `json:"warehouse_code"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MinimumTemperature float64 `json:"minimum_temperature"`
}
