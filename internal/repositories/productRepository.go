package repositories

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

type ProductRepositoryI interface {
	GetAll() ([]models.Product, error)
	GetByID(id int) (*models.Product, error)
	Create(product models.Product) (models.Product, error)
	CreateByBatch(products []models.Product) ([]models.Product, error)
	UpdateById(id int, product models.Product) (models.Product, error)
	DeleteById(id int) error
	ExistById(id int) bool
}

type productRepository struct {
	Storage loader.StorageJSON[models.Product]
}

func NewProductRepository(storage loader.StorageJSON[models.Product]) ProductRepositoryI {
	return &productRepository{
		Storage: storage,
	}
}

func (pr *productRepository) GetAll() ([]models.Product, error) {
	productsMap, err := pr.Storage.ReadAll()

	if err != nil {
		return nil, err
	}

	if len(productsMap) == 0 {
		return nil, error_message.ErrNotFound
	}

	productSlice := pr.Storage.MapToSlice(productsMap)

	return productSlice, nil
}

func (pr *productRepository) GetByID(id int) (*models.Product, error) {
	products, err := pr.Storage.ReadAll()

	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, error_message.ErrNotFound
	}

	productFound := products[id]

	if productFound == (models.Product{}) {
		return nil, error_message.ErrNotFound
	}

	return &productFound, nil

}

func (pr *productRepository) Create(newProduct models.Product) (models.Product, error) {

	if pr.ExistById(newProduct.Id) {
		return models.Product{}, error_message.ErrAlreadyExists
	}

	productsMap, err := pr.Storage.ReadAll()

	if err != nil {
		return models.Product{}, err
	}

	productsMap[newProduct.Id] = newProduct

	productSlice := pr.Storage.MapToSlice(productsMap)

	err = pr.Storage.WriteAll(productSlice)

	if err != nil {
		return models.Product{}, err
	}

	return newProduct, nil

}

func (pr *productRepository) CreateByBatch(products []models.Product) ([]models.Product, error) {
	for _, currentProduct := range products {
		_, err := pr.Create(currentProduct)
		if err != nil {
			return nil, err
		}
	}

	return products, nil
}

func (pr *productRepository) UpdateById(id int, productToUpdate models.Product) (models.Product, error) {
	if !pr.ExistById(id) {
		return models.Product{}, error_message.ErrNotFound
	}

	productsMap, err := pr.Storage.ReadAll()

	if err != nil {
		return models.Product{}, err
	}

	productToUpdate.Id = id

	productsMap[id] = productToUpdate
	productsSlice := pr.Storage.MapToSlice(productsMap)

	err = pr.Storage.WriteAll(productsSlice)

	if err != nil {
		return models.Product{}, err
	}

	return productToUpdate, nil

}

func (pr *productRepository) DeleteById(id int) error {
	if !pr.ExistById(id) {
		return error_message.ErrNotFound
	}

	productsMap, _ := pr.Storage.ReadAll()

	delete(productsMap, id)

	productsSlice := pr.Storage.MapToSlice(productsMap)

	err := pr.Storage.WriteAll(productsSlice)

	return err
}

func (pr *productRepository) ExistById(id int) bool {
	products, err := pr.Storage.ReadAll()

	if err != nil {
		return false
	}

	if len(products) == 0 {
		return false
	}

	_, exist := products[id]

	return exist

}
