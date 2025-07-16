package mappers

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

func MapCarryToCreateCarryResponse(carry models.Carry) responses.CreateCarryResponse {
	return responses.CreateCarryResponse{
		Id:          carry.Id,
		Cid:         carry.Cid,
		CompanyName: carry.CompanyName,
		Address:     carry.Address,
		Telephone:   carry.Telephone,
		LocalityId:  carry.LocalityId,
	}
}

func MapCarryRequestToCarry(request requests.CarryRequest) models.Carry {
	return models.Carry{
		Cid:         request.Cid,
		CompanyName: request.CompanyName,
		Address:     request.Address,
		Telephone:   request.Telephone,
		LocalityId:  request.LocalityId,
	}
}
