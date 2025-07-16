package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"net/http"
	"strconv"
	"time"
)

type LocalityHandler struct {
	service services.LocalityService
}

func NewLocalityHandler(service services.LocalityService) *LocalityHandler {
	return &LocalityHandler{service: service}
}

func (h *LocalityHandler) Save(w http.ResponseWriter, r *http.Request) {
	var localityToCreate requests.LocalityRequest
	errorBody := json.NewDecoder(r.Body).Decode(&localityToCreate)
	if errorBody != nil {
		response.Error(w, http.StatusBadRequest, errorBody.Error())
		return
	}
	data := localityToCreate.Data
	fmt.Println(data)
	if err := validations.ValidateLocalityRequestStruct(data); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	parentCtx := context.Background()
	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	defer cancel()

	localityCreated, err := h.service.Save(ctx, data)
	if err != nil {
		if errors.Is(err, error_message.ErrQuery) {
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		if errors.Is(err, error_message.ErrAlreadyExists) {
			response.Error(w, http.StatusConflict, err.Error())
		}
		if errors.Is(err, context.DeadlineExceeded) {
			response.Error(w, http.StatusGatewayTimeout, err.Error())
		}
		return
	}
	response.JSON(w, http.StatusOK, responses.DataResponse{Data: localityCreated})

}

func (h *LocalityHandler) GetSellerReportByLocality(w http.ResponseWriter, r *http.Request) {
	parentCtx := context.Background()
	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	defer cancel()

	localityIdStr := r.URL.Query().Get("id")
	var localityId int

	if r.URL.RawQuery != "" && localityIdStr == "" {
		response.Error(w, http.StatusInternalServerError, "Invalid query parameter")
		return
	}

	if localityIdStr == "" {
		// Si no viene el parámetro, usar 0 por defecto
		localityId = 0
	} else {
		// Si viene, intentar convertirlo a int
		var err error
		localityId, err = strconv.Atoi(localityIdStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	result, err := h.service.GetSellerReports(ctx, localityId)

	if errors.Is(err, error_message.ErrFailedCheckingExistence) || errors.Is(err, error_message.ErrQueryingReport) || errors.Is(err, error_message.ErrFailedToScan) {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if errors.Is(err, error_message.ErrNotFound) {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		response.Error(w, http.StatusGatewayTimeout, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{result})
}
