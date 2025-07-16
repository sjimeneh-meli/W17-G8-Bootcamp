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

// GetPurchaseOrderHandler creates and returns a new instance of PurchaseOrderHandler with the required service
// GetPurchaseOrderHandler crea y retorna una nueva instancia de PurchaseOrderHandler con el servicio requerido
func GetPurchaseOrderHandler(service services.PurchaseOrderServiceI) PurchaseOrderHandlerI {
	return &PurchaseOrderHandler{
		service: service,
	}
}

// PurchaseOrderHandlerI defines the contract for purchase order HTTP handlers
// PurchaseOrderHandlerI define el contrato para los manejadores HTTP de órdenes de compra
type PurchaseOrderHandlerI interface {
	GetAll() http.HandlerFunc
	GetPurchaseOrdersReport() http.HandlerFunc
	PostPurchaseOrder() http.HandlerFunc
}

// PurchaseOrderHandler implements PurchaseOrderHandlerI and handles HTTP requests for purchase order operations
// PurchaseOrderHandler implementa PurchaseOrderHandlerI y maneja las solicitudes HTTP para operaciones de órdenes de compra
type PurchaseOrderHandler struct {
	service services.PurchaseOrderServiceI // Service layer for purchase order business logic / Capa de servicio para lógica de negocio de órdenes de compra
}

// PostPurchaseOrder handles HTTP POST requests to create a new purchase order
// Validates the request body and returns appropriate HTTP status codes
// PostPurchaseOrder maneja las solicitudes HTTP POST para crear una nueva orden de compra
// Valida el cuerpo de la solicitud y retorna códigos de estado HTTP apropiados
func (h *PurchaseOrderHandler) PostPurchaseOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var (
			requestResponse *responses.DataResponse       = &responses.DataResponse{}
			requestOrder    requests.PurchaseOrderRequest = requests.PurchaseOrderRequest{}
		)

		// Parse and validate JSON request body / Parsear y validar cuerpo de solicitud JSON
		err := request.JSON(r, &requestOrder)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Validate request structure and business rules / Validar estructura de solicitud y reglas de negocio
		err = validations.ValidatePurchaseOrderRequestStruct(requestOrder)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		// Map request to model and create through service / Mapear solicitud a modelo y crear a través del servicio
		modelOrder := mappers.GetModelPurchaseOrderFromRequest(requestOrder)
		orderDb, err := h.service.Create(ctx, *modelOrder)
		if err != nil {
			// Handle specific error types / Manejar tipos de error específicos
			if errors.Is(err, error_message.ErrAlreadyExists) {
				response.Error(w, http.StatusConflict, err.Error())
				return
			}

			if errors.Is(err, error_message.ErrNotFound) {
				response.Error(w, http.StatusConflict, err.Error())
				return
			}

			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Map model to response format / Mapear modelo a formato de respuesta
		orderResponse := mappers.GetResponsePurchaseOrderFromModel(&orderDb)
		requestResponse.Data = orderResponse

		response.JSON(w, http.StatusCreated, requestResponse)
	}
}

// GetAll handles HTTP GET requests to retrieve all purchase orders
// Returns a JSON response with all purchase orders or appropriate error codes
// GetAll maneja las solicitudes HTTP GET para recuperar todas las órdenes de compra
// Retorna una respuesta JSON con todas las órdenes de compra o códigos de error apropiados
func (h *PurchaseOrderHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		var (
			requestResponse       *responses.DataResponse = &responses.DataResponse{}
			purchaseOrderResponse []*responses.PurchaseOrderResponse
			purchaseOrders        []*models.PurchaseOrder
		)

		// Get all purchase orders from service layer / Obtener todas las órdenes de compra de la capa de servicio
		purchaseOrdersMap, err := h.service.GetAll(ctx)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
		}

		// Convert map to slice and map to response format / Convertir mapa a slice y mapear a formato de respuesta
		purchaseOrders = purchaseOrderMapToPurchaseOrderList(purchaseOrdersMap)
		purchaseOrderResponse = mappers.GetListPurchaseOrderResponseFromListModel(purchaseOrders)

		requestResponse.Data = purchaseOrderResponse
		response.JSON(w, http.StatusOK, requestResponse)
	}
}

// GetPurchaseOrdersReport handles HTTP GET requests to retrieve purchase order reports
// Accepts an optional 'id' query parameter to filter by buyer ID
// GetPurchaseOrdersReport maneja las solicitudes HTTP GET para recuperar reportes de órdenes de compra
// Acepta un parámetro de consulta 'id' opcional para filtrar por ID de comprador
func (h *PurchaseOrderHandler) GetPurchaseOrdersReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		var requestResponse *responses.DataResponse = &responses.DataResponse{}
		var idRequest *int = nil

		// Get and validate optional buyer ID query parameter / Obtener y validar parámetro opcional de consulta ID de comprador
		idParam := r.URL.Query().Get("id")
		if idParam != "" {
			// Parse and validate buyer ID / Parsear y validar ID de comprador
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				response.Error(w, http.StatusBadRequest, err.Error())
				return
			}
			idRequest = &id
		}

		// Get purchase order reports from service layer / Obtener reportes de órdenes de compra de la capa de servicio
		report, err := h.service.GetPurchaseOrdersReport(ctx, idRequest)
		if err != nil {
			// Handle specific error types / Manejar tipos de error específicos
			if errors.Is(err, error_message.ErrNotFound) {
				response.Error(w, http.StatusNotFound, err.Error())
				return
			}
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		requestResponse.Data = report
		response.JSON(w, http.StatusOK, requestResponse)
	}
}

// purchaseOrderMapToPurchaseOrderList converts a map of purchase orders to a slice of purchase order pointers
// Helper function for data transformation in the handler layer
// purchaseOrderMapToPurchaseOrderList convierte un mapa de órdenes de compra a un slice de punteros de órdenes de compra
// Función auxiliar para transformación de datos en la capa de manejadores
func purchaseOrderMapToPurchaseOrderList(orders map[int]models.PurchaseOrder) []*models.PurchaseOrder {
	ordersList := []*models.PurchaseOrder{}

	for _, order := range orders {
		ordersList = append(ordersList, &order)
	}

	return ordersList
}
