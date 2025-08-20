// Package services - Locality Service Unit Tests / Tests Unitarios del Servicio Locality
// Business logic testing for locality management with country/province creation and seller reporting
// Testing de lógica de negocio para gestión de localidades con creación de país/provincia y reportes de vendedores
package services

import (
	"context"
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
)

// setupLocalityServiceTest - Helper function to setup test dependencies
// setupLocalityServiceTest - Función helper para configurar dependencias de test
// Creates mock repository and service instance for testing
// Crea repositorio mock e instancia de servicio para testing
func setupLocalityServiceTest() (*tests.LocalityRepositoryMock, LocalityService) {
	mockLocalityRepo := &tests.LocalityRepositoryMock{}
	// Create direct service instance to avoid singleton
	// Crear una instancia directa del service para evitar el singleton
	service := &SQLLocalityService{repo: mockLocalityRepo}
	return mockLocalityRepo, service
}

// TestLocalityService_Save - Tests for Save method / Tests para método Save
// Cases: successful creation, repository error
// Casos: creación exitosa, error de repositorio
func TestLocalityService_Save(t *testing.T) {
	t.Run("save_ok - Creates locality successfully when repository works", func(t *testing.T) {
		// Test: Valid locality data saves successfully through service delegation / Datos válidos de localidad se guardan exitosamente a través de delegación del service
		mockRepo, service := setupLocalityServiceTest()

		inputLocality := models.Locality{
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		expectedLocality := models.Locality{
			Id:           1,
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		// Mock repository save success / Mock de éxito en save del repository
		mockRepo.SaveFunc = func(ctx context.Context, locality models.Locality) (models.Locality, error) {
			assert.Equal(t, inputLocality, locality)
			return expectedLocality, nil
		}

		result, err := service.Save(context.Background(), inputLocality)

		assert.NoError(t, err)
		assert.Equal(t, expectedLocality, result)
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, "Test City", result.LocalityName)
		assert.Equal(t, "Test Province", result.ProvinceName)
		assert.Equal(t, "Test Country", result.CountryName)
	})

	t.Run("save_error - Returns error when repository save fails", func(t *testing.T) {
		// Test: Service propagates repository errors correctly / Service propaga errores del repository correctamente
		mockRepo, service := setupLocalityServiceTest()

		inputLocality := models.Locality{
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		// Mock repository save error / Mock de error en save del repository
		mockRepo.SaveFunc = func(ctx context.Context, locality models.Locality) (models.Locality, error) {
			assert.Equal(t, inputLocality, locality)
			return models.Locality{}, error_message.ErrAlreadyExists
		}

		result, err := service.Save(context.Background(), inputLocality)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrAlreadyExists, err)
		assert.Equal(t, models.Locality{}, result)
	})

	t.Run("save_error_query - Returns error when repository query fails", func(t *testing.T) {
		// Test: Service handles repository query errors / Service maneja errores de consulta del repository
		mockRepo, service := setupLocalityServiceTest()

		inputLocality := models.Locality{
			LocalityName: "Test City",
			ProvinceName: "Test Province",
			CountryName:  "Test Country",
		}

		// Mock repository query error / Mock de error de consulta del repository
		mockRepo.SaveFunc = func(ctx context.Context, locality models.Locality) (models.Locality, error) {
			assert.Equal(t, inputLocality, locality)
			return models.Locality{}, error_message.ErrQuery
		}

		result, err := service.Save(context.Background(), inputLocality)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQuery, err)
		assert.Equal(t, models.Locality{}, result)
	})
}

