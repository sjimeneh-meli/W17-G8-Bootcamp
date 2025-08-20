package mappers_test

import (
	"testing"
	"time"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetModelProductBatchFromRequest(t *testing.T) {
	t.Run("Successfully maps a ProductBatchRequest to ProductBatch", func(t *testing.T) {
		var exampleDate = time.Date(
			2025,
			time.August,
			19,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		expectedProductBatch := &models.ProductBatch{
			Id:                 0,
			BatchNumber:        "AAA",
			CurrentQuantity:    1,
			CurrentTemperature: 3.43,
			DueDate:            exampleDate,
			InitialQuantity:    1,
			ManufacturingDate:  exampleDate,
			ManufacturingHour:  exampleDate,
			MinimumTemperature: 3.43,
			ProductID:          1,
			SectionID:          1,
		}

		productBatchRequest := &requests.ProductBatchRequest{
			BatchNumber:        "AAA",
			CurrentQuantity:    1,
			CurrentTemperature: 3.43,
			DueDate:            "2025-08-19",
			InitialQuantity:    1,
			ManufacturingDate:  "2025-08-19",
			ManufacturingHour:  0,
			MinimumTemperature: 3.43,
			ProductID:          1,
			SectionID:          1,
		}

		result, _ := mappers.GetProductBatchModelFromRequest(productBatchRequest)

		assert.Equal(t, expectedProductBatch, result)
	})
}

func TestGetResponseProductBatchFromModel(t *testing.T) {
	t.Run("Successfully maps a ProductBatchModel to Response", func(t *testing.T) {
		var exampleDate = time.Date(
			2025,
			time.August,
			19,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		productBatch := &models.ProductBatch{
			Id:                 0,
			BatchNumber:        "AAA",
			CurrentQuantity:    1,
			CurrentTemperature: 3.43,
			DueDate:            exampleDate,
			InitialQuantity:    1,
			ManufacturingDate:  exampleDate,
			ManufacturingHour:  exampleDate,
			MinimumTemperature: 3.43,
			ProductID:          1,
			SectionID:          1,
		}

		expectedProductBatchResponse := &responses.ProductBatchResponse{
			Id:                 0,
			BatchNumber:        "AAA",
			CurrentQuantity:    1,
			CurrentTemperature: 3.43,
			DueDate:            "2025-08-19 00:00:00 +0000 UTC",
			InitialQuantity:    1,
			ManufacturingDate:  "2025-08-19 00:00:00 +0000 UTC",
			ManufacturingHour:  0,
			MinimumTemperature: 3.43,
			ProductID:          1,
			SectionID:          1,
		}

		result := mappers.GetProductBatchResponseFromModel(productBatch)

		assert.Equal(t, expectedProductBatchResponse, result)
	})
}
