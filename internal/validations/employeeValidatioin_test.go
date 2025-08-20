package validations_test

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"github.com/stretchr/testify/assert"
)

func TestGetEmployeeValidation(t *testing.T) {
	t.Run("Successfully creates EmployeeValidation instance", func(t *testing.T) {
		result := validations.GetEmployeeValidation()

		assert.NotNil(t, result)
		assert.IsType(t, &validations.EmployeeValidation{}, result)
	})
}

func TestEmployeeValidations_ValidateEmployeeRequestStruct(t *testing.T) {
	t.Run("EmployeeRequest doesn't have required fields CardNumberID, FirstName, LastName and WarehouseID", func(t *testing.T) {
		expectedErrorMessage := "first_name: cannot be blank; id_card_number: cannot be blank; last_name: cannot be blank; warehouse_id: cannot be blank."
		employeeRequest := requests.EmployeeRequest{}

		err := validations.ValidateEmployeeRequestStruct(employeeRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("EmployeeRequest doesn't have required field CardNumberID", func(t *testing.T) {
		expectedErrorMessage := "id_card_number: cannot be blank."
		employeeRequest := requests.EmployeeRequest{
			FirstName:   "Juan",
			LastName:    "Perez",
			WarehouseID: 1,
		}

		err := validations.ValidateEmployeeRequestStruct(employeeRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("EmployeeRequest doesn't have required field FirstName", func(t *testing.T) {
		expectedErrorMessage := "first_name: cannot be blank."
		employeeRequest := requests.EmployeeRequest{
			CardNumberID: "EMP-001",
			LastName:     "Perez",
			WarehouseID:  1,
		}

		err := validations.ValidateEmployeeRequestStruct(employeeRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("EmployeeRequest doesn't have required field LastName", func(t *testing.T) {
		expectedErrorMessage := "last_name: cannot be blank."
		employeeRequest := requests.EmployeeRequest{
			CardNumberID: "EMP-001",
			FirstName:    "Juan",
			WarehouseID:  1,
		}

		err := validations.ValidateEmployeeRequestStruct(employeeRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("EmployeeRequest doesn't have required field WarehouseID", func(t *testing.T) {
		expectedErrorMessage := "warehouse_id: cannot be blank."
		employeeRequest := requests.EmployeeRequest{
			CardNumberID: "EMP-001",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  0, // Zero value is considered blank
		}

		err := validations.ValidateEmployeeRequestStruct(employeeRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("EmployeeRequest missing CardNumberID and FirstName", func(t *testing.T) {
		expectedErrorMessage := "first_name: cannot be blank; id_card_number: cannot be blank."
		employeeRequest := requests.EmployeeRequest{
			LastName:    "Perez",
			WarehouseID: 1,
		}

		err := validations.ValidateEmployeeRequestStruct(employeeRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("EmployeeRequest missing LastName and WarehouseID", func(t *testing.T) {
		expectedErrorMessage := "last_name: cannot be blank; warehouse_id: cannot be blank."
		employeeRequest := requests.EmployeeRequest{
			CardNumberID: "EMP-001",
			FirstName:    "Juan",
		}

		err := validations.ValidateEmployeeRequestStruct(employeeRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("EmployeeRequest has all fields correctly", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			CardNumberID: "EMP-001",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		}

		err := validations.ValidateEmployeeRequestStruct(employeeRequest)

		assert.Nil(t, err, "err should be nil")
	})

	t.Run("EmployeeRequest has all fields with valid positive WarehouseID", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			CardNumberID: "EMP-999",
			FirstName:    "Maria",
			LastName:     "Garcia",
			WarehouseID:  999,
		}

		err := validations.ValidateEmployeeRequestStruct(employeeRequest)

		assert.Nil(t, err, "err should be nil")
	})
}

func TestEmployeeValidations_IsNotAnEmptyEmployee(t *testing.T) {
	t.Run("EmployeeRequest is completely empty returns error", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{}
		expectedErrorMessage := "at least one of id_card_number, first_name, or last_name is required"

		err := validations.IsNotAnEmptyEmployee(employeeRequest)

		assert.Error(t, err, "err should be expected")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("EmployeeRequest with only WarehouseID returns error", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			WarehouseID: 1,
		}
		expectedErrorMessage := "at least one of id_card_number, first_name, or last_name is required"

		err := validations.IsNotAnEmptyEmployee(employeeRequest)

		assert.Error(t, err, "err should be expected")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("EmployeeRequest has CardNumberID", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			CardNumberID: "EMP-001",
		}

		err := validations.IsNotAnEmptyEmployee(employeeRequest)

		assert.Nil(t, err, "err should be nil")
	})

	t.Run("EmployeeRequest has FirstName", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			FirstName: "Juan",
		}

		err := validations.IsNotAnEmptyEmployee(employeeRequest)

		assert.Nil(t, err, "err should be nil")
	})

	t.Run("EmployeeRequest has LastName", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			LastName: "Perez",
		}

		err := validations.IsNotAnEmptyEmployee(employeeRequest)

		assert.Nil(t, err, "err should be nil")
	})

	t.Run("EmployeeRequest has CardNumberID and FirstName", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			CardNumberID: "EMP-001",
			FirstName:    "Juan",
		}

		err := validations.IsNotAnEmptyEmployee(employeeRequest)

		assert.Nil(t, err, "err should be nil")
	})

	t.Run("EmployeeRequest has FirstName and LastName", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			FirstName: "Juan",
			LastName:  "Perez",
		}

		err := validations.IsNotAnEmptyEmployee(employeeRequest)

		assert.Nil(t, err, "err should be nil")
	})

	t.Run("EmployeeRequest has CardNumberID and LastName", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			CardNumberID: "EMP-001",
			LastName:     "Perez",
		}

		err := validations.IsNotAnEmptyEmployee(employeeRequest)

		assert.Nil(t, err, "err should be nil")
	})

	t.Run("EmployeeRequest is complete", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			CardNumberID: "EMP-001",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		}

		err := validations.IsNotAnEmptyEmployee(employeeRequest)

		assert.Nil(t, err, "err should be nil")
	})

	t.Run("EmployeeRequest has only CardNumberID with spaces", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			CardNumberID: " ",
		}

		err := validations.IsNotAnEmptyEmployee(employeeRequest)

		assert.Nil(t, err, "err should be nil - spaces are considered non-empty")
	})

	t.Run("EmployeeRequest has only FirstName with spaces", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			FirstName: "   ",
		}

		err := validations.IsNotAnEmptyEmployee(employeeRequest)

		assert.Nil(t, err, "err should be nil - spaces are considered non-empty")
	})

	t.Run("EmployeeRequest has only LastName with spaces", func(t *testing.T) {
		employeeRequest := requests.EmployeeRequest{
			LastName: "  \t  ",
		}

		err := validations.IsNotAnEmptyEmployee(employeeRequest)

		assert.Nil(t, err, "err should be nil - spaces are considered non-empty")
	})
}
