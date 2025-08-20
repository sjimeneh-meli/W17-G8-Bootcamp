package tests

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/stretchr/testify/mock"
)

func GetProductBatchValidationMock() *ProductBatchValidationMock {
	return &ProductBatchValidationMock{}
}

type ProductBatchValidationMock struct {
	mock.Mock
}

func (m *ProductBatchValidationMock) ValidateProductBatchRequestStruc(r requests.ProductBatchRequest) error {
	args := m.Called(r)
	return args.Error(0)
}
