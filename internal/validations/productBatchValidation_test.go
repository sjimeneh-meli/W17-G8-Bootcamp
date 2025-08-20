package validations_test

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/seeders"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"github.com/stretchr/testify/assert"
)

func TestValidateProductBatchRequestStruct(t *testing.T) {
	t.Run("ProductBatchRequest doesn't have required fields", func(t *testing.T) {
		expectedErrorMessage := "batch_number: cannot be blank; current_quantity: cannot be blank; current_temperature: cannot be blank; due_date: cannot be blank; initial_quantity: cannot be blank; manufacturing_date: cannot be blank; manufacturing_hour: cannot be blank; minimum_temperature: cannot be blank; product_id: cannot be blank; section_id: cannot be blank."
		productBatchRequest := requests.ProductBatchRequest{}

		vld := validations.GetProductBatchValidation()
		err := vld.ValidateProductBatchRequestStruc(productBatchRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("ProductBatchRequest has all fields correctly", func(t *testing.T) {
		vld := validations.GetProductBatchValidation()
		err := vld.ValidateProductBatchRequestStruc(seeders.NewProductBatchRequest)

		assert.Nil(t, err, "err should be nil")
	})
}
