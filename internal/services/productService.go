package services

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

// ProductServiceI define la interfaz para el servicio de productos
// ProductServiceI defines the interface for the product service
// Esta interfaz implementa el patrón Repository/Service, separando la lógica de negocio del acceso a datos
// This interface implements the Repository/Service pattern, separating business logic from data access
type ProductServiceI interface {
	GetAll() ([]models.Product, error)                                // Obtiene todos los productos / Gets all products
	GetByID(id int) (*models.Product, error)                           // Obtiene un producto por ID / Gets a product by ID
	Create(product models.Product) (models.Product, error)             // Crea un nuevo producto / Creates a new product
	CreateByBatch(products []models.Product) ([]models.Product, error) // Crea múltiples productos en lote / Creates multiple products in batch
	UpdateById(id int, product models.Product) (models.Product, error) // Actualiza un producto por ID / Updates a product by ID
	DeleteById(id int) error                                           // Elimina un producto por ID / Deletes a product by ID
	ExistById(id int) bool                                             // Verifica si existe un producto por ID / Checks if a product exists by ID
}

// productService es la implementación concreta de ProductServiceI
// productService is the concrete implementation of ProductServiceI
// Utiliza composición para encapsular el repositorio y aplicar inyección de dependencias
// Uses composition to encapsulate the repository and apply dependency injection
type productService struct {
	repository repositories.ProductRepositoryI // Dependencia del repositorio / Repository dependency
}

// NewProductService es un constructor que implementa el patrón de inyección de dependencias
// NewProductService is a constructor that implements the dependency injection pattern
// Recibe una interfaz del repositorio, permitiendo testeo y flexibilidad
// Receives a repository interface, allowing testing and flexibility
func NewProductService(repository repositories.ProductRepositoryI) ProductServiceI {
	return &productService{
		repository: repository,
	}
}

// GetAll obtiene todos los productos delegando la responsabilidad al repositorio
// GetAll gets all products by delegating responsibility to the repository
func (ps *productService) GetAll() ([]models.Product, error) {
	return ps.repository.GetAll()
}

// GetByID obtiene un producto específico por su ID
// GetByID gets a specific product by its ID
func (ps *productService) GetByID(id int) (*models.Product, error) {
	return ps.repository.GetByID(id)
}

// Create crea un nuevo producto en el sistema
// Create creates a new product in the system
func (ps *productService) Create(newProduct models.Product) (models.Product, error) {
	return ps.repository.Create(newProduct)
}

// CreateByBatch permite crear múltiples productos de una vez (operación en lote)
// CreateByBatch allows creating multiple products at once (batch operation)
func (ps *productService) CreateByBatch(products []models.Product) ([]models.Product, error) {
	return ps.repository.CreateByBatch(products)
}

// UpdateById actualiza un producto existente identificado por su ID
// UpdateById updates an existing product identified by its ID
func (ps *productService) UpdateById(id int, updateProduct models.Product) (models.Product, error) {
	return ps.repository.UpdateById(id, updateProduct)
}

// DeleteById elimina un producto del sistema por su ID
// DeleteById deletes a product from the system by its ID
func (ps *productService) DeleteById(id int) error {
	return ps.repository.DeleteById(id)
}

// ExistById verifica si un producto existe en el sistema
// ExistById checks if a product exists in the system
func (ps *productService) ExistById(id int) bool {
	return ps.repository.ExistById(id)
}
