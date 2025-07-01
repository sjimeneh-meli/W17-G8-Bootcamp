package services

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

type ProductServiceI interface {
	GetAll() ([]models.Product, error)
	GetByID(id int) (*models.Product, error)
	Create(product models.Product) (models.Product, error)
	CreateByBatch(products []models.Product) ([]models.Product, error)
	UpdateById(id int, product models.Product) (models.Product, error)
	DeleteById(id int) error
	ExistById(id int) bool
}

type productService struct {
	repository repositories.ProductRepositoryI
}

func NewProductService(repository repositories.ProductRepositoryI) ProductServiceI {
	return &productService{
		repository: repository,
	}
}

func (ps *productService) GetAll() ([]models.Product, error) {
	return ps.repository.GetAll()
}

func (ps *productService) GetByID(id int) (*models.Product, error) {
	return ps.repository.GetByID(id)
}

func (ps *productService) Create(newProduct models.Product) (models.Product, error) {
	return ps.repository.Create(newProduct)
}

func (ps *productService) CreateByBatch(products []models.Product) ([]models.Product, error) {
	return ps.repository.CreateByBatch(products)
}

func (ps *productService) UpdateById(id int, updateProduct models.Product) (models.Product, error) {
	return ps.repository.UpdateById(id, updateProduct)
}

func (ps *productService) DeleteById(id int) error {
	return ps.repository.DeleteById(id)
}

func (ps *productService) ExistById(id int) bool {
	return ps.ExistById(id)

}
