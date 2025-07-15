package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

func GetEmployeeHandler(service services.EmployeeServiceI, validation *validations.EmployeeValidation) EmployeeHandlerI {
	return &EmployeeHandler{
		service:    service,
		validation: validation,
	}
}

type EmployeeHandlerI interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	DeleteById(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type EmployeeHandler struct {
	service    services.EmployeeServiceI
	validation *validations.EmployeeValidation
}

func (h *EmployeeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	{
		var (
			responseJson     *responses.DataResponse = &responses.DataResponse{}
			employeeResponse []*responses.EmployeeResponse
			employee         []*models.Employee
		)

		employee = h.service.GetAll()
		employeeResponse = mappers.GetListEmployeeResponseFromListModel(employee)
		responseJson.Data = employeeResponse

		response.JSON(w, http.StatusOK, responseJson)
	}
}
func (h *EmployeeHandler) GetById(w http.ResponseWriter, r *http.Request) {

	var (
		responseJson     *responses.DataResponse = &responses.DataResponse{}
		employeeResponse *responses.EmployeeResponse
	)

	idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
	if convErr != nil {
		response.Error(w, http.StatusExpectationFailed, convErr.Error())
	}

	employee, srvErr := h.service.GetById(idParam)
	if srvErr != nil {
		response.Error(w, http.StatusNotFound, srvErr.Error())
	}

	employeeResponse = mappers.GetEmployeeResponseFromModel(employee)
	responseJson.Data = employeeResponse

	response.JSON(w, http.StatusOK, responseJson)

}
func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		request      *requests.EmployeeRequest = &requests.EmployeeRequest{}
		responseJson *responses.DataResponse   = &responses.DataResponse{}
		employee     *models.Employee
	)

	if reqErr := json.NewDecoder(r.Body).Decode(&request); reqErr != nil {
		response.Error(w, http.StatusExpectationFailed, reqErr.Error())
		return
	}

	if valErr := h.validation.ValidateEmployeeRequestStruct(*request); valErr != nil {
		response.Error(w, http.StatusUnprocessableEntity, valErr.Error())
		return
	}

	employee = mappers.GetEmployeeModelFromRequest(request)

	if h.service.ExistsWhCardNumber(employee.Id, employee.CardNumberID) {
		response.Error(w, http.StatusConflict, "Already exist with same id_card_number")
		return
	}

	if srvErr := h.service.Create(employee); srvErr != nil {
		response.Error(w, http.StatusExpectationFailed, srvErr.Error())
		return
	}

	responseJson.Data = mappers.GetEmployeeResponseFromModel(employee)
	response.JSON(w, http.StatusCreated, responseJson)

}
func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {

	var (
		request      *requests.EmployeeRequest = &requests.EmployeeRequest{}
		responseJson *responses.DataResponse   = &responses.DataResponse{}
	)

	idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
	if convErr != nil {
		response.Error(w, http.StatusExpectationFailed, convErr.Error())
		return
	}

	employee, srvErr := h.service.GetById(idParam)
	if srvErr != nil {
		response.Error(w, http.StatusNotFound, srvErr.Error())
		return
	}

	if reqErr := json.NewDecoder(r.Body).Decode(request); reqErr != nil {
		response.Error(w, http.StatusExpectationFailed, reqErr.Error())
		return
	}

	if valErr := h.validation.ValidateEmployeeRequestStruct(*request); valErr != nil {
		response.Error(w, http.StatusUnprocessableEntity, valErr.Error())
		return
	}

	if h.service.ExistsWhCardNumber(employee.Id, request.CardNumberID) {
		response.Error(w, http.StatusConflict, "")
		return
	}

	mappers.UpdateEmployeeModelFromRequest(employee, request)

	responseJson.Data = mappers.GetEmployeeResponseFromModel(employee)
	response.JSON(w, http.StatusOK, responseJson)
}
func (h *EmployeeHandler) DeleteById(w http.ResponseWriter, r *http.Request) {

	idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
	if convErr != nil {
		response.Error(w, http.StatusExpectationFailed, convErr.Error())
		return
	}

	srvErr := h.service.DeleteById(idParam)
	if srvErr != nil {
		response.Error(w, http.StatusNotFound, srvErr.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
