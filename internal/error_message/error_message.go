package error_message

import (
	"errors"
)

var (
	// ErrNotFound is returned when a requested resource does not exist (HTTP 404 Not Found).
	// ErrNotFound se devuelve cuando un recurso solicitado no existe (HTTP 404 Not Found).
	ErrNotFound = errors.New("error: the requested resource was not found")

	// ErrAlreadyExists is returned when an attempt is made to create a resource that already exists (HTTP 409 Conflict).
	// ErrAlreadyExists se devuelve cuando se intenta crear un recurso que ya existe (HTTP 409 Conflict).
	ErrAlreadyExists = errors.New("error: resource with the provided identifier already exists")

	// ErrInvalidInput is returned when the request payload is malformed or
	// contains invalid/missing fields (HTTP 422 Unprocessable Entity).
	// ErrInvalidInput se devuelve cuando el payload de la solicitud está mal formado o
	// contiene campos inválidos/faltantes (HTTP 422 Unprocessable Entity).
	ErrInvalidInput = errors.New("error: the provided input is invalid or missing required fields")

	// ErrInternalServerError is a generic error for unexpected server issues (HTTP 500 Internal Server Error).
	// ErrInternalServerError es un error genérico para problemas inesperados del servidor (HTTP 500 Internal Server Error).
	ErrInternalServerError = errors.New("error: an unexpected internal server error occurred")

	// ErrDependencyNotFound is returned when an entity depends on another entity that does not exist.
	// For example, trying to create an order with a non-existent product ID.
	// ErrDependencyNotFound se devuelve cuando una entidad depende de otra entidad que no existe.
	// Por ejemplo, intentar crear un pedido con un ID de producto inexistente.
	ErrDependencyNotFound = errors.New("error: a required dependent entity was not found")

	ErrFailedCheckingExistence = errors.New("error: failed checking locality existence")
	ErrQueryingReport          = errors.New("error: querying report failed")
	ErrFailedToScan            = errors.New("error: failed to scan record row")
)
