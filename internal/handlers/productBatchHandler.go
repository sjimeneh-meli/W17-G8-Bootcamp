package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

// GetProductBatchHandler creates and returns a new instance of ProductBatchHandler with required services and validation
// GetProductBatchHandler crea y retorna una nueva instancia de ProductBatchHandler con los servicios y validación requeridos
func GetProductBatchHandler(service services.ProductBatchServiceI,
	sectionService services.SectionServiceI,
	productService services.ProductService,
	validation validations.ProductBatchValidation) ProductBatchHandlerI {

	return &ProductBatchHandler{
		service:        service,
		sectionService: sectionService,
		productService: productService,
		validation:     &validation,
	}
}

// ProductBatchHandlerI defines the contract for product batch HTTP handlers
// ProductBatchHandlerI define el contrato para los manejadores HTTP de lotes de productos
type ProductBatchHandlerI interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetReportProduct(w http.ResponseWriter, r *http.Request)
}

// ProductBatchHandler implements ProductBatchHandlerI and handles HTTP requests for product batch operations
// ProductBatchHandler implementa ProductBatchHandlerI y maneja las solicitudes HTTP para operaciones de lotes de productos
type ProductBatchHandler struct {
	service        services.ProductBatchServiceI       // Service layer for product batch business logic / Capa de servicio para lógica de negocio de lotes de productos
	sectionService services.SectionServiceI            // Service layer for section validation / Capa de servicio para validación de secciones
	productService services.ProductService             // Service layer for product validation / Capa de servicio para validación de productos
	validation     *validations.ProductBatchValidation // Validation layer for product batch requests / Capa de validación para solicitudes de lotes de productos
}

// Create handles HTTP POST requests to create a new product batch
// Validates dependencies (section and product) and creates the batch with business rules validation
// Create maneja las solicitudes HTTP POST para crear un nuevo lote de productos
// Valida dependencias (sección y producto) y crea el lote con validación de reglas de negocio
func (h *ProductBatchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		request      *requests.ProductBatchRequest = &requests.ProductBatchRequest{}
		responseJson *responses.DataResponse       = &responses.DataResponse{}
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
	if valErr := h.validation.ValidateProductBatchRequestStruc(*request); valErr != nil {
		response.Error(w, http.StatusUnprocessableEntity, valErr.Error())
		return
	}

	// Validate that section exists / Validar que la sección exista
	if !h.sectionService.ExistWithID(ctx, request.SectionID) {
		response.Error(w, http.StatusNotFound, "section not found")
		return
	}

	// Validate that product exists / Validar que el producto exista
	exists, _ := h.productService.ExistById(ctx, int64(request.ProductID))
	if !exists {
		response.Error(w, http.StatusNotFound, "product not found")
		return
	}

	// Map request to product batch model / Mapear solicitud a modelo de lote de productos
	productBatch, mapErr := mappers.GetProductBatchModelFromRequest(request)
	if mapErr != nil {
		response.Error(w, http.StatusExpectationFailed, mapErr.Error())
		return
	}

	// Validate batch number uniqueness / Validar unicidad del número de lote
	if h.service.ExistsWithBatchNumber(ctx, productBatch.Id, productBatch.BatchNumber) {
		response.Error(w, http.StatusConflict, "already exist a batch with the same number")
		return
	}

	// Create product batch through service layer / Crear lote de productos a través de la capa de servicio
	if srvErr := h.service.Create(ctx, productBatch); srvErr != nil {
		response.Error(w, http.StatusExpectationFailed, srvErr.Error())
		return
	}

	// Map model to response format / Mapear modelo a formato de respuesta
	responseJson.Data = mappers.GetProductBatchResponseFromModel(productBatch)
	response.JSON(w, http.StatusCreated, responseJson)
}

// GetReportProduct handles HTTP GET requests to retrieve product quantity reports by section
// Accepts an optional 'id' query parameter to filter by section ID
// GetReportProduct maneja las solicitudes HTTP GET para recuperar reportes de cantidad de productos por sección
// Acepta un parámetro de consulta 'id' opcional para filtrar por ID de sección
func (h *ProductBatchHandler) GetReportProduct(w http.ResponseWriter, r *http.Request) {
	var responseJson *responses.DataResponse = &responses.DataResponse{}

	// Set timeout context for the request / Establecer contexto con timeout para la solicitud
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Get and validate optional section ID query parameter / Obtener y validar parámetro opcional de consulta ID de sección
	idParamString := r.URL.Query().Get("id")
	if idParamString != "" {
		// Parse and validate section ID / Parsear y validar ID de sección
		idParam, convErr := strconv.Atoi(idParamString)
		if convErr != nil {
			response.Error(w, http.StatusExpectationFailed, convErr.Error())
			return
		}

		// Get section by ID to validate existence / Obtener sección por ID para validar existencia
		section, srvErr := h.sectionService.GetByID(ctx, idParam)
		if srvErr != nil {
			response.Error(w, http.StatusNotFound, srvErr.Error())
			return
		}

		// Get product quantity for specific section / Obtener cantidad de productos para sección específica
		quantity := h.service.GetProductQuantityBySectionId(ctx, section.Id)
		responseJson.Data = map[string]any{
			"section_id":     section.Id,
			"section_number": section.SectionNumber,
			"products_count": quantity,
		}
		response.JSON(w, http.StatusCreated, responseJson)

	} else {
		// Get reports for all sections / Obtener reportes para todas las secciones
		sections, srvErr := h.sectionService.GetAll(ctx)
		data := make([]map[string]any, 0, len(sections))
		if srvErr != nil {
			response.Error(w, http.StatusExpectationFailed, srvErr.Error())
			return
		}

		// Build report data for each section / Construir datos de reporte para cada sección
		for _, s := range sections {
			quantity := h.service.GetProductQuantityBySectionId(ctx, s.Id)
			data = append(data, map[string]any{
				"section_id":     s.Id,
				"section_number": s.SectionNumber,
				"products_count": quantity,
			})
		}
		responseJson.Data = data

		response.JSON(w, http.StatusCreated, responseJson)
	}
}
