package mappers_test

import (
	"testing"
	"time"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_PurchaseOrder_GetModelPurchaseOrderFromRequest(t *testing.T) {
	t.Run("Successfully maps a PurchaseOrderRequest to PurchaseOrder", func(t *testing.T) {
		expectedPurchaseOrder := &models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		purchaseOrderRequest := requests.PurchaseOrderRequest{
			Data: requests.PurchaseOrderAttributes{
				OrderNumber:     "ORDER-001",
				OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				TrackingCode:    "TRACK-001",
				BuyerId:         1,
				ProductRecordId: 1,
			},
		}

		result := mappers.GetModelPurchaseOrderFromRequest(purchaseOrderRequest)

		assert.Equal(t, expectedPurchaseOrder, result)
	})
}

func Test_PurchaseOrder_GetResponsePurchaseOrderFromModel(t *testing.T) {
	t.Run("Successfully maps a PurchaseOrder to PurchaseOrderResponse", func(t *testing.T) {
		purchaseOrderDb := models.PurchaseOrder{
			Id:              10,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		expectedPurchaseOrder := &responses.PurchaseOrderResponse{
			Id:              10,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		result := mappers.GetResponsePurchaseOrderFromModel(&purchaseOrderDb)

		assert.Equal(t, expectedPurchaseOrder, result)
	})
}

func Test_PurchaseOrder_GetListPurchaseOrderResponseFromListModel(t *testing.T) {
	t.Run("Successfully maps a list of PurchaseOrder to a list of PurchaseOrderResponse", func(t *testing.T) {
		purchaseOrdersList := []*models.PurchaseOrder{
			{
				Id:              10,
				OrderNumber:     "ORDER-001",
				OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				TrackingCode:    "TRACK-001",
				BuyerId:         1,
				ProductRecordId: 1,
			},
			{
				Id:              11,
				OrderNumber:     "ORDER-002",
				OrderDate:       time.Date(2024, 1, 16, 14, 45, 0, 0, time.UTC),
				TrackingCode:    "TRACK-002",
				BuyerId:         2,
				ProductRecordId: 2,
			},
		}

		expectedResponsePurchaseOrderList := []*responses.PurchaseOrderResponse{
			{
				Id:              10,
				OrderNumber:     "ORDER-001",
				OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				TrackingCode:    "TRACK-001",
				BuyerId:         1,
				ProductRecordId: 1,
			},
			{
				Id:              11,
				OrderNumber:     "ORDER-002",
				OrderDate:       time.Date(2024, 1, 16, 14, 45, 0, 0, time.UTC),
				TrackingCode:    "TRACK-002",
				BuyerId:         2,
				ProductRecordId: 2,
			},
		}

		result := mappers.GetListPurchaseOrderResponseFromListModel(purchaseOrdersList)

		assert.Equal(t, expectedResponsePurchaseOrderList, result)
	})
}
