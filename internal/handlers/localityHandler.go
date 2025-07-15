package handlers

import (
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

	localityCreated, err := h.service.Save(data)
	if err != nil {
		response.Error(w, http.StatusConflict, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, responses.DataResponse{Data: localityCreated})

}

func (h *LocalityHandler) GetSellerReportByLocality(w http.ResponseWriter, r *http.Request) {
	localityIdStr := r.URL.Query().Get("id")
	var localityId int

	if localityIdStr == "" {
		// Si no viene el parámetro, usar 0 por defecto
		localityId = 0
	} else {
		// Si viene, intentar convertirlo a int
		var err error
		localityId, err = strconv.Atoi(localityIdStr)
		if err != nil {
			http.Error(w, "Invalid 'id' parameter: must be an integer", http.StatusUnprocessableEntity)
			return
		}
	}
	result, err := h.service.GetSellerReports(localityId)

	if errors.Is(err, error_message.ErrFailedCheckingExistence) || errors.Is(err, error_message.ErrQueryingReport) || errors.Is(err, error_message.ErrFailedToScan) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, error_message.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	/*resultFormated := responses.LocalitySellerReportResponse{
		Data: result,
	}*/

	response.JSON(w, http.StatusOK, responses.DataResponse{result})
}
