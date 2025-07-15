package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

func GetProductBatchHandler(service services.ProductBatchServiceI, validation validations.ProductBatchValidation) ProductBatchHandlerI {
	return &ProductBatchHandler{
		service:    service,
		validation: &validation,
	}
}

type ProductBatchHandlerI interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type ProductBatchHandler struct {
	service    services.ProductBatchServiceI
	validation *validations.ProductBatchValidation
}

func (h *ProductBatchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		request      *requests.ProductBatchRequest = &requests.ProductBatchRequest{}
		responseJson *responses.DataResponse       = &responses.DataResponse{}
	)

	if reqErr := json.NewDecoder(r.Body).Decode(request); reqErr != nil {
		response.Error(w, http.StatusExpectationFailed, reqErr.Error())
		return
	}

	if valErr := h.validation.ValidateProductBatchRequestStruc(*request); valErr != nil {
		response.Error(w, http.StatusUnprocessableEntity, valErr.Error())
		return
	}

	productBatch, mapErr := mappers.GetProductBatchModelFromRequest(request)
	if mapErr != nil {
		response.Error(w, http.StatusExpectationFailed, mapErr.Error())
		return
	}

	if h.service.ExistsWithBatchNumber(productBatch.Id, productBatch.BatchNumber) {
		response.Error(w, http.StatusConflict, "already exist a batch with the same number")
		return
	}

	if srvErr := h.service.Create(productBatch); srvErr != nil {
		response.Error(w, http.StatusExpectationFailed, srvErr.Error())
		return
	}

	responseJson.Data = mappers.GetProductBatchResponseFromModel(productBatch)
	response.JSON(w, http.StatusCreated, responseJson)
}
