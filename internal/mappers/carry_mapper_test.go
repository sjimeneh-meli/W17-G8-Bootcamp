package mappers

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCarryMapper_MapCarryToCreateCarryResponse(t *testing.T) {
	t.Run("success - maps carry model to response correctly", func(t *testing.T) {
		// Arrange
		carry := models.Carry{
			Id:          1,
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		// Act
		result := MapCarryToCreateCarryResponse(carry)

		// Assert
		assert.Equal(t, carry.Id, result.Id)
		assert.Equal(t, carry.Cid, result.Cid)
		assert.Equal(t, carry.CompanyName, result.CompanyName)
		assert.Equal(t, carry.Address, result.Address)
		assert.Equal(t, carry.Telephone, result.Telephone)
		assert.Equal(t, carry.LocalityId, result.LocalityId)
	})

	t.Run("success - maps carry with zero values", func(t *testing.T) {
		// Arrange
		carry := models.Carry{
			Id:          0,
			Cid:         "",
			CompanyName: "",
			Address:     "",
			Telephone:   "",
			LocalityId:  0,
		}

		// Act
		result := MapCarryToCreateCarryResponse(carry)

		// Assert
		assert.Equal(t, 0, result.Id)
		assert.Equal(t, "", result.Cid)
		assert.Equal(t, "", result.CompanyName)
		assert.Equal(t, "", result.Address)
		assert.Equal(t, "", result.Telephone)
		assert.Equal(t, 0, result.LocalityId)
	})

	t.Run("success - maps carry with negative values", func(t *testing.T) {
		// Arrange
		carry := models.Carry{
			Id:          -1,
			Cid:         "CID-123",
			CompanyName: "Negative Company",
			Address:     "Negative Address",
			Telephone:   "-123456789",
			LocalityId:  -1,
		}

		// Act
		result := MapCarryToCreateCarryResponse(carry)

		// Assert
		assert.Equal(t, -1, result.Id)
		assert.Equal(t, "CID-123", result.Cid)
		assert.Equal(t, "Negative Company", result.CompanyName)
		assert.Equal(t, "Negative Address", result.Address)
		assert.Equal(t, "-123456789", result.Telephone)
		assert.Equal(t, -1, result.LocalityId)
	})

	t.Run("success - maps carry with large values", func(t *testing.T) {
		// Arrange
		carry := models.Carry{
			Id:          999999999,
			Cid:         "CID999999999",
			CompanyName: "Very Long Company Name That Exceeds Normal Length Expectations",
			Address:     "Very Long Address That Exceeds Normal Length Expectations",
			Telephone:   "12345678901234567890",
			LocalityId:  999999999,
		}

		// Act
		result := MapCarryToCreateCarryResponse(carry)

		// Assert
		assert.Equal(t, 999999999, result.Id)
		assert.Equal(t, "CID999999999", result.Cid)
		assert.Equal(t, "Very Long Company Name That Exceeds Normal Length Expectations", result.CompanyName)
		assert.Equal(t, "Very Long Address That Exceeds Normal Length Expectations", result.Address)
		assert.Equal(t, "12345678901234567890", result.Telephone)
		assert.Equal(t, 999999999, result.LocalityId)
	})

	t.Run("success - maps carry with special characters", func(t *testing.T) {
		// Arrange
		carry := models.Carry{
			Id:          1,
			Cid:         "CID@#$%",
			CompanyName: "Company & Sons, Ltd.",
			Address:     "123 Main St., Apt. #4B",
			Telephone:   "+1-555-123-4567",
			LocalityId:  1,
		}

		// Act
		result := MapCarryToCreateCarryResponse(carry)

		// Assert
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, "CID@#$%", result.Cid)
		assert.Equal(t, "Company & Sons, Ltd.", result.CompanyName)
		assert.Equal(t, "123 Main St., Apt. #4B", result.Address)
		assert.Equal(t, "+1-555-123-4567", result.Telephone)
		assert.Equal(t, 1, result.LocalityId)
	})
}

