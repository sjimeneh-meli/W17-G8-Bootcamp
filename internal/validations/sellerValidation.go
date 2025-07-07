package validations

import "github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
import validation "github.com/go-ozzo/ozzo-validation/v4"

func ValidateSellerRequestStruct(r requests.SellerRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CID, validation.Required),
		validation.Field(&r.CompanyName, validation.Required),
		validation.Field(&r.Address, validation.Required),
		validation.Field(&r.Telephone, validation.Required),
	)
}
