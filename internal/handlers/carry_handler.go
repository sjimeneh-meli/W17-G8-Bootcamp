package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bootcamp-go/web/response"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

type CarryHandler struct {
	carryService services.CarryService
}

func NewCarryHandler(carryService services.CarryService) *CarryHandler {
	return &CarryHandler{carryService: carryService}
}

func (h *CarryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	var request requests.CarryRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := validations.ValidateCarryRequest(request); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Mapear request a modelo
	carry := mappers.MapCarryRequestToCarry(request)

	// Crear carry a través del servicio
	newCarry, err := h.carryService.CreateCarry(ctx, carry)

	if err != nil {
		if ctx.Err() != nil {
			response.Error(w, http.StatusRequestTimeout, "Request timeout cancelled")
			return
		}
		// Manejar error de CID duplicado o locality no encontrada
		if errors.Is(err, error_message.ErrAlreadyExists) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error creating carry")
		return
	}

	// Mapear modelo a respuesta
	carryResponse := mappers.MapCarryToCreateCarryResponse(newCarry)

	response.JSON(w, http.StatusCreated, responses.DataResponse{
		Data: carryResponse,
	})
}

func (h *CarryHandler) GetCarryReportByLocality(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Obtener el parámetro id de la query string
	localityIDStr := r.URL.Query().Get("id")
	var localityID int

	if localityIDStr != "" {
		// Si se proporciona un ID, validarlo
		var err error
		localityID, err = strconv.Atoi(localityIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid locality ID format")
			return
		}
		if localityID <= 0 {
			response.Error(w, http.StatusBadRequest, "Locality ID must be a positive number")
			return
		}
	}
	// Si localityIDStr está vacío, localityID será 0 (reporte de todas las localidades)

	reports, err := h.carryService.GetCarryReportByLocality(ctx, localityID)

	if err != nil {
		if ctx.Err() != nil {
			response.Error(w, http.StatusRequestTimeout, "Request timeout cancelled")
			return
		}
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "Locality not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error getting carry report by locality")
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{
		Data: reports,
	})
}
