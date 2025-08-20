package seeders

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

var Sections = []*models.Section{
	{
		Id:                 1,
		SectionNumber:      "A-01",
		CurrentCapacity:    1,
		CurrentTemperature: 3.43,
		MaximumCapacity:    1,
		MinimumCapacity:    1,
		MinimumTemperature: 1,
		ProductTypeID:      1,
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
		WarehouseID:        2,
	},
	{
		Id:                 3,
		SectionNumber:      "C-01",
		CurrentCapacity:    3,
		CurrentTemperature: 3.43,
		MaximumCapacity:    3,
		MinimumCapacity:    3,
		MinimumTemperature: 3,
		ProductTypeID:      3,
		WarehouseID:        3,
	},
}

var NewSectionModel = &models.Section{
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

var NewSectionRequest = requests.SectionRequest{
	SectionNumber:      "A-01",
	CurrentCapacity:    2,
	CurrentTemperature: 3.43,
	MaximumCapacity:    2,
	MinimumCapacity:    2,
	MinimumTemperature: 2,
	ProductTypeID:      2,
	WarehouseID:        1,
}

var NewSectionRequestTwo = requests.SectionRequest{
	SectionNumber:      "B-01",
	CurrentCapacity:    2,
	CurrentTemperature: 3.43,
	MaximumCapacity:    2,
	MinimumCapacity:    2,
	MinimumTemperature: 2,
	ProductTypeID:      2,
	WarehouseID:        1,
}

var NewSectionRequestWithoutNumber = requests.SectionRequest{
	CurrentCapacity:    2,
	CurrentTemperature: 3.43,
	MaximumCapacity:    2,
	MinimumCapacity:    2,
	MinimumTemperature: 2,
	ProductTypeID:      2,
	WarehouseID:        1,
}
