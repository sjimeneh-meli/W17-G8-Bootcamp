package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/bootcamp-go/web/response"
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

func purchaseOrderMapToPurchaseOrderList(orders map[int]models.PurchaseOrder) []*models.PurchaseOrder {
	ordersList := []*models.PurchaseOrder{}

	for _, order := range orders {
		ordersList = append(ordersList, &order)
	}

	return ordersList
}
