package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

// GetSectionHandler creates and returns a new instance of SectionHandler with required services and validation
// GetSectionHandler crea y retorna una nueva instancia de SectionHandler con los servicios y validación requeridos
func GetSectionHandler(service services.SectionServiceI, warehouseService services.WarehouseService, validation *validations.SectionValidation) SectionHandlerI {
	return &SectionHandler{
		service:          service,
		warehouseService: warehouseService,
		validation:       validation,
	}
}

// SectionHandlerI defines the contract for section HTTP handlers with RESTful operations
// SectionHandlerI define el contrato para los manejadores HTTP de secciones con operaciones RESTful
type SectionHandlerI interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	DeleteByID(w http.ResponseWriter, r *http.Request)
}

// SectionHandler implements SectionHandlerI and handles HTTP requests for section operations
// SectionHandler implementa SectionHandlerI y maneja las solicitudes HTTP para operaciones de secciones
type SectionHandler struct {
	service          services.SectionServiceI       // Service layer for section business logic / Capa de servicio para lógica de negocio de secciones
	warehouseService services.WarehouseService      // Service layer for warehouse validation / Capa de servicio para validación de almacenes
	validation       *validations.SectionValidation // Validation layer for section requests / Capa de validación para solicitudes de secciones
}

// GetAll handles HTTP GET requests to retrieve all sections
// Returns a JSON response with all sections or appropriate error codes
// GetAll maneja las solicitudes HTTP GET para recuperar todas las secciones
// Retorna una respuesta JSON con todas las secciones o códigos de error apropiados
func (h *SectionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var responseJson *responses.DataResponse = &responses.DataResponse{}

	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Get all sections from service layer / Obtener todas las secciones de la capa de servicio
	sections, srvErr := h.service.GetAll(ctx)
	if srvErr != nil {
		response.Error(w, http.StatusNotFound, srvErr.Error())
		return
	}

	// Map models to response format / Mapear modelos a formato de respuesta
	responseJson.Data = mappers.GetListSectionResponseFromListModel(sections)
	response.JSON(w, http.StatusOK, responseJson)
}

// GetByID handles HTTP GET requests to retrieve a section by ID
// Extracts the ID from the URL parameter and returns the section data or appropriate error codes
// GetByID maneja las solicitudes HTTP GET para recuperar una sección por ID
// Extrae el ID del parámetro de URL y retorna los datos de la sección o códigos de error apropiados
func (h *SectionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	var responseJson *responses.DataResponse = &responses.DataResponse{}

	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
	idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
	if convErr != nil {
		response.Error(w, http.StatusExpectationFailed, convErr.Error())
		return
	}

	// Get section by ID from service layer / Obtener sección por ID de la capa de servicio
	section, srvErr := h.service.GetByID(ctx, idParam)
	if srvErr != nil {
		response.Error(w, http.StatusNotFound, srvErr.Error())
		return
	}

	// Map model to response format / Mapear modelo a formato de respuesta
	responseJson.Data = mappers.GetSectionResponseFromModel(section)
	response.JSON(w, http.StatusOK, responseJson)
}

