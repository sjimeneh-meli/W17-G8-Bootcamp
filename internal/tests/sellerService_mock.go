package tests

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockSellerService is a mock implementation of SellerService using testify/mock
// MockSellerService es una implementación mock de SellerService usando testify/mock
type MockSellerService struct {
	mock.Mock
}

// NewMockSellerService creates a new instance of MockSellerService
// NewMockSellerService crea una nueva instancia de MockSellerService
func NewMockSellerService() *MockSellerService {
	return &MockSellerService{}
}

// GetAll mocks the GetAll method of SellerService
// GetAll mockea el método GetAll de SellerService
func (m *MockSellerService) GetAll() ([]models.Seller, error) {
	args := m.Called()
	return args.Get(0).([]models.Seller), args.Error(1)
}

// GetById mocks the GetById method of SellerService
// GetById mockea el método GetById de SellerService
func (m *MockSellerService) GetById(id int) (models.Seller, error) {
	args := m.Called(id)
	return args.Get(0).(models.Seller), args.Error(1)
}

// Save mocks the Save method of SellerService
// Save mockea el método Save de SellerService
func (m *MockSellerService) Save(seller models.Seller) ([]models.Seller, error) {
	args := m.Called(seller)
	return args.Get(0).([]models.Seller), args.Error(1)
}

// Update mocks the Update method of SellerService
// Update mockea el método Update de SellerService
func (m *MockSellerService) Update(id int, seller models.Seller) ([]models.Seller, error) {
	args := m.Called(id, seller)
	return args.Get(0).([]models.Seller), args.Error(1)
}

// Delete mocks the Delete method of SellerService
// Delete mockea el método Delete de SellerService
func (m *MockSellerService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
