package validations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

// ValidatePurchaseOrderRequestStruct validates that all required fields in PurchaseOrderRequest are present
// Uses ozzo-validation to ensure OrderNumber, OrderDate, TrackingCode, BuyerId, and ProductRecordId are not empty
func ValidatePurchaseOrderRequestStruct(r requests.PurchaseOrderRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OrderNumber, validation.Required),
		validation.Field(&r.OrderDate, validation.Required),
		validation.Field(&r.TrackingCode, validation.Required),
		validation.Field(&r.BuyerId, validation.Required),
		validation.Field(&r.ProductRecordId, validation.Required),
	)
}
