// Package services_test - Purchase Order Service Unit Tests / Tests Unitarios del Servicio Purchase Order
// Comprehensive business logic testing for purchase order management operations
// Testing comprehensivo de lógica de negocio para operaciones de gestión de órdenes de compra
package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPurchaseOrderService_GetAll - Tests for GetAll method / Tests para método GetAll
// Cases: successful retrieval, empty list, repository error
// Casos: recuperación exitosa, lista vacía, error de repositorio
func TestPurchaseOrderService_GetAll(t *testing.T) {
	t.Run("GetAll successfully returns all purchase orders", func(t *testing.T) {
		// Test: Repository returns data successfully / Repositorio retorna datos exitosamente
		// Arrange
		ctx := context.Background()

		mockOrders := map[int]models.PurchaseOrder{
			1: {
				Id:              1,
				OrderNumber:     "ORDER-001",
				OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				TrackingCode:    "TRACK-001",
				BuyerId:         1,
				ProductRecordId: 1,
			},
			2: {
				Id:              2,
				OrderNumber:     "ORDER-002",
				OrderDate:       time.Date(2024, 1, 16, 14, 45, 0, 0, time.UTC),
				TrackingCode:    "TRACK-002",
				BuyerId:         2,
				ProductRecordId: 2,
			},
		}
		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("GetAll", ctx).Return(mockOrders, nil).Once()

		// Act
		result, err := service.GetAll(ctx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, mockOrders, result)
		mockPurchaseOrderRepo.AssertExpectations(t)
	})

	t.Run("GetAll with empty list returns zero elements", func(t *testing.T) {
		// Test: Empty repository returns empty list / Repositorio vacío retorna lista vacía
		// Arrange
		ctx := context.Background()
		emptyOrders := map[int]models.PurchaseOrder{}
		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("GetAll", ctx).Return(emptyOrders, nil).Once()

		// Act
		result, err := service.GetAll(ctx)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 0, len(result))
		assert.Equal(t, emptyOrders, result)
		mockPurchaseOrderRepo.AssertExpectations(t)
	})

	t.Run("GetAll fails because of repository error", func(t *testing.T) {
		// Test: Repository error propagation / Propagación de error del repositorio
		// Arrange
		ctx := context.Background()
		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("GetAll", ctx).Return(map[int]models.PurchaseOrder{}, error_message.ErrInternalServerError).Once()

		// Act
		result, err := service.GetAll(ctx)

		// Assert
		require.Error(t, err)
		assert.Equal(t, error_message.ErrInternalServerError, err)
		assert.Equal(t, 0, len(result))
		mockPurchaseOrderRepo.AssertExpectations(t)
	})
}

// TestPurchaseOrderService_GetPurchaseOrdersReport - Tests for GetPurchaseOrdersReport method / Tests para método GetPurchaseOrdersReport
// Cases: successful retrieval for specific buyer, successful retrieval for all buyers, repository errors
// Casos: recuperación exitosa para comprador específico, recuperación exitosa para todos los compradores, errores de repositorio
func TestPurchaseOrderService_GetPurchaseOrdersReport(t *testing.T) {
	t.Run("GetPurchaseOrdersReport successfully returns report for specific buyer", func(t *testing.T) {
		// Test: Valid buyer ID returns specific buyer report / ID de comprador válido retorna reporte específico del comprador
		// Arrange
		ctx := context.Background()
		buyerId := 1
		expectedReport := models.PurchaseOrderReport{
			Id:                 1,
			IdCardNumber:       "CARD-001",
			FirstName:          "Ignacio",
			LastName:           "Garcia",
			PurchaseOrderCount: 3,
		}
		expectedReports := []models.PurchaseOrderReport{expectedReport}

		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("GetPurchaseOrdersReportByBuyerId", ctx, buyerId).Return(expectedReport, nil).Once()

		// Act
		result, err := service.GetPurchaseOrdersReport(ctx, &buyerId)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, expectedReports, result)
		mockPurchaseOrderRepo.AssertExpectations(t)
	})

	t.Run("GetPurchaseOrdersReport successfully returns reports for all buyers", func(t *testing.T) {
		// Test: No buyer ID provided returns all buyers reports / Sin ID de comprador retorna reportes de todos los compradores
		// Arrange
		ctx := context.Background()
		expectedReports := []models.PurchaseOrderReport{
			{
				Id:                 1,
				IdCardNumber:       "CARD-001",
				FirstName:          "Ignacio",
				LastName:           "Garcia",
				PurchaseOrderCount: 3,
			},
			{
				Id:                 2,
				IdCardNumber:       "CARD-002",
				FirstName:          "Jesus",
				LastName:           "Ortega",
				PurchaseOrderCount: 2,
			},
		}

		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("GetAllPurchaseOrdersReports", ctx).Return(expectedReports, nil).Once()

		// Act
		result, err := service.GetPurchaseOrdersReport(ctx, nil)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, expectedReports, result)
		mockPurchaseOrderRepo.AssertExpectations(t)
	})

	t.Run("GetPurchaseOrdersReport fails because of repository error for specific buyer", func(t *testing.T) {
		// Test: Repository error when getting specific buyer report / Error de repositorio al obtener reporte específico del comprador
		// Arrange
		ctx := context.Background()
		buyerId := 999
		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("GetPurchaseOrdersReportByBuyerId", ctx, buyerId).Return(models.PurchaseOrderReport{}, error_message.ErrNotFound).Once()

		// Act
		result, err := service.GetPurchaseOrdersReport(ctx, &buyerId)

		// Assert
		require.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)
		assert.Equal(t, 0, len(result))
		mockPurchaseOrderRepo.AssertExpectations(t)
	})

	t.Run("GetPurchaseOrdersReport fails because of repository error for all buyers", func(t *testing.T) {
		// Test: Repository error when getting all buyers reports / Error de repositorio al obtener reportes de todos los compradores
		// Arrange
		ctx := context.Background()
		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("GetAllPurchaseOrdersReports", ctx).Return([]models.PurchaseOrderReport{}, error_message.ErrInternalServerError).Once()

		// Act
		result, err := service.GetPurchaseOrdersReport(ctx, nil)

		// Assert
		require.Error(t, err)
		assert.Equal(t, error_message.ErrInternalServerError, err)
		assert.Equal(t, 0, len(result))
		mockPurchaseOrderRepo.AssertExpectations(t)
	})
}

