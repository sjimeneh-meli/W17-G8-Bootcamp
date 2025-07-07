package validations

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

func GetEmployeeValidation() *EmployeeValidation {
	return &EmployeeValidation{}
}

type EmployeeValidation struct {
}

func (v EmployeeValidation) ValidateEmployeeRequestStruct(r requests.EmployeeRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CardNumberID, validation.Required),
		validation.Field(&r.FirstName, validation.Required),
		validation.Field(&r.LastName, validation.Required),
		validation.Field(&r.WarehouseID, validation.Required),
	)
}
