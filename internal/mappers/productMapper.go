package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func GetProductResponseFromModel(model *models.Product) *responses.ProductResponse {
	return &responses.ProductResponse{
		ProductCode:                    model.ProductCode,
		Description:                    model.Description,
		Width:                          model.Width,
		Height:                         model.Height,
		Length:                         model.Length,
		NetWeight:                      model.NetWeight,
		ExpirationRate:                 model.ExpirationRate,
		RecommendedFreezingTemperature: model.RecommendedFreezingTemperature,
		FreezingRate:                   model.FreezingRate,
		ProductTypeID:                  model.ProductTypeID,
		SellerID:                       model.SellerID,
	}
}

func GetProductFromRequest(request requests.ProductRequest) models.Product {
	return models.Product{
		ProductCode:                    request.ProductCode,
		Description:                    request.Description,
		Width:                          request.Width,
		Height:                         request.Height,
		Length:                         request.Length,
		NetWeight:                      request.NetWeight,
		ExpirationRate:                 request.ExpirationRate,
		RecommendedFreezingTemperature: request.RecommendedFreezingTemperature,
		FreezingRate:                   request.FreezingRate,
		ProductTypeID:                  request.ProductTypeID,
		SellerID:                       request.SellerID,
	}
}

func GetListProductResponseFromListModel(models []*models.Product) []*responses.ProductResponse {
	var listProductResponse []*responses.ProductResponse
	if len(models) > 0 {
		for _, s := range models {
			listProductResponse = append(listProductResponse, GetProductResponseFromModel(s))
		}
	}
	return listProductResponse
}
