package models

type Warehouse struct {
	Id            int    `json:"id"`
	Address       string `json:"address"`
	Telephone     string `json:"telephone"`
	WareHouseCode string `json:"warehouse_code"`
}
