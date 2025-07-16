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

// GetEmployeeHandler creates and returns a new instance of EmployeeHandler with the required service and validation
// GetEmployeeHandler crea y retorna una nueva instancia de EmployeeHandler con el servicio y validación requeridos
func GetEmployeeHandler(service services.EmployeeServiceI, validation *validations.EmployeeValidation) EmployeeHandlerI {
	return &EmployeeHandler{
		service: service,
	}
}

// EmployeeHandlerI defines the contract for employee HTTP handlers with RESTful operations
// EmployeeHandlerI define el contrato para los manejadores HTTP de empleados con operaciones RESTful
type EmployeeHandlerI interface {
	GetAllEmployee() http.HandlerFunc
	PostEmployee() http.HandlerFunc
	GetByIdEmployee() http.HandlerFunc
	DeleteByIdEmployee() http.HandlerFunc
	PatchEmployee() http.HandlerFunc
}

// EmployeeHandler implements EmployeeHandlerI and handles HTTP requests for employee operations
// EmployeeHandler implementa EmployeeHandlerI y maneja las solicitudes HTTP para operaciones de empleados
type EmployeeHandler struct {
	service services.EmployeeServiceI // Service layer for employee business logic / Capa de servicio para lógica de negocio de empleados
}

// GetAllEmployee handles HTTP GET requests to retrieve all employees
// Returns a JSON response with all employees or appropriate error codes
// GetAllEmployee maneja las solicitudes HTTP GET para recuperar todos los empleados
// Retorna una respuesta JSON con todos los empleados o códigos de error apropiados
func (h *EmployeeHandler) GetAllEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var (
			requestResponse  *responses.DataResponse = &responses.DataResponse{}
			employeeResponse []*responses.EmployeeResponse
			employee         []*models.Employee
		)

		// Get all employees from service layer / Obtener todos los empleados de la capa de servicio
		employeeMap, err := h.service.GetAll(ctx)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Convert map to slice and map to response format / Convertir mapa a slice y mapear a formato de respuesta
		employee = employeeMapToList(employeeMap)
		employeeResponse = mappers.GetListEmployeeResponseFromListModel(employee)
		requestResponse.Data = employeeResponse

		response.JSON(w, http.StatusOK, requestResponse)
	}
}

// GetByIdEmployee handles HTTP GET requests to retrieve an employee by ID
// Extracts the ID from the URL parameter and returns the employee data or appropriate error codes
// GetByIdEmployee maneja las solicitudes HTTP GET para recuperar un empleado por ID
// Extrae el ID del parámetro de URL y retorna los datos del empleado o códigos de error apropiados
func (h *EmployeeHandler) GetByIdEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var (
			requestResponse  *responses.DataResponse = &responses.DataResponse{}
			employeeResponse *responses.EmployeeResponse
			employee         models.Employee
		)

		// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Get employee by ID from service layer / Obtener empleado por ID de la capa de servicio
		employee, err = h.service.GetById(ctx, id)
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
		employeeResponse = mappers.GetEmployeeResponseFromModel(&employee)
		requestResponse.Data = employeeResponse
		response.JSON(w, http.StatusOK, requestResponse)
	}
}

// PostEmployee handles HTTP POST requests to create a new employee
// Validates the request body and returns appropriate HTTP status codes
// PostEmployee maneja las solicitudes HTTP POST para crear un nuevo empleado
// Valida el cuerpo de la solicitud y retorna códigos de estado HTTP apropiados
func (h *EmployeeHandler) PostEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestResponse *responses.DataResponse = &responses.DataResponse{}

		// Parse and validate request body / Parsear y validar cuerpo de la solicitud
		requestEmployee := requests.EmployeeRequest{}
		request.JSON(r, &requestEmployee)

		err := validations.ValidateEmployeeRequestStruct(requestEmployee)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Map request to model and create through service / Mapear solicitud a modelo y crear a través del servicio
		modelEmployee := mappers.GetEmployeeModelFromRequest(requestEmployee)
		employeeDb, err := h.service.Create(ctx, *modelEmployee)
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
		employeeResponse := mappers.GetEmployeeResponseFromModel(&employeeDb)
		requestResponse.Data = employeeResponse

		response.JSON(w, http.StatusCreated, requestResponse)
	}
}

// PatchEmployee handles HTTP PATCH requests to update an existing employee
// Extracts the ID from the URL parameter and updates the employee with partial data
// PatchEmployee maneja las solicitudes HTTP PATCH para actualizar un empleado existente
// Extrae el ID del parámetro de URL y actualiza el empleado con datos parciales
func (h *EmployeeHandler) PatchEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		var requestResponse *responses.DataResponse = &responses.DataResponse{}

		// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Parse and validate request body for partial update / Parsear y validar cuerpo de solicitud para actualización parcial
		requestEmployee := requests.EmployeeRequest{}
		request.JSON(r, &requestEmployee)

		err = validations.IsNotAnEmptyEmployee(requestEmployee)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Map request to model and update through service / Mapear solicitud a modelo y actualizar a través del servicio
		modelEmployee := mappers.GetEmployeeModelFromRequest(requestEmployee)
		employeeDb, err := h.service.Update(ctx, id, *modelEmployee)
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
		employeeResponse := mappers.GetEmployeeResponseFromModel(&employeeDb)
		requestResponse.Data = employeeResponse

		response.JSON(w, http.StatusOK, requestResponse)
	}
}

// DeleteByIdEmployee handles HTTP DELETE requests to remove an employee by ID
// Extracts the ID from the URL parameter and deletes the employee
// DeleteByIdEmployee maneja las solicitudes HTTP DELETE para eliminar un empleado por ID
// Extrae el ID del parámetro de URL y elimina el empleado
func (h *EmployeeHandler) DeleteByIdEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set timeout context for the request / Establecer contexto con timeout para la solicitud
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		// Extract and validate ID parameter from URL / Extraer y validar parámetro ID de la URL
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Delete employee through service layer / Eliminar empleado a través de la capa de servicio
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

// employeeMapToList converts a map of employees to a slice of employee pointers
// Helper function for data transformation in the handler layer
// employeeMapToList convierte un mapa de empleados a un slice de punteros de empleados
// Función auxiliar para transformación de datos en la capa de manejadores
func employeeMapToList(employee map[int]models.Employee) []*models.Employee {
	employeeList := []*models.Employee{}
	for _, empl := range employee {
		employeeList = append(employeeList, &empl)
	}
	return employeeList
}
