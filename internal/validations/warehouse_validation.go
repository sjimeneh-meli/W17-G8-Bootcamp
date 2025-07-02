package validations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/dto"
)

func ValidateWarehouseRequestStruct(r dto.WarehouseRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Address, validation.Required),
		validation.Field(&r.Telephone, validation.Required),
		validation.Field(&r.WareHouseCode, validation.Required),
		validation.Field(&r.MinimumCapacity, validation.Required, validation.Min(1)),
		validation.Field(&r.MinimumTemperature, validation.Required),
	)
}
