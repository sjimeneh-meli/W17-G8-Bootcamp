package validations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

func ValidateWarehouseRequestStruct(r requests.WarehouseRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Address, validation.Required),
		validation.Field(&r.Telephone, validation.Required),
		validation.Field(&r.WareHouseCode, validation.Required),
		validation.Field(&r.MinimumCapacity, validation.Required, validation.Min(1)),
		validation.Field(&r.MinimumTemperature, validation.Required),
	)
}

func ValidateWarehousePatchRequest(r requests.WarehousePatchRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Address, validation.When(r.Address != nil, validation.Required)),
		validation.Field(&r.Telephone, validation.When(r.Telephone != nil, validation.Required)),
		validation.Field(&r.WareHouseCode, validation.When(r.WareHouseCode != nil, validation.Required)),
		validation.Field(&r.MinimumCapacity, validation.When(r.MinimumCapacity != nil, validation.Required, validation.Min(1))),
		validation.Field(&r.MinimumTemperature, validation.When(r.MinimumTemperature != nil, validation.Required)),
	)
}
