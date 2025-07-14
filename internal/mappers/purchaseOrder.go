package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func GetModelPurchaseOrderFromRequest(por requests.PurchaseOrderRequest) *models.PurchaseOrder {
	return &models.PurchaseOrder{
		Id:              0,
		OrderNumber:     por.OrderNumber,
		OrderDate:       por.OrderDate,
		TrackingCode:    por.TrackingCode,
		BuyerId:         por.BuyerId,
		ProductRecordId: por.ProductRecordId,
	}
}

func GetResponsePurchaseOrderFromModel(po *models.PurchaseOrder) *responses.PurchaseOrderResponse {
	return &responses.PurchaseOrderResponse{
		OrderNumber:     po.OrderNumber,
		OrderDate:       po.OrderDate,
		TrackingCode:    po.TrackingCode,
		BuyerId:         po.BuyerId,
		ProductRecordId: po.ProductRecordId,
	}
}

func GetListPurchaseOrderResponseFromListModel(models []*models.PurchaseOrder) []*responses.PurchaseOrderResponse {
	var listPurchaseOrderResponse []*responses.PurchaseOrderResponse

	for _, purchaseOrder := range models {
		listPurchaseOrderResponse = append(listPurchaseOrderResponse, GetResponsePurchaseOrderFromModel(purchaseOrder))
	}

	return listPurchaseOrderResponse
}
