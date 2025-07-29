package validations_test

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"github.com/stretchr/testify/assert"
)

func TestValidateSectionRequestStruct(t *testing.T) {
	t.Run("SectionRequest doesn't have required fields", func(t *testing.T) {
		expectedErrorMessage := "current_capacity: cannot be blank; current_temperature: cannot be blank; maximum_capacity: cannot be blank; minimum_capacity: cannot be blank; minimum_temperature: cannot be blank; product_type_id: cannot be blank; section_number: cannot be blank; warehouse_id: cannot be blank."
		sectionRequest := requests.SectionRequest{}

		vld := validations.GetSectionValidation()
		err := vld.ValidateSectionRequestStruct(sectionRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("SectionRequest has all fields correctly", func(t *testing.T) {
		sectionRequest := requests.SectionRequest{
			SectionNumber:      "AAA",
			CurrentCapacity:    1,
			CurrentTemperature: 3.43,
			MaximumCapacity:    1,
			MinimumCapacity:    1,
			MinimumTemperature: 3.43,
			ProductTypeID:      1,
			WarehouseID:        1,
		}

		vld := validations.GetSectionValidation()
		err := vld.ValidateSectionRequestStruct(sectionRequest)

		assert.Nil(t, err, "err should be nil")
	})
}
