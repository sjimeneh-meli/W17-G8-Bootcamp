package validations

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/stretchr/testify/assert"
)

func TestCarryValidation_ValidateCarryRequest(t *testing.T) {
	t.Run("success - valid carry request", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		// Act
		err := ValidateCarryRequest(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("fail - missing required fields", func(t *testing.T) {
		// Arrange - Multiple missing fields to cover all validation rules
		request := requests.CarryRequest{
			Cid:         "", // Missing required field
			CompanyName: "", // Missing required field
			Address:     "", // Missing required field
			Telephone:   "", // Missing required field
			LocalityId:  0,  // Missing required field (zero value)
		}

		// Act
		err := ValidateCarryRequest(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cid")
		assert.Contains(t, err.Error(), "company_name")
		assert.Contains(t, err.Error(), "address")
		assert.Contains(t, err.Error(), "telephone")
		assert.Contains(t, err.Error(), "locality_id")
	})

	t.Run("fail - CID length validation", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "CID1234567890", // Above maximum length of 10
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		// Act
		err := ValidateCarryRequest(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cid")
	})

	t.Run("fail - company name length validation", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "CID123",
			CompanyName: "Very Long Company Name That Exceeds The Maximum Allowed Length Of One Hundred Characters And Should Cause A Validation Error Because It Is Too Long",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		// Act
		err := ValidateCarryRequest(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "company_name")
	})

	t.Run("fail - telephone length validation", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "12345678901", // Above maximum length of 10
			LocalityId:  1,
		}

		// Act
		err := ValidateCarryRequest(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "telephone")
	})

	t.Run("fail - negative locality ID", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "CID123",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  -1, // Negative value
		}

		// Act
		err := ValidateCarryRequest(request)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "locality_id")
	})

	t.Run("success - minimum valid values", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "A", // Minimum length (1 character)
			CompanyName: "A", // Minimum length (1 character)
			Address:     "A", // Minimum length (1 character)
			Telephone:   "A", // Minimum length (1 character)
			LocalityId:  1,   // Minimum value
		}

		// Act
		err := ValidateCarryRequest(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - maximum valid values", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "1234567890", // Maximum length (10 characters)
			CompanyName: "Very Long Company Name That Exceeds Normal Length Expectations And Reaches The Maximum Allowed",
			Address:     "Very Long Address That Exceeds Normal Length Expectations And Reaches The Maximum Allowed",
			Telephone:   "1234567890", // Maximum length (10 characters)
			LocalityId:  999999999,    // Large positive value
		}

		// Act
		err := ValidateCarryRequest(request)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("success - special characters", func(t *testing.T) {
		// Arrange
		request := requests.CarryRequest{
			Cid:         "CID@123",
			CompanyName: "Company & Sons, Ltd.",
			Address:     "123 Main St., Apt. #4B",
			Telephone:   "+1-555-123",
			LocalityId:  1,
		}

		// Act
		err := ValidateCarryRequest(request)

		// Assert
		assert.NoError(t, err)
	})
}
