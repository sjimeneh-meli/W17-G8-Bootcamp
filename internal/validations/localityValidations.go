package validations

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)
import validation "github.com/go-ozzo/ozzo-validation/v4"

func ValidateLocalityRequestStruct(r models.Locality) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.LocalityName, validation.Required),
		validation.Field(&r.ProvinceName, validation.Required),
		validation.Field(&r.CountryName, validation.Required),
	)
}
