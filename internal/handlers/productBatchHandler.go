package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

func GetProductBatchHandler(service services.ProductBatchServiceI,
	sectionService services.SectionServiceI,
	productService services.ProductService,
	validation validations.ProductBatchValidation) ProductBatchHandlerI {

	return &ProductBatchHandler{
		service:        service,
		sectionService: sectionService,
		productService: productService,
		validation:     &validation,
	}
}

type ProductBatchHandlerI interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetReportProduct(w http.ResponseWriter, r *http.Request)
}

type ProductBatchHandler struct {
	service        services.ProductBatchServiceI
	sectionService services.SectionServiceI
	productService services.ProductService
	validation     *validations.ProductBatchValidation
}

func (h *ProductBatchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		request      *requests.ProductBatchRequest = &requests.ProductBatchRequest{}
		responseJson *responses.DataResponse       = &responses.DataResponse{}
	)

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if reqErr := json.NewDecoder(r.Body).Decode(request); reqErr != nil {
		response.Error(w, http.StatusExpectationFailed, reqErr.Error())
		return
	}

	if valErr := h.validation.ValidateProductBatchRequestStruc(*request); valErr != nil {
		response.Error(w, http.StatusUnprocessableEntity, valErr.Error())
		return
	}

	if !h.sectionService.ExistWithID(request.SectionID) {
		response.Error(w, http.StatusNotFound, "section not found")
	}

	exists, _ := h.productService.ExistById(ctx, int64(request.ProductID))
	if !exists {
		response.Error(w, http.StatusNotFound, "product not found")
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

func (h *ProductBatchHandler) GetReportProduct(w http.ResponseWriter, r *http.Request) {
	var responseJson *responses.DataResponse = &responses.DataResponse{}

	idParamString := r.URL.Query().Get("id")
	if idParamString != "" {

		idParam, convErr := strconv.Atoi(idParamString)
		if convErr != nil {
			response.Error(w, http.StatusExpectationFailed, convErr.Error())
			return
		}

		section, srvErr := h.sectionService.GetByID(idParam)
		if srvErr != nil {
			response.Error(w, http.StatusNotFound, srvErr.Error())
			return
		}

		quantity := h.service.GetProductQuantityBySectionId(section.Id)
		responseJson.Data = map[string]any{
			"section_id":     section.Id,
			"section_number": section.SectionNumber,
			"products_count": quantity,
		}
		response.JSON(w, http.StatusCreated, responseJson)

	} else {
		sections, srvErr := h.sectionService.GetAll()
		data := make([]map[string]any, 0, len(sections))
		if srvErr != nil {
			response.Error(w, http.StatusExpectationFailed, srvErr.Error())
			return
		}

		for _, s := range sections {
			quantity := h.service.GetProductQuantityBySectionId(s.Id)
			data = append(data, map[string]any{
				"section_id":     s.Id,
				"section_number": s.SectionNumber,
				"products_count": quantity,
			})
		}
		responseJson.Data = data

		response.JSON(w, http.StatusCreated, responseJson)

	}

}