// Create handles HTTP POST requests to create a new section
// Validates dependencies (warehouse) and creates the section with business rules validation
// Create maneja las solicitudes HTTP POST para crear una nueva sección
// Valida dependencias (almacén) y crea la sección con validación de reglas de negocio
func (h *SectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		request      *requests.SectionRequest = &requests.SectionRequest{}
		responseJson *responses.DataResponse  = &responses.DataResponse{}
		section      *models.Section
	)

	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Parse and validate JSON request body / Parsear y validar cuerpo de solicitud JSON
	if reqErr := json.NewDecoder(r.Body).Decode(request); reqErr != nil {
		response.Error(w, http.StatusExpectationFailed, reqErr.Error())
		return
	}

	// Validate request structure and business rules / Validar estructura de solicitud y reglas de negocio
	if valErr := h.validation.ValidateSectionRequestStruct(*request); valErr != nil {
		response.Error(w, http.StatusUnprocessableEntity, valErr.Error())
		return
	}

	// Validate that warehouse exists / Validar que el almacén exista
	_, srvErr := h.warehouseService.GetById(ctx, request.WarehouseID)
	if srvErr != nil {
		response.Error(w, http.StatusNotFound, srvErr.Error())
		return
	}

	// Map request to section model / Mapear solicitud a modelo de sección
	section = mappers.GetSectionModelFromRequest(request)

	// Validate section number uniqueness / Validar unicidad del número de sección
	if h.service.ExistsWithSectionNumber(ctx, section.Id, section.SectionNumber) {
		response.Error(w, http.StatusConflict, "already exist a section with the same number")
		return
	}

	// Create section through service layer / Crear sección a través de la capa de servicio
	if srvErr := h.service.Create(ctx, section); srvErr != nil {
		response.Error(w, http.StatusExpectationFailed, srvErr.Error())
		return
	}

	// Map model to response format / Mapear modelo a formato de respuesta
	responseJson.Data = mappers.GetSectionResponseFromModel(section)
	response.JSON(w, http.StatusCreated, responseJson)
}

// Update handles HTTP PUT requests to update an existing section
// Extracts the ID from the URL parameter and updates the section with validation
// Update maneja las solicitudes HTTP PUT para actualizar una sección existente
// Extrae el ID del parámetro de URL y actualiza la sección con validación
func (h *SectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		request      *requests.SectionRequest = &requests.SectionRequest{}
		responseJson *responses.DataResponse  = &responses.DataResponse{}
	)

	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
	idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
	if convErr != nil {
		response.Error(w, http.StatusExpectationFailed, convErr.Error())
		return
	}

	// Get existing section by ID / Obtener sección existente por ID
	section, srvErr := h.service.GetByID(ctx, idParam)
	if srvErr != nil {
		response.Error(w, http.StatusNotFound, srvErr.Error())
		return
	}

	// Parse and validate JSON request body / Parsear y validar cuerpo de solicitud JSON
	if reqErr := json.NewDecoder(r.Body).Decode(request); reqErr != nil {
		response.Error(w, http.StatusExpectationFailed, reqErr.Error())
		return
	}

	// Validate request structure and business rules / Validar estructura de solicitud y reglas de negocio
	if valErr := h.validation.ValidateSectionRequestStruct(*request); valErr != nil {
		response.Error(w, http.StatusUnprocessableEntity, valErr.Error())
		return
	}

	// Validate section number uniqueness for update / Validar unicidad del número de sección para actualización
	if h.service.ExistsWithSectionNumber(ctx, section.Id, request.SectionNumber) {
		response.Error(w, http.StatusConflict, "already exist a section with the same number")
		return
	}

	// Update section model with request data / Actualizar modelo de sección con datos de la solicitud
	mappers.UpdateSectionModelFromRequest(section, request)
	if srvErr := h.service.Update(ctx, section); srvErr != nil {
		response.Error(w, http.StatusExpectationFailed, srvErr.Error())
		return
	}

	// Map model to response format / Mapear modelo a formato de respuesta
	responseJson.Data = mappers.GetSectionResponseFromModel(section)
	response.JSON(w, http.StatusOK, responseJson)
}

// DeleteByID handles HTTP DELETE requests to remove a section by ID
// Extracts the ID from the URL parameter and deletes the section
// DeleteByID maneja las solicitudes HTTP DELETE para eliminar una sección por ID
// Extrae el ID del parámetro de URL y elimina la sección
func (h *SectionHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
	idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
	if convErr != nil {
		response.Error(w, http.StatusExpectationFailed, convErr.Error())
		return
	}

	// Delete section through service layer / Eliminar sección a través de la capa de servicio
	srvErr := h.service.DeleteByID(ctx, idParam)
	if srvErr != nil {
		response.Error(w, http.StatusNotFound, srvErr.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
