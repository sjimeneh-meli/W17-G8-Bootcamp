// Package mappers - Seller Mapper Unit Tests / Tests Unitarios del Mapper Seller
// Data transformation testing for seller mapping between request, model, and response structures
// Testing de transformación de datos para mapeo de vendedor entre estructuras de request, modelo y response
package mappers

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestToRequestToSellerStruct - Tests for ToRequestToSellerStruct function
// Tests data mapping from SellerRequest to Seller model
// Tests para función ToRequestToSellerStruct
// Prueba el mapeo de datos de SellerRequest a modelo Seller
func TestToRequestToSellerStruct(t *testing.T) {
	t.Run("success - maps seller request to model correctly", func(t *testing.T) {
		// Test: Valid seller request data maps successfully to model / Datos válidos de solicitud de vendedor se mapean exitosamente al modelo
		// Arrange
		request := requests.SellerRequest{
			CID:         "SEL-001",
			CompanyName: "Test Company Ltd",
			Address:     "123 Main Street, Downtown",
			Telephone:   "+1-555-123-4567",
			LocalityID:  1,
		}

		expectedSeller := models.Seller{
			Id:          0, // ID should not be mapped from request / ID no debería ser mapeado del request
			CID:         "SEL-001",
			CompanyName: "Test Company Ltd",
			Address:     "123 Main Street, Downtown",
			Telephone:   "+1-555-123-4567",
			LocalityID:  1,
		}

		// Act
		result := ToRequestToSellerStruct(request)

		// Assert
		assert.Equal(t, expectedSeller, result)
		assert.Equal(t, "SEL-001", result.CID)
		assert.Equal(t, "Test Company Ltd", result.CompanyName)
		assert.Equal(t, "123 Main Street, Downtown", result.Address)
		assert.Equal(t, "+1-555-123-4567", result.Telephone)
		assert.Equal(t, 1, result.LocalityID)
		assert.Equal(t, 0, result.Id) // ID should be zero value / ID debería ser valor cero
	})

	t.Run("success - maps seller request with empty values", func(t *testing.T) {
		// Test: Empty string values in request map correctly / Valores de string vacíos en request se mapean correctamente
		// Arrange
		request := requests.SellerRequest{
			CID:         "",
			CompanyName: "",
			Address:     "",
			Telephone:   "",
			LocalityID:  0,
		}

		expectedSeller := models.Seller{
			Id:          0,
			CID:         "",
			CompanyName: "",
			Address:     "",
			Telephone:   "",
			LocalityID:  0,
		}

		// Act
		result := ToRequestToSellerStruct(request)

		// Assert
		assert.Equal(t, expectedSeller, result)
		assert.Equal(t, "", result.CID)
		assert.Equal(t, "", result.CompanyName)
		assert.Equal(t, "", result.Address)
		assert.Equal(t, "", result.Telephone)
		assert.Equal(t, 0, result.LocalityID)
	})

	t.Run("success - maps seller request with long values", func(t *testing.T) {
		// Test: Very long string values map correctly / Valores de string muy largos se mapean correctamente
		// Arrange
		longCID := "VERY-LONG-CID-THAT-EXCEEDS-NORMAL-EXPECTATIONS"
		longCompanyName := "Very Long Company Name That Exceeds Normal Length Expectations And Reaches The Maximum Allowed Length"
		longAddress := "Very Long Address That Exceeds Normal Length Expectations And Contains Multiple Street Names And Numbers"
		longTelephone := "12345678901234567890123456789012345678901234567890"

		request := requests.SellerRequest{
			CID:         longCID,
			CompanyName: longCompanyName,
			Address:     longAddress,
			Telephone:   longTelephone,
			LocalityID:  999999999,
		}

		expectedSeller := models.Seller{
			Id:          0,
			CID:         longCID,
			CompanyName: longCompanyName,
			Address:     longAddress,
			Telephone:   longTelephone,
			LocalityID:  999999999,
		}

		// Act
		result := ToRequestToSellerStruct(request)

		// Assert
		assert.Equal(t, expectedSeller, result)
		assert.Equal(t, longCID, result.CID)
		assert.Equal(t, longCompanyName, result.CompanyName)
		assert.Equal(t, longAddress, result.Address)
		assert.Equal(t, longTelephone, result.Telephone)
		assert.Equal(t, 999999999, result.LocalityID)
	})
}

