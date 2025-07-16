package requests

type WarehouseRequest struct {
	Address            string  `json:"address"`
	Telephone          string  `json:"telephone"`
	WareHouseCode      string  `json:"warehouse_code"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	LocalityId         int     `json:"locality_id"`
}

type WarehousePatchRequest struct {
	Address            *string  `json:"address,omitempty"`
	Telephone          *string  `json:"telephone,omitempty"`
	WareHouseCode      *string  `json:"warehouse_code,omitempty"`
	MinimumCapacity    *int     `json:"minimum_capacity,omitempty"`
	MinimumTemperature *float64 `json:"minimum_temperature,omitempty"`
}
