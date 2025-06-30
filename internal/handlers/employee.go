package handlers

import (
	"encoding/json"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"net/http"
)

func GetEmployeeHandler() EmployeeHandlerI {
	return &EmployeeHandler{
		service: services.GetEmployeeService(),
	}
}

type EmployeeHandlerI interface {
	GetAll() http.HandlerFunc
}

type EmployeeHandler struct {
	service services.EmployeeServiceI
}

func (h *EmployeeHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			response         *responses.DataResponse = &responses.DataResponse{}
			employeeResponse []*responses.EmployeeResponse
			employee         []*models.Employee
		)

		w.Header().Set("Content-Type", "application/json")

		employee = h.service.GetAll()
		employeeResponse = mappers.GetListEmployeeResponseFromListModel(employee)
		response.Data = employeeResponse

		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(response)
		if encodeErr != nil {
			response.SetError(encodeErr.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}
