package handlers

import (
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
)

type WarehouseHandler struct {
	warehouseService services.WarehouseService
}
type WarehouseResponse struct {
	ID            int    `json:"id"`
	Address       string `json:"address"`
	Telephone     string `json:"telephone"`
	WareHouseCode string `json:"warehouse_code"`
}
type dataResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewWarehouseHandler(warehouseService services.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{warehouseService: warehouseService}
}

func (h *WarehouseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	warehouses, err := h.warehouseService.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(warehouses) == 0 {
		response.Error(w, http.StatusNotFound, "No se encontraron almacenes")
		return
	}

	response.JSON(w, http.StatusOK, dataResponse{
		Message: "Almacenes obtenidos correctamente",
		Data:    warehouses,
	})
}
