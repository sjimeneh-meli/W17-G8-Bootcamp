package seeders

import (
	"time"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

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

var NewProductBatchRequest = requests.ProductBatchRequest{
	BatchNumber:        "AAA",
	CurrentQuantity:    1,
	CurrentTemperature: 3.43,
	DueDate:            "2025-08-19",
	InitialQuantity:    1,
	ManufacturingDate:  "2025-08-19",
	ManufacturingHour:  1,
	MinimumTemperature: 3.43,
	ProductID:          1,
	SectionID:          1,
}

var NewProductBatchRequestTwo = requests.ProductBatchRequest{
	BatchNumber:        "AAA",
	CurrentQuantity:    1,
	CurrentTemperature: 3.43,
	DueDate:            "2025-08-19",
	InitialQuantity:    1,
	ManufacturingDate:  "2025-08-19",
	ManufacturingHour:  1,
	MinimumTemperature: 3.43,
	ProductID:          1,
	SectionID:          5,
}

var NewProductBatchRequestThree = requests.ProductBatchRequest{
	BatchNumber:        "AAA",
	CurrentQuantity:    1,
	CurrentTemperature: 3.43,
	DueDate:            "2025-08-19",
	InitialQuantity:    1,
	ManufacturingDate:  "2025-08-19",
	ManufacturingHour:  0,
	MinimumTemperature: 3.43,
	ProductID:          1,
	SectionID:          5,
}

var NewProductBatchRequestWithoutNumber = requests.ProductBatchRequest{
	CurrentQuantity:    1,
	CurrentTemperature: 3.43,
	DueDate:            "2025-08-19",
	InitialQuantity:    1,
	ManufacturingDate:  "2025-08-19",
	ManufacturingHour:  1,
	MinimumTemperature: 3.43,
	ProductID:          1,
	SectionID:          1,
}

var NewProductBatchModel = models.ProductBatch{
	Id:                 1,
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

var NewProductBatchModelTwo = &models.ProductBatch{
	Id:                 1,
	BatchNumber:        "AAA",
	CurrentQuantity:    1,
	CurrentTemperature: 3.43,
	DueDate:            exampleDate,
	InitialQuantity:    1,
	ManufacturingDate:  exampleDate,
	ManufacturingHour:  exampleDate,
	MinimumTemperature: 3.43,
	ProductID:          1,
	SectionID:          5,
}

var NewProductBatchModelThree = &models.ProductBatch{
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
	SectionID:          5,
}