// TestToSellerStructToResponse - Tests for ToSellerStructToResponse function
// Tests data mapping from Seller model to SellerResponse
// Tests para función ToSellerStructToResponse
// Prueba el mapeo de datos de modelo Seller a SellerResponse
func TestToSellerStructToResponse(t *testing.T) {
	t.Run("success - maps seller model to response correctly", func(t *testing.T) {
		// Test: Valid seller model data maps successfully to response / Datos válidos de modelo de vendedor se mapean exitosamente a response
		// Arrange
		seller := models.Seller{
			Id:          123,
			CID:         "SEL-123",
			CompanyName: "Response Test Company",
			Address:     "456 Response Street",
			Telephone:   "+1-555-987-6543",
			LocalityID:  2,
		}

		expectedResponse := responses.SellerResponse{
			Id:          123,
			CID:         "SEL-123",
			CompanyName: "Response Test Company",
			Address:     "456 Response Street",
			Telephone:   "+1-555-987-6543",
			LocalityID:  2,
		}

		// Act
		result := ToSellerStructToResponse(seller)

		// Assert
		assert.Equal(t, expectedResponse, result)
		assert.Equal(t, 123, result.Id)
		assert.Equal(t, "SEL-123", result.CID)
		assert.Equal(t, "Response Test Company", result.CompanyName)
		assert.Equal(t, "456 Response Street", result.Address)
		assert.Equal(t, "+1-555-987-6543", result.Telephone)
		assert.Equal(t, 2, result.LocalityID)
	})

	t.Run("success - maps seller model with zero values", func(t *testing.T) {
		// Test: Zero and empty values in model map correctly / Valores cero y vacíos en modelo se mapean correctamente
		// Arrange
		seller := models.Seller{
			Id:          0,
			CID:         "",
			CompanyName: "",
			Address:     "",
			Telephone:   "",
			LocalityID:  0,
		}

		expectedResponse := responses.SellerResponse{
			Id:          0,
			CID:         "",
			CompanyName: "",
			Address:     "",
			Telephone:   "",
			LocalityID:  0,
		}

		// Act
		result := ToSellerStructToResponse(seller)

		// Assert
		assert.Equal(t, expectedResponse, result)
		assert.Equal(t, 0, result.Id)
		assert.Equal(t, "", result.CID)
		assert.Equal(t, "", result.CompanyName)
		assert.Equal(t, "", result.Address)
		assert.Equal(t, "", result.Telephone)
		assert.Equal(t, 0, result.LocalityID)
	})

	t.Run("success - maps seller model with special characters", func(t *testing.T) {
		// Test: Special characters and symbols in model map correctly / Caracteres especiales y símbolos en modelo se mapean correctamente
		// Arrange
		seller := models.Seller{
			Id:          456,
			CID:         "SPEC!@L-CH@RS",
			CompanyName: "Special & Characters Co., Inc.",
			Address:     "123 Special St., Unit #456",
			Telephone:   "+1 (555) 123-4567 x890",
			LocalityID:  789,
		}

		expectedResponse := responses.SellerResponse{
			Id:          456,
			CID:         "SPEC!@L-CH@RS",
			CompanyName: "Special & Characters Co., Inc.",
			Address:     "123 Special St., Unit #456",
			Telephone:   "+1 (555) 123-4567 x890",
			LocalityID:  789,
		}

		// Act
		result := ToSellerStructToResponse(seller)

		// Assert
		assert.Equal(t, expectedResponse, result)
		assert.Equal(t, 456, result.Id)
		assert.Equal(t, "SPEC!@L-CH@RS", result.CID)
		assert.Equal(t, "Special & Characters Co., Inc.", result.CompanyName)
		assert.Equal(t, "123 Special St., Unit #456", result.Address)
		assert.Equal(t, "+1 (555) 123-4567 x890", result.Telephone)
		assert.Equal(t, 789, result.LocalityID)
	})

}
