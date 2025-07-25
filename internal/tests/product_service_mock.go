package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

// ProductServiceMock es un mock del ProductService para pruebas unitarias
// ProductServiceMock is a mock of ProductService for unit testing
type ProductServiceMock struct {
	mock.Mock
}

// GetAll simula la obtención de todos los productos
// GetAll simulates retrieving all products
func (m *ProductServiceMock) GetAll(ctx context.Context) ([]models.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Product), args.Error(1)
}

// GetByID simula la obtención de un producto por su ID
// GetByID simulates retrieving a product by its ID
func (m *ProductServiceMock) GetByID(ctx context.Context, id int64) (models.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Product), args.Error(1)
}

// Create simula la creación de un nuevo producto
// Create simulates creating a new product
func (m *ProductServiceMock) Create(ctx context.Context, product models.Product) (models.Product, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(models.Product), args.Error(1)
}

// CreateByBatch simula la creación de múltiples productos en lote
// CreateByBatch simulates creating multiple products in batch
func (m *ProductServiceMock) CreateByBatch(ctx context.Context, products []models.Product) ([]models.Product, error) {
	args := m.Called(ctx, products)
	return args.Get(0).([]models.Product), args.Error(1)
}

// Update simula la actualización de un producto existente
// Update simulates updating an existing product
func (m *ProductServiceMock) Update(ctx context.Context, id int64, product models.Product) (models.Product, error) {
	args := m.Called(ctx, id, product)
	return args.Get(0).(models.Product), args.Error(1)
}

// Delete simula la eliminación de un producto
// Delete simulates deleting a product
func (m *ProductServiceMock) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ExistById simula la verificación de existencia de un producto por ID
// ExistById simulates checking if a product exists by ID
func (m *ProductServiceMock) ExistById(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}
