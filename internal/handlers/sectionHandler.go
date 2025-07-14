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
	var responseJson *responses.DataResponse = &responses.DataResponse{}

	sections, srvErr := h.service.GetAll()
	if srvErr != nil {
		response.Error(w, http.StatusNotFound, srvErr.Error())
		return
	}

	responseJson.Data = mappers.GetListSectionResponseFromListModel(sections)
	response.JSON(w, http.StatusOK, responseJson)
}

func (h *SectionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	var responseJson *responses.DataResponse = &responses.DataResponse{}

	idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
	if convErr != nil {
		response.Error(w, http.StatusExpectationFailed, convErr.Error())
		return
	}

	section, srvErr := h.service.GetByID(idParam)
	if srvErr != nil {
		response.Error(w, http.StatusNotFound, srvErr.Error())
		return
	}

	responseJson.Data = mappers.GetSectionResponseFromModel(section)
	response.JSON(w, http.StatusOK, responseJson)

}

func (h *SectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	//VALIDATE IF PRODUCT AND WAREHOUSE EXIST
	var (
		request      *requests.SectionRequest = &requests.SectionRequest{}
		responseJson *responses.DataResponse  = &responses.DataResponse{}
		section      *models.Section
	)

	if reqErr := json.NewDecoder(r.Body).Decode(request); reqErr != nil {
		response.Error(w, http.StatusExpectationFailed, reqErr.Error())
		return
	}

	if valErr := h.validation.ValidateSectionRequestStruct(*request); valErr != nil {
		response.Error(w, http.StatusUnprocessableEntity, valErr.Error())
		return
	}

	section = mappers.GetSectionModelFromRequest(request)

	if h.service.ExistsWithSectionNumber(section.Id, section.SectionNumber) {
		response.Error(w, http.StatusConflict, "already exist a section with the same number")
		return
	}

	if srvErr := h.service.Create(section); srvErr != nil {
		response.Error(w, http.StatusExpectationFailed, srvErr.Error())
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

	idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
	if convErr != nil {
		response.Error(w, http.StatusExpectationFailed, convErr.Error())
		return
	}

	section, srvErr := h.service.GetByID(idParam)
	if srvErr != nil {
		response.Error(w, http.StatusNotFound, srvErr.Error())
		return
	}

	if reqErr := json.NewDecoder(r.Body).Decode(request); reqErr != nil {
		response.Error(w, http.StatusExpectationFailed, reqErr.Error())
		return
	}

	if valErr := h.validation.ValidateSectionRequestStruct(*request); valErr != nil {
		response.Error(w, http.StatusUnprocessableEntity, valErr.Error())
		return
	}

	if h.service.ExistsWithSectionNumber(section.Id, request.SectionNumber) {
		response.Error(w, http.StatusConflict, "already exist a section with the same number")
		return
	}

	mappers.UpdateSectionModelFromRequest(section, request)
	if srvErr := h.service.Update(section); srvErr != nil {
		response.Error(w, http.StatusExpectationFailed, srvErr.Error())
		return
	}

	responseJson.Data = mappers.GetSectionResponseFromModel(section)
	response.JSON(w, http.StatusOK, responseJson)

}

func (h *SectionHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {

	idParam, convErr := strconv.Atoi(chi.URLParam(r, "id"))
	if convErr != nil {
		response.Error(w, http.StatusExpectationFailed, convErr.Error())
		return
	}

	srvErr := h.service.DeleteByID(idParam)
	if srvErr != nil {
		response.Error(w, http.StatusNotFound, srvErr.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
