package validations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

type ProductRecordValidation struct {
}

func GetProductRecordValidation() *ProductValidation {
	return &ProductValidation{}
}

func (v *ProductRecordValidation) ValidateProductRecordRequestStruct(r requests.ProductRecordRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.LastUpdateDate, validation.Required),
		validation.Field(&r.PurchasePrice, validation.Required),
		validation.Field(&r.SalePrice, validation.Required),
		validation.Field(&r.ProductID, validation.Required),
	)
}
