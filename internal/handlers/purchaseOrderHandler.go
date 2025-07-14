package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
)

func GetPurchaseOrderHandler(service services.PurchaseOrderServiceI) PurchaseOrderHandlerI {
	return &PurchaseOrderHandler{
		service: service,
	}
}

type PurchaseOrderHandlerI interface {
	GetAll() http.HandlerFunc
	GetPurchaseOrdersReport() http.HandlerFunc
}

type PurchaseOrderHandler struct {
	service services.PurchaseOrderServiceI
}

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

func purchaseOrderMapToPurchaseOrderList(orders map[int]models.PurchaseOrder) []*models.PurchaseOrder {
	ordersList := []*models.PurchaseOrder{}

	for _, order := range orders {
		ordersList = append(ordersList, &order)
	}

	return ordersList
}
