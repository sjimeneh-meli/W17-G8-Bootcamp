// Package services - Carry Service Unit Tests / Tests Unitarios del Servicio Carry
// Business logic testing for carry management with CID uniqueness and locality validation
// Testing de lógica de negocio para gestión de transportistas con validación de unicidad de CID y localidad
package services

import (
	"context"
	"errors"
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
)

// setupCarryServiceTest - Helper function to setup test dependencies
// setupCarryServiceTest - Función helper para configurar dependencias de test
// Creates mock repositories and service instance for testing
// Crea repositorios mock e instancia de servicio para testing
func setupCarryServiceTest() (*tests.CarryRepositoryMock, *tests.LocalityRepositoryMock, CarryService) {
	mockCarryRepo := &tests.CarryRepositoryMock{}
	mockLocalityRepo := &tests.LocalityRepositoryMock{}
	// Create direct service instance to avoid singleton
	// Crear una instancia directa del service para evitar el singleton
	service := &CarryServiceImpl{carryRepository: mockCarryRepo, localityRepository: mockLocalityRepo}
	return mockCarryRepo, mockLocalityRepo, service
}

// TestCarryService_CreateCarry - Tests for CreateCarry method / Tests para método CreateCarry
// Cases: successful creation, CID conflict, locality not found
// Casos: creación exitosa, conflicto de CID, localidad no encontrada
func TestCarryService_CreateCarry(t *testing.T) {
	t.Run("create_ok - Creates carry if all validations pass", func(t *testing.T) {
		// Test: Valid carry data creates successfully / Datos válidos de transportista se crean exitosamente
		mockCarryRepo, mockLocalityRepo, service := setupCarryServiceTest()

		expectedCarry := models.Carry{
			Cid:         "CID001",
			CompanyName: "Empresa Test",
			Address:     "Dirección Test",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		// Mock locality exists validation / Mock validación de existencia de localidad
		mockLocalityRepo.ExistByIdFunc = func(ctx context.Context, id int) (bool, error) {
			return true, nil
		}

		// Mock CID uniqueness validation / Mock validación de unicidad de CID
		mockCarryRepo.ExistsByCidFunc = func(ctx context.Context, cid string) (bool, error) {
			return false, nil
		}

		// Mock carry creation / Mock creación de transportista
		mockCarryRepo.CreateFunc = func(ctx context.Context, carry models.Carry) (models.Carry, error) {
			carry.Id = 1 // Simulate generated ID / Simular ID generado
			return carry, nil
		}

		result, err := service.CreateCarry(context.Background(), expectedCarry)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, expectedCarry.Cid, result.Cid)
		assert.Equal(t, expectedCarry.CompanyName, result.CompanyName)
	})

	t.Run("create_conflict_cid - Cannot be created if CID already exists", func(t *testing.T) {
		// Test: Duplicate CID prevents creation / CID duplicado previene creación
		mockCarryRepo, mockLocalityRepo, service := setupCarryServiceTest()

		carry := models.Carry{
			Cid:         "CID001",
			CompanyName: "Empresa Test",
			Address:     "Dirección Test",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		// Mock locality exists validation / Mock validación de existencia de localidad
		mockLocalityRepo.ExistByIdFunc = func(ctx context.Context, id int) (bool, error) {
			return true, nil
		}

		// Mock CID already exists / Mock CID ya existe
		mockCarryRepo.ExistsByCidFunc = func(ctx context.Context, cid string) (bool, error) {
			return true, nil
		}

		result, err := service.CreateCarry(context.Background(), carry)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrAlreadyExists))
		assert.Empty(t, result)
	})

	t.Run("create_locality_not_found - Cannot be created if locality doesn't exist", func(t *testing.T) {
		// Test: Non-existent locality prevents creation / Localidad inexistente previene creación
		_, mockLocalityRepo, service := setupCarryServiceTest()

		carry := models.Carry{
			Cid:         "CID001",
			CompanyName: "Empresa Test",
			Address:     "Dirección Test",
			Telephone:   "123456789",
			LocalityId:  999,
		}

		// Mock locality doesn't exist / Mock localidad no existe
		mockLocalityRepo.ExistByIdFunc = func(ctx context.Context, id int) (bool, error) {
			return false, nil
		}

		result, err := service.CreateCarry(context.Background(), carry)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrNotFound))
		assert.Empty(t, result)
	})

	t.Run("create_locality_repository_error - Returns internal error if locality repository fails", func(t *testing.T) {
		// Test: Locality repository error during validation / Error de repositorio de localidad durante validación
		_, mockLocalityRepo, service := setupCarryServiceTest()

		carry := models.Carry{
			Cid:         "CID001",
			CompanyName: "Empresa Test",
			Address:     "Dirección Test",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		// Mock locality repository error / Mock error del repositorio de localidad
		mockLocalityRepo.ExistByIdFunc = func(ctx context.Context, id int) (bool, error) {
			return false, error_message.ErrInternalServerError
		}

		result, err := service.CreateCarry(context.Background(), carry)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
		assert.Empty(t, result)
	})

	t.Run("create_cid_repository_error - Returns internal error if CID repository fails", func(t *testing.T) {
		// Test: CID repository error during validation / Error de repositorio de CID durante validación
		mockCarryRepo, mockLocalityRepo, service := setupCarryServiceTest()

		carry := models.Carry{
			Cid:         "CID001",
			CompanyName: "Empresa Test",
			Address:     "Dirección Test",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		// Mock locality exists validation / Mock validación de existencia de localidad
		mockLocalityRepo.ExistByIdFunc = func(ctx context.Context, id int) (bool, error) {
			return true, nil
		}

		// Mock CID repository error / Mock error del repositorio de CID
		mockCarryRepo.ExistsByCidFunc = func(ctx context.Context, cid string) (bool, error) {
			return false, error_message.ErrInternalServerError
		}

		result, err := service.CreateCarry(context.Background(), carry)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
		assert.Empty(t, result)
	})

	t.Run("create_carry_repository_error - Returns internal error if carry repository fails", func(t *testing.T) {
		// Test: Carry repository error during creation / Error de repositorio de transportista durante creación
		mockCarryRepo, mockLocalityRepo, service := setupCarryServiceTest()

		carry := models.Carry{
			Cid:         "CID001",
			CompanyName: "Empresa Test",
			Address:     "Dirección Test",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		// Mock locality exists validation / Mock validación de existencia de localidad
		mockLocalityRepo.ExistByIdFunc = func(ctx context.Context, id int) (bool, error) {
			return true, nil
		}

		// Mock CID uniqueness validation / Mock validación de unicidad de CID
		mockCarryRepo.ExistsByCidFunc = func(ctx context.Context, cid string) (bool, error) {
			return false, nil
		}

		// Mock carry repository error / Mock error del repositorio de transportista
		mockCarryRepo.CreateFunc = func(ctx context.Context, c models.Carry) (models.Carry, error) {
			return models.Carry{}, error_message.ErrInternalServerError
		}

		result, err := service.CreateCarry(context.Background(), carry)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
		assert.Empty(t, result)
	})
}

// TestCarryService_GetCarryReportByLocality - Tests for GetCarryReportByLocality method
// TestCarryService_GetCarryReportByLocality - Tests para método GetCarryReportByLocality
// Cases: all localities report, specific locality report, locality not found
// Casos: reporte de todas las localidades, reporte de localidad específica, localidad no encontrada
func TestCarryService_GetCarryReportByLocality(t *testing.T) {
	t.Run("get_all_localities_report - Returns reports for all localities when localityID is 0", func(t *testing.T) {
		// Test: Zero locality ID returns all localities report / ID de localidad cero retorna reporte de todas las localidades
		mockCarryRepo, _, service := setupCarryServiceTest()

		expectedReports := []responses.LocalityCarryReport{
			{LocalityId: 1, LocalityName: "Localidad 1", CarriersCount: 5},
			{LocalityId: 2, LocalityName: "Localidad 2", CarriersCount: 3},
		}

		// Mock carry repository returns all reports / Mock repositorio de transportista retorna todos los reportes
		mockCarryRepo.GetCarryReportsByLocalityFunc = func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
			if localityID == 0 {
				return expectedReports, nil
			}
			return []responses.LocalityCarryReport{}, nil
		}

		result, err := service.GetCarryReportByLocality(context.Background(), 0)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedReports, result)
	})

	t.Run("get_specific_locality_report - Returns report for specific locality when localityID is provided", func(t *testing.T) {
		// Test: Specific locality ID returns locality report / ID de localidad específico retorna reporte de localidad
		mockCarryRepo, mockLocalityRepo, service := setupCarryServiceTest()

		expectedReports := []responses.LocalityCarryReport{
			{LocalityId: 1, LocalityName: "Localidad 1", CarriersCount: 5},
		}

		// Mock locality exists validation / Mock validación de existencia de localidad
		mockLocalityRepo.ExistByIdFunc = func(ctx context.Context, id int) (bool, error) {
			return id == 1, nil
		}

		// Mock carry repository returns specific locality report / Mock repositorio de transportista retorna reporte de localidad específica
		mockCarryRepo.GetCarryReportsByLocalityFunc = func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
			if localityID == 1 {
				return expectedReports, nil
			}
			return []responses.LocalityCarryReport{}, nil
		}

		result, err := service.GetCarryReportByLocality(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, expectedReports, result)
	})

	t.Run("get_specific_locality_not_found - Returns error if specific locality doesn't exist", func(t *testing.T) {
		// Test: Non-existent locality ID returns error / ID de localidad inexistente retorna error
		_, mockLocalityRepo, service := setupCarryServiceTest()

		// Mock locality doesn't exist / Mock localidad no existe
		mockLocalityRepo.ExistByIdFunc = func(ctx context.Context, id int) (bool, error) {
			return false, nil
		}

		result, err := service.GetCarryReportByLocality(context.Background(), 999)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrNotFound))
		assert.Empty(t, result)
	})

	t.Run("get_specific_locality_repository_error - Returns internal error if locality repository fails", func(t *testing.T) {
		// Test: Locality repository error during validation / Error de repositorio de localidad durante validación
		_, mockLocalityRepo, service := setupCarryServiceTest()

		// Mock locality repository error / Mock error del repositorio de localidad
		mockLocalityRepo.ExistByIdFunc = func(ctx context.Context, id int) (bool, error) {
			return false, error_message.ErrInternalServerError
		}

		result, err := service.GetCarryReportByLocality(context.Background(), 1)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
		assert.Empty(t, result)
	})

	t.Run("get_carry_reports_repository_error - Returns internal error if carry repository fails", func(t *testing.T) {
		// Test: Carry repository error during report retrieval / Error de repositorio de transportista durante recuperación de reportes
		mockCarryRepo, mockLocalityRepo, service := setupCarryServiceTest()

		// Mock locality exists validation / Mock validación de existencia de localidad
		mockLocalityRepo.ExistByIdFunc = func(ctx context.Context, id int) (bool, error) {
			return true, nil
		}

		// Mock carry repository error / Mock error del repositorio de transportista
		mockCarryRepo.GetCarryReportsByLocalityFunc = func(ctx context.Context, localityID int) ([]responses.LocalityCarryReport, error) {
			return nil, error_message.ErrInternalServerError
		}

		result, err := service.GetCarryReportByLocality(context.Background(), 1)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
		assert.Empty(t, result)
	})
}
