package repositories

import (
	"errors"
	"strings"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

// ProductRepositoryI defines the contract for product repository operations
// ProductRepositoryI define el contrato para las operaciones del repositorio de productos
type ProductRepositoryI interface {
	// GetAll retrieves all products from storage
	// GetAll obtiene todos los productos del almacenamiento
	GetAll() ([]models.Product, error)

	// GetByID retrieves a product by its ID
	// GetByID obtiene un producto por su ID
	GetByID(id int) (*models.Product, error)

	// Create adds a new product to storage
	// Create agrega un nuevo producto al almacenamiento
	Create(product models.Product) (models.Product, error)

	// CreateByBatch adds multiple products to storage in a single operation
	// CreateByBatch agrega múltiples productos al almacenamiento en una sola operación
	CreateByBatch(products []models.Product) ([]models.Product, error)

	// UpdateById updates an existing product by its ID
	// UpdateById actualiza un producto existente por su ID
	UpdateById(id int, product models.Product) (models.Product, error)

	// DeleteById removes a product from storage by its ID
	// DeleteById elimina un producto del almacenamiento por su ID
	DeleteById(id int) error

	// ExistById checks if a product exists by its ID
	// ExistById verifica si un producto existe por su ID
	ExistById(id int) bool

	// ExistByProductCode checks if a product exists by its product code
	// ExistByProductCode verifica si un producto existe por su código de producto
	ExistByProductCode(id string) bool
}

// productRepository implements ProductRepositoryI interface using JSON storage
// productRepository implementa la interfaz ProductRepositoryI usando almacenamiento JSON
type productRepository struct {
	Storage loader.StorageJSON[models.Product]
}

// NewProductRepository creates a new instance of product repository
// NewProductRepository crea una nueva instancia del repositorio de productos
func NewProductRepository(storage loader.StorageJSON[models.Product]) ProductRepositoryI {
	return &productRepository{
		Storage: storage,
	}
}

// GetAll retrieves all products from storage and returns them as a slice
// GetAll obtiene todos los productos del almacenamiento y los devuelve como un slice
func (pr *productRepository) GetAll() ([]models.Product, error) {
	// Read all products from storage
	// Leer todos los productos del almacenamiento
	productsMap, err := pr.Storage.ReadAll()

	if err != nil {
		return nil, err
	}

	// Check if there are no products
	// Verificar si no hay productos
	if len(productsMap) == 0 {
		return nil, error_message.ErrNotFound
	}

	// Convert the map to a slice
	// Convertir el mapa a un slice
	productSlice := pr.Storage.MapToSlice(productsMap)

	return productSlice, nil
}

// GetByID retrieves a specific product by its ID
// GetByID obtiene un producto específico por su ID
func (pr *productRepository) GetByID(id int) (*models.Product, error) {
	products, err := pr.Storage.ReadAll()

	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, error_message.ErrNotFound
	}

	// Get the product by ID
	// Obtener el producto por ID
	productFound := products[id]

	if productFound == (models.Product{}) {
		return nil, error_message.ErrNotFound
	}

	return &productFound, nil

}

// Create adds a new product to storage after validating it doesn't already exist
// Create agrega un nuevo producto al almacenamiento después de validar que no exista ya
func (pr *productRepository) Create(newProduct models.Product) (models.Product, error) {

	// Check if product already exists by ID
	// Verificar si el producto ya existe por ID
	if pr.ExistById(newProduct.Id) {
		return models.Product{}, error_message.ErrAlreadyExists
	}

	// Check if product already exists by product code
	// Verificar si el producto ya existe por código de producto
	if pr.ExistByProductCode(newProduct.ProductCode) {
		return models.Product{}, errors.Join(error_message.ErrAlreadyExists, errors.New("error: product with the same product_code exists"))

	}

	// Read all products from storage
	// Leer todos los productos del almacenamiento
	productsMap, err := pr.Storage.ReadAll()

	if err != nil {
		return models.Product{}, err
	}

	// Get the new ID for the product
	// Obtener el nuevo ID para el producto
	idCount := len(productsMap) + 1
	newProduct.Id = idCount

	// Add the new product to the map
	// Agregar el nuevo producto al mapa
	productsMap[newProduct.Id] = newProduct

	// Convert the map to a slice
	// Convertir el mapa a un slice
	productSlice := pr.Storage.MapToSlice(productsMap)

	err = pr.Storage.WriteAll(productSlice)

	if err != nil {
		return models.Product{}, err
	}

	return newProduct, nil

}

// CreateByBatch creates multiple products in a single batch operation
// CreateByBatch crea múltiples productos en una sola operación por lotes
func (pr *productRepository) CreateByBatch(products []models.Product) ([]models.Product, error) {
	// Create each product individually
	// Crear cada producto individualmente
	for _, currentProduct := range products {
		_, err := pr.Create(currentProduct)
		if err != nil {
			return nil, err
		}
	}

	return products, nil
}

// UpdateById updates an existing product by its ID
// UpdateById actualiza un producto existente por su ID
func (pr *productRepository) UpdateById(id int, productToUpdate models.Product) (models.Product, error) {
	if !pr.ExistById(id) {
		return models.Product{}, error_message.ErrNotFound
	}

	productsMap, err := pr.Storage.ReadAll()

	if err != nil {
		return models.Product{}, err
	}

	// Ensure the product maintains the correct ID
	// Asegurar que el producto mantenga el ID correcto
	productToUpdate.Id = id

	// Update the product in the map
	// Actualizar el producto en el mapa
	productsMap[id] = productToUpdate
	productsSlice := pr.Storage.MapToSlice(productsMap)

	err = pr.Storage.WriteAll(productsSlice)

	if err != nil {
		return models.Product{}, err
	}

	return productToUpdate, nil

}

// DeleteById removes a product from storage by its ID
// DeleteById elimina un producto del almacenamiento por su ID
func (pr *productRepository) DeleteById(id int) error {
	if !pr.ExistById(id) {
		return error_message.ErrNotFound
	}

	productsMap, _ := pr.Storage.ReadAll()

	// Remove product from map
	// Eliminar producto del mapa
	delete(productsMap, id)

	productsSlice := pr.Storage.MapToSlice(productsMap)

	err := pr.Storage.WriteAll(productsSlice)

	return err
}

// ExistById checks if a product exists in storage by its ID
// ExistById verifica si un producto existe en el almacenamiento por su ID
func (pr *productRepository) ExistById(id int) bool {
	products, err := pr.Storage.ReadAll()

	if err != nil {
		return false
	}

	if len(products) == 0 {
		return false
	}

	// Check if the product exists in the map
	// Verificar si el producto existe en el mapa
	_, exist := products[id]

	return exist

}

// ExistByProductCode checks if a product exists in storage by its product code (case-insensitive)
// ExistByProductCode verifica si un producto existe en el almacenamiento por su código de producto (insensible a mayúsculas)
func (pr *productRepository) ExistByProductCode(productCode string) bool {
	products, err := pr.Storage.ReadAll()

	if err != nil {
		return false
	}

	// Check if there are no products
	// Verificar si no hay productos
	if len(products) == 0 {
		return false
	}

	// Compare product codes case-insensitively using EqualFold for better performance
	// Comparar códigos de productos sin distinguir mayúsculas usando EqualFold para mejor rendimiento
	for _, currentProduct := range products {
		if strings.EqualFold(currentProduct.ProductCode, productCode) {
			return true
		}
	}

	return false
}
