package validations

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
)

type EmployeeValidation struct {
}

func ValidateEmployeeRequestStruct(r requests.EmployeeRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CardNumberID, validation.Required),
		validation.Field(&r.FirstName, validation.Required),
		validation.Field(&r.LastName, validation.Required),
		validation.Field(&r.WarehouseID, validation.Required),
	)
}

func IsNotAnEmptyEmployee(r requests.EmployeeRequest) error {
	if r.CardNumberID != "" || r.FirstName != "" || r.LastName != "" {
		return nil
	}
	return errors.New("at least one of id_card_number, first_name, or last_name is required")
}
