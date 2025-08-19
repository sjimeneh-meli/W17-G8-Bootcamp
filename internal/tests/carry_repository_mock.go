package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type CarryRepositoryMock struct {
	CreateFunc                    func(ctx context.Context, carry models.Carry) (models.Carry, error)
	ExistsByCidFunc               func(ctx context.Context, cid string) (bool, error)
	GetCarryReportsByLocalityFunc func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error)
}

func (m *CarryRepositoryMock) Create(ctx context.Context, carry models.Carry) (models.Carry, error) {
	return m.CreateFunc(ctx, carry)
}

func (m *CarryRepositoryMock) ExistsByCid(ctx context.Context, cid string) (bool, error) {
	return m.ExistsByCidFunc(ctx, cid)
}

func (m *CarryRepositoryMock) GetCarryReportsByLocality(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
	return m.GetCarryReportsByLocalityFunc(ctx, localityID)
}
