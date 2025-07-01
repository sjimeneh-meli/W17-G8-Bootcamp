package validations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

type ProductValidation struct {
}

func GetProductValidation() *ProductValidation {
	return &ProductValidation{}
}

func (v ProductValidation) ValidateProductRequestStruct(r requests.ProductRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Id, validation.Required),
		validation.Field(&r.ProductCode, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.Width, validation.Required),
		validation.Field(&r.Height, validation.Required),
		validation.Field(&r.Length, validation.Required),
		validation.Field(&r.ProductTypeID, validation.Required),
		validation.Field(&r.NetWeight, validation.Required),
		validation.Field(&r.ExpirationRate, validation.Required),
		validation.Field(&r.RecommendedFreezingTemperature, validation.Required),
		validation.Field(&r.FreezingRate, validation.Required),
		validation.Field(&r.ProductTypeID, validation.Required),
		validation.Field(&r.SellerID, validation.Required),
	)
}
