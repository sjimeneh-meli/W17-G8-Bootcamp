package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

// WarehouseHandler handles HTTP requests for warehouse operations
// WarehouseHandler maneja las solicitudes HTTP para operaciones de almacenes
type WarehouseHandler struct {
	warehouseService services.WarehouseService // Service layer for warehouse business logic / Capa de servicio para lógica de negocio de almacenes
}

// NewWarehouseHandler creates and returns a new instance of WarehouseHandler with the required service
// NewWarehouseHandler crea y retorna una nueva instancia de WarehouseHandler con el servicio requerido
func NewWarehouseHandler(warehouseService services.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{warehouseService: warehouseService}
}

// GetAll handles HTTP GET requests to retrieve all warehouses
// Returns a JSON response with all warehouses or appropriate error codes
// GetAll maneja las solicitudes HTTP GET para recuperar todos los almacenes
// Retorna una respuesta JSON con todos los almacenes o códigos de error apropiados
func (h *WarehouseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Get all warehouses from service layer / Obtener todos los almacenes de la capa de servicio
	warehouses, err := h.warehouseService.GetAll(ctx)
	if err != nil {
		// Handle timeout errors / Manejar errores de timeout
		if ctx.Err() != nil {
			response.Error(w, http.StatusRequestTimeout, "Request timeout cancelled")
			return
		}
		// Handle specific error types / Manejar tipos de error específicos
		if errors.Is(err, error_message.ErrInternalServerError) {
			response.Error(w, http.StatusInternalServerError, "Error al leer la base de datos de warehouses")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener los warehouses")
		return
	}

	// Check if any warehouses were found / Verificar si se encontraron almacenes
	if len(warehouses) == 0 {
		response.Error(w, http.StatusNotFound, "No se encontraron almacenes")
		return
	}

	// Map models to response format / Mapear modelos a formato de respuesta
	warehouseResponses := make([]responses.WarehouseResponse, 0)
	for _, warehouse := range warehouses {
		warehouseResponses = append(warehouseResponses, mappers.ToResponse(warehouse))
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{
		Data: warehouseResponses,
	})
}

// Create handles HTTP POST requests to create a new warehouse
// Validates the request body, code uniqueness, and returns appropriate HTTP status codes
// Create maneja las solicitudes HTTP POST para crear un nuevo almacén
// Valida el cuerpo de la solicitud, unicidad del código y retorna códigos de estado HTTP apropiados
func (h *WarehouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	var warehouseRequest requests.WarehouseRequest

	// Parse and validate JSON request body / Parsear y validar cuerpo de solicitud JSON
	if err := json.NewDecoder(r.Body).Decode(&warehouseRequest); err != nil {
		response.Error(w, http.StatusBadRequest, "Formato JSON inválido")
		return
	}

	// Validate request structure and business rules / Validar estructura de solicitud y reglas de negocio
	if err := validations.ValidateWarehouseRequestStruct(warehouseRequest); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate warehouse code uniqueness / Validar unicidad del código de almacén
	if err := h.warehouseService.ValidateCodeUniqueness(ctx, warehouseRequest.WareHouseCode); err != nil {
		// Handle timeout errors / Manejar errores de timeout
		if ctx.Err() != nil {
			response.Error(w, http.StatusRequestTimeout, "Request timeout cancelled")
			return
		}
		// Handle specific error types / Manejar tipos de error específicos
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

	// Map request to warehouse model and create through service / Mapear solicitud a modelo de almacén y crear a través del servicio
	warehouse := mappers.ToRequest(warehouseRequest)
	createdWarehouse, err := h.warehouseService.Create(ctx, warehouse)
	if err != nil {
		// Handle timeout errors / Manejar errores de timeout
		if ctx.Err() != nil {
			response.Error(w, http.StatusRequestTimeout, "Request timeout cancelled")
			return
		}
		// Handle specific error types / Manejar tipos de error específicos
		if errors.Is(err, error_message.ErrInternalServerError) {
			response.Error(w, http.StatusInternalServerError, "Error al guardar el warehouse en la base de datos")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al crear el warehouse")
		return
	}

	// Map model to response format / Mapear modelo a formato de respuesta
	warehouseResponse := mappers.ToResponse(createdWarehouse)
	response.JSON(w, http.StatusCreated, responses.DataResponse{
		Data: warehouseResponse,
	})
}

// GetById handles HTTP GET requests to retrieve a warehouse by ID
// Extracts the ID from the URL parameter and returns the warehouse data or appropriate error codes
// GetById maneja las solicitudes HTTP GET para recuperar un almacén por ID
// Extrae el ID del parámetro de URL y retorna los datos del almacén o códigos de error apropiados
func (h *WarehouseHandler) GetById(w http.ResponseWriter, r *http.Request) {
	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.Error(w, http.StatusBadRequest, "El ID del almacén es requerido")
		return
	}

	// Parse and validate ID parameter / Parsear y validar parámetro ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "El ID del almacén debe ser un número")
		return
	}

	// Get warehouse by ID from service layer / Obtener almacén por ID de la capa de servicio
	warehouse, err := h.warehouseService.GetById(ctx, id)
	if err != nil {
		// Handle timeout errors / Manejar errores de timeout
		if ctx.Err() != nil {
			response.Error(w, http.StatusRequestTimeout, "Request timeout cancelled")
			return
		}
		// Handle specific error types / Manejar tipos de error específicos
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

	// Map model to response format / Mapear modelo a formato de respuesta
	warehouseResponse := mappers.ToResponse(warehouse)
	response.JSON(w, http.StatusOK, responses.DataResponse{
		Data: warehouseResponse,
	})
}

// Delete handles HTTP DELETE requests to remove a warehouse by ID
// Extracts the ID from the URL parameter and deletes the warehouse
// Delete maneja las solicitudes HTTP DELETE para eliminar un almacén por ID
// Extrae el ID del parámetro de URL y elimina el almacén
func (h *WarehouseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.Error(w, http.StatusBadRequest, "El ID del almacén es requerido")
		return
	}

	// Parse and validate ID parameter / Parsear y validar parámetro ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "El ID del almacén debe ser un número")
		return
	}

	// Delete warehouse through service layer / Eliminar almacén a través de la capa de servicio
	if err := h.warehouseService.Delete(ctx, id); err != nil {
		// Handle timeout errors / Manejar errores de timeout
		if ctx.Err() != nil {
			response.Error(w, http.StatusRequestTimeout, "Request timeout cancelled")
			return
		}
		// Handle specific error types / Manejar tipos de error específicos
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

// Update handles HTTP PUT requests to update an existing warehouse
// Extracts the ID from the URL parameter, validates code uniqueness, and updates the warehouse
// Update maneja las solicitudes HTTP PUT para actualizar un almacén existente
// Extrae el ID del parámetro de URL, valida unicidad del código y actualiza el almacén
func (h *WarehouseHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.Error(w, http.StatusBadRequest, "El ID del almacén es requerido")
		return
	}

	// Parse and validate ID parameter / Parsear y validar parámetro ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "El ID del almacén debe ser un número")
		return
	}

	// Get existing warehouse by ID for comparison / Obtener almacén existente por ID para comparación
	existingWarehouse, err := h.warehouseService.GetById(ctx, id)
	if err != nil {
		// Handle timeout errors / Manejar errores de timeout
		if ctx.Err() != nil {
			response.Error(w, http.StatusRequestTimeout, "Request timeout cancelled")
			return
		}
		// Handle specific error types / Manejar tipos de error específicos
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

	// Parse and validate JSON request body / Parsear y validar cuerpo de solicitud JSON
	if err := json.NewDecoder(r.Body).Decode(&warehousePatchRequest); err != nil {
		response.Error(w, http.StatusBadRequest, "Formato JSON inválido")
		return
	}

	// Validate request structure for patch operation / Validar estructura de solicitud para operación patch
	if err := validations.ValidateWarehousePatchRequest(warehousePatchRequest); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Apply patch to existing warehouse / Aplicar patch al almacén existente
	updatedWarehouse := mappers.ApplyPatch(existingWarehouse, warehousePatchRequest)

	// Validate code uniqueness only if code has changed / Validar unicidad del código solo si el código ha cambiado
	if warehousePatchRequest.WareHouseCode != nil && *warehousePatchRequest.WareHouseCode != existingWarehouse.WareHouseCode {
		if err := h.warehouseService.ValidateCodeUniqueness(ctx, *warehousePatchRequest.WareHouseCode); err != nil {
			// Handle timeout errors / Manejar errores de timeout
			if ctx.Err() != nil {
				response.Error(w, http.StatusRequestTimeout, "Request timeout cancelled")
				return
			}
			// Handle specific error types / Manejar tipos de error específicos
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

	// Update warehouse through service layer / Actualizar almacén a través de la capa de servicio
	updatedWarehouse, err = h.warehouseService.Update(ctx, id, updatedWarehouse)
	if err != nil {
		// Handle timeout errors / Manejar errores de timeout
		if ctx.Err() != nil {
			response.Error(w, http.StatusRequestTimeout, "Request timeout cancelled")
			return
		}
		// Handle specific error types / Manejar tipos de error específicos
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrInternalServerError) {
			response.Error(w, http.StatusInternalServerError, "Error al actualizar el warehouse en la base de datos")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al actualizar el warehouse")
		return
	}

	// Map model to response format / Mapear modelo a formato de respuesta
	warehouseResponse := mappers.ToResponse(updatedWarehouse)
	response.JSON(w, http.StatusOK, responses.DataResponse{
		Data: warehouseResponse,
	})
}
