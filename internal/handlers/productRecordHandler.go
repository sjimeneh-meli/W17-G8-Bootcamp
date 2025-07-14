package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

// ProductRecordHandlerI - Interface defining the contract for product record HTTP handlers
// ProductRecordHandlerI - Interfaz que define el contrato para los manejadores HTTP de registros de productos
type ProductRecordHandlerI interface {
	// Create - Handles HTTP POST requests for creating product records
	// Create - Maneja las peticiones HTTP POST para crear registros de productos
	Create(w http.ResponseWriter, r *http.Request)

	// GetReportByIdProduct - Handles HTTP GET requests for product record reports
	// GetReportByIdProduct - Maneja las peticiones HTTP GET para reportes de registros de productos
	GetReportByIdProduct(w http.ResponseWriter, r *http.Request)
}

// productRecordHandler - Handler layer implementation for product record HTTP operations
// productRecordHandler - Implementación de la capa de handler para operaciones HTTP de registros de productos
type productRecordHandler struct {
	Service services.ProductRecordServiceI // Service layer dependency for business logic / Dependencia de la capa de servicio para lógica de negocio
}

// NewProductRecordHandler - Constructor function that creates a new handler instance with service dependency injection
// NewProductRecordHandler - Función constructora que crea una nueva instancia del handler con inyección de dependencias del servicio
func NewProductRecordHandler(service services.ProductRecordServiceI) ProductRecordHandlerI {
	return &productRecordHandler{Service: service}
}

// Create - HTTP handler for creating product records with full request processing pipeline
// Create - Manejador HTTP para crear registros de productos con pipeline completo de procesamiento de peticiones
func (prh *productRecordHandler) Create(w http.ResponseWriter, r *http.Request) {
	// CONTEXT MANAGEMENT: Create context with timeout for request processing
	// GESTIÓN DE CONTEXTO: Crear contexto con timeout para procesamiento de peticiones
	parentCtx := context.Background()
	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	defer cancel() // Ensure context is cancelled to prevent resource leaks / Asegurar que el contexto se cancele para prevenir fugas de recursos

	var productRecordRequest requests.ProductRecordRequest

	// REQUEST PARSING: Parse JSON request body into structured data
	// ANÁLISIS DE PETICIÓN: Analizar el cuerpo de la petición JSON en datos estructurados
	err := request.JSON(r, &productRecordRequest)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// INPUT VALIDATION: Validate request structure and business rules
	// VALIDACIÓN DE ENTRADA: Validar estructura de petición y reglas de negocio
	v := validations.GetProductRecordValidation()
	err = v.ValidateProductRecordRequestStruct(productRecordRequest)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// DATA TRANSFORMATION: Convert request DTO to domain model
	// TRANSFORMACIÓN DE DATOS: Convertir DTO de petición a modelo de dominio
	productRecord := mappers.GetProductRecordFromRequest(productRecordRequest)

	// BUSINESS LOGIC DELEGATION: Call service layer for business processing
	// DELEGACIÓN DE LÓGICA DE NEGOCIO: Llamar a la capa de servicio para procesamiento de negocio
	result, err := prh.Service.CreateProductRecord(ctx, productRecord)

	if err != nil {
		// ERROR MAPPING: Map business errors to appropriate HTTP status codes
		// MAPEO DE ERRORES: Mapear errores de negocio a códigos de estado HTTP apropiados
		if errors.Is(err, error_message.ErrDependencyNotFound) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// RESPONSE FORMATTING: Transform domain model to response DTO and return success
	// FORMATEO DE RESPUESTA: Transformar modelo de dominio a DTO de respuesta y retornar éxito
	productRecordResponse := mappers.GetProductRecordResponseFromModel(result)
	response.JSON(w, http.StatusCreated, responses.DataResponse{
		Data: productRecordResponse,
	})
}

// GetReportByIdProduct - HTTP handler for retrieving product record reports with URL parameter processing
// GetReportByIdProduct - Manejador HTTP para obtener reportes de registros de productos con procesamiento de parámetros URL
func (prh *productRecordHandler) GetReportByIdProduct(w http.ResponseWriter, r *http.Request) {
	// CONTEXT MANAGEMENT: Create context with longer timeout for report generation
	// GESTIÓN DE CONTEXTO: Crear contexto con timeout más largo para generación de reportes
	parentCtx := context.Background()
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel() // Ensure context cleanup / Asegurar limpieza del contexto

	// URL PARAMETER EXTRACTION: Extract and validate ID parameter from query string
	// EXTRACCIÓN DE PARÁMETROS URL: Extraer y validar parámetro ID del query string
	idString := strings.TrimSpace(r.URL.Query().Get("id"))

	if idString == "" {
		report, err := prh.Service.GetReport(ctx)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.JSON(w, http.StatusOK, responses.DataResponse{Data: report})

	} else {

		// INPUT SANITIZATION: Convert string ID to integer with validation
		// SANITIZACIÓN DE ENTRADA: Convertir ID string a entero con validación
		idInt, err := strconv.Atoi(idString)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "error: id not is a number")
			return
		}

		// BUSINESS LOGIC DELEGATION: Call service layer for report generation
		// DELEGACIÓN DE LÓGICA DE NEGOCIO: Llamar a la capa de servicio para generación de reportes
		productRecordReport, err := prh.Service.GetReportByIdProduct(ctx, int64(idInt))

		if err != nil {
			// ERROR MAPPING: Map service layer errors to HTTP status codes
			// MAPEO DE ERRORES: Mapear errores de capa de servicio a códigos de estado HTTP
			if errors.Is(err, error_message.ErrDependencyNotFound) {
				response.Error(w, http.StatusConflict, err.Error())
				return
			}
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		// RESPONSE FORMATTING: Format successful response with report data as array
		// FORMATEO DE RESPUESTA: Formatear respuesta exitosa con datos del reporte como array
		response.JSON(w, http.StatusOK, responses.DataResponse{
			Data: []*models.ProductRecordReport{productRecordReport},
		})

	}
}
