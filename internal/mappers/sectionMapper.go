package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func GetSectionResponseFromModel(model *models.Section) *responses.SectionResponse {
	return &responses.SectionResponse{
		ID:                 model.Id,
		SectionNumber:      model.SectionNumber,
		CurrentCapacity:    model.CurrentCapacity,
		CurrentTemperature: model.CurrentTemperature,
		MaximumCapacity:    model.MaximumCapacity,
		MinimumCapacity:    model.MinimumCapacity,
		MinimumTemperature: model.MinimumTemperature,
		ProductTypeID:      model.ProductTypeID,
		WarehouseID:        model.WarehouseID,
	}
}

func GetListSectionResponseFromListModel(models []*models.Section) []*responses.SectionResponse {
	var listSectionResponse []*responses.SectionResponse
	if len(models) > 0 {
		for _, s := range models {
			listSectionResponse = append(listSectionResponse, GetSectionResponseFromModel(s))
		}
	}
	return listSectionResponse
}
