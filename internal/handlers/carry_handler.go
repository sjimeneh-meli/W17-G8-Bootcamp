package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

// CarryHandler handles HTTP requests for carry operations
// CarryHandler maneja las solicitudes HTTP para operaciones de transporte
type CarryHandler struct {
	carryService services.CarryService // Service layer for carry business logic / Capa de servicio para lógica de negocio de transporte
}

// NewCarryHandler creates and returns a new instance of CarryHandler with the required service
// NewCarryHandler crea y retorna una nueva instancia de CarryHandler con el servicio requerido
func NewCarryHandler(carryService services.CarryService) *CarryHandler {
	return &CarryHandler{carryService: carryService}
}

// Create handles HTTP POST requests to create a new carry
// Validates the request body, creates the carry, and returns appropriate HTTP status codes
// Create maneja las solicitudes HTTP POST para crear un nuevo transporte
// Valida el cuerpo de la solicitud, crea el transporte y retorna códigos de estado HTTP apropiados
func (h *CarryHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	var request requests.CarryRequest

	// Parse and validate JSON request body / Parsear y validar cuerpo de solicitud JSON
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate request structure and business rules / Validar estructura de solicitud y reglas de negocio
	if err := validations.ValidateCarryRequest(request); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Map request to carry model / Mapear solicitud a modelo de transporte
	carry := mappers.MapCarryRequestToCarry(request)

	// Create carry through service layer / Crear transporte a través de la capa de servicio
	newCarry, err := h.carryService.CreateCarry(ctx, carry)

	if err != nil {
		// Handle timeout errors / Manejar errores de timeout
		if ctx.Err() != nil {
			response.Error(w, http.StatusRequestTimeout, "Request timeout cancelled")
			return
		}
		// Handle specific business logic errors / Manejar errores específicos de lógica de negocio
		if errors.Is(err, error_message.ErrAlreadyExists) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error creating carry")
		return
	}

	// Map model to response format / Mapear modelo a formato de respuesta
	carryResponse := mappers.MapCarryToCreateCarryResponse(newCarry)

	response.JSON(w, http.StatusCreated, responses.DataResponse{
		Data: carryResponse,
	})
}

// GetCarryReportByLocality handles HTTP GET requests to retrieve carry reports by locality
// Accepts an optional 'id' query parameter to filter by locality ID
// GetCarryReportByLocality maneja las solicitudes HTTP GET para recuperar reportes de transporte por localidad
// Acepta un parámetro de consulta 'id' opcional para filtrar por ID de localidad
func (h *CarryHandler) GetCarryReportByLocality(w http.ResponseWriter, r *http.Request) {
	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Get and validate optional locality ID parameter / Obtener y validar parámetro opcional de ID de localidad
	localityIDStr := r.URL.Query().Get("id")
	var localityID int

	if localityIDStr != "" {
		// Validate and convert locality ID if provided / Validar y convertir ID de localidad si se proporciona
		var err error
		localityID, err = strconv.Atoi(localityIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid locality ID format")
			return
		}
		if localityID <= 0 {
			response.Error(w, http.StatusBadRequest, "Locality ID must be a positive number")
			return
		}
	}
	// If localityIDStr is empty, localityID will be 0 (report for all localities)
	// Si localityIDStr está vacío, localityID será 0 (reporte para todas las localidades)

	// Get reports from service layer / Obtener reportes de la capa de servicio
	reports, err := h.carryService.GetCarryReportByLocality(ctx, localityID)

	if err != nil {
		// Handle timeout errors / Manejar errores de timeout
		if ctx.Err() != nil {
			response.Error(w, http.StatusRequestTimeout, "Request timeout cancelled")
			return
		}
		// Handle specific business logic errors / Manejar errores específicos de lógica de negocio
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "Locality not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error getting carry report by locality")
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{
		Data: reports,
	})
}
