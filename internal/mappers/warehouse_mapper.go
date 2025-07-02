package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/dto"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func ToResponse(warehouse models.Warehouse) dto.WarehouseResponse {
	return dto.WarehouseResponse{
		ID:            warehouse.Id,
		Address:       warehouse.Address,
		Telephone:     warehouse.Telephone,
		WareHouseCode: warehouse.WareHouseCode,
	}
}

func ToRequest(warehouseRequest dto.WarehouseRequest) models.Warehouse {
	return models.Warehouse{
		Address:       warehouseRequest.Address,
		Telephone:     warehouseRequest.Telephone,
		WareHouseCode: warehouseRequest.WareHouseCode,
	}
}
