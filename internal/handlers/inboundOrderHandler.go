package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
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

// GetInboundOrderHandler creates and returns a new instance of InboundOrderHandler with the required service
// GetInboundOrderHandler crea y retorna una nueva instancia de InboundOrderHandler con el servicio requerido
func GetInboundOrderHandler(service services.InboundOrdersServiceI) InboundOrderHandlerI {
	return &InboundOrderHandler{
		service: service,
	}
}

// InboundOrderHandlerI defines the contract for inbound order HTTP handlers
// InboundOrderHandlerI define el contrato para los manejadores HTTP de órdenes de entrada
type InboundOrderHandlerI interface {
	GetInboundOrdersReport() http.HandlerFunc
	PostInboundOrder() http.HandlerFunc
}

// InboundOrderHandler implements InboundOrderHandlerI and handles HTTP requests for inbound order operations
// InboundOrderHandler implementa InboundOrderHandlerI y maneja las solicitudes HTTP para operaciones de órdenes de entrada
type InboundOrderHandler struct {
	service services.InboundOrdersServiceI // Service layer for inbound order business logic / Capa de servicio para lógica de negocio de órdenes de entrada
}

// GetInboundOrdersReport handles HTTP GET requests to retrieve inbound order reports
// Accepts an optional 'id' query parameter to filter by employee ID
// GetInboundOrdersReport maneja las solicitudes HTTP GET para recuperar reportes de órdenes de entrada
// Acepta un parámetro de consulta 'id' opcional para filtrar por ID de empleado
func (h *InboundOrderHandler) GetInboundOrdersReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestResponse *responses.DataResponse = &responses.DataResponse{}
		var responseData []models.InboundOrderReport

		// Get and validate optional employee ID query parameter / Obtener y validar parámetro opcional de consulta ID de empleado
		queryId := r.URL.Query().Get("id")
		if queryId != "" {
			// Parse and validate employee ID / Parsear y validar ID de empleado
			employeeId, err := strconv.Atoi(queryId)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Invalid employee ID")
				return
			}

			// Get report for specific employee / Obtener reporte para empleado específico
			report, err := h.service.GetInboundOrdersReportByEmployeeId(ctx, employeeId)
			if err != nil {
				response.Error(w, http.StatusInternalServerError, err.Error())
				return
			}

			responseData = append(responseData, report)
		} else {
			// Get reports for all employees / Obtener reportes para todos los empleados
			reports, err := h.service.GetAllInboundOrdersReports(ctx)
			if err != nil {
				response.Error(w, http.StatusInternalServerError, err.Error())
				return
			}
			responseData = reports
		}

		requestResponse.Data = responseData
		response.JSON(w, http.StatusOK, requestResponse)
	}
}

// PostInboundOrder handles HTTP POST requests to create a new inbound order
// Validates the request body and returns appropriate HTTP status codes
// PostInboundOrder maneja las solicitudes HTTP POST para crear una nueva orden de entrada
// Valida el cuerpo de la solicitud y retorna códigos de estado HTTP apropiados
func (h *InboundOrderHandler) PostInboundOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var (
			reqResponse    *responses.DataResponse      = &responses.DataResponse{}
			requestInbound requests.InboundOrderRequest = requests.InboundOrderRequest{}
		)

		// Parse and validate JSON request body / Parsear y validar cuerpo de solicitud JSON
		err := request.JSON(r, &requestInbound)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Validate request structure and business rules / Validar estructura de solicitud y reglas de negocio
		err = validations.ValidateInboundOrderRequestStruct(requestInbound)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		// Map request to model and create through service / Mapear solicitud a modelo y crear a través del servicio
		modelInbound := mappers.GetModelInboundOrderFromRequest(requestInbound)
		order, err := h.service.Create(ctx, *modelInbound)
		if err != nil {
			// Handle specific error types / Manejar tipos de error específicos
			switch {
			case errors.Is(err, error_message.ErrAlreadyExists):
				response.Error(w, http.StatusConflict, err.Error())
				return
			case errors.Is(err, error_message.ErrDependencyNotFound):
				response.Error(w, http.StatusConflict, err.Error())
				return
			default:
				response.Error(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		// Map model to response format / Mapear modelo a formato de respuesta
		orderResponse := mappers.GetResponseInboundOrderFromModel(&order)
		reqResponse.Data = orderResponse

		response.JSON(w, http.StatusCreated, reqResponse)
	}
}
