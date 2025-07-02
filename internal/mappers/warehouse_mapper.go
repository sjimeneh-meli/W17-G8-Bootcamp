package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func ToResponse(warehouse models.Warehouse) handlers.WarehouseResponse {
	return handlers.WarehouseResponse{
		ID:            warehouse.Id,
		Address:       warehouse.Address,
		Telephone:     warehouse.Telephone,
		WareHouseCode: warehouse.WareHouseCode,
	}
}
