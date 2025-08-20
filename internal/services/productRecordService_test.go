package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestNewProductRecordService(t *testing.T) {
	// Reset singleton instance for testing
	productRecordServiceInstance = nil
	defer func() { productRecordServiceInstance = nil }()

	t.Run("should create new instance when instance is nil", func(t *testing.T) {
		// Arrange
		mockRepo := tests.GetNewProductRecordRepositoryMock()
		mockProductService := &tests.ProductServiceMock{}

		// Act
		service := NewProductRecordService(mockRepo, mockProductService)

		// Assert
		assert.NotNil(t, service)
		assert.Equal(t, productRecordServiceInstance, service)
	})

	t.Run("should return existing instance when instance already exists", func(t *testing.T) {
		// Arrange
		mockRepo1 := tests.GetNewProductRecordRepositoryMock()
		mockProductService1 := &tests.ProductServiceMock{}
		mockRepo2 := tests.GetNewProductRecordRepositoryMock()
		mockProductService2 := &tests.ProductServiceMock{}

		// Act
		service1 := NewProductRecordService(mockRepo1, mockProductService1)
		service2 := NewProductRecordService(mockRepo2, mockProductService2)

		// Assert
		assert.Equal(t, service1, service2)
	})
}

func TestCreateProductRecord(t *testing.T) {
	// Reset singleton instance for testing
	originalInstance := productRecordServiceInstance
	defer func() { productRecordServiceInstance = originalInstance }()

	ctx := context.Background()

	t.Run("should create product record successfully when product exists", func(t *testing.T) {
		// Arrange
		productRecordServiceInstance = nil
		mockRepo := tests.GetNewProductRecordRepositoryMock()
		mockProductService := &tests.ProductServiceMock{}
		service := NewProductRecordService(mockRepo, mockProductService)

		productRecord := models.ProductRecord{
			ID:             1,
			LastUpdateDate: time.Now(),
			PurchasePrice:  10.50,
			SalePrice:      15.75,
			ProductID:      100,
		}

		expectedRecord := &models.ProductRecord{
			ID:             1,
			LastUpdateDate: productRecord.LastUpdateDate,
			PurchasePrice:  10.50,
			SalePrice:      15.75,
			ProductID:      100,
		}

		mockProductService.On("ExistById", ctx, int64(100)).Return(true, nil)
		mockRepo.On("Create", ctx, &productRecord).Return(expectedRecord, nil)

		// Act
		result, err := service.CreateProductRecord(ctx, productRecord)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedRecord, result)
		mockProductService.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when ProductService.ExistById fails", func(t *testing.T) {
		// Arrange
		productRecordServiceInstance = nil
		mockRepo := tests.GetNewProductRecordRepositoryMock()
		mockProductService := &tests.ProductServiceMock{}
		service := NewProductRecordService(mockRepo, mockProductService)

		productRecord := models.ProductRecord{
			ProductID: 100,
		}

		expectedError := errors.New("database connection error")
		mockProductService.On("ExistById", ctx, int64(100)).Return(false, expectedError)

		// Act
		result, err := service.CreateProductRecord(ctx, productRecord)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, &models.ProductRecord{}, result)
		mockProductService.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should return error when product does not exist", func(t *testing.T) {
		// Arrange
		productRecordServiceInstance = nil
		mockRepo := tests.GetNewProductRecordRepositoryMock()
		mockProductService := &tests.ProductServiceMock{}
		service := NewProductRecordService(mockRepo, mockProductService)

		productRecord := models.ProductRecord{
			ProductID: 100,
		}

		mockProductService.On("ExistById", ctx, int64(100)).Return(false, nil)

		// Act
		result, err := service.CreateProductRecord(ctx, productRecord)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product by id : 100 does not exist")
		assert.True(t, errors.Is(err, error_message.ErrDependencyNotFound))
		assert.Equal(t, &models.ProductRecord{}, result)
		mockProductService.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Create")
	})
}

func TestGetReportByIdProduct(t *testing.T) {
	// Reset singleton instance for testing
	originalInstance := productRecordServiceInstance
	defer func() { productRecordServiceInstance = originalInstance }()

	ctx := context.Background()

	t.Run("should get report successfully when product exists", func(t *testing.T) {
		// Arrange
		productRecordServiceInstance = nil
		mockRepo := tests.GetNewProductRecordRepositoryMock()
		mockProductService := &tests.ProductServiceMock{}
		service := NewProductRecordService(mockRepo, mockProductService)

		productID := int64(100)
		expectedReport := &models.ProductRecordReport{
			ProductId:    100,
			Description:  "Test Product",
			RecordsCount: 5,
		}

		mockProductService.On("ExistById", ctx, productID).Return(true, nil)
		mockRepo.On("GetReportByIdProduct", ctx, productID).Return(expectedReport, nil)

		// Act
		result, err := service.GetReportByIdProduct(ctx, productID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedReport, result)
		mockProductService.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when ProductService.ExistById fails", func(t *testing.T) {
		// Arrange
		productRecordServiceInstance = nil
		mockRepo := tests.GetNewProductRecordRepositoryMock()
		mockProductService := &tests.ProductServiceMock{}
		service := NewProductRecordService(mockRepo, mockProductService)

		productID := int64(100)
		expectedError := errors.New("database connection error")
		mockProductService.On("ExistById", ctx, productID).Return(false, expectedError)

		// Act
		result, err := service.GetReportByIdProduct(ctx, productID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, &models.ProductRecordReport{}, result)
		mockProductService.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "GetReportByIdProduct")
	})

	t.Run("should return error when product does not exist", func(t *testing.T) {
		// Arrange
		productRecordServiceInstance = nil
		mockRepo := tests.GetNewProductRecordRepositoryMock()
		mockProductService := &tests.ProductServiceMock{}
		service := NewProductRecordService(mockRepo, mockProductService)

		productID := int64(100)
		mockProductService.On("ExistById", ctx, productID).Return(false, nil)

		// Act
		result, err := service.GetReportByIdProduct(ctx, productID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product by id : 100 does not exist")
		assert.True(t, errors.Is(err, error_message.ErrDependencyNotFound))
		assert.Equal(t, &models.ProductRecordReport{}, result)
		mockProductService.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "GetReportByIdProduct")
	})
}

