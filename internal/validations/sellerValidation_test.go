// Package validations - Seller Validation Unit Tests / Tests Unitarios de Validación Seller
// Request validation testing for seller data with required field checks
// Testing de validación de request para datos de vendedor con verificación de campos requeridos
package validations

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/stretchr/testify/assert"
)

// TestValidateSellerRequestStruct - Tests for ValidateSellerRequestStruct function
// Tests validation of seller request data structure
// Tests para función ValidateSellerRequestStruct
// Prueba la validación de estructura de datos de request de vendedor
func TestValidateSellerRequestStruct(t *testing.T) {
	t.Run("success - valid seller request with all required fields", func(t *testing.T) {
		// Test: Valid seller request data passes validation / Datos válidos de solicitud de vendedor pasan validación
		// Arrange
		request := requests.SellerRequest{
			CID:         "SEL-001",
			CompanyName: "Test Company Ltd",
			Address:     "123 Main Street, Downtown",
			Telephone:   "+1-555-123-4567",
			LocalityID:  1,
		}

		// Act
		err := ValidateSellerRequestStruct(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("fail - missing all required fields", func(t *testing.T) {
		// Test: Empty required fields fail validation / Campos requeridos vacíos fallan validación
		// Arrange
		request := requests.SellerRequest{
			CID:         "", // Missing required field
			CompanyName: "", // Missing required field
			Address:     "", // Missing required field
			Telephone:   "", // Missing required field
			LocalityID:  0,  // LocalityID is not validated as required
		}

		// Act
		err := ValidateSellerRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cid")
		assert.Contains(t, err.Error(), "company_name")
		assert.Contains(t, err.Error(), "address")
		assert.Contains(t, err.Error(), "telephone")
		// LocalityID should not be in error message as it's not required
		assert.NotContains(t, err.Error(), "locality_id")
	})

	t.Run("fail - missing CID only", func(t *testing.T) {
		// Test: Missing CID field fails validation / Campo CID faltante falla validación
		// Arrange
		request := requests.SellerRequest{
			CID:         "", // Missing required field
			CompanyName: "Test Company Ltd",
			Address:     "123 Main Street",
			Telephone:   "+1-555-123-4567",
			LocalityID:  1,
		}

		// Act
		err := ValidateSellerRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cid")
		assert.NotContains(t, err.Error(), "company_name")
		assert.NotContains(t, err.Error(), "address")
		assert.NotContains(t, err.Error(), "telephone")
	})

	t.Run("fail - missing CompanyName only", func(t *testing.T) {
		// Test: Missing CompanyName field fails validation / Campo CompanyName faltante falla validación
		// Arrange
		request := requests.SellerRequest{
			CID:         "SEL-002",
			CompanyName: "", // Missing required field
			Address:     "123 Main Street",
			Telephone:   "+1-555-123-4567",
			LocalityID:  1,
		}

		// Act
		err := ValidateSellerRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "company_name")
		assert.NotContains(t, err.Error(), "cid")
		assert.NotContains(t, err.Error(), "address")
		assert.NotContains(t, err.Error(), "telephone")
	})

	t.Run("fail - missing Address only", func(t *testing.T) {
		// Test: Missing Address field fails validation / Campo Address faltante falla validación
		// Arrange
		request := requests.SellerRequest{
			CID:         "SEL-003",
			CompanyName: "Test Company Ltd",
			Address:     "", // Missing required field
			Telephone:   "+1-555-123-4567",
			LocalityID:  1,
		}

		// Act
		err := ValidateSellerRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "address")
		assert.NotContains(t, err.Error(), "cid")
		assert.NotContains(t, err.Error(), "company_name")
		assert.NotContains(t, err.Error(), "telephone")
	})

	t.Run("fail - missing Telephone only", func(t *testing.T) {
		// Test: Missing Telephone field fails validation / Campo Telephone faltante falla validación
		// Arrange
		request := requests.SellerRequest{
			CID:         "SEL-004",
			CompanyName: "Test Company Ltd",
			Address:     "123 Main Street",
			Telephone:   "", // Missing required field
			LocalityID:  1,
		}

		// Act
		err := ValidateSellerRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "telephone")
		assert.NotContains(t, err.Error(), "cid")
		assert.NotContains(t, err.Error(), "company_name")
		assert.NotContains(t, err.Error(), "address")
	})

	t.Run("success - LocalityID is zero (not required)", func(t *testing.T) {
		// Test: Zero LocalityID passes validation as it's not required / LocalityID cero pasa validación ya que no es requerido
		// Arrange
		request := requests.SellerRequest{
			CID:         "SEL-005",
			CompanyName: "Test Company Ltd",
			Address:     "123 Main Street",
			Telephone:   "+1-555-123-4567",
			LocalityID:  0, // Zero value should be valid
		}

		// Act
		err := ValidateSellerRequestStruct(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - negative LocalityID (not validated)", func(t *testing.T) {
		// Test: Negative LocalityID passes validation as field is not validated / LocalityID negativo pasa validación ya que el campo no es validado
		// Arrange
		request := requests.SellerRequest{
			CID:         "SEL-006",
			CompanyName: "Test Company Ltd",
			Address:     "123 Main Street",
			Telephone:   "+1-555-123-4567",
			LocalityID:  -1, // Negative value should be valid (no validation on this field)
		}

		// Act
		err := ValidateSellerRequestStruct(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - minimum valid values", func(t *testing.T) {
		// Test: Single character values pass validation / Valores de un solo carácter pasan validación
		// Arrange
		request := requests.SellerRequest{
			CID:         "A", // Minimum length (1 character)
			CompanyName: "B", // Minimum length (1 character)
			Address:     "C", // Minimum length (1 character)
			Telephone:   "D", // Minimum length (1 character)
			LocalityID:  1,
		}

		// Act
		err := ValidateSellerRequestStruct(request)

		// Assert
		assert.NoError(t, err)
	})

}
