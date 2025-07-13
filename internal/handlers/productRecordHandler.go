package handlers

import (
	"context"
	"errors"
	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ProductRecordHandlerI interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetReportByIdProduct(w http.ResponseWriter, r *http.Request)
}

type productRecordHandler struct {
	Service services.ProductRecordServiceI
}

func NewProductRecordHandler(service services.ProductRecordServiceI) ProductRecordHandlerI {
	return &productRecordHandler{Service: service}
}

func (prh *productRecordHandler) Create(w http.ResponseWriter, r *http.Request) {
	parentCtx := context.Background()

	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)

	defer cancel()

	var productRecordRequest requests.ProductRecordRequest

	err := request.JSON(r, &productRecordRequest)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	v := validations.GetProductRecordValidation()

	err = v.ValidateProductRecordRequestStruct(productRecordRequest)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	productRecord := mappers.GetProductRecordFromRequest(productRecordRequest)

	result, err := prh.Service.CreateProductRecord(ctx, productRecord)

	if err != nil {
		if errors.Is(err, error_message.ErrDependencyNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	productRecordResponse := mappers.GetProductRecordResponseFromModel(result)

	response.JSON(w, http.StatusCreated, responses.DataResponse{
		Data: productRecordResponse,
	})

}
func (prh *productRecordHandler) GetReportByIdProduct(w http.ResponseWriter, r *http.Request) {
	parentCtx := context.Background()

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)

	defer cancel()

	idString := strings.TrimSpace(r.URL.Query().Get("id"))

	if idString == "" {
		response.Error(w, http.StatusBadRequest, "error: id is required")
		return
	}

	idInt, err := strconv.Atoi(idString)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "error: id not is a number")
		return
	}

	productRecordReport, err := prh.Service.GetReportByIdProduct(ctx, int64(idInt))

	if err != nil {
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{
		Data: []*models.ProductRecordReport{productRecordReport},
	})

}
