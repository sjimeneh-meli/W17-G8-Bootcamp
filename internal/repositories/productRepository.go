package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
)

// Constantes SQL para operaciones de productos
// SQL constants for product operations
const (
	// productColumns define las columnas principales de la tabla products
	// productColumns defines the main columns of the products table
	productColumns = `
		id, description, expiration_rate, freezing_rate, height, 
		length, net_weight, product_code, recommended_freezing_temperature, 
		width, product_type_id, seller_id
	`
	queryGetAllProducts = "SELECT " + productColumns + " FROM products"
	queryGetProductByID = "SELECT " + productColumns + " FROM products WHERE id = ?"
	queryCreateProduct  = `
		INSERT INTO products (description, expiration_rate, freezing_rate, height, 
							  length, net_weight, product_code, recommended_freezing_temperature, 
							  width, product_type_id, seller_id) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	queryUpdateProduct = `
		UPDATE products SET description = ?, expiration_rate = ?, freezing_rate = ?, height = ?, 
		length = ?, net_weight = ?, product_code = ?, recommended_freezing_temperature = ?, 
		width = ?, product_type_id = ?, seller_id = ?
		WHERE id = ?`
	queryDeleteProduct    = "DELETE FROM products WHERE id = ?"
	queryExistProductID   = "SELECT 1 FROM products WHERE id = ? LIMIT 1"
	queryExistProductCode = "SELECT 1 FROM products WHERE product_code = ? LIMIT 1"
)

// ProductRepository define la interfaz para operaciones de productos en la base de datos
// ProductRepository defines the interface for product database operations
type ProductRepository interface {
	// GetAll obtiene todos los productos de la base de datos
	// GetAll retrieves all products from the database
	GetAll(ctx context.Context) ([]models.Product, error)

	// GetByID obtiene un producto por su ID
	// GetByID retrieves a product by its ID
	GetByID(ctx context.Context, id int64) (models.Product, error)

	// Create crea un nuevo producto en la base de datos
	// Create creates a new product in the database
	Create(ctx context.Context, newProduct models.Product) (models.Product, error)

	// CreateByBatch crea múltiples productos en una sola transacción
	// CreateByBatch creates multiple products in a single transaction
	CreateByBatch(ctx context.Context, products []models.Product) ([]models.Product, error)

	// Update actualiza un producto existente
	// Update updates an existing product
	Update(ctx context.Context, id int64, product models.Product) (models.Product, error)

	// Delete elimina un producto por su ID
	// Delete removes a product by its ID
	Delete(ctx context.Context, id int64) error

	// Exists verifica si un producto existe por su ID
	// Exists checks if a product exists by its ID
	Exists(ctx context.Context, id int64) (bool, error)

	// ExistsByProductCode verifica si existe un producto con el código dado
	// ExistsByProductCode checks if a product exists with the given product code
	ExistsByProductCode(ctx context.Context, productCode string) (bool, error)
}

// service implementa la interfaz ProductRepository
// service implements the ProductRepository interface
type service struct {
	db *sql.DB
}

// NewProductRepository crea una nueva instancia del repositorio de productos
// NewProductRepository creates a new instance of the product repository
func NewProductRepository(db *sql.DB) ProductRepository {
	return &service{
		db: db,
	}
}

// GetAll obtiene todos los productos de la base de datos
// GetAll retrieves all products from the database
func (pr *service) GetAll(ctx context.Context) ([]models.Product, error) {
	rows, err := pr.db.QueryContext(ctx, queryGetAllProducts)

	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product

		if err := rows.Scan(&p.Id, &p.Description, &p.ExpirationRate, &p.FreezingRate, &p.Height,
			&p.Length, &p.NetWeight, &p.ProductCode, &p.RecommendedFreezingTemperature,
			&p.Width, &p.ProductTypeID, &p.SellerID); err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product rows: %w", err)
	}

	return products, nil
}

// GetByID obtiene un producto específico por su ID
// GetByID retrieves a specific product by its ID
func (pr *service) GetByID(ctx context.Context, id int64) (models.Product, error) {
	var p models.Product
	err := pr.db.QueryRowContext(ctx, queryGetProductByID, id).Scan(
		&p.Id, &p.Description, &p.ExpirationRate, &p.FreezingRate, &p.Height,
		&p.Length, &p.NetWeight, &p.ProductCode, &p.RecommendedFreezingTemperature,
		&p.Width, &p.ProductTypeID, &p.SellerID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Product{}, error_message.ErrNotFound
		}
		return models.Product{}, fmt.Errorf("failed to scan product with id %d: %w", id, err)
	}

	return p, nil
}

// Create crea un nuevo producto en la base de datos y retorna el producto con su ID asignado
// Create creates a new product in the database and returns the product with its assigned ID
func (pr *service) Create(ctx context.Context, newProduct models.Product) (models.Product, error) {
	result, err := pr.db.ExecContext(ctx, queryCreateProduct, newProduct.Description, newProduct.ExpirationRate,
		newProduct.FreezingRate, newProduct.Height, newProduct.Length, newProduct.NetWeight,
		newProduct.ProductCode, newProduct.RecommendedFreezingTemperature, newProduct.Width,
		newProduct.ProductTypeID, newProduct.SellerID)

	if err != nil {
		return models.Product{}, fmt.Errorf("failed to create product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Product{}, fmt.Errorf("failed to get last insert id for product: %w", err)
	}

	newProduct.Id = id
	return newProduct, nil
}

// CreateByBatch crea múltiples productos en una sola transacción para mejorar el rendimiento
// CreateByBatch creates multiple products in a single transaction to improve performance
func (pr *service) CreateByBatch(ctx context.Context, products []models.Product) ([]models.Product, error) {
	// Inicia una transacción para asegurar atomicidad
	// Start a transaction to ensure atomicity
	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	// Prepara la declaración SQL para reutilizar
	// Prepare the SQL statement for reuse
	stmt, err := tx.PrepareContext(ctx, queryCreateProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Ejecuta la inserción para cada producto
	// Execute the insertion for each product
	for i := range products {
		product := &products[i]
		result, err := stmt.ExecContext(ctx, product.Description, product.ExpirationRate,
			product.FreezingRate, product.Height, product.Length, product.NetWeight,
			product.ProductCode, product.RecommendedFreezingTemperature, product.Width,
			product.ProductTypeID, product.SellerID)

		if err != nil {

			return nil, fmt.Errorf("failed to execute statement for product %s: %w", product.ProductCode, err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get last insert id for product %s: %w", product.ProductCode, err)
		}
		product.Id = id
	}

	// Confirma la transacción
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return products, nil
}

// Update actualiza un producto existente en la base de datos
// Update updates an existing product in the database
func (pr *service) Update(ctx context.Context, id int64, product models.Product) (models.Product, error) {
	result, err := pr.db.ExecContext(ctx, queryUpdateProduct,
		product.Description, product.ExpirationRate, product.FreezingRate, product.Height,
		product.Length, product.NetWeight, product.ProductCode, product.RecommendedFreezingTemperature,
		product.Width, product.ProductTypeID, product.SellerID, id)

	if err != nil {
		return models.Product{}, fmt.Errorf("failed to update product %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Product{}, fmt.Errorf("failed to get rows affected for product %d: %w", id, err)
	}

	// Verifica si el producto existe
	// Check if the product exists
	if rowsAffected == 0 {
		return models.Product{}, error_message.ErrNotFound
	}

	product.Id = id
	return product, nil
}

// Delete elimina un producto de la base de datos por su ID
// Delete removes a product from the database by its ID
func (pr *service) Delete(ctx context.Context, id int64) error {
	result, err := pr.db.ExecContext(ctx, queryDeleteProduct, id)
	if err != nil {
		return fmt.Errorf("failed to delete product %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected on delete for product %d: %w", id, err)
	}

	// Verifica si el producto existía
	// Check if the product existed
	if rowsAffected == 0 {
		return error_message.ErrNotFound
	}

	return nil
}

// Exists verifica si un producto existe en la base de datos por su ID
// Exists checks if a product exists in the database by its ID
func (pr *service) Exists(ctx context.Context, id int64) (bool, error) {
	var exists int
	err := pr.db.QueryRowContext(ctx, queryExistProductID, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error checking existence of product %d: %w", id, err)
	}
	return true, nil
}

// ExistsByProductCode verifica si existe un producto con el código de producto dado
// ExistsByProductCode checks if a product exists with the given product code
func (pr *service) ExistsByProductCode(ctx context.Context, productCode string) (bool, error) {
	var exists int
	err := pr.db.QueryRowContext(ctx, queryExistProductCode, productCode).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("error checking existence of product code %s: %w", productCode, err)
	}
	return true, nil
}
