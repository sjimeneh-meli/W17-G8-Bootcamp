package handlers

import (
	"encoding/json"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
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
	GetById() http.HandlerFunc
	DeleteById() http.HandlerFunc
}

type EmployeeHandler struct {
	service services.EmployeeServiceI
	validation *validations.EmployeeValidation
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

func (h *EmployeeHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			request  *requests.EmployeeRequest = &requests.EmployeeRequest{}
			response *responses.DataResponse  = &responses.DataResponse{}
			section  *models.Employee
		)

		w.Header().Set("Content-Type", "application/json")

		if reqErr := json.NewDecoder(r.Body).Decode(request); reqErr != nil {
			response.SetError(error_message.ErrInvalidInput.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}

		if valErr := h.validation.ValidateEmployeeRequestStruct(*request); valErr != nil {
			response.SetError(valErr.Error())
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response)
			return
		}

		employee = mappers.GetEmployeeModelFromRequest(request)

		if h.service.ExistsWithSectionNumber(section.Id, section.SectionNumber) {
			response.SetError("already exist a section with the same number")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		}

		if srvErr := h.service.Create(section); srvErr != nil {
			response.SetError(error_message.ErrInvalidInput.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}

		response.Data = mappers.GetSectionResponseFromModel(section)
		w.WriteHeader(http.StatusCreated)
		encodeErr := json.NewEncoder(w).Encode(response)
		if encodeErr != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			response.SetError(encodeErr.Error())
			json.NewEncoder(w).Encode(response)
			return
		}

	}
}
}