// TestPurchaseOrderService_Create - Tests for Create method with comprehensive validation
// Tests para método Create con validación comprehensiva
// Cases: order number exists, buyer doesn't exist, product record doesn't exist, successful creation, repository errors
// Casos: número de orden existe, comprador no existe, registro de producto no existe, creación exitosa, errores de repositorio
func TestPurchaseOrderService_Create(t *testing.T) {
	t.Run("Create fails because order number already exists", func(t *testing.T) {
		// Test: Duplicate order number validation / Validación de número de orden duplicado
		// Arrange
		ctx := context.Background()
		newOrder := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("ExistPurchaseOrderByOrderNumber", ctx, "ORDER-001").Return(true, nil).Once()

		// Act
		result, err := service.Create(ctx, newOrder)

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, error_message.ErrAlreadyExists)
		assert.Equal(t, models.PurchaseOrder{}, result)
		mockPurchaseOrderRepo.AssertExpectations(t)
	})

	t.Run("Create fails because of repository error checking order number existence", func(t *testing.T) {
		// Test: Repository error when checking order number existence / Error de repositorio al verificar existencia de número de orden
		// Arrange
		ctx := context.Background()
		newOrder := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("ExistPurchaseOrderByOrderNumber", ctx, "ORDER-001").Return(false, error_message.ErrInternalServerError).Once()

		// Act
		result, err := service.Create(ctx, newOrder)

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, error_message.ErrInternalServerError)
		assert.Equal(t, models.PurchaseOrder{}, result)
		mockPurchaseOrderRepo.AssertExpectations(t)
	})

	t.Run("Create fails because buyer doesn't exist", func(t *testing.T) {
		// Test: Buyer validation failure / Falla en validación de comprador
		// Arrange
		ctx := context.Background()
		newOrder := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         999,
			ProductRecordId: 1,
		}

		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("ExistPurchaseOrderByOrderNumber", ctx, "ORDER-001").Return(false, nil).Once()
		mockBuyerRepo.On("ExistBuyerById", ctx, 999).Return(false, nil).Once()

		// Act
		result, err := service.Create(ctx, newOrder)

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, error_message.ErrNotFound)
		assert.Equal(t, models.PurchaseOrder{}, result)
		mockPurchaseOrderRepo.AssertExpectations(t)
		mockBuyerRepo.AssertExpectations(t)
	})

	t.Run("Create fails because of repository error checking buyer existence", func(t *testing.T) {
		// Test: Repository error when checking buyer existence / Error de repositorio al verificar existencia de comprador
		// Arrange
		ctx := context.Background()
		newOrder := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("ExistPurchaseOrderByOrderNumber", ctx, "ORDER-001").Return(false, nil).Once()
		mockBuyerRepo.On("ExistBuyerById", ctx, 1).Return(false, error_message.ErrInternalServerError).Once()

		// Act
		result, err := service.Create(ctx, newOrder)

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, error_message.ErrInternalServerError)
		assert.Equal(t, models.PurchaseOrder{}, result)
		mockPurchaseOrderRepo.AssertExpectations(t)
		mockBuyerRepo.AssertExpectations(t)
	})

	t.Run("Create fails because product record doesn't exist", func(t *testing.T) {
		// Test: Product record validation failure / Falla en validación de registro de producto
		// Arrange
		ctx := context.Background()
		newOrder := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 999,
		}

		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("ExistPurchaseOrderByOrderNumber", ctx, "ORDER-001").Return(false, nil).Once()
		mockBuyerRepo.On("ExistBuyerById", ctx, 1).Return(true, nil).Once()
		mockProductRecordRepo.On("ExistProductRecordByID", ctx, int64(999)).Return(false).Once()

		// Act
		result, err := service.Create(ctx, newOrder)

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, error_message.ErrNotFound)
		assert.Equal(t, models.PurchaseOrder{}, result)
		mockPurchaseOrderRepo.AssertExpectations(t)
		mockBuyerRepo.AssertExpectations(t)
		mockProductRecordRepo.AssertExpectations(t)
	})

	t.Run("Create successfully creates a new purchase order", func(t *testing.T) {
		// Test: Successful purchase order creation / Creación exitosa de orden de compra
		// Arrange
		ctx := context.Background()
		newOrder := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 1,
		}
		expectedOrder := models.PurchaseOrder{
			Id:              17,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("ExistPurchaseOrderByOrderNumber", ctx, "ORDER-001").Return(false, nil).Once()
		mockBuyerRepo.On("ExistBuyerById", ctx, 1).Return(true, nil).Once()
		mockProductRecordRepo.On("ExistProductRecordByID", ctx, int64(1)).Return(true).Once()
		mockPurchaseOrderRepo.On("Create", ctx, newOrder).Return(expectedOrder, nil).Once()

		// Act
		result, err := service.Create(ctx, newOrder)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedOrder, result)
		assert.Equal(t, 17, result.Id)
		mockPurchaseOrderRepo.AssertExpectations(t)
		mockBuyerRepo.AssertExpectations(t)
		mockProductRecordRepo.AssertExpectations(t)
	})

	t.Run("Create fails because of repository error during creation", func(t *testing.T) {
		// Test: Repository error during creation / Error de repositorio durante creación
		// Arrange
		ctx := context.Background()
		newOrder := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-001",
			OrderDate:       time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			TrackingCode:    "TRACK-001",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()
		service := services.PurchaseOrderService{
			PurchaseOrderRepository: mockPurchaseOrderRepo,
			BuyerRepository:         mockBuyerRepo,
			ProductRecordRepository: mockProductRecordRepo,
		}
		mockPurchaseOrderRepo.On("ExistPurchaseOrderByOrderNumber", ctx, "ORDER-001").Return(false, nil).Once()
		mockBuyerRepo.On("ExistBuyerById", ctx, 1).Return(true, nil).Once()
		mockProductRecordRepo.On("ExistProductRecordByID", ctx, int64(1)).Return(true).Once()
		mockPurchaseOrderRepo.On("Create", ctx, newOrder).Return(models.PurchaseOrder{}, error_message.ErrInternalServerError).Once()

		// Act
		result, err := service.Create(ctx, newOrder)

		// Assert
		require.Error(t, err)
		assert.Equal(t, error_message.ErrInternalServerError, err)
		assert.Equal(t, models.PurchaseOrder{}, result)
		mockPurchaseOrderRepo.AssertExpectations(t)
		mockBuyerRepo.AssertExpectations(t)
		mockProductRecordRepo.AssertExpectations(t)
	})
}

// TestGetPurchaseOrderService - Tests for service singleton pattern / Tests para patrón singleton del servicio
func TestPurchaseOrderService_GetPurchaseOrderService(t *testing.T) {
	t.Run("GetPurchaseOrderService returns same instance when called multiple times", func(t *testing.T) {
		// Test: Singleton pattern validation / Validación de patrón singleton
		mockPurchaseOrderRepo := tests.GetNewPurchaseOrderRepositoryMock()
		mockBuyerRepo := tests.GetNewBuyerRepositoryMock()
		mockProductRecordRepo := tests.GetNewProductRecordRepositoryMock()

		service1 := services.GetPurchaseOrderService(mockPurchaseOrderRepo, mockBuyerRepo, mockProductRecordRepo)
		service2 := services.GetPurchaseOrderService(nil, nil, nil)

		assert.Equal(t, service1, service2)
	})
}
