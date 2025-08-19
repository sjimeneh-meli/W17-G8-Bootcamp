package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type CarryServiceMock struct {
	CreateCarryFunc              func(ctx context.Context, carry models.Carry) (models.Carry, error)
	GetCarryReportByLocalityFunc func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error)
}

func (m *CarryServiceMock) CreateCarry(ctx context.Context, carry models.Carry) (models.Carry, error) {
	return m.CreateCarryFunc(ctx, carry)
}

func (m *CarryServiceMock) GetCarryReportByLocality(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
	return m.GetCarryReportByLocalityFunc(ctx, localityID)
}
