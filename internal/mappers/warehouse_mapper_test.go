package mappers

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWarehouseMapper_ToResponse(t *testing.T) {
	t.Run("success - maps warehouse model to response correctly", func(t *testing.T) {
		// Arrange
		warehouse := models.Warehouse{
			Id:                 1,
			Address:            "Test Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		// Act
		result := ToResponse(warehouse)

		// Assert
		assert.Equal(t, warehouse.Id, result.ID)
		assert.Equal(t, warehouse.Address, result.Address)
		assert.Equal(t, warehouse.Telephone, result.Telephone)
		assert.Equal(t, warehouse.WareHouseCode, result.WareHouseCode)
		assert.Equal(t, warehouse.MinimumCapacity, result.MinimumCapacity)
		assert.Equal(t, warehouse.MinimumTemperature, result.MinimumTemperature)
	})

	t.Run("success - maps warehouse with zero values", func(t *testing.T) {
		// Arrange
		warehouse := models.Warehouse{
			Id:                 0,
			Address:            "",
			Telephone:          "",
			WareHouseCode:      "",
			MinimumCapacity:    0,
			MinimumTemperature: 0.0,
			LocalityId:         0,
		}

		// Act
		result := ToResponse(warehouse)

		// Assert
		assert.Equal(t, 0, result.ID)
		assert.Equal(t, "", result.Address)
		assert.Equal(t, "", result.Telephone)
		assert.Equal(t, "", result.WareHouseCode)
		assert.Equal(t, 0, result.MinimumCapacity)
		assert.Equal(t, 0.0, result.MinimumTemperature)
	})

	t.Run("success - maps warehouse with negative values", func(t *testing.T) {
		// Arrange
		warehouse := models.Warehouse{
			Id:                 -1,
			Address:            "Negative Address",
			Telephone:          "-123456789",
			WareHouseCode:      "WH-001",
			MinimumCapacity:    -100,
			MinimumTemperature: -5.0,
			LocalityId:         -1,
		}

		// Act
		result := ToResponse(warehouse)

		// Assert
		assert.Equal(t, -1, result.ID)
		assert.Equal(t, "Negative Address", result.Address)
		assert.Equal(t, "-123456789", result.Telephone)
		assert.Equal(t, "WH-001", result.WareHouseCode)
		assert.Equal(t, -100, result.MinimumCapacity)
		assert.Equal(t, -5.0, result.MinimumTemperature)
	})

	t.Run("success - maps warehouse with large values", func(t *testing.T) {
		// Arrange
		warehouse := models.Warehouse{
			Id:                 999999999,
			Address:            "Very Long Address That Exceeds Normal Length Expectations",
			Telephone:          "12345678901234567890",
			WareHouseCode:      "WH999999999",
			MinimumCapacity:    999999999,
			MinimumTemperature: 999.99,
			LocalityId:         999999999,
		}

		// Act
		result := ToResponse(warehouse)

		// Assert
		assert.Equal(t, 999999999, result.ID)
		assert.Equal(t, "Very Long Address That Exceeds Normal Length Expectations", result.Address)
		assert.Equal(t, "12345678901234567890", result.Telephone)
		assert.Equal(t, "WH999999999", result.WareHouseCode)
		assert.Equal(t, 999999999, result.MinimumCapacity)
		assert.Equal(t, 999.99, result.MinimumTemperature)
	})
}

func TestWarehouseMapper_ToRequest(t *testing.T) {
	t.Run("success - maps warehouse request to model correctly", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "Test Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		// Act
		result := ToRequest(request)

		// Assert
		assert.Equal(t, request.Address, result.Address)
		assert.Equal(t, request.Telephone, result.Telephone)
		assert.Equal(t, request.WareHouseCode, result.WareHouseCode)
		assert.Equal(t, request.MinimumCapacity, result.MinimumCapacity)
		assert.Equal(t, request.MinimumTemperature, result.MinimumTemperature)
		assert.Equal(t, request.LocalityId, result.LocalityId)
		// ID should not be set from request
		assert.Equal(t, 0, result.Id)
	})

	t.Run("success - maps warehouse request with zero values", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "",
			Telephone:          "",
			WareHouseCode:      "",
			MinimumCapacity:    0,
			MinimumTemperature: 0.0,
			LocalityId:         0,
		}

		// Act
		result := ToRequest(request)

		// Assert
		assert.Equal(t, "", result.Address)
		assert.Equal(t, "", result.Telephone)
		assert.Equal(t, "", result.WareHouseCode)
		assert.Equal(t, 0, result.MinimumCapacity)
		assert.Equal(t, 0.0, result.MinimumTemperature)
		assert.Equal(t, 0, result.LocalityId)
		assert.Equal(t, 0, result.Id)
	})

	t.Run("success - maps warehouse request with negative values", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "Negative Address",
			Telephone:          "-123456789",
			WareHouseCode:      "WH-001",
			MinimumCapacity:    -100,
			MinimumTemperature: -5.0,
			LocalityId:         -1,
		}

		// Act
		result := ToRequest(request)

		// Assert
		assert.Equal(t, "Negative Address", result.Address)
		assert.Equal(t, "-123456789", result.Telephone)
		assert.Equal(t, "WH-001", result.WareHouseCode)
		assert.Equal(t, -100, result.MinimumCapacity)
		assert.Equal(t, -5.0, result.MinimumTemperature)
		assert.Equal(t, -1, result.LocalityId)
		assert.Equal(t, 0, result.Id)
	})

	t.Run("success - maps warehouse request with large values", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "Very Long Address That Exceeds Normal Length Expectations",
			Telephone:          "12345678901234567890",
			WareHouseCode:      "WH999999999",
			MinimumCapacity:    999999999,
			MinimumTemperature: 999.99,
			LocalityId:         999999999,
		}

		// Act
		result := ToRequest(request)

		// Assert
		assert.Equal(t, "Very Long Address That Exceeds Normal Length Expectations", result.Address)
		assert.Equal(t, "12345678901234567890", result.Telephone)
		assert.Equal(t, "WH999999999", result.WareHouseCode)
		assert.Equal(t, 999999999, result.MinimumCapacity)
		assert.Equal(t, 999.99, result.MinimumTemperature)
		assert.Equal(t, 999999999, result.LocalityId)
		assert.Equal(t, 0, result.Id)
	})
}

