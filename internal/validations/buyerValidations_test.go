package validations_test

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"github.com/stretchr/testify/assert"
)

func TestValidateBuyerRequestStruct(t *testing.T) {
	t.Run("BuyerRequest doesn't have required fields CardNumberId, FirstName and LastName", func(t *testing.T) {
		expectedErrorMessage := "first_name: cannot be blank; id_card_number: cannot be blank; last_name: cannot be blank."
		buyerRequest := requests.BuyerRequest{}

		err := validations.ValidateBuyerRequestStruct(buyerRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("BuyerRequest doesn't have require field CardNumberId", func(t *testing.T) {
		expectedErrorMessage := "id_card_number: cannot be blank."
		buyerRequest := requests.BuyerRequest{
			FirstName: "Ignacio",
			LastName:  "Garcia",
		}

		err := validations.ValidateBuyerRequestStruct(buyerRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("BuyerRequest doesn't have require field FirstName", func(t *testing.T) {
		expectedErrorMessage := "first_name: cannot be blank."
		buyerRequest := requests.BuyerRequest{
			CardNumberId: "CARD-01",
			LastName:     "Garcia",
		}

		err := validations.ValidateBuyerRequestStruct(buyerRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("BuyerRequest doesn't have require field LastName", func(t *testing.T) {
		expectedErrorMessage := "last_name: cannot be blank."
		buyerRequest := requests.BuyerRequest{
			CardNumberId: "CARD-01",
			FirstName:    "Ignacio",
		}

		err := validations.ValidateBuyerRequestStruct(buyerRequest)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Error(t, err, "err should have an error")
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("BuyerRequest has all fields correctly", func(t *testing.T) {
		buyerRequest := requests.BuyerRequest{
			CardNumberId: "CARD-01",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		err := validations.ValidateBuyerRequestStruct(buyerRequest)

		assert.Nil(t, err, "err should be nil")
	})
}

func TestIsNotAnEmptyBuyer(t *testing.T) {
	t.Run("BuyerRequest is empty returns error", func(t *testing.T) {
		buyerRequest := requests.BuyerRequest{}

		err := validations.IsNotAnEmptyBuyer(buyerRequest)

		assert.Error(t, err, "err should be expected")
	})

	t.Run("BuyerRequest has CardNumberId", func(t *testing.T) {
		buyerRequest := requests.BuyerRequest{
			CardNumberId: "CARD-01",
		}

		err := validations.IsNotAnEmptyBuyer(buyerRequest)

		assert.Nil(t, err, "err should be nil")
	})

	t.Run("BuyerRequest has FirstName", func(t *testing.T) {
		buyerRequest := requests.BuyerRequest{
			FirstName: "Ignacio",
		}

		err := validations.IsNotAnEmptyBuyer(buyerRequest)

		assert.Nil(t, err, "err should be nil")
	})

	t.Run("BuyerRequest has LastName", func(t *testing.T) {
		buyerRequest := requests.BuyerRequest{
			LastName: "Garcia",
		}

		err := validations.IsNotAnEmptyBuyer(buyerRequest)

		assert.Nil(t, err, "err should be nil")
	})

	t.Run("BuyerRequest is complete", func(t *testing.T) {
		buyerRequest := requests.BuyerRequest{
			CardNumberId: "CARD-01",
			FirstName:    "Ignacio",
			LastName:     "Garcia",
		}

		err := validations.IsNotAnEmptyBuyer(buyerRequest)

		assert.Nil(t, err, "err should be nil")
	})
}
