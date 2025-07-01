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
	AddEmployee() http.HandlerFunc
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

func (h *EmployeeHandler) AddEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			// Variables auxiliares para manejar la solicitud entrante, el modelo interno,
			//la conversión a respuesta y la estructura final que se envia

			response   *responses.DataResponse          = &responses.DataResponse{}
			requestDTO *responses.EmployeeCreateRequest = &responses.EmployeeCreateRequest{}
			newEmp     *models.Employee
			result     *responses.EmployeeResponse
		)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewDecoder(r.Body).Decode(&requestDTO); err != nil {
			response.SetError("invalid JSON format")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		//Validate basic
		if requestDTO.CardNumberID == "" || requestDTO.FirstName == "" || requestDTO.LastName == "" {
			response.SetError("missing required fields")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		//Map
		newEmp = &models.Employee{
			CardNumberID: requestDTO.CardNumberID,
			FirstName:    requestDTO.FirstName,
			LastName:     requestDTO.LastName,
			WarehouseID:  requestDTO.WarehouseID,
		}

		newEmp, err := h.service.Create(newEmp)
		if err != nil {
			response.SetError(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		result = mappers.GetEmployeeResponseFromModel(newEmp)
		response.Data = result

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
