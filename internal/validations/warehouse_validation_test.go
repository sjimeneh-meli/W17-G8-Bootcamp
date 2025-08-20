package validations

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/stretchr/testify/assert"
)

func TestWarehouseValidation_ValidateWarehouseRequestStruct(t *testing.T) {
	t.Run("success - valid warehouse request", func(t *testing.T) {
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
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("fail - missing address", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "", // Missing required field
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		// Act
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "address")
	})

	t.Run("fail - missing telephone", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "Test Address",
			Telephone:          "", // Missing required field
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		// Act
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "telephone")
	})

	t.Run("fail - missing warehouse code", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "Test Address",
			Telephone:          "123456789",
			WareHouseCode:      "", // Missing required field
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		// Act
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "warehouse_code")
	})

	t.Run("fail - missing minimum capacity", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "Test Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    0, // Missing required field (zero value)
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		// Act
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "minimum_capacity")
	})

	t.Run("fail - minimum capacity below minimum value", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "Test Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    0, // Below minimum value of 1
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		// Act
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "minimum_capacity")
	})

	t.Run("fail - minimum capacity negative", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "Test Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    -100, // Negative value
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		// Act
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "minimum_capacity")
	})

	t.Run("fail - missing minimum temperature", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "Test Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 0.0, // Missing required field (zero value)
			LocalityId:         1,
		}

		// Act
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "minimum_temperature")
	})

	t.Run("fail - missing locality ID", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "Test Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         0, // Missing required field (zero value)
		}

		// Act
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "locality_id")
	})

	t.Run("success - minimum valid values", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "A",
			Telephone:          "A",
			WareHouseCode:      "A",
			MinimumCapacity:    1, // Minimum allowed value
			MinimumTemperature: 0.1, // Small positive value
			LocalityId:         1,
		}

		// Act
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - large valid values", func(t *testing.T) {
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
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - negative temperature (valid)", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "Test Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: -50.0, // Negative temperature is valid
			LocalityId:         1,
		}

		// Act
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - decimal temperature", func(t *testing.T) {
		// Arrange
		request := requests.WarehouseRequest{
			Address:            "Test Address",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.5, // Decimal temperature
			LocalityId:         1,
		}

		// Act
		err := ValidateWarehouseRequestStruct(request)

		// Assert
		assert.NoError(t, err)
	})
}

func TestWarehouseValidation_ValidateWarehousePatchRequest(t *testing.T) {
	t.Run("success - valid patch request with all fields", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			Address:            stringPtr("Updated Address"),
			Telephone:          stringPtr("987654321"),
			WareHouseCode:      stringPtr("WH002"),
			MinimumCapacity:    intPtr(200),
			MinimumTemperature: float64Ptr(10.0),
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - valid patch request with partial fields", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			Address: stringPtr("Updated Address"),
			// Other fields are nil, so they should not be validated
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - valid patch request with no fields (all nil)", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			// All fields are nil
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("fail - address is empty string when provided", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			Address:            stringPtr(""), // Empty string when field is provided
			Telephone:          stringPtr("987654321"),
			WareHouseCode:      stringPtr("WH002"),
			MinimumCapacity:    intPtr(200),
			MinimumTemperature: float64Ptr(10.0),
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "address")
	})

	t.Run("fail - telephone is empty string when provided", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			Address:            stringPtr("Updated Address"),
			Telephone:          stringPtr(""), // Empty string when field is provided
			WareHouseCode:      stringPtr("WH002"),
			MinimumCapacity:    intPtr(200),
			MinimumTemperature: float64Ptr(10.0),
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "telephone")
	})

	t.Run("fail - warehouse code is empty string when provided", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			Address:            stringPtr("Updated Address"),
			Telephone:          stringPtr("987654321"),
			WareHouseCode:      stringPtr(""), // Empty string when field is provided
			MinimumCapacity:    intPtr(200),
			MinimumTemperature: float64Ptr(10.0),
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "warehouse_code")
	})

	t.Run("fail - minimum capacity is zero when provided", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			Address:            stringPtr("Updated Address"),
			Telephone:          stringPtr("987654321"),
			WareHouseCode:      stringPtr("WH002"),
			MinimumCapacity:    intPtr(0), // Zero when field is provided
			MinimumTemperature: float64Ptr(10.0),
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "minimum_capacity")
	})

	t.Run("fail - minimum capacity is negative when provided", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			Address:            stringPtr("Updated Address"),
			Telephone:          stringPtr("987654321"),
			WareHouseCode:      stringPtr("WH002"),
			MinimumCapacity:    intPtr(-100), // Negative when field is provided
			MinimumTemperature: float64Ptr(10.0),
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "minimum_capacity")
	})

	t.Run("fail - minimum temperature is zero when provided", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			Address:            stringPtr("Updated Address"),
			Telephone:          stringPtr("987654321"),
			WareHouseCode:      stringPtr("WH002"),
			MinimumCapacity:    intPtr(200),
			MinimumTemperature: float64Ptr(0.0), // Zero when field is provided
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "minimum_temperature")
	})

	t.Run("success - minimum capacity at boundary value (1)", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			Address:            stringPtr("Updated Address"),
			Telephone:          stringPtr("987654321"),
			WareHouseCode:      stringPtr("WH002"),
			MinimumCapacity:    intPtr(1), // Boundary value
			MinimumTemperature: float64Ptr(10.0),
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - negative temperature when provided", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			Address:            stringPtr("Updated Address"),
			Telephone:          stringPtr("987654321"),
			WareHouseCode:      stringPtr("WH002"),
			MinimumCapacity:    intPtr(200),
			MinimumTemperature: float64Ptr(-50.0), // Negative temperature is valid
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - decimal temperature when provided", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			Address:            stringPtr("Updated Address"),
			Telephone:          stringPtr("987654321"),
			WareHouseCode:      stringPtr("WH002"),
			MinimumCapacity:    intPtr(200),
			MinimumTemperature: float64Ptr(10.5), // Decimal temperature
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - large values when provided", func(t *testing.T) {
		// Arrange
		request := requests.WarehousePatchRequest{
			Address:            stringPtr("Very Long Address That Exceeds Normal Length Expectations"),
			Telephone:          stringPtr("12345678901234567890"),
			WareHouseCode:      stringPtr("WH999999999"),
			MinimumCapacity:    intPtr(999999999),
			MinimumTemperature: float64Ptr(999.99),
		}

		// Act
		err := ValidateWarehousePatchRequest(request)

		// Assert
		assert.NoError(t, err)
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
