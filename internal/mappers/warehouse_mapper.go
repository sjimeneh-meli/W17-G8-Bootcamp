package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func ToResponse(warehouse models.Warehouse) responses.WarehouseResponse {
	return responses.WarehouseResponse{
		ID:                 warehouse.Id,
		Address:            warehouse.Address,
		Telephone:          warehouse.Telephone,
		WareHouseCode:      warehouse.WareHouseCode,
		MinimumCapacity:    warehouse.MinimumCapacity,
		MinimumTemperature: warehouse.MinimumTemperature,
	}
}

func ToRequest(warehouseRequest requests.WarehouseRequest) models.Warehouse {
	return models.Warehouse{
		Address:            warehouseRequest.Address,
		Telephone:          warehouseRequest.Telephone,
		WareHouseCode:      warehouseRequest.WareHouseCode,
		MinimumCapacity:    warehouseRequest.MinimumCapacity,
		MinimumTemperature: warehouseRequest.MinimumTemperature,
	}
}

func ApplyPatch(existing models.Warehouse, patch requests.WarehousePatchRequest) models.Warehouse {
	if patch.Address != nil {
		existing.Address = *patch.Address
	}
	if patch.Telephone != nil {
		existing.Telephone = *patch.Telephone
	}
	if patch.WareHouseCode != nil {
		existing.WareHouseCode = *patch.WareHouseCode
	}
	if patch.MinimumCapacity != nil {
		existing.MinimumCapacity = *patch.MinimumCapacity
	}
	if patch.MinimumTemperature != nil {
		existing.MinimumTemperature = *patch.MinimumTemperature
	}
	return existing
}
