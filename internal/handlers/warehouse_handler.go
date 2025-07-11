package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
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
		if errors.Is(err, error_message.ErrInternalServerError) {
			response.Error(w, http.StatusInternalServerError, "Error al leer la base de datos de warehouses")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener los warehouses")
		return
	}

	if len(warehouses) == 0 {
		response.Error(w, http.StatusNotFound, "No se encontraron almacenes")
		return
	}

	warehouseResponses := make([]responses.WarehouseResponse, 0)
	for _, warehouse := range warehouses {
		warehouseResponses = append(warehouseResponses, mappers.ToResponse(warehouse))
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{
		Data: warehouseResponses,
	})
}

func (h *WarehouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var warehouseRequest requests.WarehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&warehouseRequest); err != nil {
		response.Error(w, http.StatusBadRequest, "Formato JSON inválido")
		return
	}

	if err := validations.ValidateWarehouseRequestStruct(warehouseRequest); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.warehouseService.ValidateCodeUniqueness(warehouseRequest.WareHouseCode); err != nil {
		if errors.Is(err, error_message.ErrAlreadyExists) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrInternalServerError) {
			response.Error(w, http.StatusInternalServerError, "Error al validar la unicidad del código en la base de datos")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al validar el código del warehouse")
		return
	}

	warehouse := mappers.ToRequest(warehouseRequest)
	createdWarehouse, err := h.warehouseService.Create(warehouse)
	if err != nil {
		if errors.Is(err, error_message.ErrInternalServerError) {
			response.Error(w, http.StatusInternalServerError, "Error al guardar el warehouse en la base de datos")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al crear el warehouse")
		return
	}
	warehouseResponse := mappers.ToResponse(createdWarehouse)
	response.JSON(w, http.StatusCreated, responses.DataResponse{
		Data: warehouseResponse,
	})
}

func (h *WarehouseHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.Error(w, http.StatusBadRequest, "El ID del almacén es requerido")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "El ID del almacén debe ser un número")
		return
	}

	warehouse, err := h.warehouseService.GetById(id)
	if err != nil {
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrInternalServerError) {
			response.Error(w, http.StatusInternalServerError, "Error al buscar el warehouse en la base de datos")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener el warehouse")
		return
	}

	warehouseResponse := mappers.ToResponse(warehouse)
	response.JSON(w, http.StatusOK, responses.DataResponse{
		Data: warehouseResponse,
	})
}

func (h *WarehouseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.Error(w, http.StatusBadRequest, "El ID del almacén es requerido")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "El ID del almacén debe ser un número")
		return
	}

	if err := h.warehouseService.Delete(id); err != nil {
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrInternalServerError) {
			response.Error(w, http.StatusInternalServerError, "Error al eliminar el warehouse de la base de datos")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar el warehouse")
		return
	}

	response.JSON(w, http.StatusNoContent, responses.DataResponse{
		Data: nil,
	})
}

func (h *WarehouseHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.Error(w, http.StatusBadRequest, "El ID del almacén es requerido")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "El ID del almacén debe ser un número")
		return
	}

	existingWarehouse, err := h.warehouseService.GetById(id)
	if err != nil {
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrInternalServerError) {
			response.Error(w, http.StatusInternalServerError, "Error al buscar el warehouse en la base de datos")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener el warehouse")
		return
	}

	var warehousePatchRequest requests.WarehousePatchRequest
	if err := json.NewDecoder(r.Body).Decode(&warehousePatchRequest); err != nil {
		response.Error(w, http.StatusBadRequest, "Formato JSON inválido")
		return
	}

	if err := validations.ValidateWarehousePatchRequest(warehousePatchRequest); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedWarehouse := mappers.ApplyPatch(existingWarehouse, warehousePatchRequest)

	if warehousePatchRequest.WareHouseCode != nil && *warehousePatchRequest.WareHouseCode != existingWarehouse.WareHouseCode {
		if err := h.warehouseService.ValidateCodeUniqueness(*warehousePatchRequest.WareHouseCode); err != nil {
			if errors.Is(err, error_message.ErrAlreadyExists) {
				response.Error(w, http.StatusConflict, err.Error())
				return
			}
			if errors.Is(err, error_message.ErrInternalServerError) {
				response.Error(w, http.StatusInternalServerError, "Error al validar la unicidad del código en la base de datos")
				return
			}
			response.Error(w, http.StatusInternalServerError, "Error al validar el código del warehouse")
			return
		}
	}

	finalWarehouse, err := h.warehouseService.Update(id, updatedWarehouse)
	if err != nil {
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrAlreadyExists) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrInternalServerError) {
			response.Error(w, http.StatusInternalServerError, "Error al actualizar el warehouse en la base de datos")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al actualizar el warehouse")
		return
	}

	warehouseResponse := mappers.ToResponse(finalWarehouse)
	response.JSON(w, http.StatusOK, responses.DataResponse{
		Data: warehouseResponse,
	})
}
