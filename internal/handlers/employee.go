package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"net/http"
	"strconv"
)

func GetEmployeeHandler(service services.EmployeeServiceI, validation *validations.EmployeeValidation) EmployeeHandlerI {
	return &EmployeeHandler{
		service:    service,
		validation: validation,
	}
}

type EmployeeHandlerI interface {
	GetAll() http.HandlerFunc
	Create() http.HandlerFunc
	GetById() http.HandlerFunc
	DeleteById() http.HandlerFunc
	PatchEmployee() http.HandlerFunc
}

type EmployeeHandler struct {
	service    services.EmployeeServiceI
	validation *validations.EmployeeValidation
}

func (h *EmployeeHandler) PatchEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			request  *requests.EmployeeRequest = &requests.EmployeeRequest{}
			response *responses.DataResponse   = &responses.DataResponse{}
		)

		w.Header().Set("Content-Type", "application/json")

		idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
		if convErr != nil {
			response.SetError(error_message.ErrInvalidInput.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}

		employee, srvErr := h.service.GetById(idParam)
		if srvErr != nil {
			response.SetError(error_message.ErrNotFound.Error())
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}

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

		if h.service.ExistsWhCardNumber(employee.Id, request.CardNumberID) {
			response.SetError("already exist an employee with the same card_number_id")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		}

		mappers.UpdateEmployeeModelFromRequest(employee, request)

		response.Data = mappers.GetEmployeeResponseFromModel(employee)
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			response.SetError(err.Error())
			json.NewEncoder(w).Encode(response)
			return
		}
	}
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
			response *responses.DataResponse   = &responses.DataResponse{}
			employee *models.Employee
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

		if h.service.ExistsWhCardNumber(employee.Id, employee.CardNumberID) {
			response.SetError("already exist a section with the same number")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		}

		if srvErr := h.service.Create(employee); srvErr != nil {
			response.SetError(error_message.ErrInvalidInput.Error())
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response)
			return
		}

		response.Data = mappers.GetEmployeeResponseFromModel(employee)
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

func (h *EmployeeHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			response         *responses.DataResponse = &responses.DataResponse{}
			employeeResponse *responses.EmployeeResponse
		)

		w.Header().Set("Content-Type", "application/json")

		idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
		if convErr != nil {
			response.SetError(error_message.ErrInvalidInput.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}

		employee, srvErr := h.service.GetById(idParam)
		if srvErr != nil {
			response.SetError(srvErr.Error())
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}

		employeeResponse = mappers.GetEmployeeResponseFromModel(employee)
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

func (h *EmployeeHandler) DeleteById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			response *responses.DataResponse = &responses.DataResponse{}
		)

		w.Header().Set("Content-Type", "application/json")

		idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
		if convErr != nil {
			response.SetError(error_message.ErrInvalidInput.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}

		srvErr := h.service.DeleteById(idParam)
		if srvErr != nil {
			response.SetError(srvErr.Error())
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
