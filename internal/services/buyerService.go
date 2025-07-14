package services

import (
	"context"
	"fmt"
	"slices"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
)

func GetBuyerService(repo repositories.BuyerRepositoryI) BuyerServiceI {
	return &BuyerService{
		repository: repo,
	}
}

type BuyerServiceI interface {
	GetAll(ctx context.Context) (map[int]models.Buyer, error)
	GetById(ctx context.Context, id int) (models.Buyer, error)
	DeleteById(ictx context.Context, d int) error
	Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error)
	Update(ctx context.Context, buyerId int, buyer models.Buyer) (models.Buyer, error)

	GetPurchaseOrdersReport(ctx context.Context, id *int) ([]models.PurchaseOrderReport, error)
}

type BuyerService struct {
	repository repositories.BuyerRepositoryI
}

func (s *BuyerService) GetAll(ctx context.Context) (map[int]models.Buyer, error) {
	return s.repository.GetAll(ctx)
}

func (s *BuyerService) GetById(ctx context.Context, id int) (models.Buyer, error) {
	return s.repository.GetById(ctx, id)
}

func (s *BuyerService) DeleteById(ctx context.Context, id int) error {
	return s.repository.DeleteById(ctx, id)
}

func (s *BuyerService) Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error) {

	existingCardNumbers, err := s.repository.GetCardNumberIds()
	if err != nil {
		return models.Buyer{}, err
	}
	if slices.Contains(existingCardNumbers, buyer.CardNumberId) {
		return models.Buyer{}, fmt.Errorf("%w - %s %s %s", error_message.ErrAlreadyExists, "card number with id:", buyer.CardNumberId, "already exists.")
	}

	return s.repository.Create(ctx, buyer)
}

func (s *BuyerService) Update(ctx context.Context, id int, buyer models.Buyer) (models.Buyer, error) {
	existingCardNumbers, err := s.repository.GetCardNumberIds()
	if err != nil {
		return models.Buyer{}, err
	}
	if slices.Contains(existingCardNumbers, buyer.CardNumberId) {
		return models.Buyer{}, fmt.Errorf("%w - %s %s %s", error_message.ErrAlreadyExists, "card number with id:", buyer.CardNumberId, "already exists.")
	}
	return s.repository.Update(ctx, id, buyer)
}

func (s *BuyerService) GetPurchaseOrdersReport(ctx context.Context, id *int) ([]models.PurchaseOrderReport, error) {
	if id != nil {
		reports := []models.PurchaseOrderReport{}
		report, err := s.repository.GetPurchaseOrdersReportByBuyerId(ctx, *id)
		if err != nil {
			return reports, err
		}
		reports = append(reports, report)
		return reports, nil
	}
	return s.repository.GetPurchaseOrdersReport(ctx)
	//Implementar el repository de múltiples reports
}
