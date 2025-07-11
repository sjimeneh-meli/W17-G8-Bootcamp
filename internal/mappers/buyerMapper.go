package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func GetModelBuyerFromRequest(br requests.BuyerRequest) *models.Buyer {
	return &models.Buyer{
		Id:           0,
		CardNumberId: br.CardNumberId,
		FirstName:    br.FirstName,
		LastName:     br.LastName,
	}
}

func GetResponseBuyerFromModel(b *models.Buyer) *responses.BuyerResponse {
	return &responses.BuyerResponse{
		Id:           b.Id,
		CardNumberId: b.CardNumberId,
		FirstName:    b.FirstName,
		LastName:     b.LastName,
	}
}

func GetListBuyerResponseFromListModel(models []*models.Buyer) []*responses.BuyerResponse {
	var listBuyerResponse []*responses.BuyerResponse

	for _, buyer := range models {
		listBuyerResponse = append(listBuyerResponse, GetResponseBuyerFromModel(buyer))
	}

	return listBuyerResponse
}
