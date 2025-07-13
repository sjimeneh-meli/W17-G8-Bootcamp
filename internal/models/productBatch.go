package models

import "time"

type ProductBatch struct {
	Id                 int       `json:"id"`
	BatchNumber        string    `json:"batch_number"`
	CurrentQuantity    int       `json:"current_quantity"`
	CurrentTemperature float64   `json:"current_temperature"`
	DueDate            time.Time `json:"due_date"`
	InitialQuantity    int       `json:"initial_quantity"`
	ManufacturingDate  time.Time `json:"manufacturing_date"`
	ManufacturingHour  time.Time `json:"manufacturing_hour"`
	MinimumTeperature  float64   `json:"minimum_temperature"`
	ProductID          int       `json:"product_id"`
	SectionID          int       `json:"section_id"`
}
