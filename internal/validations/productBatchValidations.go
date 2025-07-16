package validations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

func GetProductBatchValidation() *ProductBatchValidation {
	return &ProductBatchValidation{}
}

type ProductBatchValidation struct{}

func (v ProductBatchValidation) ValidateProductBatchRequestStruc(r requests.ProductBatchRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BatchNumber, validation.Required),
		validation.Field(&r.CurrentQuantity, validation.Required),
		validation.Field(&r.CurrentTemperature, validation.Required),
		validation.Field(&r.DueDate, validation.Required),
		validation.Field(&r.InitialQuantity, validation.Required),
		validation.Field(&r.ManufacturingDate, validation.Required),
		validation.Field(&r.ManufacturingHour, validation.Required),
		validation.Field(&r.MinimumTemperature, validation.Required),
		validation.Field(&r.ProductID, validation.Required),
		validation.Field(&r.SectionID, validation.Required),
	)
}
