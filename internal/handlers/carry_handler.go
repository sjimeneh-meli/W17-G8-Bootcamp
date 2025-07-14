package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
)

type CarryHandler struct {
	carryService services.CarryService
}

func NewCarryHandler(carryService services.CarryService) *CarryHandler {
	return &CarryHandler{carryService: carryService}
}

func (h *CarryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request requests.CarryRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Mapear request a modelo
	carry := mappers.GetCarryModelFromRequest(&request)

	// Crear carry a través del servicio
	newCarry, err := h.carryService.Create(*carry)
	if err != nil {
		// Manejar error de CID duplicado
		if errors.Is(err, error_message.ErrAlreadyExists) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error creating carry")
		return
	}

	// Mapear modelo a respuesta
	carryResponse := mappers.GetCarryResponseFromModel(&newCarry)

	response.JSON(w, http.StatusCreated, responses.DataResponse{
		Data: carryResponse,
	})
}
