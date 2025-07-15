package validations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

func ValidateCarryRequest(request requests.CarryRequest) error {
	return validation.ValidateStruct(&request,
		validation.Field(&request.Cid, validation.Required, validation.Length(1, 10)),
		validation.Field(&request.CompanyName, validation.Required, validation.Length(1, 100)),
		validation.Field(&request.Address, validation.Required, validation.Length(1, 100)),
		validation.Field(&request.Telephone, validation.Required, validation.Length(1, 10)),
		validation.Field(&request.LocalityId, validation.Required, validation.Min(1)),
	)
}
