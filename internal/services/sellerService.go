package services

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var sellerServiceInstance SellerService

// NewJSONSellerService creates and returns a singleton instance of JsonSellerService with the required repository
// NewJSONSellerService crea y retorna una instancia singleton de JsonSellerService con el repositorio requerido
func NewJSONSellerService(repo repositories.SellerRepository) SellerService {
	if sellerServiceInstance != nil {
		return sellerServiceInstance
	}
	sellerServiceInstance = &JsonSellerService{
		repo: repo,
	}
	return sellerServiceInstance
}

// SellerService defines the contract for seller service operations with business logic
// SellerService define el contrato para las operaciones de servicio de vendedores con lógica de negocio
type SellerService interface {
	GetAll() ([]models.Seller, error)
	GetById(id int) (models.Seller, error)
	Save(seller models.Seller) ([]models.Seller, error)
	Update(id int, seller models.Seller) ([]models.Seller, error)
	Delete(id int) error
}

// JsonSellerService implements SellerService and contains business logic for seller operations using JSON storage
// JsonSellerService implementa SellerService y contiene la lógica de negocio para operaciones de vendedores usando almacenamiento JSON
type JsonSellerService struct {
	repo repositories.SellerRepository // Repository for seller data access / Repositorio para acceso a datos de vendedores
}

// GetAll retrieves all sellers from the repository
// GetAll recupera todos los vendedores del repositorio
func (s *JsonSellerService) GetAll() ([]models.Seller, error) {
	sellers, err := s.repo.GetAll()
	return sellers, err
}

// GetById retrieves a seller by its ID with business logic to search through all sellers
// Returns ErrNotFound if the seller doesn't exist
// GetById recupera un vendedor por su ID con lógica de negocio para buscar entre todos los vendedores
// Retorna ErrNotFound si el vendedor no existe
func (s *JsonSellerService) GetById(id int) (models.Seller, error) {
	// Get all sellers from repository / Obtener todos los vendedores del repositorio
	sellers, err := s.repo.GetAll()
	if err != nil {
		return models.Seller{}, error_message.ErrNotFound
	}

	// Search for seller with matching ID / Buscar vendedor con ID coincidente
	for _, seller := range sellers {
		if seller.Id == id {
			return seller, nil
		}
	}

	// Return error if seller not found / Retornar error si el vendedor no se encuentra
	return models.Seller{}, error_message.ErrNotFound
}

// Save creates a new seller in the repository
// Save crea un nuevo vendedor en el repositorio
func (s *JsonSellerService) Save(seller models.Seller) ([]models.Seller, error) {
	sellerCreated, err := s.repo.Save(seller)
	return sellerCreated, err
}

// Update modifies an existing seller by ID in the repository
// Update modifica un vendedor existente por ID en el repositorio
func (s *JsonSellerService) Update(id int, seller models.Seller) ([]models.Seller, error) {
	sellerFounded, err := s.repo.Update(id, seller)
	return sellerFounded, err
}

// Delete removes a seller by ID from the repository
// Delete elimina un vendedor por ID del repositorio
func (s *JsonSellerService) Delete(id int) error {
	err := s.repo.Delete(id)
	return err
}
