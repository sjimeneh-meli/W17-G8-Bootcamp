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

// GetBuyerHandler creates and returns a new instance of BuyerHandler with the required service
// GetBuyerHandler crea y retorna una nueva instancia de BuyerHandler con el servicio requerido
func GetBuyerHandler(service services.BuyerServiceI) BuyerHandlerI {
	return &BuyerHandler{
		service: service,
	}
}

// BuyerHandlerI defines the contract for buyer HTTP handlers with RESTful operations
// BuyerHandlerI define el contrato para los manejadores HTTP de compradores con operaciones RESTful
type BuyerHandlerI interface {
	GetAll() http.HandlerFunc
	GetById() http.HandlerFunc
	DeleteById() http.HandlerFunc
	PostBuyer() http.HandlerFunc
	PatchBuyer() http.HandlerFunc
}

// BuyerHandler implements BuyerHandlerI and handles HTTP requests for buyer operations
// BuyerHandler implementa BuyerHandlerI y maneja las solicitudes HTTP para operaciones de compradores
type BuyerHandler struct {
	service services.BuyerServiceI // Service layer for buyer business logic / Capa de servicio para lógica de negocio de compradores
}

// GetAll handles HTTP GET requests to retrieve all buyers
// Returns a JSON response with all buyers or appropriate error codes
// GetAll maneja las solicitudes HTTP GET para recuperar todos los compradores
// Retorna una respuesta JSON con todos los compradores o códigos de error apropiados
func (h *BuyerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var (
			requestResponse *responses.DataResponse = &responses.DataResponse{}
			buyerResponse   []*responses.BuyerResponse
			buyers          []*models.Buyer
		)

		// Get all buyers from service layer / Obtener todos los compradores de la capa de servicio
		buyersMap, err := h.service.GetAll(ctx)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Convert map to slice and map to response format / Convertir mapa a slice y mapear a formato de respuesta
		buyers = buyerMapToBuyerList(buyersMap)
		buyerResponse = mappers.GetListBuyerResponseFromListModel(buyers)
		requestResponse.Data = buyerResponse

		response.JSON(w, http.StatusOK, requestResponse)
	}
}

// GetById handles HTTP GET requests to retrieve a buyer by ID
// Extracts the ID from the URL parameter and returns the buyer data or appropriate error codes
// GetById maneja las solicitudes HTTP GET para recuperar un comprador por ID
// Extrae el ID del parámetro de URL y retorna los datos del comprador o códigos de error apropiados
func (h *BuyerHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var (
			requestResponse *responses.DataResponse = &responses.DataResponse{}
			buyerResponse   *responses.BuyerResponse
			buyer           models.Buyer
		)

		// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, error_message.ErrInvalidInput.Error())
			return
		}

		// Get buyer by ID from service layer / Obtener comprador por ID de la capa de servicio
		buyer, err = h.service.GetById(ctx, id)
		if err != nil {
			// Handle specific error types / Manejar tipos de error específicos
			if errors.Is(err, error_message.ErrNotFound) {
				response.Error(w, http.StatusNotFound, err.Error())
				return
			}

			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Map model to response format / Mapear modelo a formato de respuesta
		buyerResponse = mappers.GetResponseBuyerFromModel(&buyer)
		requestResponse.Data = buyerResponse
		response.JSON(w, http.StatusOK, requestResponse)
	}
}

// DeleteById handles HTTP DELETE requests to remove a buyer by ID
// Extracts the ID from the URL parameter and deletes the buyer
// DeleteById maneja las solicitudes HTTP DELETE para eliminar un comprador por ID
// Extrae el ID del parámetro de URL y elimina el comprador
func (h *BuyerHandler) DeleteById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, error_message.ErrInvalidInput.Error())
			return
		}

		// Delete buyer through service layer / Eliminar comprador a través de la capa de servicio
		err = h.service.DeleteById(ctx, id)
		if err != nil {
			// Handle specific error types / Manejar tipos de error específicos
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
// PostBuyer maneja las solicitudes HTTP POST para crear un nuevo comprador
// Valida el cuerpo de la solicitud y retorna códigos de estado HTTP apropiados
func (h *BuyerHandler) PostBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestResponse *responses.DataResponse = &responses.DataResponse{}

		// Parse and validate request body / Parsear y validar cuerpo de la solicitud
		requestBuyer := requests.BuyerRequest{}
		request.JSON(r, &requestBuyer)

		err := validations.ValidateBuyerRequestStruct(requestBuyer)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Map request to model and create through service / Mapear solicitud a modelo y crear a través del servicio
		modelBuyer := mappers.GetModelBuyerFromRequest(requestBuyer)
		buyerDb, err := h.service.Create(ctx, *modelBuyer)
		if err != nil {
			// Handle specific error types / Manejar tipos de error específicos
			if errors.Is(err, error_message.ErrAlreadyExists) {
				response.Error(w, http.StatusConflict, err.Error())
				return
			}

			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Map model to response format / Mapear modelo a formato de respuesta
		buyerResponse := mappers.GetResponseBuyerFromModel(&buyerDb)
		requestResponse.Data = buyerResponse

		response.JSON(w, http.StatusCreated, requestResponse)
	}
}

// PatchBuyer handles HTTP PATCH requests to update an existing buyer
// Extracts the ID from the URL parameter and updates the buyer with partial data
// PatchBuyer maneja las solicitudes HTTP PATCH para actualizar un comprador existente
// Extrae el ID del parámetro de URL y actualiza el comprador con datos parciales
func (h *BuyerHandler) PatchBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestResponse *responses.DataResponse = &responses.DataResponse{}

		// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, error_message.ErrInvalidInput.Error())
			return
		}

		// Parse and validate request body for partial update / Parsear y validar cuerpo de solicitud para actualización parcial
		requestBuyer := requests.BuyerRequest{}
		request.JSON(r, &requestBuyer)

		err = validations.IsNotAnEmptyBuyer(requestBuyer)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Map request to model and update through service / Mapear solicitud a modelo y actualizar a través del servicio
		modelBuyer := mappers.GetModelBuyerFromRequest(requestBuyer)
		buyerDb, err := h.service.Update(ctx, id, *modelBuyer)
		if err != nil {
			// Handle specific error types / Manejar tipos de error específicos
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

		// Map model to response format / Mapear modelo a formato de respuesta
		buyerResponse := mappers.GetResponseBuyerFromModel(&buyerDb)
		requestResponse.Data = buyerResponse

		response.JSON(w, http.StatusOK, requestResponse)
	}
}

// buyerMapToBuyerList converts a map of buyers to a slice of buyer pointers
// Helper function for data transformation in the handler layer
// buyerMapToBuyerList convierte un mapa de compradores a un slice de punteros de compradores
// Función auxiliar para transformación de datos en la capa de manejadores
func buyerMapToBuyerList(buyers map[int]models.Buyer) []*models.Buyer {
	buyersList := []*models.Buyer{}
	for _, buyer := range buyers {
		buyersList = append(buyersList, &buyer)
	}
	return buyersList
}
