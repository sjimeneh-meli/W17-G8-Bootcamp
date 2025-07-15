package validations

import (
	"fmt"
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

// ValidatePurchaseOrderRequestStruct validates that all required fields in PurchaseOrderRequest are present
// Uses ozzo-validation to ensure Data, OrderNumber, OrderDate, TrackingCode, BuyerId, and ProductRecordId are not empty
func ValidatePurchaseOrderRequestStruct(r requests.PurchaseOrderRequest) error {
	if isPurchaseOrderAttributesEmpty(r.Data) {
		fields := []string{}
		val := reflect.ValueOf(r.Data)
		for i := 0; i < val.Type().NumField(); i++ {
			fields = append(fields, val.Type().Field(i).Tag.Get("json"))
		}

		return fmt.Errorf("data: cannot be blank. fields %v are neccesary inside of data", strings.Join(fields, ", "))
	}

	// validation that internal fields of data are present
	return validation.ValidateStruct(&r.Data,
		validation.Field(&r.Data.OrderNumber, validation.Required),
		validation.Field(&r.Data.OrderDate, validation.Required),
		validation.Field(&r.Data.TrackingCode, validation.Required),
		validation.Field(&r.Data.BuyerId, validation.Required),
		validation.Field(&r.Data.ProductRecordId, validation.Required),
	)
}

func isPurchaseOrderAttributesEmpty(d requests.PurchaseOrderAttributes) bool {
	return d.OrderNumber == "" &&
		d.OrderDate.IsZero() &&
		d.TrackingCode == "" &&
		d.BuyerId == 0 &&
		d.ProductRecordId == 0
}
