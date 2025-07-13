package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
)

// ProductHandler handles HTTP requests for product operations.
// It acts as a bridge between HTTP requests and business logic services.
//
// ProductHandler maneja las peticiones HTTP para operaciones de productos.
// Actúa como un puente entre las peticiones HTTP y los servicios de lógica de negocio.
type ProductHandler struct {
	service services.ProductServiceI
}

// NewProductHandler creates a new product handler instance with the provided service.
// It follows the dependency injection pattern for better testability and modularity.
//
// NewProductHandler crea una nueva instancia del manejador de productos con el servicio proporcionado.
// Sigue el patrón de inyección de dependencias para mejor testabilidad y modularidad.
func NewProductHandler(service services.ProductServiceI) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// GetAll retrieves all products from the system.
// Returns 200 OK with products data or appropriate error status codes.
//
// GetAll obtiene todos los productos del sistema.
// Retorna 200 OK con los datos de productos o códigos de estado de error apropiados.
func (ph *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := ph.service.GetAll()

	if err != nil {
		// Handle not found case / Manejar caso de no encontrado
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		// Handle internal server errors / Manejar errores internos del servidor
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{Data: products})
}

// Save creates a new product in the system.
// Validates the request payload and returns 201 Created on success.
// Returns 400 Bad Request for validation errors or 409 Conflict if product already exists.
//
// Save crea un nuevo producto en el sistema.
// Valida la carga útil de la petición y retorna 201 Created en caso de éxito.
// Retorna 400 Bad Request para errores de validación o 409 Conflict si el producto ya existe.
func (ph *ProductHandler) Save(w http.ResponseWriter, r *http.Request) {
	var productRequest requests.ProductRequest

	// Parse JSON request body / Parsear el cuerpo de la petición JSON
	err := request.JSON(r, &productRequest)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request data / Validar datos de la petición
	validation := validations.GetProductValidation()

	err = validation.ValidateProductRequestStruct(productRequest)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Map request to domain model / Mapear petición a modelo de dominio
	product := mappers.GetProductFromRequest(productRequest)

	// Create product through service layer / Crear producto a través de la capa de servicio
	newProduct, err := ph.service.Create(product)

	if err != nil {
		// Handle duplicate product case / Manejar caso de producto duplicado
		if errors.Is(err, error_message.ErrAlreadyExists) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map domain model to response / Mapear modelo de dominio a respuesta
	productResponse := mappers.GetProductResponseFromModel(&newProduct)

	response.JSON(w, http.StatusCreated, responses.DataResponse{Data: productResponse})

}

// GetById retrieves a specific product by its ID.
// Validates the ID parameter and returns 200 OK with product data or 404 Not Found.
//
// GetById obtiene un producto específico por su ID.
// Valida el parámetro ID y retorna 200 OK con los datos del producto o 404 Not Found.
func (ph *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {

	// Extract and validate ID from URL path / Extraer y validar ID de la ruta URL
	idString := strings.TrimSpace(chi.URLParam(r, "id"))

	if idString == "" {
		response.Error(w, http.StatusBadRequest, "error: id is required")
		return
	}

	// Convert ID to integer / Convertir ID a entero
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "error: id is not a number")
		return
	}

	// Retrieve product from service / Obtener producto del servicio
	product, err := ph.service.GetByID(idInt)

	if err != nil {
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, err.Error())
		return

	}

	// Map to response format / Mapear a formato de respuesta
	productResponse := mappers.GetProductResponseFromModel(product)

	response.JSON(w, http.StatusOK, responses.DataResponse{Data: productResponse})

}

// Update modifies an existing product by its ID.
// Validates the ID parameter and request payload, returns 200 OK on success.
//
// Update modifica un producto existente por su ID.
// Valida el parámetro ID y la carga útil de la petición, retorna 200 OK en caso de éxito.
func (ph *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extract and validate ID from URL path / Extraer y validar ID de la ruta URL
	idString := strings.TrimSpace(chi.URLParam(r, "id"))

	if idString == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	idInt, err := strconv.Atoi(idString)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "id is not a number")
		return
	}

	var productRequest requests.ProductRequest

	// Parse update request payload / Parsear carga útil de petición de actualización
	err = request.JSON(r, &productRequest)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Map request to domain model / Mapear petición a modelo de dominio
	productToUpdate := mappers.GetProductFromRequest(productRequest)

	// Update product through service layer / Actualizar producto a través de la capa de servicio
	productUpdated, err := ph.service.UpdateById(idInt, productToUpdate)

	if err != nil {
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, responses.DataResponse{Data: productUpdated})

}

// DeleteById removes a product from the system by its ID.
// Validates the ID parameter and returns 204 No Content on successful deletion.
//
// DeleteById elimina un producto del sistema por su ID.
// Valida el parámetro ID y retorna 204 No Content en caso de eliminación exitosa.
func (ph *ProductHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	// Extract and validate ID from URL path / Extraer y validar ID de la ruta URL
	idString := strings.TrimSpace(chi.URLParam(r, "id"))

	if idString == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	idInt, err := strconv.Atoi(idString)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "id is not a number")
		return
	}

	// Delete product through service layer / Eliminar producto a través de la capa de servicio
	err = ph.service.DeleteById(idInt)

	if err != nil {
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
