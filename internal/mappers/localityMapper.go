package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

/*type LocalitySellerReportResponse struct {
	Data []responses.LocalitySellerResponse `json:"data"`
}*/

func ToRequestToLocalityStruct(locality requests.LocalityRequest) models.Locality {
	localityFormated := models.Locality{
		LocalityName: locality.Data.LocalityName,
		ProvinceName: locality.Data.ProvinceName,
		CountryName:  locality.Data.CountryName,
	}
	return localityFormated
}
