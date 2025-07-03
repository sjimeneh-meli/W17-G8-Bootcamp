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

func GetSectionHandler(service services.SectionServiceI, validation *validations.SectionValidation) SectionHandlerI {
	return &SectionHandler{
		service:    service,
		validation: validation,
	}
}

type SectionHandlerI interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	DeleteByID(w http.ResponseWriter, r *http.Request)
}

type SectionHandler struct {
	service    services.SectionServiceI
	validation *validations.SectionValidation
}

func (h *SectionHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	var (
		responseJson     *responses.DataResponse = &responses.DataResponse{}
		sectionsResponse []*responses.SectionResponse
		sections         []*models.Section
	)

	w.Header().Set("Content-Type", "application/json")

	sections = h.service.GetAll()
	sectionsResponse = mappers.GetListSectionResponseFromListModel(sections)
	responseJson.Data = sectionsResponse

	response.JSON(w, http.StatusOK, responseJson)

}

func (h *SectionHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	var (
		responseJson    *responses.DataResponse = &responses.DataResponse{}
		sectionResponse *responses.SectionResponse
	)

	w.Header().Set("Content-Type", "application/json")

	idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
	if convErr != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		json.NewEncoder(w).Encode(responseJson)
		return
	}

	section, srvErr := h.service.GetByID(idParam)
	if srvErr != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseJson)
		return
	}

	sectionResponse = mappers.GetSectionResponseFromModel(section)
	responseJson.Data = sectionResponse

	response.JSON(w, http.StatusOK, responseJson)

}

func (h *SectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	//VALIDATE IF PRODUCT AND WAREHOUSE EXIST
	var (
		request      *requests.SectionRequest = &requests.SectionRequest{}
		responseJson *responses.DataResponse  = &responses.DataResponse{}
		section      *models.Section
	)

	w.Header().Set("Content-Type", "application/json")

	if reqErr := json.NewDecoder(r.Body).Decode(request); reqErr != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		json.NewEncoder(w).Encode(responseJson)
		return
	}

	if valErr := h.validation.ValidateSectionRequestStruct(*request); valErr != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(responseJson)
		return
	}

	section = mappers.GetSectionModelFromRequest(request)

	if h.service.ExistsWithSectionNumber(section.Id, section.SectionNumber) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(responseJson)
		return
	}

	if srvErr := h.service.Create(section); srvErr != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		json.NewEncoder(w).Encode(responseJson)
		return
	}

	responseJson.Data = mappers.GetSectionResponseFromModel(section)
	response.JSON(w, http.StatusCreated, responseJson)

}

func (h *SectionHandler) Update(w http.ResponseWriter, r *http.Request) {

	var (
		request      *requests.SectionRequest = &requests.SectionRequest{}
		responseJson *responses.DataResponse  = &responses.DataResponse{}
	)

	w.Header().Set("Content-Type", "application/json")

	idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
	if convErr != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		json.NewEncoder(w).Encode(responseJson)
		return
	}

	section, srvErr := h.service.GetByID(idParam)
	if srvErr != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responseJson)
		return
	}

	if reqErr := json.NewDecoder(r.Body).Decode(request); reqErr != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		json.NewEncoder(w).Encode(responseJson)
		return
	}

	if valErr := h.validation.ValidateSectionRequestStruct(*request); valErr != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(responseJson)
		return
	}

	if h.service.ExistsWithSectionNumber(section.Id, request.SectionNumber) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(responseJson)
		return
	}

	mappers.UpdateSectionModelFromRequest(section, request)

	response.Data = mappers.GetSectionResponseFromModel(section)
	response.JSON(w, http.StatusOK, responseJson)

}

func (h *SectionHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {

	var (
		response *responses.DataResponse = &responses.DataResponse{}
	)

	w.Header().Set("Content-Type", "application/json")

	idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
	if convErr != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		json.NewEncoder(w).Encode(response)
		return
	}

	srvErr := h.service.DeleteByID(idParam)
	if srvErr != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
