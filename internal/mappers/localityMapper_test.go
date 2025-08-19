// Package mappers - Locality Mapper Unit Tests / Tests Unitarios del Mapper Locality
// Data transformation testing for locality mapping between request and model structures
// Testing de transformación de datos para mapeo de localidad entre estructuras de request y modelo
package mappers

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestToRequestToLocalityStruct - Tests for ToRequestToLocalityStruct function
// Tests data mapping from LocalityRequest to Locality model
// Tests para función ToRequestToLocalityStruct
// Prueba el mapeo de datos de LocalityRequest a modelo Locality
func TestToRequestToLocalityStruct(t *testing.T) {
	t.Run("success - maps locality request to model correctly", func(t *testing.T) {
		// Test: Valid locality request data maps successfully to model / Datos válidos de solicitud de localidad se mapean exitosamente al modelo
		// Arrange
		request := requests.LocalityRequest{
			Data: models.Locality{
				Id:           1, // This should be ignored in mapping / Esto debería ser ignorado en el mapeo
				LocalityName: "Buenos Aires",
				ProvinceName: "Ciudad Autónoma de Buenos Aires",
				CountryName:  "Argentina",
			},
		}

		expectedLocality := models.Locality{
			Id:           0, // ID should not be mapped from request / ID no debería ser mapeado del request
			LocalityName: "Buenos Aires",
			ProvinceName: "Ciudad Autónoma de Buenos Aires",
			CountryName:  "Argentina",
		}

		// Act
		result := ToRequestToLocalityStruct(request)

		// Assert
		assert.Equal(t, expectedLocality, result)
		assert.Equal(t, "Buenos Aires", result.LocalityName)
		assert.Equal(t, "Ciudad Autónoma de Buenos Aires", result.ProvinceName)
		assert.Equal(t, "Argentina", result.CountryName)
		assert.Equal(t, 0, result.Id) // ID should be zero value / ID debería ser valor cero
	})

	t.Run("success - maps locality request with empty values", func(t *testing.T) {
		// Test: Empty string values in request map correctly / Valores de string vacíos en request se mapean correctamente
		// Arrange
		request := requests.LocalityRequest{
			Data: models.Locality{
				Id:           999, // This should be ignored / Esto debería ser ignorado
				LocalityName: "",
				ProvinceName: "",
				CountryName:  "",
			},
		}

		expectedLocality := models.Locality{
			Id:           0,
			LocalityName: "",
			ProvinceName: "",
			CountryName:  "",
		}

		// Act
		result := ToRequestToLocalityStruct(request)

		// Assert
		assert.Equal(t, expectedLocality, result)
		assert.Equal(t, "", result.LocalityName)
		assert.Equal(t, "", result.ProvinceName)
		assert.Equal(t, "", result.CountryName)
		assert.Equal(t, 0, result.Id)
	})

	t.Run("success - maps locality request with special characters", func(t *testing.T) {
		// Test: Special characters and symbols map correctly / Caracteres especiales y símbolos se mapean correctamente
		// Arrange
		request := requests.LocalityRequest{
			Data: models.Locality{
				Id:           123,
				LocalityName: "São Paulo & Co.",
				ProvinceName: "São Paulo - S.P.",
				CountryName:  "Brasil (BR)",
			},
		}

		expectedLocality := models.Locality{
			Id:           0,
			LocalityName: "São Paulo & Co.",
			ProvinceName: "São Paulo - S.P.",
			CountryName:  "Brasil (BR)",
		}

		// Act
		result := ToRequestToLocalityStruct(request)

		// Assert
		assert.Equal(t, expectedLocality, result)
		assert.Equal(t, "São Paulo & Co.", result.LocalityName)
		assert.Equal(t, "São Paulo - S.P.", result.ProvinceName)
		assert.Equal(t, "Brasil (BR)", result.CountryName)
	})

}
