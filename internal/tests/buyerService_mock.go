package tests

import (
	"context"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/mock"
)

type BuyerServiceMock struct {
	mock.Mock
}

func GetNewBuyerServiceMock() *BuyerServiceMock {
	return &BuyerServiceMock{}
}

func (r *BuyerServiceMock) GetAll(ctx context.Context) (map[int]models.Buyer, error) {
	args := r.Called(ctx)
	return args.Get(0).(map[int]models.Buyer), args.Error(1)
}

func (r *BuyerServiceMock) GetById(ctx context.Context, id int) (models.Buyer, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (r *BuyerServiceMock) DeleteById(ictx context.Context, id int) error {
	args := r.Called(ictx, id)
	return args.Error(0)
}

func (r *BuyerServiceMock) Create(ctx context.Context, buyer models.Buyer) (models.Buyer, error) {
	args := r.Called(ctx, buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (r *BuyerServiceMock) Update(ctx context.Context, buyerId int, buyer models.Buyer) (models.Buyer, error) {
	args := r.Called(ctx, buyerId, buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}