// TestLocalityService_GetSellerReports - Tests for GetSellerReports method / Tests para método GetSellerReports
// Cases: successful retrieval for all localities, successful retrieval for specific locality, repository error
// Casos: recuperación exitosa para todas las localidades, recuperación exitosa para localidad específica, error de repositorio
func TestLocalityService_GetSellerReports(t *testing.T) {
	t.Run("get_all_reports_ok - Returns all locality seller reports successfully", func(t *testing.T) {
		// Test: Service returns all locality reports when ID is 0 / Service retorna todos los reportes de localidad cuando ID es 0
		mockRepo, service := setupLocalityServiceTest()

		expectedReports := []responses.LocalitySellerReport{
			{
				LocalityID:   1,
				LocalityName: "City One",
				SellerCount:  5,
			},
			{
				LocalityID:   2,
				LocalityName: "City Two",
				SellerCount:  3,
			},
			{
				LocalityID:   3,
				LocalityName: "City Three",
				SellerCount:  0,
			},
		}

		// Mock repository get all reports success / Mock de éxito en obtener todos los reportes del repository
		mockRepo.GetSellerReportsFunc = func(ctx context.Context, localityID int) ([]responses.LocalitySellerReport, error) {
			assert.Equal(t, 0, localityID) // 0 means all localities / 0 significa todas las localidades
			return expectedReports, nil
		}

		result, err := service.GetSellerReports(context.Background(), 0)

		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, expectedReports, result)
		assert.Equal(t, 5, result[0].SellerCount)
		assert.Equal(t, 0, result[2].SellerCount) // Verify zero sellers case / Verificar caso de cero vendedores
	})

	t.Run("get_specific_report_ok - Returns specific locality seller report successfully", func(t *testing.T) {
		// Test: Service returns specific locality report when valid ID provided / Service retorna reporte específico cuando se proporciona ID válido
		mockRepo, service := setupLocalityServiceTest()

		localityID := 1
		expectedReports := []responses.LocalitySellerReport{
			{
				LocalityID:   1,
				LocalityName: "Specific City",
				SellerCount:  2,
			},
		}

		// Mock repository get specific report success / Mock de éxito en obtener reporte específico del repository
		mockRepo.GetSellerReportsFunc = func(ctx context.Context, id int) ([]responses.LocalitySellerReport, error) {
			assert.Equal(t, localityID, id)
			return expectedReports, nil
		}

		result, err := service.GetSellerReports(context.Background(), localityID)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, expectedReports, result)
		assert.Equal(t, localityID, result[0].LocalityID)
		assert.Equal(t, "Specific City", result[0].LocalityName)
		assert.Equal(t, 2, result[0].SellerCount)
	})

	t.Run("get_reports_empty - Returns empty slice when no localities exist", func(t *testing.T) {
		// Test: Service returns empty slice when repository returns no data / Service retorna slice vacío cuando repository no retorna datos
		mockRepo, service := setupLocalityServiceTest()

		// Mock repository empty response / Mock de respuesta vacía del repository
		mockRepo.GetSellerReportsFunc = func(ctx context.Context, localityID int) ([]responses.LocalitySellerReport, error) {
			assert.Equal(t, 0, localityID)
			return nil, nil // Empty response / Respuesta vacía
		}

		result, err := service.GetSellerReports(context.Background(), 0)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
		assert.Nil(t, result)
	})

	t.Run("get_reports_error_not_found - Returns error when specific locality not found", func(t *testing.T) {
		// Test: Service propagates not found error from repository / Service propaga error not found del repository
		mockRepo, service := setupLocalityServiceTest()

		nonExistentLocalityID := 999

		// Mock repository not found error / Mock de error not found del repository
		mockRepo.GetSellerReportsFunc = func(ctx context.Context, localityID int) ([]responses.LocalitySellerReport, error) {
			assert.Equal(t, nonExistentLocalityID, localityID)
			return nil, error_message.ErrNotFound
		}

		result, err := service.GetSellerReports(context.Background(), nonExistentLocalityID)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrNotFound, err)
		assert.Nil(t, result)
	})

	t.Run("get_reports_error_query - Returns error when repository query fails", func(t *testing.T) {
		// Test: Service handles repository query errors correctly / Service maneja errores de consulta del repository correctamente
		mockRepo, service := setupLocalityServiceTest()

		localityID := 1

		// Mock repository query error / Mock de error de consulta del repository
		mockRepo.GetSellerReportsFunc = func(ctx context.Context, id int) ([]responses.LocalitySellerReport, error) {
			assert.Equal(t, localityID, id)
			return nil, error_message.ErrQueryingReport
		}

		result, err := service.GetSellerReports(context.Background(), localityID)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrQueryingReport, err)
		assert.Nil(t, result)
	})

	t.Run("get_reports_error_existence_check - Returns error when existence check fails", func(t *testing.T) {
		// Test: Service propagates existence check errors from repository / Service propaga errores de verificación de existencia del repository
		mockRepo, service := setupLocalityServiceTest()

		localityID := 1

		// Mock repository existence check error / Mock de error de verificación de existencia del repository
		mockRepo.GetSellerReportsFunc = func(ctx context.Context, id int) ([]responses.LocalitySellerReport, error) {
			assert.Equal(t, localityID, id)
			return nil, error_message.ErrFailedCheckingExistence
		}

		result, err := service.GetSellerReports(context.Background(), localityID)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrFailedCheckingExistence, err)
		assert.Nil(t, result)
	})

	t.Run("get_reports_error_scan - Returns error when scan fails", func(t *testing.T) {
		// Test: Service handles scan errors from repository / Service maneja errores de escaneo del repository
		mockRepo, service := setupLocalityServiceTest()

		// Mock repository scan error / Mock de error de escaneo del repository
		mockRepo.GetSellerReportsFunc = func(ctx context.Context, localityID int) ([]responses.LocalitySellerReport, error) {
			assert.Equal(t, 0, localityID)
			return nil, error_message.ErrFailedToScan
		}

		result, err := service.GetSellerReports(context.Background(), 0)

		assert.Error(t, err)
		assert.Equal(t, error_message.ErrFailedToScan, err)
		assert.Nil(t, result)
	})
}
