package requests

type SectionRequest struct {
	SectionNumber      string  `json:"section_number"`
	CurrentCapacity    int     `json:"current_capacity"`
	CurrentTemperature float64 `json:"current_temperature"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	ProductTypeID      int     `json:"product_type_id"`
	WarehouseID        int     `json:"warehouse_id"`
}