func TestCarryMapper_MapCarryRequestToCarry(t *testing.T) {
	t.Run("success - maps carry request to model correctly", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		// Act
		result := MapCarryRequestToCarry(request)

		// Assert
		assert.Equal(t, request.Cid, result.Cid)
		assert.Equal(t, request.CompanyName, result.CompanyName)
		assert.Equal(t, request.Address, result.Address)
		assert.Equal(t, request.Telephone, result.Telephone)
		assert.Equal(t, request.LocalityId, result.LocalityId)
		// ID should not be set from request
		assert.Equal(t, 0, result.Id)
	})

	t.Run("success - maps carry request with zero values", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "",
			CompanyName: "",
			Address:     "",
			Telephone:   "",
			LocalityId:  0,
		}

		// Act
		result := MapCarryRequestToCarry(request)

		// Assert
		assert.Equal(t, "", result.Cid)
		assert.Equal(t, "", result.CompanyName)
		assert.Equal(t, "", result.Address)
		assert.Equal(t, "", result.Telephone)
		assert.Equal(t, 0, result.LocalityId)
		assert.Equal(t, 0, result.Id)
	})

	t.Run("success - maps carry request with negative values", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "CID-123",
			CompanyName: "Negative Company",
			Address:     "Negative Address",
			Telephone:   "-123456789",
			LocalityId:  -1,
		}

		// Act
		result := MapCarryRequestToCarry(request)

		// Assert
		assert.Equal(t, "CID-123", result.Cid)
		assert.Equal(t, "Negative Company", result.CompanyName)
		assert.Equal(t, "Negative Address", result.Address)
		assert.Equal(t, "-123456789", result.Telephone)
		assert.Equal(t, -1, result.LocalityId)
		assert.Equal(t, 0, result.Id)
	})

	t.Run("success - maps carry request with large values", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "CID999999999",
			CompanyName: "Very Long Company Name That Exceeds Normal Length Expectations",
			Address:     "Very Long Address That Exceeds Normal Length Expectations",
			Telephone:   "12345678901234567890",
			LocalityId:  999999999,
		}

		// Act
		result := MapCarryRequestToCarry(request)

		// Assert
		assert.Equal(t, "CID999999999", result.Cid)
		assert.Equal(t, "Very Long Company Name That Exceeds Normal Length Expectations", result.CompanyName)
		assert.Equal(t, "Very Long Address That Exceeds Normal Length Expectations", result.Address)
		assert.Equal(t, "12345678901234567890", result.Telephone)
		assert.Equal(t, 999999999, result.LocalityId)
		assert.Equal(t, 0, result.Id)
	})

	t.Run("success - maps carry request with special characters", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "CID@#$%",
			CompanyName: "Company & Sons, Ltd.",
			Address:     "123 Main St., Apt. #4B",
			Telephone:   "+1-555-123-4567",
			LocalityId:  1,
		}

		// Act
		result := MapCarryRequestToCarry(request)

		// Assert
		assert.Equal(t, "CID@#$%", result.Cid)
		assert.Equal(t, "Company & Sons, Ltd.", result.CompanyName)
		assert.Equal(t, "123 Main St., Apt. #4B", result.Address)
		assert.Equal(t, "+1-555-123-4567", result.Telephone)
		assert.Equal(t, 1, result.LocalityId)
		assert.Equal(t, 0, result.Id)
	})

	t.Run("success - maps carry request with minimum valid values", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "A",
			CompanyName: "A",
			Address:     "A",
			Telephone:   "A",
			LocalityId:  1,
		}

		// Act
		result := MapCarryRequestToCarry(request)

		// Assert
		assert.Equal(t, "A", result.Cid)
		assert.Equal(t, "A", result.CompanyName)
		assert.Equal(t, "A", result.Address)
		assert.Equal(t, "A", result.Telephone)
		assert.Equal(t, 1, result.LocalityId)
		assert.Equal(t, 0, result.Id)
	})

	t.Run("success - maps carry request with maximum valid values", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "1234567890", // 10 characters (max length)
			CompanyName: "Very Long Company Name That Exceeds Normal Length Expectations And Reaches The Maximum Allowed Length For This Field",
			Address:     "Very Long Address That Exceeds Normal Length Expectations And Reaches The Maximum Allowed Length For This Field",
			Telephone:   "1234567890", // 10 characters (max length)
			LocalityId:  999999999,
		}

		// Act
		result := MapCarryRequestToCarry(request)

		// Assert
		assert.Equal(t, "1234567890", result.Cid)
		assert.Equal(t, "Very Long Company Name That Exceeds Normal Length Expectations And Reaches The Maximum Allowed Length For This Field", result.CompanyName)
		assert.Equal(t, "Very Long Address That Exceeds Normal Length Expectations And Reaches The Maximum Allowed Length For This Field", result.Address)
		assert.Equal(t, "1234567890", result.Telephone)
		assert.Equal(t, 999999999, result.LocalityId)
		assert.Equal(t, 0, result.Id)
	})
}
