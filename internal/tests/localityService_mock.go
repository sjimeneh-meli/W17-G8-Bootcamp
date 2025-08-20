package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockLocalityService is a mock implementation of LocalityService using testify/mock
// MockLocalityService es una implementación mock de LocalityService usando testify/mock
type MockLocalityService struct {
	mock.Mock
}

// NewMockLocalityService creates a new instance of MockLocalityService
// NewMockLocalityService crea una nueva instancia de MockLocalityService
func NewMockLocalityService() *MockLocalityService {
	return &MockLocalityService{}
}

// Save mocks the Save method of LocalityService
// Save mockea el método Save de LocalityService
func (m *MockLocalityService) Save(ctx context.Context, locality models.Locality) (models.Locality, error) {
	args := m.Called(ctx, locality)
	return args.Get(0).(models.Locality), args.Error(1)
}

// GetSellerReports mocks the GetSellerReports method of LocalityService
// GetSellerReports mockea el método GetSellerReports de LocalityService
func (m *MockLocalityService) GetSellerReports(ctx context.Context, id int) ([]responses.LocalitySellerReport, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]responses.LocalitySellerReport), args.Error(1)
}
