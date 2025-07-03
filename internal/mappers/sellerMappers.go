package mappers

import "github.com/sajimenezher_meli/meli-frescos-8/internal/models"

func ToRequestToSellerStruct(seller models.SellerRequest) models.Seller {
	sellerFormated := models.Seller{
		Id:          0,
		CID:         seller.CID,
		CompanyName: seller.CompanyName,
		Address:     seller.Address,
		Telephone:   seller.Telephone,
	}
	return sellerFormated
}
