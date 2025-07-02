package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/dto"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
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

	// Validar que el código de almacén sea único
	if err := h.warehouseService.ValidateCodeUniqueness(warehouseRequest.WareHouseCode); err != nil {
		if errors.Is(err, error_message.ErrEntityExists) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrDatabaseError) {
			response.Error(w, http.StatusInternalServerError, "Error interno del servidor")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	warehouse := mappers.ToRequest(warehouseRequest)
	createdWarehouse, err := h.warehouseService.Create(warehouse)
	if err != nil {
		if errors.Is(err, error_message.ErrDatabaseError) {
			response.Error(w, http.StatusInternalServerError, "Error interno del servidor")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	warehouseResponse := mappers.ToResponse(createdWarehouse)
	response.JSON(w, http.StatusCreated, responses.DataResponse{
		Message: "Almacen creado correctamente",
		Data:    warehouseResponse,
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
		// Determinar el código HTTP basado en el tipo de error
		if errors.Is(err, error_message.ErrEntityNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrDatabaseError) {
			response.Error(w, http.StatusInternalServerError, "Error interno del servidor")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener el warehouse")
		return
	}

	warehouseResponse := mappers.ToResponse(warehouse)
	response.JSON(w, http.StatusOK, responses.DataResponse{
		Message: "Almacen encontrado correctamente",
		Data:    warehouseResponse,
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
		// Determinar el código HTTP basado en el tipo de error
		if errors.Is(err, error_message.ErrEntityNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrDatabaseError) {
			response.Error(w, http.StatusInternalServerError, "Error interno del servidor")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar el warehouse")
		return
	}

	response.JSON(w, http.StatusNoContent, responses.DataResponse{
		Message: "Almacen eliminado correctamente",
		Data:    nil,
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
	updatedWarehouse, err := h.warehouseService.Update(id, warehouse)
	if err != nil {
		if errors.Is(err, error_message.ErrEntityNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrEntityExists) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrDatabaseError) {
			response.Error(w, http.StatusInternalServerError, "Error interno del servidor")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al actualizar el warehouse")
		return
	}

	warehouseResponse := mappers.ToResponse(updatedWarehouse)
	response.JSON(w, http.StatusOK, responses.DataResponse{
		Message: "Almacen actualizado correctamente",
		Data:    warehouseResponse,
	})
}
