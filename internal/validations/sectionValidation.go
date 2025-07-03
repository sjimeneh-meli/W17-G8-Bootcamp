package validations

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

func GetSectionValidation() *SectionValidation {
	return &SectionValidation{}
}

type SectionValidation struct {
}

func (v SectionValidation) ValidateSectionRequestStruct(r requests.SectionRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.SectionNumber, validation.Required),
		validation.Field(&r.CurrentCapacity, validation.Required),
		validation.Field(&r.CurrentTemperature, validation.Required),
		validation.Field(&r.MaximumCapacity, validation.Required),
		validation.Field(&r.MinimumCapacity, validation.Required),
		validation.Field(&r.MinimumTemperature, validation.Required),
		validation.Field(&r.ProductTypeID, validation.Required),
		validation.Field(&r.WarehouseID, validation.Required),
	)
}
