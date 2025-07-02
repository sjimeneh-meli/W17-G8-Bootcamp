package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/dto"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

type WarehouseHandler struct {
	warehouseService services.WarehouseService
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

	// Mapear a DTOs para la respuesta
	warehouseResponses := make(map[int]dto.WarehouseResponse)
	for id, warehouse := range warehouses {
		warehouseResponses[id] = mappers.ToResponse(warehouse)
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{
		Message: "Almacenes obtenidos correctamente",
		Data:    warehouseResponses,
	})
}

func (h *WarehouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var warehouseRequest dto.WarehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&warehouseRequest); err != nil {
		response.Error(w, http.StatusBadRequest, "Formato JSON inválido")
		return
	}

	if err := validations.ValidateWarehouseRequestStruct(warehouseRequest); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	warehouse := mappers.ToRequest(warehouseRequest)
	createdWarehouse, err := h.warehouseService.Create(warehouse)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	warehouseResponse := mappers.ToResponse(createdWarehouse)
	response.JSON(w, http.StatusCreated, responses.DataResponse{
		Message: "Almacen creado correctamente",
		Data:    warehouseResponse,
	})
}
