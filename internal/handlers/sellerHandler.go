package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
)

// SellerHandler handles HTTP requests for seller operations
// SellerHandler maneja las solicitudes HTTP para operaciones de vendedores
type SellerHandler struct {
	service services.SellerService // Service layer for seller business logic / Capa de servicio para lógica de negocio de vendedores
}

// NewSellerHandler creates and returns a new instance of SellerHandler with the required service
// NewSellerHandler crea y retorna una nueva instancia de SellerHandler con el servicio requerido
func NewSellerHandler(service services.SellerService) *SellerHandler {
	return &SellerHandler{service: service}
}

// GetAll handles HTTP GET requests to retrieve all sellers
// Returns a JSON response with all sellers or appropriate error codes
// GetAll maneja las solicitudes HTTP GET para recuperar todos los vendedores
// Retorna una respuesta JSON con todos los vendedores o códigos de error apropiados
func (h *SellerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Get all sellers from service layer / Obtener todos los vendedores de la capa de servicio
	sellers, err := h.service.GetAll()
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{Data: sellers})
}

// GetById handles HTTP GET requests to retrieve a seller by ID
// Extracts the ID from the URL parameter and returns the seller data or appropriate error codes
// GetById maneja las solicitudes HTTP GET para recuperar un vendedor por ID
// Extrae el ID del parámetro de URL y retorna los datos del vendedor o códigos de error apropiados
func (h *SellerHandler) GetById(w http.ResponseWriter, r *http.Request) {
	// Extract ID parameter from URL / Extraer parámetro ID de la URL
	id := chi.URLParam(r, "id")

	// Parse and validate ID parameter / Parsear y validar parámetro ID
	idFormated, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get seller by ID from service layer / Obtener vendedor por ID de la capa de servicio
	seller, err1 := h.service.GetById(idFormated)
	if err1 != nil {
		response.Error(w, http.StatusNotFound, err1.Error())
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{Data: seller})
}

// Save handles HTTP POST requests to create a new seller
// Validates the request body and returns appropriate HTTP status codes
// Save maneja las solicitudes HTTP POST para crear un nuevo vendedor
// Valida el cuerpo de la solicitud y retorna códigos de estado HTTP apropiados
func (h *SellerHandler) Save(w http.ResponseWriter, r *http.Request) {
	var sellerToCreate requests.SellerRequest
	data := r.Body

	// Parse and validate JSON request body / Parsear y validar cuerpo de solicitud JSON
	errorBody := json.NewDecoder(data).Decode(&sellerToCreate)
	if errorBody != nil {
		response.Error(w, http.StatusBadRequest, errorBody.Error())
		return
	}

	// Validate request structure and business rules / Validar estructura de solicitud y reglas de negocio
	if err := validations.ValidateSellerRequestStruct(sellerToCreate); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Map request to seller model / Mapear solicitud a modelo de vendedor
	sellerParced := mappers.ToRequestToSellerStruct(sellerToCreate)

	// Create seller through service layer / Crear vendedor a través de la capa de servicio
	sellerCreated, err := h.service.Save(sellerParced)
	if err != nil {
		response.Error(w, http.StatusConflict, err.Error())
		return
	}

	// Map model to response format / Mapear modelo a formato de respuesta
	sellerResponse := mappers.ToSellerStructToResponse(sellerCreated[0])

	response.JSON(w, http.StatusOK, responses.DataResponse{Data: sellerResponse})
}

// Update handles HTTP PUT requests to update an existing seller
// Extracts the ID from the URL parameter and updates the seller data
// Update maneja las solicitudes HTTP PUT para actualizar un vendedor existente
// Extrae el ID del parámetro de URL y actualiza los datos del vendedor
func (h *SellerHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	// Parse and validate ID parameter / Parsear y validar parámetro ID
	idFormated, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse and validate JSON request body / Parsear y validar cuerpo de solicitud JSON
	data := r.Body
	var bodyFormated requests.SellerRequest
	errBody := json.NewDecoder(data).Decode(&bodyFormated)
	if errBody != nil {
		response.Error(w, http.StatusBadRequest, errBody.Error())
		return
	}

	// Map request to seller model and update through service / Mapear solicitud a modelo de vendedor y actualizar a través del servicio
	sellerToUpdate := mappers.ToRequestToSellerStruct(bodyFormated)
	sellerUpdated, errUpdate := h.service.Update(idFormated, sellerToUpdate)
	if errUpdate != nil {
		response.Error(w, http.StatusNotFound, errUpdate.Error())
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{Data: sellerUpdated})
}

// Delete handles HTTP DELETE requests to remove a seller by ID
// Extracts the ID from the URL parameter and deletes the seller
// Delete maneja las solicitudes HTTP DELETE para eliminar un vendedor por ID
// Extrae el ID del parámetro de URL y elimina el vendedor
func (h *SellerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	// Parse and validate ID parameter / Parsear y validar parámetro ID
	idFormated, err := strconv.Atoi(id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Delete seller through service layer / Eliminar vendedor a través de la capa de servicio
	errDelete := h.service.Delete(idFormated)
	if errDelete != nil {
		response.Error(w, http.StatusNotFound, errDelete.Error())
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
