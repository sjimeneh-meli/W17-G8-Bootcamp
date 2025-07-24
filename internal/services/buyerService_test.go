package services_test

import (
	"context"
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuyerService_GetAll(t *testing.T) {
	t.Run("GetAll successfully returns all buyers", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		mockBuyers := map[int]models.Buyer{
			1: {Id: 1, CardNumberId: "CARD-1001", FirstName: "Juan", LastName: "Pérez"},
			2: {Id: 2, CardNumberId: "CARD-1002", FirstName: "María", LastName: "Gómez"},
			3: {Id: 3, CardNumberId: "CARD-1003", FirstName: "Carlos", LastName: "López"},
		}
		mockRepo := tests.GetNewBuyerRepositoryMock()
		service := services.BuyerService{
			Repository: mockRepo,
		}
		mockRepo.On("GetAll", ctx).Return(mockBuyers, nil).Once()

		// Act

		result, err := service.GetAll(ctx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 3, len(result))
		assert.Equal(t, mockBuyers, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetAll with empty list returns zero elements", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		emptyBuyers := map[int]models.Buyer{}
		mockRepo := tests.GetNewBuyerRepositoryMock()
		service := services.BuyerService{
			Repository: mockRepo,
		}
		mockRepo.On("GetAll", ctx).Return(emptyBuyers, nil).Once()

		// Act
		result, err := service.GetAll(ctx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 0, len(result))
		assert.Equal(t, emptyBuyers, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetAll fails because of repository error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		mockRepo := tests.GetNewBuyerRepositoryMock()
		service := services.BuyerService{
			Repository: mockRepo,
		}
		mockRepo.On("GetAll", ctx).Return(map[int]models.Buyer{}, error_message.ErrInternalServerError).Once()

		// Act
		result, err := service.GetAll(ctx)

		// Assert
		require.Error(t, err)
		assert.Equal(t, error_message.ErrInternalServerError, err)
		assert.Equal(t, 0, len(result))
		mockRepo.AssertExpectations(t)
	})
}

func TestBuyerService_GetById(t *testing.T) {
	t.Run("GetById successfully returns buyer information", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		expectedBuyer := models.Buyer{
			Id:           1,
			CardNumberId: "CARD-1001",
			FirstName:    "Juan",
			LastName:     "Pérez",
		}
		mockRepo := tests.GetNewBuyerRepositoryMock()
		service := services.BuyerService{
			Repository: mockRepo,
		}
		mockRepo.On("GetById", ctx, 1).Return(expectedBuyer, nil).Once()

		// Act
		result, err := service.GetById(ctx, 1)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedBuyer, result)
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, "CARD-1001", result.CardNumberId)
		assert.Equal(t, "Juan", result.FirstName)
		assert.Equal(t, "Pérez", result.LastName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetById fails because buyer doesn't exist", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		mockRepo := tests.GetNewBuyerRepositoryMock()
		service := services.BuyerService{
			Repository: mockRepo,
		}
		mockRepo.On("GetById", ctx, 999).Return(models.Buyer{}, error_message.ErrNotFound).Once()

		// Act
		result, err := service.GetById(ctx, 999)

		// Assert
		require.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)
		assert.Equal(t, models.Buyer{}, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetById fails because of repository error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		mockRepo := tests.GetNewBuyerRepositoryMock()
		service := services.BuyerService{
			Repository: mockRepo,
		}
		mockRepo.On("GetById", ctx, 1).Return(models.Buyer{}, error_message.ErrInternalServerError).Once()

		// Act
		result, err := service.GetById(ctx, 1)

		// Assert
		require.Error(t, err)
		assert.Equal(t, error_message.ErrInternalServerError, err)
		assert.Equal(t, models.Buyer{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestBuyerService_DeleteById(t *testing.T) {
	t.Run("DeleteById successfully deletes buyer", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		mockRepo := tests.GetNewBuyerRepositoryMock()
		service := services.BuyerService{
			Repository: mockRepo,
		}
		mockRepo.On("DeleteById", ctx, 1).Return(nil).Once()

		// Act
		err := service.DeleteById(ctx, 1)

		// Assert
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteById fails because buyer doesn't exist", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		mockRepo := tests.GetNewBuyerRepositoryMock()
		service := services.BuyerService{
			Repository: mockRepo,
		}
		mockRepo.On("DeleteById", ctx, 999).Return(error_message.ErrNotFound).Once()

		// Act
		err := service.DeleteById(ctx, 999)

		// Assert
		require.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteById fails because of repository error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		mockRepo := tests.GetNewBuyerRepositoryMock()
		service := services.BuyerService{
			Repository: mockRepo,
		}
		mockRepo.On("DeleteById", ctx, 1).Return(error_message.ErrInternalServerError).Once()

		// Act
		err := service.DeleteById(ctx, 1)

		// Assert
		require.Error(t, err)
		assert.Equal(t, error_message.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetBuyerService(t *testing.T) {
	t.Run("GetBuyerService returns same instance when called multiple times", func(t *testing.T) {
		mockRepository := tests.GetNewBuyerRepositoryMock()
		service1 := services.GetBuyerService(mockRepository)
		service2 := services.GetBuyerService(nil)

		assert.Equal(t, service1, service2)
	})
}
