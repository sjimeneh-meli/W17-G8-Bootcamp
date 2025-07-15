package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

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

// GetBuyerHandler creates and returns a new instance of BuyerHandler
func GetBuyerHandler(service services.BuyerServiceI) BuyerHandlerI {
	return &BuyerHandler{
		service: service,
	}
}

// BuyerHandlerI defines the interface for buyer HTTP handlers
type BuyerHandlerI interface {
	GetAll() http.HandlerFunc
	GetById() http.HandlerFunc
	DeleteById() http.HandlerFunc
	PostBuyer() http.HandlerFunc
	PatchBuyer() http.HandlerFunc
}

// BuyerHandler implements BuyerHandlerI and handles HTTP requests for buyer operations
type BuyerHandler struct {
	service services.BuyerServiceI
}

// GetAll handles HTTP GET requests to retrieve all buyers
// Returns a JSON response with all buyers
func (h *BuyerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var (
			requestResponse *responses.DataResponse = &responses.DataResponse{}
			buyerResponse   []*responses.BuyerResponse
			buyers          []*models.Buyer
		)

		buyersMap, err := h.service.GetAll(ctx)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		buyers = buyerMapToBuyerList(buyersMap)
		buyerResponse = mappers.GetListBuyerResponseFromListModel(buyers)
		requestResponse.Data = buyerResponse

		response.JSON(w, http.StatusOK, requestResponse)
	}
}

// GetById handles HTTP GET requests to retrieve a buyer by ID
// Extracts the ID from the URL parameter and returns the buyer data
func (h *BuyerHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var (
			requestResponse *responses.DataResponse = &responses.DataResponse{}
			buyerResponse   *responses.BuyerResponse
			buyer           models.Buyer
		)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		buyer, err = h.service.GetById(ctx, id)
		if err != nil {

			if errors.Is(err, error_message.ErrNotFound) {
				response.Error(w, http.StatusNotFound, err.Error())
				return
			}

			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		buyerResponse = mappers.GetResponseBuyerFromModel(&buyer)
		requestResponse.Data = buyerResponse
		response.JSON(w, http.StatusOK, requestResponse)
	}
}

// DeleteById handles HTTP DELETE requests to remove a buyer by ID
// Extracts the ID from the URL parameter and deletes the buyer
func (h *BuyerHandler) DeleteById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.service.DeleteById(ctx, id)
		if err != nil {

			if errors.Is(err, error_message.ErrNotFound) {
				response.Error(w, http.StatusNotFound, err.Error())
				return
			}

			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// PostBuyer handles HTTP POST requests to create a new buyer
// Validates the request body and returns appropriate HTTP status codes
func (h *BuyerHandler) PostBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestResponse *responses.DataResponse = &responses.DataResponse{}

		requestBuyer := requests.BuyerRequest{}
		request.JSON(r, &requestBuyer)

		err := validations.ValidateBuyerRequestStruct(requestBuyer)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		modelBuyer := mappers.GetModelBuyerFromRequest(requestBuyer)
		buyerDb, err := h.service.Create(ctx, *modelBuyer)
		if err != nil {

			if errors.Is(err, error_message.ErrAlreadyExists) {
				response.Error(w, http.StatusConflict, err.Error())
				return
			}

			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		buyerResponse := mappers.GetResponseBuyerFromModel(&buyerDb)
		requestResponse.Data = buyerResponse

		response.JSON(w, http.StatusCreated, requestResponse)
	}
}

// PatchBuyer handles HTTP PATCH requests to update an existing buyer
// Extracts the ID from the URL parameter and updates the buyer with partial data
func (h *BuyerHandler) PatchBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestResponse *responses.DataResponse = &responses.DataResponse{}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		requestBuyer := requests.BuyerRequest{}
		request.JSON(r, &requestBuyer)

		err = validations.IsNotAnEmptyBuyer(requestBuyer)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		modelBuyer := mappers.GetModelBuyerFromRequest(requestBuyer)
		buyerDb, err := h.service.Update(ctx, id, *modelBuyer)
		if err != nil {

			if errors.Is(err, error_message.ErrNotFound) {
				response.Error(w, http.StatusNotFound, err.Error())
				return
			}

			if errors.Is(err, error_message.ErrAlreadyExists) {
				response.Error(w, http.StatusConflict, err.Error())
				return
			}

			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		buyerResponse := mappers.GetResponseBuyerFromModel(&buyerDb)
		requestResponse.Data = buyerResponse

		response.JSON(w, http.StatusOK, requestResponse)

	}
}

// buyerMapToBuyerList converts a map of buyers to a slice of buyer pointers
func buyerMapToBuyerList(buyers map[int]models.Buyer) []*models.Buyer {
	buyersList := []*models.Buyer{}
	for _, buyer := range buyers {
		buyersList = append(buyersList, &buyer)
	}
	return buyersList
}
