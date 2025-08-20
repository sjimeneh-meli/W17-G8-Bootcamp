package mappers_test

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetModelSectionFromRequest(t *testing.T) {
	t.Run("Successfully maps a SectionRequest to Section", func(t *testing.T) {
		expectedSection := &models.Section{
			Id:                 0,
			SectionNumber:      "A-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		sectionRequest := &requests.SectionRequest{
			SectionNumber:      "A-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		result := mappers.GetSectionModelFromRequest(sectionRequest)

		assert.Equal(t, expectedSection, result)
	})
}

func TestGetResponseSectionFromModel(t *testing.T) {
	t.Run("Successfully maps a Section to SectionResponse", func(t *testing.T) {
		section := &models.Section{
			Id:                 1,
			SectionNumber:      "A-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		expectedSectionResponse := &responses.SectionResponse{
			ID:                 1,
			SectionNumber:      "A-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		}

		result := mappers.GetSectionResponseFromModel(section)

		assert.Equal(t, expectedSectionResponse, result)
	})
}

func TestTestGetListSectionResponseFromListModel(t *testing.T) {
	t.Run("Successfully maps a list of Section to a list of SectionResponse", func(t *testing.T) {
		sections := []*models.Section{{
			Id:                 1,
			SectionNumber:      "A-01",
			CurrentCapacity:    2,
			CurrentTemperature: 3.43,
			MaximumCapacity:    2,
			MinimumCapacity:    2,
			MinimumTemperature: 2,
			ProductTypeID:      2,
			WarehouseID:        1,
		},
			{
				Id:                 2,
				SectionNumber:      "B-01",
				CurrentCapacity:    2,
				CurrentTemperature: 3.43,
				MaximumCapacity:    2,
				MinimumCapacity:    2,
				MinimumTemperature: 2,
				ProductTypeID:      2,
				WarehouseID:        1,
			},
		}
		expectedSectionResponseList := []*responses.SectionResponse{
			{
				ID:                 1,
				SectionNumber:      "A-01",
				CurrentCapacity:    2,
				CurrentTemperature: 3.43,
				MaximumCapacity:    2,
				MinimumCapacity:    2,
				MinimumTemperature: 2,
				ProductTypeID:      2,
				WarehouseID:        1,
			},
			{
				ID:                 2,
				SectionNumber:      "B-01",
				CurrentCapacity:    2,
				CurrentTemperature: 3.43,
				MaximumCapacity:    2,
				MinimumCapacity:    2,
				MinimumTemperature: 2,
				ProductTypeID:      2,
				WarehouseID:        1,
			},
		}

		result := mappers.GetListSectionResponseFromListModel(sections)

		assert.Equal(t, expectedSectionResponseList, result)

	})

}
