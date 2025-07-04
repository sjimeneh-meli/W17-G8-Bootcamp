package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func ToRequestToSellerStruct(seller requests.SellerRequest) models.Seller {
	sellerFormated := models.Seller{
		CID:         seller.CID,
		CompanyName: seller.CompanyName,
		Address:     seller.Address,
		Telephone:   seller.Telephone,
	}
	return sellerFormated
}

func ToSellerStructToResponse(seller models.Seller) responses.SellerResponse {
	sellerFormated := responses.SellerResponse{
		Id:          seller.Id,
		CID:         seller.CID,
		CompanyName: seller.CompanyName,
		Address:     seller.Address,
		Telephone:   seller.Telephone,
	}
	return sellerFormated
}
