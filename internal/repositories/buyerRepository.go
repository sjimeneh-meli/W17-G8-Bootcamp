package repositories

import (
	"context"
	"fmt"
	"slices"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

type BuyerRepositoryI interface {
	GetAll(ctx context.Context) (map[int]models.Buyer, error)
	GetById(ctx context.Context, id int) (models.Buyer, error)
	DeleteById(ctx context.Context, id int) error
	Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error)
	Update(ctx context.Context, buyerId int, buyer models.Buyer) (models.Buyer, error)

	GetPurchaseOrdersReportByBuyerId(ctx context.Context, buyerId int) (models.PurchaseOrderReport, error)
	GetPurchaseOrdersReport(ctx context.Context) ([]models.PurchaseOrderReport, error)
	GetCardNumberIds() ([]string, error)
}

type BuyerRepository struct {
	storage map[int]models.Buyer
	loader  loader.Storage[models.Buyer]
}

func (r *BuyerRepository) Update(ctx context.Context, id int, buyer models.Buyer) (models.Buyer, error) {
	_, exists := r.storage[id]
	if !exists {
		return models.Buyer{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", id, "doesn't exists.")
	}

	updatedBuyer := r.storage[id]
	if buyer.CardNumberId != "" {
		existingCardNumberIds, _ := r.GetCardNumberIds()
		if slices.Contains(existingCardNumberIds, buyer.CardNumberId) {
			return models.Buyer{}, fmt.Errorf("%w. %s %s %s", error_message.ErrAlreadyExists, "Card number with Id", buyer.CardNumberId, "already exists.")
		}
		updatedBuyer.CardNumberId = buyer.CardNumberId
	}
	if buyer.FirstName != "" {
		updatedBuyer.FirstName = buyer.FirstName
	}
	if buyer.LastName != "" {
		updatedBuyer.LastName = buyer.LastName
	}

	r.storage[id] = updatedBuyer
	err := r.Save()
	if err != nil {
		return models.Buyer{}, err
	}

	return updatedBuyer, nil
}

func (r *BuyerRepository) GetAll(ctx context.Context) (map[int]models.Buyer, error) {
	return r.storage, nil
}

func (r *BuyerRepository) GetById(ctx context.Context, id int) (models.Buyer, error) {
	_, exists := r.storage[id]
	if !exists {
		return models.Buyer{}, fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", id, "doesn't exists.")
	}

	return r.storage[id], nil
}

func (r *BuyerRepository) DeleteById(ctx context.Context, id int) error {
	_, exists := r.storage[id]
	if !exists {
		return fmt.Errorf("%w. %s %d %s", error_message.ErrNotFound, "Buyer with Id", id, "doesn't exists.")
	}
	delete(r.storage, id)

	err := r.Save()
	if err != nil {
		return err
	}
	return nil
}

func (r *BuyerRepository) Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error) {
	newId := r.GetNewId()
	buyer.Id = newId

	_, exists := r.storage[buyer.Id]
	if exists {
		return models.Buyer{}, fmt.Errorf("%w. %s %d %s", error_message.ErrAlreadyExists, "Buyer with Id", buyer.Id, "already exists.")
	}

	r.storage[buyer.Id] = buyer

	err := r.Save()
	if err != nil {
		return models.Buyer{}, fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err.Error())
	}

	return buyer, nil
}

func (r *BuyerRepository) GetNewId() int {
	lastId := 0
	for _, buyer := range r.storage {
		if buyer.Id > lastId {
			lastId = buyer.Id
		}
	}
	return (lastId + 1)
}

func (r *BuyerRepository) Save() error {
	buyerArray := []models.Buyer{}

	for _, buyer := range r.storage {
		buyerArray = append(buyerArray, buyer)
	}

	err := r.loader.WriteAll(buyerArray)
	if err != nil {
		return fmt.Errorf("%w. %s", error_message.ErrInternalServerError, err.Error())
	}
	return nil
}

func GetNewBuyerRepository(loader loader.Storage[models.Buyer]) (BuyerRepositoryI, error) {
	storage, err := loader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w:%v", error_message.ErrInternalServerError, err)
	}

	return &BuyerRepository{
		storage: storage,
		loader:  loader,
	}, nil
}

func (r *BuyerRepository) GetCardNumberIds() ([]string, error) {
	cardNumbers := []string{}
	for _, buyer := range r.storage {
		cardNumbers = append(cardNumbers, buyer.CardNumberId)
	}
	return cardNumbers, nil
}

func (r *BuyerRepository) GetPurchaseOrdersReportByBuyerId(ctx context.Context, buyerId int) (models.PurchaseOrderReport, error) {
	return models.PurchaseOrderReport{}, nil
}

func (r *BuyerRepository) GetPurchaseOrdersReport(ctx context.Context) ([]models.PurchaseOrderReport, error) {
	return []models.PurchaseOrderReport{}, nil
}
