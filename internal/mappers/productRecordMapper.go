package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func GetProductRecordResponseFromModel(model *models.ProductRecord) *responses.ProductRecordResponse {
	return &responses.ProductRecordResponse{
		LastUpdateDate: model.LastUpdateDate,
		PurchasePrice:  model.PurchasePrice,
		SalePrice:      model.SalePrice,
		ProductID:      model.ProductID,
	}
}

func GetProductRecordFromRequest(request requests.ProductRecordRequest) models.ProductRecord {
	return models.ProductRecord{
		LastUpdateDate: request.LastUpdateDate,
		PurchasePrice:  request.PurchasePrice,
		SalePrice:      request.SalePrice,
		ProductID:      request.ProductID,
	}
}
