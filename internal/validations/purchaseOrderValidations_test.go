package validations_test

import (
	"testing"
	"time"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"github.com/stretchr/testify/assert"
)

func TestPurchaseOrderValidations_ValidatePurchaseOrderRequestStruct(t *testing.T) {
	t.Run("PurchaseOrderRequest has all fields empty - should return custom error", func(t *testing.T) {
		expectedErrorMessage := "data: cannot be blank. fields order_number, order_date, tracking_code, buyer_id, product_record_id are neccesary inside of data"
		purchaseOrderRequest := requests.PurchaseOrderRequest{
			Data: requests.PurchaseOrderAttributes{},
		}

		err := validations.ValidatePurchaseOrderRequestStruct(purchaseOrderRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("PurchaseOrderRequest doesn't have required field OrderNumber", func(t *testing.T) {
		expectedErrorMessage := "order_number: cannot be blank."
		purchaseOrderRequest := requests.PurchaseOrderRequest{
			Data: requests.PurchaseOrderAttributes{
				OrderDate:       time.Now(),
				TrackingCode:    "TRACK-001",
				BuyerId:         1,
				ProductRecordId: 1,
			},
		}

		err := validations.ValidatePurchaseOrderRequestStruct(purchaseOrderRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("PurchaseOrderRequest doesn't have required field OrderDate", func(t *testing.T) {
		expectedErrorMessage := "order_date: cannot be blank."
		purchaseOrderRequest := requests.PurchaseOrderRequest{
			Data: requests.PurchaseOrderAttributes{
				OrderNumber:     "ORDER-001",
				TrackingCode:    "TRACK-001",
				BuyerId:         1,
				ProductRecordId: 1,
			},
		}

		err := validations.ValidatePurchaseOrderRequestStruct(purchaseOrderRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("PurchaseOrderRequest doesn't have required field TrackingCode", func(t *testing.T) {
		expectedErrorMessage := "tracking_code: cannot be blank."
		purchaseOrderRequest := requests.PurchaseOrderRequest{
			Data: requests.PurchaseOrderAttributes{
				OrderNumber:     "ORDER-001",
				OrderDate:       time.Now(),
				BuyerId:         1,
				ProductRecordId: 1,
			},
		}

		err := validations.ValidatePurchaseOrderRequestStruct(purchaseOrderRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("PurchaseOrderRequest doesn't have required field BuyerId", func(t *testing.T) {
		expectedErrorMessage := "buyer_id: cannot be blank."
		purchaseOrderRequest := requests.PurchaseOrderRequest{
			Data: requests.PurchaseOrderAttributes{
				OrderNumber:     "ORDER-001",
				OrderDate:       time.Now(),
				TrackingCode:    "TRACK-001",
				ProductRecordId: 1,
			},
		}

		err := validations.ValidatePurchaseOrderRequestStruct(purchaseOrderRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("PurchaseOrderRequest doesn't have required field ProductRecordId", func(t *testing.T) {
		expectedErrorMessage := "product_record_id: cannot be blank."
		purchaseOrderRequest := requests.PurchaseOrderRequest{
			Data: requests.PurchaseOrderAttributes{
				OrderNumber:  "ORDER-001",
				OrderDate:    time.Now(),
				TrackingCode: "TRACK-001",
				BuyerId:      1,
			},
		}

		err := validations.ValidatePurchaseOrderRequestStruct(purchaseOrderRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("PurchaseOrderRequest doesn't have multiple required fields", func(t *testing.T) {
		expectedErrorMessage := "buyer_id: cannot be blank; order_number: cannot be blank; product_record_id: cannot be blank."
		purchaseOrderRequest := requests.PurchaseOrderRequest{
			Data: requests.PurchaseOrderAttributes{
				OrderDate:    time.Now(),
				TrackingCode: "TRACK-001",
			},
		}

		err := validations.ValidatePurchaseOrderRequestStruct(purchaseOrderRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("PurchaseOrderRequest has all fields correctly", func(t *testing.T) {
		purchaseOrderRequest := requests.PurchaseOrderRequest{
			Data: requests.PurchaseOrderAttributes{
				OrderNumber:     "ORDER-001",
				OrderDate:       time.Now(),
				TrackingCode:    "TRACK-001",
				BuyerId:         1,
				ProductRecordId: 1,
			},
		}

		err := validations.ValidatePurchaseOrderRequestStruct(purchaseOrderRequest)

		assert.Nil(t, err, "err should be nil")
	})
}