func TestGetReport(t *testing.T) {
	// Reset singleton instance for testing
	originalInstance := productRecordServiceInstance
	defer func() { productRecordServiceInstance = originalInstance }()

	ctx := context.Background()

	t.Run("should get all reports successfully", func(t *testing.T) {
		// Arrange
		productRecordServiceInstance = nil
		mockRepo := tests.GetNewProductRecordRepositoryMock()
		mockProductService := &tests.ProductServiceMock{}
		service := NewProductRecordService(mockRepo, mockProductService)

		expectedReports := []*models.ProductRecordReport{
			{
				ProductId:    100,
				Description:  "Product 1",
				RecordsCount: 5,
			},
			{
				ProductId:    200,
				Description:  "Product 2",
				RecordsCount: 3,
			},
		}

		mockRepo.On("GetReport", ctx).Return(expectedReports, nil)

		// Act
		result, err := service.GetReport(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedReports, result)
		mockRepo.AssertExpectations(t)
		// ProductService should not be called for this method
		mockProductService.AssertNotCalled(t, "ExistById")
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		productRecordServiceInstance = nil
		mockRepo := tests.GetNewProductRecordRepositoryMock()
		mockProductService := &tests.ProductServiceMock{}
		service := NewProductRecordService(mockRepo, mockProductService)

		expectedError := errors.New("database query error")
		mockRepo.On("GetReport", ctx).Return([]*models.ProductRecordReport(nil), expectedError)

		// Act
		result, err := service.GetReport(ctx)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
		// ProductService should not be called for this method
		mockProductService.AssertNotCalled(t, "ExistById")
	})
}

func TestExistProductRecordByID(t *testing.T) {
	// Reset singleton instance for testing
	originalInstance := productRecordServiceInstance
	defer func() { productRecordServiceInstance = originalInstance }()

	ctx := context.Background()

	t.Run("should return true when product record exists", func(t *testing.T) {
		// Arrange
		productRecordServiceInstance = nil
		mockRepo := tests.GetNewProductRecordRepositoryMock()
		mockProductService := &tests.ProductServiceMock{}
		service := NewProductRecordService(mockRepo, mockProductService)

		recordID := int64(100)
		mockRepo.On("ExistProductRecordByID", ctx, recordID).Return(true)

		// Act
		result := service.ExistProductRecordByID(ctx, recordID)

		// Assert
		assert.True(t, result)
		mockRepo.AssertExpectations(t)
		// ProductService should not be called for this method
		mockProductService.AssertNotCalled(t, "ExistById")
	})

	t.Run("should return false when product record does not exist", func(t *testing.T) {
		// Arrange
		productRecordServiceInstance = nil
		mockRepo := tests.GetNewProductRecordRepositoryMock()
		mockProductService := &tests.ProductServiceMock{}
		service := NewProductRecordService(mockRepo, mockProductService)

		recordID := int64(100)
		mockRepo.On("ExistProductRecordByID", ctx, recordID).Return(false)

		// Act
		result := service.ExistProductRecordByID(ctx, recordID)

		// Assert
		assert.False(t, result)
		mockRepo.AssertExpectations(t)
		// ProductService should not be called for this method
		mockProductService.AssertNotCalled(t, "ExistById")
	})
}

// Test to ensure the service implements the interface correctly
func TestProductRecordServiceImplementsInterface(t *testing.T) {
	// This test ensures that productRecordService implements ProductRecordServiceI
	var _ ProductRecordServiceI = (*productRecordService)(nil)
}

// Test for constructor with singleton behavior edge cases
func TestNewProductRecordServiceSingletonBehavior(t *testing.T) {
	// Save original instance
	originalInstance := productRecordServiceInstance

	t.Run("should maintain singleton across multiple calls", func(t *testing.T) {
		// Reset singleton
		productRecordServiceInstance = nil

		// Create multiple instances
		mockRepo1 := tests.GetNewProductRecordRepositoryMock()
		mockProductService1 := &tests.ProductServiceMock{}
		service1 := NewProductRecordService(mockRepo1, mockProductService1)

		mockRepo2 := tests.GetNewProductRecordRepositoryMock()
		mockProductService2 := &tests.ProductServiceMock{}
		service2 := NewProductRecordService(mockRepo2, mockProductService2)

		mockRepo3 := tests.GetNewProductRecordRepositoryMock()
		mockProductService3 := &tests.ProductServiceMock{}
		service3 := NewProductRecordService(mockRepo3, mockProductService3)

		// All should be the same instance
		assert.Equal(t, service1, service2)
		assert.Equal(t, service2, service3)
		assert.Equal(t, service1, service3)
	})

	// Restore original instance
	productRecordServiceInstance = originalInstance
}
