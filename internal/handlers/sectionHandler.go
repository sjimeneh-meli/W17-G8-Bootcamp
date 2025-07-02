package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

func GetSectionHandler(service services.SectionServiceI, validation *validations.SectionValidation) SectionHandlerI {
	return &SectionHandler{
		service:    service,
		validation: validation,
	}
}

type SectionHandlerI interface {
	GetAll() http.HandlerFunc
	GetByID() http.HandlerFunc
	Create() http.HandlerFunc
	Update() http.HandlerFunc
	DeleteByID() http.HandlerFunc
}

type SectionHandler struct {
	service    services.SectionServiceI
	validation *validations.SectionValidation
}

func (h *SectionHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			response         *responses.DataResponse = &responses.DataResponse{}
			sectionsResponse []*responses.SectionResponse
			sections         []*models.Section
		)

		w.Header().Set("Content-Type", "application/json")

		sections = h.service.GetAll()
		sectionsResponse = mappers.GetListSectionResponseFromListModel(sections)
		response.Data = sectionsResponse

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

func (h *SectionHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			response        *responses.DataResponse = &responses.DataResponse{}
			sectionResponse *responses.SectionResponse
		)

		w.Header().Set("Content-Type", "application/json")

		idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
		if convErr != nil {
			response.SetError(error_message.ErrInvalidInput.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}

		section, srvErr := h.service.GetByID(idParam)
		if srvErr != nil {
			response.SetError(srvErr.Error())
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}

		sectionResponse = mappers.GetSectionResponseFromModel(section)
		response.Data = sectionResponse

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

func (h *SectionHandler) Create() http.HandlerFunc {
	//VALIDATE IF PRODUCT AND WAREHOUSE EXIST
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			request  *requests.SectionRequest = &requests.SectionRequest{}
			response *responses.DataResponse  = &responses.DataResponse{}
			section  *models.Section
		)

		w.Header().Set("Content-Type", "application/json")

		if reqErr := json.NewDecoder(r.Body).Decode(request); reqErr != nil {
			response.SetError(error_message.ErrInvalidInput.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}

		if valErr := h.validation.ValidateSectionRequestStruct(*request); valErr != nil {
			response.SetError(valErr.Error())
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response)
			return
		}

		section = mappers.GetSectionModelFromRequest(request)

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

func (h *SectionHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			request  *requests.SectionRequest = &requests.SectionRequest{}
			response *responses.DataResponse  = &responses.DataResponse{}
		)

		w.Header().Set("Content-Type", "application/json")

		idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
		if convErr != nil {
			response.SetError(error_message.ErrInvalidInput.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}

		section, srvErr := h.service.GetByID(idParam)
		if srvErr != nil {
			response.SetError(srvErr.Error())
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

		if valErr := h.validation.ValidateSectionRequestStruct(*request); valErr != nil {
			response.SetError(valErr.Error())
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response)
			return
		}

		if h.service.ExistsWithSectionNumber(section.Id, request.SectionNumber) {
			response.SetError("already exist a section with the same number")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response)
			return
		}

		mappers.UpdateSectionModelFromRequest(section, request)

		response.Data = mappers.GetSectionResponseFromModel(section)
		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(response)
		if encodeErr != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			response.SetError(encodeErr.Error())
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}

func (h *SectionHandler) DeleteByID() http.HandlerFunc {
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

		srvErr := h.service.DeleteByID(idParam)
		if srvErr != nil {
			response.SetError(srvErr.Error())
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
