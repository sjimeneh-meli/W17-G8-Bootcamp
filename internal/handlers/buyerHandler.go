package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
)

func GetBuyerHandler(service services.BuyerServiceI) BuyerHandlerI {
	return &BuyerHandler{
		service: service,
	}
}

type BuyerHandlerI interface {
	GetAll() http.HandlerFunc
	GetById() http.HandlerFunc
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
		w.Header().Set("Content-Type", "application/json")

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

func buyerMapToBuyerList(buyers map[int]models.Buyer) []*models.Buyer {
	buyersList := []*models.Buyer{}
	for _, buyer := range buyers {
		buyersList = append(buyersList, &buyer)
	}
	return buyersList
}
