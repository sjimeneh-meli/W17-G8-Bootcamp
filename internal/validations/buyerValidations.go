package validations

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

// ValidateBuyerRequestStruct validates that all required fields in BuyerRequest are present
// Uses ozzo-validation to ensure CardNumberId, FirstName, and LastName are not empty
func ValidateBuyerRequestStruct(r requests.BuyerRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CardNumberId, validation.Required),
		validation.Field(&r.FirstName, validation.Required),
		validation.Field(&r.LastName, validation.Required),
	)
}

// IsNotAnEmptyBuyer validates that at least one field in BuyerRequest is provided
// Used for PATCH operations where partial updates are allowed
func IsNotAnEmptyBuyer(r requests.BuyerRequest) error {
	if r.CardNumberId != "" || r.FirstName != "" || r.LastName != "" {
		return nil
	}
	return errors.New("at least one of id_card_number, first_name, or last_name is required")
}
