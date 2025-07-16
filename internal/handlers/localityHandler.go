package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

// LocalityHandler handles HTTP requests for locality operations
// LocalityHandler maneja las solicitudes HTTP para operaciones de localidad
type LocalityHandler struct {
	service services.LocalityService // Service layer for locality business logic / Capa de servicio para lógica de negocio de localidad
}

// NewLocalityHandler creates and returns a new instance of LocalityHandler with the required service
// NewLocalityHandler crea y retorna una nueva instancia de LocalityHandler con el servicio requerido
func NewLocalityHandler(service services.LocalityService) *LocalityHandler {
	return &LocalityHandler{service: service}
}

// Save handles HTTP POST requests to create a new locality
// Validates the request body and returns appropriate HTTP status codes
// Save maneja las solicitudes HTTP POST para crear una nueva localidad
// Valida el cuerpo de la solicitud y retorna códigos de estado HTTP apropiados
func (h *LocalityHandler) Save(w http.ResponseWriter, r *http.Request) {
	var localityToCreate requests.LocalityRequest

	// Parse and validate JSON request body / Parsear y validar cuerpo de solicitud JSON
	errorBody := json.NewDecoder(r.Body).Decode(&localityToCreate)
	if errorBody != nil {
		response.Error(w, http.StatusBadRequest, errorBody.Error())
		return
	}

	// Extract data from request / Extraer datos de la solicitud
	data := localityToCreate.Data
	fmt.Println(data)

	// Validate request structure and business rules / Validar estructura de solicitud y reglas de negocio
	if err := validations.ValidateLocalityRequestStruct(data); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	parentCtx := context.Background()
	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	defer cancel()

	// Create locality through service layer / Crear localidad a través de la capa de servicio
	localityCreated, err := h.service.Save(ctx, data)
	if err != nil {
		// Handle specific error types / Manejar tipos de error específicos
		if errors.Is(err, error_message.ErrQuery) {
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		if errors.Is(err, error_message.ErrAlreadyExists) {
			response.Error(w, http.StatusConflict, err.Error())
		}
		if errors.Is(err, context.DeadlineExceeded) {
			response.Error(w, http.StatusGatewayTimeout, err.Error())
		}
		return
	}
	response.JSON(w, http.StatusOK, responses.DataResponse{Data: localityCreated})
}

// GetSellerReportByLocality handles HTTP GET requests to retrieve seller reports by locality
// Accepts an optional 'id' query parameter to filter by locality ID
// GetSellerReportByLocality maneja las solicitudes HTTP GET para recuperar reportes de vendedores por localidad
// Acepta un parámetro de consulta 'id' opcional para filtrar por ID de localidad
func (h *LocalityHandler) GetSellerReportByLocality(w http.ResponseWriter, r *http.Request) {
	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	parentCtx := context.Background()
	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	defer cancel()

	// Get and validate optional locality ID query parameter / Obtener y validar parámetro opcional de consulta ID de localidad
	localityIdStr := r.URL.Query().Get("id")
	var localityId int

	// Validate query parameter format / Validar formato del parámetro de consulta
	if r.URL.RawQuery != "" && localityIdStr == "" {
		response.Error(w, http.StatusInternalServerError, "Invalid query parameter")
		return
	}

	if localityIdStr == "" {
		// If no parameter provided, use 0 as default (all localities) / Si no se proporciona parámetro, usar 0 por defecto (todas las localidades)
		localityId = 0
	} else {
		// If parameter provided, attempt to convert to int / Si se proporciona parámetro, intentar convertir a int
		var err error
		localityId, err = strconv.Atoi(localityIdStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	// Get seller reports from service layer / Obtener reportes de vendedores de la capa de servicio
	result, err := h.service.GetSellerReports(ctx, localityId)

	// Handle specific error types / Manejar tipos de error específicos
	if errors.Is(err, error_message.ErrFailedCheckingExistence) || errors.Is(err, error_message.ErrQueryingReport) || errors.Is(err, error_message.ErrFailedToScan) {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if errors.Is(err, error_message.ErrNotFound) {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		response.Error(w, http.StatusGatewayTimeout, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{result})
}
