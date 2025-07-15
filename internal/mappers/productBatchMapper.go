package mappers

import (
	"errors"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	tools "github.com/sajimenezher_meli/meli-frescos-8/pkg"
)

func GetProductBatchModelFromRequest(request *requests.ProductBatchRequest) (*models.ProductBatch, error) {
	dueDate, convErr := tools.ConvertStringToDate(request.DueDate)
	if convErr != nil {
		return nil, errors.New("cannot convert the field dueDate")
	}

	manufacturingDate, convErr := tools.ConvertStringToDate(request.ManufacturingDate)
	if convErr != nil {
		return nil, errors.New("cannot convert the field manufacturingDate")
	}

	manufacturingHour, convErr := tools.ConvertFloatToDuration(request.ManufacturingHour)
	if convErr != nil {
		return nil, errors.New("cannot convert the field manufacturingHour")
	}

	return &models.ProductBatch{
		BatchNumber:        request.BatchNumber,
		CurrentQuantity:    request.CurrentQuantity,
		CurrentTemperature: request.CurrentTemperature,
		DueDate:            dueDate,
		InitialQuantity:    request.InitialQuantity,
		ManufacturingDate:  manufacturingDate,
		ManufacturingHour:  manufacturingHour,
		MinimumTemperature: request.MinimumTemperature,
		ProductID:          request.ProductID,
		SectionID:          request.SectionID,
	}, nil
}

func GetProductBatchResponseFromModel(model *models.ProductBatch) *responses.ProductBatchResponse {
	return &responses.ProductBatchResponse{
		Id:                 model.Id,
		BatchNumber:        model.BatchNumber,
		CurrentQuantity:    model.CurrentQuantity,
		CurrentTemperature: model.CurrentTemperature,
		DueDate:            model.DueDate.String(),
		InitialQuantity:    model.InitialQuantity,
		ManufacturingDate:  model.ManufacturingDate.String(),
		ManufacturingHour:  model.ManufacturingHour.String(),
		MinimumTemperature: model.MinimumTemperature,
		ProductID:          model.ProductID,
		SectionID:          model.SectionID,
	}
}
