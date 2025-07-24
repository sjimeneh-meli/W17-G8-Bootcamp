package services

import (
	"context"
	"fmt"
	"slices"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

var buyerServiceInstance BuyerServiceI

// GetBuyerService - Creates and returns a new instance of BuyerService with the required repository using singleton pattern
// GetBuyerService - Crea y retorna una nueva instancia de BuyerService con el repositorio requerido usando patrón singleton
func GetBuyerService(repo repositories.BuyerRepositoryI) BuyerServiceI {
	if buyerServiceInstance != nil {
		return buyerServiceInstance
	}

	buyerServiceInstance = &BuyerService{
		Repository: repo,
	}
	return buyerServiceInstance
}

// BuyerServiceI - Interface defining the contract for buyer service operations with business logic
// BuyerServiceI - Interfaz que define el contrato para las operaciones del servicio de compradores con lógica de negocio
type BuyerServiceI interface {
	// GetAll - Retrieves all buyers from the system
	// GetAll - Obtiene todos los compradores del sistema
	GetAll(ctx context.Context) (map[int]models.Buyer, error)

	// GetById - Retrieves a specific buyer by their ID
	// GetById - Obtiene un comprador específico por su ID
	GetById(ctx context.Context, id int) (models.Buyer, error)

	// DeleteById - Removes a buyer from the system by their ID
	// DeleteById - Elimina un comprador del sistema por su ID
	DeleteById(ictx context.Context, d int) error

	// Create - Creates a new buyer with business validation (card number uniqueness)
	// Create - Crea un nuevo comprador con validación de negocio (unicidad del número de tarjeta)
	Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error)

	// Update - Updates an existing buyer with business validation (card number uniqueness)
	// Update - Actualiza un comprador existente con validación de negocio (unicidad del número de tarjeta)
	Update(ctx context.Context, buyerId int, buyer models.Buyer) (models.Buyer, error)
}

// BuyerService - Implementation of BuyerServiceI containing business logic for buyer operations
// BuyerService - Implementación de BuyerServiceI que contiene la lógica de negocio para operaciones de compradores
type BuyerService struct {
	Repository repositories.BuyerRepositoryI // Repository dependency for data access / Dependencia del repositorio para acceso a datos
}

// GetAll - Delegates retrieving all buyers to the repository
// GetAll - Delega la obtención de todos los compradores al repositorio
func (s *BuyerService) GetAll(ctx context.Context) (map[int]models.Buyer, error) {
	return s.Repository.GetAll(ctx)
}

// GetById - Delegates retrieving a buyer by their ID to the repository
// GetById - Delega la obtención de un comprador por su ID al repositorio
func (s *BuyerService) GetById(ctx context.Context, id int) (models.Buyer, error) {
	return s.Repository.GetById(ctx, id)
}

// DeleteById - Delegates removing a buyer from the repository by their ID
// DeleteById - Delega la eliminación de un comprador del repositorio por su ID
func (s *BuyerService) DeleteById(ctx context.Context, id int) error {
	return s.Repository.DeleteById(ctx, id)
}

// Create - Creates a new buyer with business validation to ensure card number uniqueness
// Create - Crea un nuevo comprador con validación de negocio para asegurar la unicidad del número de tarjeta
func (s *BuyerService) Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error) {
	// Business validation: Get all existing card numbers to check for duplicates
	// Validación de negocio: Obtener todos los números de tarjeta existentes para verificar duplicados
	existingCardNumbers, err := s.Repository.GetCardNumberIds()
	if err != nil {
		return models.Buyer{}, err
	}

	// Business rule: Card number must be unique across all buyers
	// Regla de negocio: El número de tarjeta debe ser único entre todos los compradores
	if slices.Contains(existingCardNumbers, buyer.CardNumberId) {
		return models.Buyer{}, fmt.Errorf("%w - %s %s %s", error_message.ErrAlreadyExists, "card number with id:", buyer.CardNumberId, "already exists.")
	}

	// If validation passes, delegate to repository for persistence
	// Si la validación pasa, delegar al repositorio para la persistencia
	return s.Repository.Create(ctx, buyer)
}

// Update - Updates an existing buyer with business validation to ensure card number uniqueness
// Update - Actualiza un comprador existente con validación de negocio para asegurar la unicidad del número de tarjeta
func (s *BuyerService) Update(ctx context.Context, id int, buyer models.Buyer) (models.Buyer, error) {
	// Business validation: Get all existing card numbers to check for duplicates
	// Validación de negocio: Obtener todos los números de tarjeta existentes para verificar duplicados
	existingCardNumbers, err := s.Repository.GetCardNumberIds()
	if err != nil {
		return models.Buyer{}, err
	}

	// Business rule: Card number must be unique across all buyers (excluding current buyer)
	// Regla de negocio: El número de tarjeta debe ser único entre todos los compradores (excluyendo el comprador actual)
	if slices.Contains(existingCardNumbers, buyer.CardNumberId) {
		return models.Buyer{}, fmt.Errorf("%w - %s %s %s", error_message.ErrAlreadyExists, "card number with id:", buyer.CardNumberId, "already exists.")
	}

	// If validation passes, delegate to repository for persistence
	// Si la validación pasa, delegar al repositorio para la persistencia
	return s.Repository.Update(ctx, id, buyer)
}
