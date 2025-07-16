package services

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var productServiceInstance ProductService

// NewProductService crea una nueva instancia del servicio de productos con inyección de dependencias
// NewProductService creates a new instance of the product service with dependency injection
func NewProductService(r repositories.ProductRepository) ProductService {
	if productServiceInstance != nil {
		return productServiceInstance
	}
	productServiceInstance = &service{
		repository: r,
	}
	return productServiceInstance
}

// ProductService define la interfaz para la lógica de negocio de productos
// ProductService defines the interface for product business logic
type ProductService interface {
	// GetAll obtiene todos los productos del sistema
	// GetAll retrieves all products from the system
	GetAll(ctx context.Context) ([]models.Product, error)

	// GetByID obtiene un producto específico por su ID
	// GetByID retrieves a specific product by its ID
	GetByID(ctx context.Context, id int64) (models.Product, error)

	// Create crea un nuevo producto en el sistema
	// Create creates a new product in the system
	Create(ctx context.Context, product models.Product) (models.Product, error)

	// CreateByBatch crea múltiples productos en lote para mejorar rendimiento
	// CreateByBatch creates multiple products in batch to improve performance
	CreateByBatch(ctx context.Context, products []models.Product) ([]models.Product, error)

	// Update actualiza un producto existente
	// Update updates an existing product
	Update(ctx context.Context, id int64, product models.Product) (models.Product, error)

	// Delete elimina un producto del sistema
	// Delete removes a product from the system
	Delete(ctx context.Context, id int64) error

	// ExistById verifica si un producto existe por su ID
	// ExistById checks if a product exists by its ID
	ExistById(ctx context.Context, id int64) (bool, error)
}

// service implementa la interfaz ProductService
// service implements the ProductService interface
type service struct {
	repository repositories.ProductRepository
}

// GetAll delega la obtención de todos los productos al repositorio
// GetAll delegates retrieving all products to the repository
func (s *service) GetAll(ctx context.Context) ([]models.Product, error) {
	return s.repository.GetAll(ctx)
}

// GetByID delega la obtención de un producto por ID al repositorio
// GetByID delegates retrieving a product by ID to the repository
func (s *service) GetByID(ctx context.Context, id int64) (models.Product, error) {
	return s.repository.GetByID(ctx, id)
}

// Create delega la creación de un nuevo producto al repositorio
// Create delegates creating a new product to the repository
func (s *service) Create(ctx context.Context, newProduct models.Product) (models.Product, error) {
	return s.repository.Create(ctx, newProduct)
}

// CreateByBatch delega la creación de múltiples productos en lote al repositorio
// CreateByBatch delegates creating multiple products in batch to the repository
func (s *service) CreateByBatch(ctx context.Context, products []models.Product) ([]models.Product, error) {
	return s.repository.CreateByBatch(ctx, products)
}

// Update delega la actualización de un producto al repositorio
// Update delegates updating a product to the repository
func (s *service) Update(ctx context.Context, id int64, updateProduct models.Product) (models.Product, error) {
	return s.repository.Update(ctx, id, updateProduct)
}

// Delete delega la eliminación de un producto al repositorio
// Delete delegates deleting a product to the repository
func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

// ExistById delega la verificación de existencia de un producto al repositorio
// ExistById delegates checking product existence to the repository
func (s *service) ExistById(ctx context.Context, id int64) (bool, error) {
	return s.repository.Exists(ctx, id)
}
