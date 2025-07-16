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
	"time"
)

func GetInboundOrderHandler(service services.InboundOrdersServiceI) InboundOrderHandlerI {
	return &InboundOrderHandler{
		service: service,
	}
}

type InboundOrderHandlerI interface {
	GetInboundOrdersReport() http.HandlerFunc
	PostInboundOrder() http.HandlerFunc
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

func (h *InboundOrderHandler) PostInboundOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var (
			reqResponse    *responses.DataResponse      = &responses.DataResponse{}
			requestInbound requests.InboundOrderRequest = requests.InboundOrderRequest{}
		)

		err := request.JSON(r, &requestInbound)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		err = validations.ValidateInboundOrderRequestStruct(requestInbound)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		modelInbound := mappers.GetModelInboundOrderFromRequest(requestInbound)
		order, err := h.service.Create(ctx, *modelInbound)
		if err != nil {
			switch {
			case errors.Is(err, error_message.ErrAlreadyExists):
				response.Error(w, http.StatusConflict, err.Error())
				return
			case errors.Is(err, error_message.ErrDependencyNotFound):
				response.Error(w, http.StatusConflict, err.Error())
				return
			default:
				response.Error(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		orderResponse := mappers.GetResponseInboundOrderFromModel(&order)
		reqResponse.Data = orderResponse

		response.JSON(w, http.StatusCreated, reqResponse)
	}
}
