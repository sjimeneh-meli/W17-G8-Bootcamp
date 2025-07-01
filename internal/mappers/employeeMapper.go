package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func GetEmployeeResponseFromModel(model *models.Employee) *responses.EmployeeResponse {
	return &responses.EmployeeResponse{
		ID:           model.Id,
		CardNumberID: model.CardNumberID,
		FirstName:    model.FirstName,
		LastName:     model.LastName,
		WarehouseID:  model.WarehouseID,
	}
}
func GetListEmployeeResponseFromListModel(models []*models.Employee) []*responses.EmployeeResponse {
	var listEmployeeResponse []*responses.EmployeeResponse
	if len(models) > 0 {
		for _, s := range models {
			listEmployeeResponse = append(listEmployeeResponse, GetEmployeeResponseFromModel(s))
		}
	}
	return listEmployeeResponse
}
func GetEmployeeModelFromRequest(model *models.Employee) *responses.EmployeeResponse {
	return &responses.EmployeeResponse{
		CardNumberID: model.CardNumberID,
		FirstName:    model.FirstName,
		LastName:     model.LastName,
		WarehouseID:  model.WarehouseID,
	}
}
