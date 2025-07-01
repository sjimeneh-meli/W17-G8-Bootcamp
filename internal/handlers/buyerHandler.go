package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

func GetBuyerHandler(service services.BuyerServiceI) BuyerHandlerI {
	return &BuyerHandler{
		service: service,
	}
}

type BuyerHandlerI interface {
	GetAll() http.HandlerFunc
	GetById() http.HandlerFunc
	DeleteById() http.HandlerFunc
	PostBuyer() http.HandlerFunc
	PatchBuyer() http.HandlerFunc
}

type BuyerHandler struct {
	service services.BuyerServiceI
}

func (h *BuyerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			requestResponse *responses.DataResponse = &responses.DataResponse{}
			buyerResponse   []*responses.BuyerResponse
			buyers          []*models.Buyer
		)

		buyersMap, err := h.service.GetAll()
		if err != nil {
			requestResponse.SetError(err.Error())
			response.JSON(w, http.StatusInternalServerError, requestResponse)
			return
		}

		buyers = buyerMapToBuyerList(buyersMap)
		buyerResponse = mappers.GetListBuyerResponseFromListModel(buyers)
		requestResponse.Data = buyerResponse
		requestResponse.Message = "success"

		response.JSON(w, http.StatusOK, requestResponse)
	}
}

func (h *BuyerHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			requestResponse *responses.DataResponse = &responses.DataResponse{}
			buyerResponse   *responses.BuyerResponse
			buyer           models.Buyer
		)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			requestResponse.SetError(err.Error())
			response.JSON(w, http.StatusInternalServerError, requestResponse)
			return
		}

		buyer, err = h.service.GetById(id)
		if err != nil {
			requestResponse.SetError(err.Error())

			if errors.Is(err, error_message.ErrNotFound) {
				response.JSON(w, http.StatusNotFound, requestResponse)
				return
			}

			response.JSON(w, http.StatusInternalServerError, requestResponse)
			return
		}

		buyerResponse = mappers.GetResponseBuyerFromModel(&buyer)
		requestResponse.Data = buyerResponse
		requestResponse.Message = "success"
		response.JSON(w, http.StatusOK, requestResponse)
	}
}

func (h *BuyerHandler) DeleteById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestResponse *responses.DataResponse = &responses.DataResponse{}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			requestResponse.SetError(err.Error())
			response.JSON(w, http.StatusInternalServerError, requestResponse)
			return
		}

		err = h.service.DeleteById(id)
		if err != nil {
			requestResponse.SetError(err.Error())

			if errors.Is(err, error_message.ErrNotFound) {
				response.JSON(w, http.StatusNotFound, requestResponse)
				return
			}

			response.JSON(w, http.StatusInternalServerError, requestResponse)
			return
		}

		requestResponse.Message = "buyer deleted successfully"
		response.JSON(w, http.StatusNoContent, requestResponse)
	}
}

func (h *BuyerHandler) PostBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestResponse *responses.DataResponse = &responses.DataResponse{}

		requestBuyer := requests.BuyerRequest{}
		request.JSON(r, &requestBuyer)

		err := validations.ValidateBuyerRequestStruct(requestBuyer)
		if err != nil {
			requestResponse.SetError(err.Error())
			response.JSON(w, http.StatusUnprocessableEntity, requestResponse)
			return
		}

		modelBuyer := mappers.GetModelBuyerFromRequest(requestBuyer)
		buyerDb, err := h.service.Create(*modelBuyer)
		if err != nil {
			requestResponse.SetError(err.Error())

			if errors.Is(err, error_message.ErrAlreadyExists) {
				response.JSON(w, http.StatusConflict, requestResponse)
				return
			}

			response.JSON(w, http.StatusInternalServerError, requestResponse)
			return
		}

		buyerResponse := mappers.GetResponseBuyerFromModel(&buyerDb)
		requestResponse.Data = buyerResponse
		requestResponse.Message = "buyer created successfully"

		response.JSON(w, http.StatusCreated, requestResponse)
	}
}

func (h *BuyerHandler) PatchBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestResponse *responses.DataResponse = &responses.DataResponse{}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			requestResponse.SetError(err.Error())
			response.JSON(w, http.StatusInternalServerError, requestResponse)
			return
		}

		requestBuyer := requests.BuyerRequest{}
		request.JSON(r, &requestBuyer)
		err = validations.IsNotAnEmptyBuyer(requestBuyer)
		if err != nil {
			requestResponse.SetError(err.Error())
			response.JSON(w, http.StatusBadRequest, requestResponse)
			return
		}

		modelBuyer := mappers.GetModelBuyerFromRequest(requestBuyer)
		buyerDb, err := h.service.Update(id, *modelBuyer)
		if err != nil {
			requestResponse.SetError(err.Error())
			if errors.Is(err, error_message.ErrNotFound) {
				response.JSON(w, http.StatusNotFound, requestResponse)
				return
			}
			response.JSON(w, http.StatusInternalServerError, requestResponse)
			return
		}

		buyerResponse := mappers.GetResponseBuyerFromModel(&buyerDb)
		requestResponse.Data = buyerResponse
		requestResponse.Message = "buyer updated successfully"

		response.JSON(w, http.StatusOK, requestResponse)

	}
}

func buyerMapToBuyerList(buyers map[int]models.Buyer) []*models.Buyer {
	buyersList := []*models.Buyer{}
	for _, buyer := range buyers {
		buyersList = append(buyersList, &buyer)
	}
	return buyersList
}
