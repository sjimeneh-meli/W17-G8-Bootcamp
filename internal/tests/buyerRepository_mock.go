package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

type BuyerRepositoryMock struct {
	mock.Mock
}

func GetNewBuyerRepositoryMock() *BuyerRepositoryMock {
	return &BuyerRepositoryMock{}
}

func (r *BuyerRepositoryMock) GetAll(ctx context.Context) (map[int]models.Buyer, error) {
	args := r.Called(ctx)
	return args.Get(0).(map[int]models.Buyer), args.Error(1)
}

func (r *BuyerRepositoryMock) GetById(ctx context.Context, id int) (models.Buyer, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (r *BuyerRepositoryMock) DeleteById(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func (r *BuyerRepositoryMock) Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error) {
	args := r.Called(ctx, buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (r *BuyerRepositoryMock) Update(ctx context.Context, buyerId int, buyer models.Buyer) (models.Buyer, error) {
	args := r.Called(ctx, buyerId, buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (r *BuyerRepositoryMock) GetCardNumberIds() ([]string, error) {
	args := r.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (r *BuyerRepositoryMock) ExistBuyerById(ctx context.Context, buyerId int) (bool, error) {
	args := r.Called(ctx, buyerId)
	return args.Bool(0), args.Error(1)
}
