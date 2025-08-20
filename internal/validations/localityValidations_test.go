// Package validations - Locality Validation Unit Tests / Tests Unitarios de Validación Locality
// Request validation testing for locality data with required field checks
// Testing de validación de request para datos de localidad con verificación de campos requeridos
package validations

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestValidateLocalityRequestStruct - Tests for ValidateLocalityRequestStruct function
// Tests validation of locality request data structure
// Tests para función ValidateLocalityRequestStruct
// Prueba la validación de estructura de datos de request de localidad
func TestValidateLocalityRequestStruct(t *testing.T) {
	t.Run("success - valid locality request with all required fields", func(t *testing.T) {
		// Test: Valid locality request data passes validation / Datos válidos de solicitud de localidad pasan validación
		// Arrange
		locality := models.Locality{
			Id:           1, // ID is not validated, can be any value
			LocalityName: "Buenos Aires",
			ProvinceName: "Ciudad Autónoma de Buenos Aires",
			CountryName:  "Argentina",
		}

		// Act
		err := ValidateLocalityRequestStruct(locality)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("fail - missing all required fields", func(t *testing.T) {
		// Test: Empty required fields fail validation / Campos requeridos vacíos fallan validación
		// Arrange
		locality := models.Locality{
			Id:           0,  // ID is not validated
			LocalityName: "", // Missing required field
			ProvinceName: "", // Missing required field
			CountryName:  "", // Missing required field
		}

		// Act
		err := ValidateLocalityRequestStruct(locality)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "locality_name")
		assert.Contains(t, err.Error(), "province_name")
		assert.Contains(t, err.Error(), "country_name")
		// ID should not be in error message as it's not required
		assert.NotContains(t, err.Error(), "id")
	})

	t.Run("fail - missing LocalityName only", func(t *testing.T) {
		// Test: Missing LocalityName field fails validation / Campo LocalityName faltante falla validación
		// Arrange
		locality := models.Locality{
			Id:           999, // ID can be any value
			LocalityName: "",  // Missing required field
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		// Act
		err := ValidateLocalityRequestStruct(locality)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "locality_name")
		assert.NotContains(t, err.Error(), "province_name")
		assert.NotContains(t, err.Error(), "country_name")
		assert.NotContains(t, err.Error(), "id")
	})

	t.Run("fail - missing ProvinceName only", func(t *testing.T) {
		// Test: Missing ProvinceName field fails validation / Campo ProvinceName faltante falla validación
		// Arrange
		locality := models.Locality{
			Id:           -1, // ID can be negative
			LocalityName: "Test City",
			ProvinceName: "", // Missing required field
			CountryName:  "Test Country",
		}

		// Act
		err := ValidateLocalityRequestStruct(locality)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "province_name")
		assert.NotContains(t, err.Error(), "locality_name")
		assert.NotContains(t, err.Error(), "country_name")
		assert.NotContains(t, err.Error(), "id")
	})

	t.Run("fail - missing CountryName only", func(t *testing.T) {
		// Test: Missing CountryName field fails validation / Campo CountryName faltante falla validación
		// Arrange
		locality := models.Locality{
			Id:           2147483647, // ID can be max int
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "", // Missing required field
		}

		// Act
		err := ValidateLocalityRequestStruct(locality)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "country_name")
		assert.NotContains(t, err.Error(), "locality_name")
		assert.NotContains(t, err.Error(), "province_name")
		assert.NotContains(t, err.Error(), "id")
	})

	t.Run("success - ID is zero (not required)", func(t *testing.T) {
		// Test: Zero ID passes validation as it's not required / ID cero pasa validación ya que no es requerido
		// Arrange
		locality := models.Locality{
			Id:           0, // Zero value should be valid
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		// Act
		err := ValidateLocalityRequestStruct(locality)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - negative ID (not validated)", func(t *testing.T) {
		// Test: Negative ID passes validation as field is not validated / ID negativo pasa validación ya que el campo no es validado
		// Arrange
		locality := models.Locality{
			Id:           -999, // Negative value should be valid (no validation on this field)
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		// Act
		err := ValidateLocalityRequestStruct(locality)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - minimum valid values", func(t *testing.T) {
		// Test: Single character values pass validation / Valores de un solo carácter pasan validación
		// Arrange
		locality := models.Locality{
			Id:           1,
			LocalityName: "A", // Minimum length (1 character)
			ProvinceName: "B", // Minimum length (1 character)
			CountryName:  "C", // Minimum length (1 character)
		}

		// Act
		err := ValidateLocalityRequestStruct(locality)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - real world examples", func(t *testing.T) {
		// Test: Real world locality examples pass validation / Ejemplos de localidades del mundo real pasan validación
		// Arrange
		localities := []models.Locality{
			{
				Id:           1,
				LocalityName: "New York City",
				ProvinceName: "New York",
				CountryName:  "United States",
			},
			{
				Id:           2,
				LocalityName: "London",
				ProvinceName: "England",
				CountryName:  "United Kingdom",
			},
			{
				Id:           3,
				LocalityName: "Tokyo",
				ProvinceName: "Tokyo Prefecture",
				CountryName:  "Japan",
			},
			{
				Id:           4,
				LocalityName: "São Paulo",
				ProvinceName: "São Paulo",
				CountryName:  "Brazil",
			},
		}

		// Act & Assert
		for _, locality := range localities {
			err := ValidateLocalityRequestStruct(locality)
			assert.NoError(t, err, "Validation should pass for locality: %+v", locality)
		}
	})

	t.Run("success - mixed case and accents", func(t *testing.T) {
		// Test: Mixed case and accented characters pass validation / Mayúsculas/minúsculas mezcladas y acentos pasan validación
		// Arrange
		locality := models.Locality{
			Id:           100,
			LocalityName: "São PaUlO",
			ProvinceName: "SãO pAuLo",
			CountryName:  "bRaSiL",
		}

		// Act
		err := ValidateLocalityRequestStruct(locality)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("fail - partial missing fields", func(t *testing.T) {
		// Test: Missing two out of three required fields / Faltando dos de tres campos requeridos
		// Arrange
		locality := models.Locality{
			Id:           200,
			LocalityName: "Only City Name", // Only this field is present
			ProvinceName: "",               // Missing
			CountryName:  "",               // Missing
		}

		// Act
		err := ValidateLocalityRequestStruct(locality)

		// Assert
		assert.Error(t, err)
		assert.NotContains(t, err.Error(), "locality_name") // This field is present
		assert.Contains(t, err.Error(), "province_name")    // Missing
		assert.Contains(t, err.Error(), "country_name")     // Missing
	})
}
