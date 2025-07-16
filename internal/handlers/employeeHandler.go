package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"

	"github.com/go-chi/chi/v5"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

func GetEmployeeHandler(service services.EmployeeServiceI, validation *validations.EmployeeValidation) EmployeeHandlerI {
	return &EmployeeHandler{
		service: service,
	}
}

type EmployeeHandlerI interface {
	GetAllEmployee() http.HandlerFunc
	PostEmployee() http.HandlerFunc
	GetByIdEmployee() http.HandlerFunc
	DeleteByIdEmployee() http.HandlerFunc
	PatchEmployee() http.HandlerFunc
}

type EmployeeHandler struct {
	service services.EmployeeServiceI
}

func (h *EmployeeHandler) GetAllEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var (
			requestResponse  *responses.DataResponse = &responses.DataResponse{}
			employeeResponse []*responses.EmployeeResponse
			employee         []*models.Employee
		)

		employeeMap, err := h.service.GetAll(ctx)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		employee = employeeMapToList(employeeMap)
		employeeResponse = mappers.GetListEmployeeResponseFromListModel(employee)
		requestResponse.Data = employeeResponse

		response.JSON(w, http.StatusOK, requestResponse)
	}
}
func (h *EmployeeHandler) GetByIdEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var (
			requestResponse  *responses.DataResponse = &responses.DataResponse{}
			employeeResponse *responses.EmployeeResponse
			employee         models.Employee
		)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		employee, err = h.service.GetById(ctx, id)
		if err != nil {

			if errors.Is(err, error_message.ErrNotFound) {
				response.Error(w, http.StatusNotFound, err.Error())
				return
			}

			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		employeeResponse = mappers.GetEmployeeResponseFromModel(&employee)
		requestResponse.Data = employeeResponse
		response.JSON(w, http.StatusOK, requestResponse)
	}
}
func (h *EmployeeHandler) PostEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestResponse *responses.DataResponse = &responses.DataResponse{}

		requestEmployee := requests.EmployeeRequest{}
		request.JSON(r, &requestEmployee)

		err := validations.ValidateEmployeeRequestStruct(requestEmployee)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		modelEmployee := mappers.GetEmployeeModelFromRequest(requestEmployee)
		employeeDb, err := h.service.Create(ctx, *modelEmployee)
		if err != nil {

			if errors.Is(err, error_message.ErrAlreadyExists) {
				response.Error(w, http.StatusConflict, err.Error())
				return
			}

			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		employeeResponse := mappers.GetEmployeeResponseFromModel(&employeeDb)
		requestResponse.Data = employeeResponse

		response.JSON(w, http.StatusCreated, requestResponse)
	}
}
func (h *EmployeeHandler) PatchEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestResponse *responses.DataResponse = &responses.DataResponse{}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		requestEmployee := requests.EmployeeRequest{}
		request.JSON(r, &requestEmployee)

		err = validations.IsNotAnEmptyEmployee(requestEmployee)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		modelEmployee := mappers.GetEmployeeModelFromRequest(requestEmployee)
		employeeDb, err := h.service.Update(ctx, id, *modelEmployee)
		if err != nil {

			if errors.Is(err, error_message.ErrNotFound) {
				response.Error(w, http.StatusNotFound, err.Error())
				return
			}

			if errors.Is(err, error_message.ErrAlreadyExists) {
				response.Error(w, http.StatusConflict, err.Error())
				return
			}

			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		employeeResponse := mappers.GetEmployeeResponseFromModel(&employeeDb)
		requestResponse.Data = employeeResponse

		response.JSON(w, http.StatusOK, requestResponse)

	}
}
func (h *EmployeeHandler) DeleteByIdEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.service.DeleteById(ctx, id)
		if err != nil {

			if errors.Is(err, error_message.ErrNotFound) {
				response.Error(w, http.StatusNotFound, err.Error())
				return
			}

			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func employeeMapToList(employee map[int]models.Employee) []*models.Employee {
	employeeList := []*models.Employee{}
	for _, empl := range employee {
		employeeList = append(employeeList, &empl)
	}
	return employeeList
}
