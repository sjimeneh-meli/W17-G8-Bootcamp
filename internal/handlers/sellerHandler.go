package handlers

import (
	"encoding/json"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"net/http"
	"strconv"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
)

type SellerHandler struct {
	service services.SellerService
}

func NewSellerHandler(service services.SellerService) *SellerHandler {
	return &SellerHandler{service: service}
}

func (h *SellerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	sellers, err := h.service.GetAll()
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{Data: sellers})
}

func (h *SellerHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	idFormated, err := strconv.Atoi(id)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	seller, err1 := h.service.GetById(idFormated)

	if err1 != nil {
		response.Error(w, http.StatusNotFound, err1.Error())
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{Data: seller})
}

func (h *SellerHandler) Save(w http.ResponseWriter, r *http.Request) {
	var sellerToCreate models.SellerRequest
	data := r.Body
	errorBody := json.NewDecoder(data).Decode(&sellerToCreate)
	if errorBody != nil {
		response.Error(w, http.StatusBadRequest, errorBody.Error())
		return
	}

	sellerParced := mappers.ToRequestToSellerStruct(sellerToCreate)
	sellerCreated, err := h.service.Save(sellerParced)

	if err != nil {
		response.Error(w, http.StatusConflict, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{Data: sellerCreated})

}

func (h *SellerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	idFormated, err := strconv.Atoi(id)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	data := r.Body
	var bodyFormated models.SellerRequest
	errBody := json.NewDecoder(data).Decode(&bodyFormated)

	if errBody != nil {
		response.Error(w, http.StatusBadRequest, errBody.Error())
		return
	}
	sellerToUpdate := mappers.ToRequestToSellerStruct(bodyFormated)
	sellerUpdated, errUpdate := h.service.Update(idFormated, sellerToUpdate)

	if errUpdate != nil {
		response.Error(w, http.StatusNotFound, errUpdate.Error())
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{Data: sellerUpdated})
}

func (h *SellerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	idFormated, err := strconv.Atoi(id)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	errDelete := h.service.Delete(idFormated)

	if errDelete != nil {
		response.Error(w, http.StatusNotFound, errDelete.Error())
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
