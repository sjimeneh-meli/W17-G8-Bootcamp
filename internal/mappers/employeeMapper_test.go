package mappers_test

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_Employee_GetEmployeeModelFromRequest(t *testing.T) {
	t.Run("Successfully maps an EmployeeRequest to Employee", func(t *testing.T) {
		expectedEmployee := &models.Employee{
			Id:           0,
			CardNumberID: "EMP-001",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		}

		employeeRequest := requests.EmployeeRequest{
			CardNumberID: "EMP-001",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		}

		result := mappers.GetEmployeeModelFromRequest(employeeRequest)

		assert.Equal(t, expectedEmployee, result)
	})

	t.Run("Successfully maps EmployeeRequest with empty fields", func(t *testing.T) {
		expectedEmployee := &models.Employee{
			Id:           0,
			CardNumberID: "",
			FirstName:    "",
			LastName:     "",
			WarehouseID:  0,
		}

		employeeRequest := requests.EmployeeRequest{
			CardNumberID: "",
			FirstName:    "",
			LastName:     "",
			WarehouseID:  0,
		}

		result := mappers.GetEmployeeModelFromRequest(employeeRequest)

		assert.Equal(t, expectedEmployee, result)
	})
}

func Test_Employee_GetEmployeeResponseFromModel(t *testing.T) {
	t.Run("Successfully maps an Employee to EmployeeResponse", func(t *testing.T) {
		employeeModel := &models.Employee{
			Id:           10,
			CardNumberID: "EMP-001",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		}

		expectedEmployeeResponse := &responses.EmployeeResponse{
			ID:           10,
			CardNumberID: "EMP-001",
			FirstName:    "Juan",
			LastName:     "Perez",
			WarehouseID:  1,
		}

		result := mappers.GetEmployeeResponseFromModel(employeeModel)

		assert.Equal(t, expectedEmployeeResponse, result)
	})

	t.Run("Successfully maps Employee with zero values", func(t *testing.T) {
		employeeModel := &models.Employee{
			Id:           0,
			CardNumberID: "",
			FirstName:    "",
			LastName:     "",
			WarehouseID:  0,
		}

		expectedEmployeeResponse := &responses.EmployeeResponse{
			ID:           0,
			CardNumberID: "",
			FirstName:    "",
			LastName:     "",
			WarehouseID:  0,
		}

		result := mappers.GetEmployeeResponseFromModel(employeeModel)

		assert.Equal(t, expectedEmployeeResponse, result)
	})
}

func Test_Employee_UpdateEmployeeModelFromRequest(t *testing.T) {
	t.Run("Successfully updates Employee model from request", func(t *testing.T) {
		// Initial employee model
		employeeModel := &models.Employee{
			Id:           10,
			CardNumberID: "OLD-001",
			FirstName:    "OldName",
			LastName:     "OldLastName",
			WarehouseID:  1,
		}

		// Update request
		updateRequest := &requests.EmployeeRequest{
			CardNumberID: "NEW-001",
			FirstName:    "NewName",
			LastName:     "NewLastName",
			WarehouseID:  2,
		}

		// Expected result after update
		expectedEmployee := &models.Employee{
			Id:           10, // ID should remain unchanged
			CardNumberID: "NEW-001",
			FirstName:    "NewName",
			LastName:     "NewLastName",
			WarehouseID:  2,
		}

		// Execute update
		mappers.UpdateEmployeeModelFromRequest(employeeModel, updateRequest)

		assert.Equal(t, expectedEmployee, employeeModel)
	})

	t.Run("Successfully updates Employee with partial data", func(t *testing.T) {
		// Initial employee model
		employeeModel := &models.Employee{
			Id:           5,
			CardNumberID: "EMP-005",
			FirstName:    "Original",
			LastName:     "Employee",
			WarehouseID:  3,
		}

		// Update request with some empty fields
		updateRequest := &requests.EmployeeRequest{
			CardNumberID: "EMP-005-UPDATED",
			FirstName:    "",
			LastName:     "UpdatedLastName",
			WarehouseID:  0,
		}

		// Expected result after update
		expectedEmployee := &models.Employee{
			Id:           5, // ID should remain unchanged
			CardNumberID: "EMP-005-UPDATED",
			FirstName:    "", // Should be updated to empty
			LastName:     "UpdatedLastName",
			WarehouseID:  0, // Should be updated to 0
		}

		// Execute update
		mappers.UpdateEmployeeModelFromRequest(employeeModel, updateRequest)

		assert.Equal(t, expectedEmployee, employeeModel)
	})
}

func Test_Employee_GetListEmployeeResponseFromListModel(t *testing.T) {
	t.Run("Successfully maps a list of Employee to a list of EmployeeResponse", func(t *testing.T) {
		employeesList := []*models.Employee{
			{
				Id:           10,
				CardNumberID: "EMP-001",
				FirstName:    "Juan",
				LastName:     "Perez",
				WarehouseID:  1,
			},
			{
				Id:           11,
				CardNumberID: "EMP-002",
				FirstName:    "Maria",
				LastName:     "Garcia",
				WarehouseID:  2,
			},
		}

		expectedResponseEmployeeList := []*responses.EmployeeResponse{
			{
				ID:           10,
				CardNumberID: "EMP-001",
				FirstName:    "Juan",
				LastName:     "Perez",
				WarehouseID:  1,
			},
			{
				ID:           11,
				CardNumberID: "EMP-002",
				FirstName:    "Maria",
				LastName:     "Garcia",
				WarehouseID:  2,
			},
		}

		result := mappers.GetListEmployeeResponseFromListModel(employeesList)

		assert.Equal(t, expectedResponseEmployeeList, result)
	})

	t.Run("Successfully handles empty list", func(t *testing.T) {
		var employeesList []*models.Employee

		result := mappers.GetListEmployeeResponseFromListModel(employeesList)

		// The function returns an empty slice, not nil
		assert.Empty(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("Successfully handles single employee list", func(t *testing.T) {
		employeesList := []*models.Employee{
			{
				Id:           1,
				CardNumberID: "EMP-SINGLE",
				FirstName:    "Single",
				LastName:     "Employee",
				WarehouseID:  99,
			},
		}

		expectedResponseEmployeeList := []*responses.EmployeeResponse{
			{
				ID:           1,
				CardNumberID: "EMP-SINGLE",
				FirstName:    "Single",
				LastName:     "Employee",
				WarehouseID:  99,
			},
		}

		result := mappers.GetListEmployeeResponseFromListModel(employeesList)

		assert.Equal(t, expectedResponseEmployeeList, result)
		assert.Len(t, result, 1)
	})
}
