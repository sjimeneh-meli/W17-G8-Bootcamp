package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func GetModelInboundOrderFromRequest(req requests.InboundOrderRequest) *models.InboundOrder {
	return &models.InboundOrder{
		OrderDate:      req.Data.OrderDate,
		OrderNumber:    req.Data.OrderNumber,
		EmployeeId:     req.Data.EmployeeId,
		ProductBatchId: req.Data.ProductBatchId,
		WarehouseId:    req.Data.WarehouseId,
	}
}
func GetResponseInboundOrderFromModel(po *models.InboundOrder) *responses.InboundOrderResponse {
	return &responses.InboundOrderResponse{
		Id:             po.Id,
		OrderNumber:    po.OrderNumber,
		OrderDate:      po.OrderDate,
		EmployeeId:     po.EmployeeId,
		ProductBatchId: po.ProductBatchId,
		WarehouseId:    po.WarehouseId,
	}
}

func GetListInboundOrderResponseFromListModel(models []*models.InboundOrder) []*responses.InboundOrderResponse {
	var listInboundOrderResponse []*responses.InboundOrderResponse

	for _, inboundOrder := range models {
		listInboundOrderResponse = append(listInboundOrderResponse, GetResponseInboundOrderFromModel(inboundOrder))
	}

	return listInboundOrderResponse
}
