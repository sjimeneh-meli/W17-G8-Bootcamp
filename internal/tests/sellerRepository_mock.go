package tests

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockSellerRepository is a mock implementation of SellerRepository using testify/mock
// MockSellerRepository es una implementación mock de SellerRepository usando testify/mock
type MockSellerRepository struct {
	mock.Mock
}

// NewMockSellerRepository creates a new instance of MockSellerRepository
// NewMockSellerRepository crea una nueva instancia de MockSellerRepository
func NewMockSellerRepository() *MockSellerRepository {
	return &MockSellerRepository{}
}

// GetAll mocks the GetAll method of SellerRepository
// GetAll mockea el método GetAll de SellerRepository
func (m *MockSellerRepository) GetAll() ([]models.Seller, error) {
	args := m.Called()
	return args.Get(0).([]models.Seller), args.Error(1)
}

// Save mocks the Save method of SellerRepository
// Save mockea el método Save de SellerRepository
func (m *MockSellerRepository) Save(seller models.Seller) ([]models.Seller, error) {
	args := m.Called(seller)
	return args.Get(0).([]models.Seller), args.Error(1)
}

// Update mocks the Update method of SellerRepository
// Update mockea el método Update de SellerRepository
func (m *MockSellerRepository) Update(id int, seller models.Seller) ([]models.Seller, error) {
	args := m.Called(id, seller)
	return args.Get(0).([]models.Seller), args.Error(1)
}

// Delete mocks the Delete method of SellerRepository
// Delete mockea el método Delete de SellerRepository
func (m *MockSellerRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
