package validations

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

func ValidateBuyerRequestStruct(r requests.BuyerRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CardNumberId, validation.Required),
		validation.Field(&r.FirstName, validation.Required),
		validation.Field(&r.LastName, validation.Required),
	)
}

func IsNotAnEmptyBuyer(r requests.BuyerRequest) error {
	if r.CardNumberId != "" || r.FirstName != "" || r.LastName != "" {
		return nil
	}
	return errors.New("at least one of card_number_id, first_name, or last_name is required")
}
