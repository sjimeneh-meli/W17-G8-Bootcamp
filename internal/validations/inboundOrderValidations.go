package validations

import (
	"fmt"
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

// ValidateInboundOrderRequestStruct valida que todos los campos necesarios estén presentes
func ValidateInboundOrderRequestStruct(r requests.InboundOrderRequest) error {
	if isInboundOrderAttributesEmpty(r.Data) {
		fields := []string{}
		val := reflect.ValueOf(r.Data)
		for i := 0; i < val.Type().NumField(); i++ {
			fields = append(fields, val.Type().Field(i).Tag.Get("json"))
		}

		return fmt.Errorf("data: cannot be blank. fields %v are required inside of data", strings.Join(fields, ", "))
	}

	return validation.ValidateStruct(&r.Data,
		validation.Field(&r.Data.OrderNumber, validation.Required),
		validation.Field(&r.Data.OrderDate, validation.Required),
		validation.Field(&r.Data.EmployeeId, validation.Required),
		validation.Field(&r.Data.ProductBatchId, validation.Required),
		validation.Field(&r.Data.WarehouseId, validation.Required),
	)
}

func isInboundOrderAttributesEmpty(d requests.InboundOrderAttributes) bool {
	return d.OrderNumber == "" &&
		d.OrderDate.IsZero() &&
		d.EmployeeId == 0 &&
		d.ProductBatchId == 0 &&
		d.WarehouseId == 0
}
