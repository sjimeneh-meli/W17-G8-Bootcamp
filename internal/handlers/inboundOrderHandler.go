package handlers

import (
	"context"
	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"net/http"
	"strconv"
	"time"
)

func GetInboundOrderHandler(service services.InboundOrdersServiceI) InboundOrderHandlerI {
	return &InboundOrderHandler{
		service: service,
	}
}

type InboundOrderHandlerI interface {
	GetInboundOrdersReport() http.HandlerFunc
}
type InboundOrderHandler struct {
	service services.InboundOrdersServiceI
}

func (h *InboundOrderHandler) GetInboundOrdersReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestResponse *responses.DataResponse = &responses.DataResponse{}
		var responseData []models.InboundOrderReport

		queryId := r.URL.Query().Get("id")
		if queryId != "" {
			employeeId, err := strconv.Atoi(queryId)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Invalid employee ID")
				return
			}

			report, err := h.service.GetInboundOrdersReportByEmployeeId(ctx, employeeId)
			if err != nil {
				response.Error(w, http.StatusInternalServerError, err.Error())
				return
			}

			responseData = append(responseData, report)
		} else {
			reports, err := h.service.GetAllInboundOrdersReports(ctx)
			if err != nil {
				response.Error(w, http.StatusInternalServerError, err.Error())
				return
			}
			responseData = reports
		}

		requestResponse.Data = responseData
		response.JSON(w, http.StatusOK, requestResponse)
	}
}
