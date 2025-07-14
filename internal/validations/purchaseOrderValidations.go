package validations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

func ValidatePurchaseOrderRequestStruct(r requests.PurchaseOrderRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OrderNumber, validation.Required),
		validation.Field(&r.OrderDate, validation.Required),
		validation.Field(&r.TrackingCode, validation.Required),
		validation.Field(&r.BuyerId, validation.Required),
		validation.Field(&r.ProductRecordId, validation.Required),
	)
}
