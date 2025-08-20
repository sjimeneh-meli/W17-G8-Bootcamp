package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

type LocalityRepositoryMock struct {
	SaveFunc             func(ctx context.Context, locality models.Locality) (models.Locality, error)
	GetSellerReportsFunc func(ctx context.Context, localityID int) ([]responses.LocalitySellerReport, error)
	ExistByIdFunc        func(ctx context.Context, localityID int) (bool, error)
}

func (m *LocalityRepositoryMock) Save(ctx context.Context, locality models.Locality) (models.Locality, error) {
	return m.SaveFunc(ctx, locality)
}

func (m *LocalityRepositoryMock) GetSellerReports(ctx context.Context, localityID int) ([]responses.LocalitySellerReport, error) {
	return m.GetSellerReportsFunc(ctx, localityID)
}

func (m *LocalityRepositoryMock) ExistById(ctx context.Context, localityID int) (bool, error) {
	return m.ExistByIdFunc(ctx, localityID)
}
