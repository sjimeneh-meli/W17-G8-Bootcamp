package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func GetEmployeeModelFromRequest(request requests.EmployeeRequest) *models.Employee {
	return &models.Employee{
		Id:           0,
		CardNumberID: request.CardNumberID,
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		WarehouseID:  request.WarehouseID,
	}
}

func GetEmployeeResponseFromModel(model *models.Employee) *responses.EmployeeResponse {
	return &responses.EmployeeResponse{
		ID:           model.Id,
		CardNumberID: model.CardNumberID,
		FirstName:    model.FirstName,
		LastName:     model.LastName,
		WarehouseID:  model.WarehouseID,
	}
}

func UpdateEmployeeModelFromRequest(model *models.Employee, request *requests.EmployeeRequest) {
	model.CardNumberID = request.CardNumberID
	model.FirstName = request.FirstName
	model.LastName = request.LastName
	model.WarehouseID = request.WarehouseID
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
