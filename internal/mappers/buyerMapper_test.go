package mappers_test

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetModelBuyerFromRequest(t *testing.T) {
	t.Run("Successfully maps a BuyerRequest to Buyer", func(t *testing.T) {
		expectedBuyer := &models.Buyer{
			Id:           0,
			CardNumberId: "CARD-01",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		buyerRequest := requests.BuyerRequest{
			CardNumberId: "CARD-01",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		result := mappers.GetModelBuyerFromRequest(buyerRequest)

		assert.Equal(t, expectedBuyer, result)
	})
}

func TestGetResponseBuyerFromModel(t *testing.T) {
	t.Run("Successfully maps a Buyer to BuyerResponse", func(t *testing.T) {
		buyerDb := models.Buyer{
			Id:           10,
			CardNumberId: "CARD-01",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		expectedBuyer := &responses.BuyerResponse{
			Id:           10,
			CardNumberId: "CARD-01",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		result := mappers.GetResponseBuyerFromModel(&buyerDb)

		assert.Equal(t, expectedBuyer, result)
	})
}

func TestTestGetListBuyerResponseFromListModel(t *testing.T) {
	t.Run("Successfully maps a list of Buyer to a list of BuyerResponse", func(t *testing.T) {
		buyersList := []*models.Buyer{{
			Id:           10,
			CardNumberId: "CARD-01",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		},
			{
				Id:           11,
				CardNumberId: "CARD-02",
				FirstName:    "Nahuel",
				LastName:     "Gomez",
			},
		}
		expectedResponseBuyerList := []*responses.BuyerResponse{
			{
				Id:           10,
				CardNumberId: "CARD-01",
				FirstName:    "Ignacio",
				LastName:     "Garcia",
			},
			{
				Id:           11,
				CardNumberId: "CARD-02",
				FirstName:    "Nahuel",
				LastName:     "Gomez",
			},
		}

		result := mappers.GetListBuyerResponseFromListModel(buyersList)

		assert.Equal(t, expectedResponseBuyerList, result)

	})

}
