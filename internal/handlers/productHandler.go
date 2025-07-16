package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time" // Import the time package

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

// ProductHandler maneja las solicitudes HTTP relacionadas con productos
// ProductHandler handles HTTP requests related to products
type ProductHandler struct {
	service services.ProductService // Dependencia del servicio de productos / Product service dependency
}

// NewProductHandler crea una nueva instancia del handler de productos con inyección de dependencias
// NewProductHandler creates a new instance of the product handler with dependency injection
func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// GetAll maneja las solicitudes GET para obtener todos los productos
// GetAll handles GET requests to retrieve all products
func (ph *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Establece timeout de 3 segundos para la operación
	// Set a 3-second timeout for the operation
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Obtiene todos los productos del servicio
	// Get all products from the service
	products, err := ph.service.GetAll(ctx)
	if err != nil {
		// Manejo de timeout
		// Handle timeout
		if errors.Is(err, context.DeadlineExceeded) {
			response.Error(w, http.StatusGatewayTimeout, "the request took too long to process")
			return
		}
		// Manejo de errores generales
		// Handle general errors
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Respuesta exitosa con los productos
	// Successful response with products
	response.JSON(w, http.StatusOK, responses.DataResponse{Data: products})
}

// Create maneja las solicitudes POST para crear un nuevo producto
// Create handles POST requests to create a new product
func (ph *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Establece timeout de 3 segundos para la operación
	// Set a 3-second timeout for the operation
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Decodifica el JSON del request
	// Decode JSON from request
	var productRequest requests.ProductRequest
	if err := request.JSON(r, &productRequest); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	// Valida la estructura del request
	// Validate request structure
	validation := validations.GetProductValidation()
	if err := validation.ValidateProductRequestStruct(productRequest); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Mapea el request a modelo de dominio
	// Map request to domain model
	product := mappers.GetProductFromRequest(productRequest)

	// Crea el producto a través del servicio
	// Create product through service
	newProduct, err := ph.service.Create(ctx, product)
	if err != nil {
		// Manejo de timeout
		// Handle timeout
		if errors.Is(err, context.DeadlineExceeded) {
			response.Error(w, http.StatusGatewayTimeout, "the request took too long to process")
			return
		}
		// Manejo de conflictos (producto ya existe)
		// Handle conflicts (product already exists)
		if errors.Is(err, error_message.ErrAlreadyExists) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		// Manejo de errores generales
		// Handle general errors
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Mapea el modelo a response y retorna
	// Map model to response and return
	productResponse := mappers.GetProductResponseFromModel(&newProduct)
	response.JSON(w, http.StatusCreated, responses.DataResponse{Data: productResponse})
}

// Get maneja las solicitudes GET para obtener un producto específico por ID
// Get handles GET requests to retrieve a specific product by ID
func (ph *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Establece timeout de 3 segundos para la operación
	// Set a 3-second timeout for the operation
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Extrae y valida el ID del parámetro de URL
	// Extract and validate ID from URL parameter
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Obtiene el producto por ID del servicio
	// Get product by ID from service
	product, err := ph.service.GetByID(ctx, id)
	if err != nil {
		// Manejo de timeout
		// Handle timeout
		if errors.Is(err, context.DeadlineExceeded) {
			response.Error(w, http.StatusGatewayTimeout, "the request took too long to process")
			return
		}
		// Manejo de producto no encontrado
		// Handle product not found
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		// Manejo de errores generales
		// Handle general errors
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Mapea el modelo a response y retorna
	// Map model to response and return
	productResponse := mappers.GetProductResponseFromModel(&product)
	response.JSON(w, http.StatusOK, responses.DataResponse{Data: productResponse})
}

// Update maneja las solicitudes PUT para actualizar un producto existente
// Update handles PUT requests to update an existing product
func (ph *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Establece timeout de 3 segundos para la operación
	// Set a 3-second timeout for the operation
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Extrae y valida el ID del parámetro de URL
	// Extract and validate ID from URL parameter
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Decodifica el JSON del request
	// Decode JSON from request
	var productRequest requests.ProductRequest
	if err := request.JSON(r, &productRequest); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	// Mapea el request a modelo de dominio
	// Map request to domain model
	productToUpdate := mappers.GetProductFromRequest(productRequest)

	// Actualiza el producto a través del servicio
	// Update product through service
	updatedProduct, err := ph.service.Update(ctx, id, productToUpdate)
	if err != nil {
		// Manejo de timeout
		// Handle timeout
		if errors.Is(err, context.DeadlineExceeded) {
			response.Error(w, http.StatusGatewayTimeout, "the request took too long to process")
			return
		}
		// Manejo de producto no encontrado
		// Handle product not found
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		// Manejo de errores generales
		// Handle general errors
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Mapea el modelo a response y retorna
	// Map model to response and return
	productResponse := mappers.GetProductResponseFromModel(&updatedProduct)
	response.JSON(w, http.StatusOK, responses.DataResponse{Data: productResponse})
}

// Delete maneja las solicitudes DELETE para eliminar un producto
// Delete handles DELETE requests to remove a product
func (ph *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Establece timeout de 3 segundos para la operación
	// Set a 3-second timeout for the operation
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Extrae y valida el ID del parámetro de URL
	// Extract and validate ID from URL parameter
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Elimina el producto a través del servicio
	// Delete product through service
	if err := ph.service.Delete(ctx, id); err != nil {
		// Manejo de timeout
		// Handle timeout
		if errors.Is(err, context.DeadlineExceeded) {
			response.Error(w, http.StatusGatewayTimeout, "the request took too long to process")
			return
		}
		// Manejo de producto no encontrado
		// Handle product not found
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		// Manejo de errores generales
		// Handle general errors
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Respuesta exitosa sin contenido
	// Successful response with no content
	w.WriteHeader(http.StatusNoContent)
}

// parseID extrae y valida el parámetro ID de la URL
// parseID extracts and validates the ID parameter from the URL
func parseID(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return 0, errors.New("id parameter is required")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id parameter: '%s' is not a valid number", idStr)
	}
	return id, nil
}
