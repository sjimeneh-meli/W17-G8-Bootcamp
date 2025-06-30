package handlers

import (
	"encoding/json"
	"net/http"

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
}

type BuyerHandler struct {
	service services.BuyerServiceI
}

func (h *BuyerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			response      *responses.DataResponse = &responses.DataResponse{}
			buyerResponse []*responses.BuyerResponse
			buyers        []*models.Buyer
		)
		w.Header().Set("Content-Type", "application/json")

		buyers = h.service.GetAll()
		buyerResponse = mappers.GetListBuyerResponseFromListModel(buyers)
		response.Data = buyerResponse

		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(response)

		if encodeErr != nil {
			response.SetError(encodeErr.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}
