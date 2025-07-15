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

// GetPurchaseOrderHandler creates and returns a new instance of PurchaseOrderHandler
func GetPurchaseOrderHandler(service services.PurchaseOrderServiceI) PurchaseOrderHandlerI {
	return &PurchaseOrderHandler{
		service: service,
	}
}

// PurchaseOrderHandlerI defines the interface for purchase order HTTP handlers
type PurchaseOrderHandlerI interface {
	GetAll() http.HandlerFunc
	GetPurchaseOrdersReport() http.HandlerFunc
	PostPurchaseOrder() http.HandlerFunc
}

// PurchaseOrderHandler implements PurchaseOrderHandlerI and handles HTTP requests for purchase order operations
type PurchaseOrderHandler struct {
	service services.PurchaseOrderServiceI
}

// PostPurchaseOrder handles HTTP POST requests to create a new purchase order
// Validates the request body and returns appropriate HTTP status codes
func (h *PurchaseOrderHandler) PostPurchaseOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var (
			requestResponse *responses.DataResponse       = &responses.DataResponse{}
			requestOrder    requests.PurchaseOrderRequest = requests.PurchaseOrderRequest{}
		)

		err := request.JSON(r, &requestOrder)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		err = validations.ValidatePurchaseOrderRequestStruct(requestOrder)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		modelOrder := mappers.GetModelPurchaseOrderFromRequest(requestOrder)
		orderDb, err := h.service.Create(ctx, *modelOrder)
		if err != nil {

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

		orderResponse := mappers.GetResponsePurchaseOrderFromModel(&orderDb)
		requestResponse.Data = orderResponse

		response.JSON(w, http.StatusCreated, requestResponse)
	}
}

// GetAll handles HTTP GET requests to retrieve all purchase orders
// Returns a JSON response with all purchase orders
func (h *PurchaseOrderHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		var (
			requestResponse       *responses.DataResponse = &responses.DataResponse{}
			purchaseOrderResponse []*responses.PurchaseOrderResponse
			purchaseOrders        []*models.PurchaseOrder
		)

		purchaseOrdersMap, err := h.service.GetAll(ctx)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
		}

		purchaseOrders = purchaseOrderMapToPurchaseOrderList(purchaseOrdersMap)
		purchaseOrderResponse = mappers.GetListPurchaseOrderResponseFromListModel(purchaseOrders)

		requestResponse.Data = purchaseOrderResponse
		response.JSON(w, http.StatusOK, requestResponse)
	}
}

// GetPurchaseOrdersReport handles HTTP GET requests to retrieve purchase order reports
// Accepts an optional 'id' query parameter to filter by buyer ID
func (h *PurchaseOrderHandler) GetPurchaseOrdersReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		var requestResponse *responses.DataResponse = &responses.DataResponse{}
		var idRequest *int = nil

		idParam := r.URL.Query().Get("id")
		if idParam != "" {
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				response.Error(w, http.StatusBadRequest, err.Error())
				return
			}
			idRequest = &id
		}

		report, err := h.service.GetPurchaseOrdersReport(ctx, idRequest)
		if err != nil {
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
func purchaseOrderMapToPurchaseOrderList(orders map[int]models.PurchaseOrder) []*models.PurchaseOrder {
	ordersList := []*models.PurchaseOrder{}

	for _, order := range orders {
		ordersList = append(ordersList, &order)
	}

	return ordersList
}