func TestWarehouseMapper_ApplyPatch(t *testing.T) {
	t.Run("success - applies all patch fields", func(t *testing.T) {
		// Arrange
		existing := models.Warehouse{
			Id:                 1,
			Address:            "Original Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		patch := requests.WarehousePatchRequest{
			Address:            stringPtr("Updated Address"),
			Telephone:          stringPtr("987654321"),
			WareHouseCode:      stringPtr("WH002"),
			MinimumCapacity:    intPtr(200),
			MinimumTemperature: float64Ptr(10.0),
		}

		// Act
		result := ApplyPatch(existing, patch)

		// Assert
		assert.Equal(t, 1, result.Id) // ID should remain unchanged
		assert.Equal(t, "Updated Address", result.Address)
		assert.Equal(t, "987654321", result.Telephone)
		assert.Equal(t, "WH002", result.WareHouseCode)
		assert.Equal(t, 200, result.MinimumCapacity)
		assert.Equal(t, 10.0, result.MinimumTemperature)
		assert.Equal(t, 1, result.LocalityId) // LocalityId should remain unchanged
	})

	t.Run("success - applies partial patch fields", func(t *testing.T) {
		// Arrange
		existing := models.Warehouse{
			Id:                 1,
			Address:            "Original Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		patch := requests.WarehousePatchRequest{
			Address: stringPtr("Updated Address"),
			// Other fields are nil, so they should remain unchanged
		}

		// Act
		result := ApplyPatch(existing, patch)

		// Assert
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, "Updated Address", result.Address) // Updated
		assert.Equal(t, "123456789", result.Telephone)     // Unchanged
		assert.Equal(t, "WH001", result.WareHouseCode)     // Unchanged
		assert.Equal(t, 100, result.MinimumCapacity)       // Unchanged
		assert.Equal(t, 5.0, result.MinimumTemperature)    // Unchanged
		assert.Equal(t, 1, result.LocalityId)
	})

	t.Run("success - applies no patch fields (all nil)", func(t *testing.T) {
		// Arrange
		existing := models.Warehouse{
			Id:                 1,
			Address:            "Original Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		patch := requests.WarehousePatchRequest{
			// All fields are nil
		}

		// Act
		result := ApplyPatch(existing, patch)

		// Assert - nothing should change
		assert.Equal(t, existing, result)
	})

	t.Run("success - applies patch with zero values", func(t *testing.T) {
		// Arrange
		existing := models.Warehouse{
			Id:                 1,
			Address:            "Original Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		patch := requests.WarehousePatchRequest{
			Address:            stringPtr(""),
			Telephone:          stringPtr(""),
			WareHouseCode:      stringPtr(""),
			MinimumCapacity:    intPtr(0),
			MinimumTemperature: float64Ptr(0.0),
		}

		// Act
		result := ApplyPatch(existing, patch)

		// Assert
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, "", result.Address)
		assert.Equal(t, "", result.Telephone)
		assert.Equal(t, "", result.WareHouseCode)
		assert.Equal(t, 0, result.MinimumCapacity)
		assert.Equal(t, 0.0, result.MinimumTemperature)
		assert.Equal(t, 1, result.LocalityId)
	})

	t.Run("success - applies patch with negative values", func(t *testing.T) {
		// Arrange
		existing := models.Warehouse{
			Id:                 1,
			Address:            "Original Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		patch := requests.WarehousePatchRequest{
			Address:            stringPtr("Negative Address"),
			Telephone:          stringPtr("-123456789"),
			WareHouseCode:      stringPtr("WH-001"),
			MinimumCapacity:    intPtr(-100),
			MinimumTemperature: float64Ptr(-5.0),
		}

		// Act
		result := ApplyPatch(existing, patch)

		// Assert
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, "Negative Address", result.Address)
		assert.Equal(t, "-123456789", result.Telephone)
		assert.Equal(t, "WH-001", result.WareHouseCode)
		assert.Equal(t, -100, result.MinimumCapacity)
		assert.Equal(t, -5.0, result.MinimumTemperature)
		assert.Equal(t, 1, result.LocalityId)
	})
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
