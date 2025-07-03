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
