package validations

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

func ValidateBuyerRequestStruct(r requests.BuyerRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CardNumberId, validation.Required),
		validation.Field(&r.FirstName, validation.Required),
		validation.Field(&r.LastName, validation.Required),
	)
}
