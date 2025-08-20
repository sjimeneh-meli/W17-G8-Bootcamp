// Package services - Warehouse Service Unit Tests / Tests Unitarios del Servicio Warehouse
// Business logic testing for warehouse management with code uniqueness validation
// Testing de lógica de negocio para gestión de almacenes con validación de unicidad de código
package services

import (
	"context"
	"errors"
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
)

// setupWarehouseServiceTest - Helper function to setup test dependencies
// setupWarehouseServiceTest - Función helper para configurar dependencias de test
// Creates mock repository and service instance for testing
// Crea repositorio mock e instancia de servicio para testing
func setupWarehouseServiceTest() (*tests.WarehouseRepositoryMock, WarehouseService) {
	mockRepo := &tests.WarehouseRepositoryMock{}
	// Create direct service instance to avoid singleton
	// Crear una instancia directa del service para evitar el singleton
	service := &WarehouseServiceImpl{warehouseRepository: mockRepo}
	return mockRepo, service
}

// TestWarehouseService_Create - Tests for Create method / Tests para método Create
// Cases: successful creation, warehouse code conflict
// Casos: creación exitosa, conflicto de código de almacén
func TestWarehouseService_Create(t *testing.T) {
	t.Run("create_ok - Creates warehouse if contains necessary fields", func(t *testing.T) {
		// Test: Valid warehouse data creates successfully / Datos válidos de almacén se crean exitosamente
		// create_ok - Si contiene los campos necesarios se creará
		mockRepo, service := setupWarehouseServiceTest()

		expectedWarehouse := models.Warehouse{
			Address:            "Dirección Test",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		mockRepo.CreateFunc = func(ctx context.Context, warehouse models.Warehouse) (models.Warehouse, error) {
			warehouse.Id = 1 // Simulate generated ID / Simular ID generado
			return warehouse, nil
		}

		result, err := service.Create(context.Background(), expectedWarehouse)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, expectedWarehouse.Address, result.Address)
		assert.Equal(t, expectedWarehouse.WareHouseCode, result.WareHouseCode)
	})

	t.Run("create_conflict - Cannot be created if warehouse_code already exists", func(t *testing.T) {
		// Test: Duplicate warehouse code prevents creation / Código de almacén duplicado previene creación
		// create_conflict - Si el warehouse_code ya existe no se podrá ser creado
		mockRepo, service := setupWarehouseServiceTest()

		warehouse := models.Warehouse{
			Address:            "Dirección Test",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		mockRepo.CreateFunc = func(ctx context.Context, w models.Warehouse) (models.Warehouse, error) {
			return models.Warehouse{}, error_message.ErrAlreadyExists
		}

		result, err := service.Create(context.Background(), warehouse)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrAlreadyExists))
		assert.Empty(t, result)
	})
}

// TestWarehouseService_GetAll - Tests for GetAll method / Tests para método GetAll
// Cases: successful retrieval with data, empty list
// Casos: recuperación exitosa con datos, lista vacía
func TestWarehouseService_GetAll(t *testing.T) {
	t.Run("find_all - Returns total number of elements if list has 'n' elements", func(t *testing.T) {
		// Test: Repository with data returns all warehouses / Repositorio con datos retorna todos los almacenes
		// find_all - Si la lista posee 'n' elementos devolverá una cantidad de los elementos totales
		mockRepo, service := setupWarehouseServiceTest()

		expectedWarehouses := []models.Warehouse{
			{Id: 1, Address: "Dirección 1", Telephone: "123456789", WareHouseCode: "WH001", MinimumCapacity: 100, MinimumTemperature: 5.0, LocalityId: 1},
			{Id: 2, Address: "Dirección 2", Telephone: "123456789", WareHouseCode: "WH002", MinimumCapacity: 150, MinimumTemperature: 3.0, LocalityId: 2},
			{Id: 3, Address: "Dirección 3", Telephone: "123456789", WareHouseCode: "WH003", MinimumCapacity: 200, MinimumTemperature: 7.0, LocalityId: 3},
		}

		mockRepo.GetAllFunc = func(ctx context.Context) ([]models.Warehouse, error) {
			return expectedWarehouses, nil
		}

		result, err := service.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, expectedWarehouses, result)
	})

	t.Run("find_all_empty - Returns empty list if list is empty", func(t *testing.T) {
		// Test: Empty repository returns empty list / Repositorio vacío retorna lista vacía
		// find_all_empty - Si la lista está vacía devuelve lista vacía
		mockRepo, service := setupWarehouseServiceTest()

		mockRepo.GetAllFunc = func(ctx context.Context) ([]models.Warehouse, error) {
			return []models.Warehouse{}, nil
		}

		result, err := service.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Len(t, result, 0)
	})
}

// TestWarehouseService_GetById - Tests for GetById method / Tests para método GetById
// Cases: existing warehouse, non-existent warehouse
// Casos: almacén existente, almacén inexistente
func TestWarehouseService_GetById(t *testing.T) {
	t.Run("find_by_id_existent - Returns element information if searched element by id exists", func(t *testing.T) {
		// Test: Valid ID returns warehouse data / ID válido retorna datos del almacén
		// find_by_id_existent - Si el elemento buscado por id existe devolverá la información del elemento solicitado
		mockRepo, service := setupWarehouseServiceTest()

		expectedWarehouse := models.Warehouse{
			Id:                 1,
			Address:            "Dirección Test",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		mockRepo.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			if id == 1 {
				return expectedWarehouse, nil
			}
			return models.Warehouse{}, error_message.ErrNotFound
		}

		result, err := service.GetById(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedWarehouse, result)
	})

	t.Run("find_by_id_non_existent - Returns null if searched element by id doesn't exist", func(t *testing.T) {
		// Test: Non-existent ID returns not found error / ID inexistente retorna error not found
		// find_by_id_non_existent - Si el elemento buscado por id no existe retorna nulo
		mockRepo, service := setupWarehouseServiceTest()

		mockRepo.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{}, error_message.ErrNotFound
		}

		result, err := service.GetById(context.Background(), 999)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrNotFound))
		assert.Empty(t, result)
	})
}

// TestWarehouseService_Update - Tests for Update method / Tests para método Update
// Cases: successful update, non-existent warehouse, code change validation, code conflict
// Casos: actualización exitosa, almacén inexistente, validación de cambio de código, conflicto de código
func TestWarehouseService_Update(t *testing.T) {
	t.Run("update_existent - Returns updated element information if fields are updated successfully", func(t *testing.T) {
		// Test: Valid update data modifies warehouse successfully / Datos de actualización válidos modifican almacén exitosamente
		// update_existent - Si los campos son actualizados exitosamente devolverá la información del elemento actualizado
		mockRepo, service := setupWarehouseServiceTest()

		currentWarehouse := models.Warehouse{
			Id:                 1,
			Address:            "Dirección Original",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		updatedWarehouse := models.Warehouse{
			Address:            "Dirección Actualizada",
			Telephone:          "987654321",
			WareHouseCode:      "WH001", // Same code, no validation needed / Mismo código, no necesita validación
			MinimumCapacity:    150,
			MinimumTemperature: 3.0,
			LocalityId:         2,
		}

		mockRepo.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			if id == 1 {
				return currentWarehouse, nil
			}
			return models.Warehouse{}, error_message.ErrNotFound
		}

		mockRepo.UpdateFunc = func(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error) {
			warehouse.Id = id
			return warehouse, nil
		}

		result, err := service.Update(context.Background(), 1, updatedWarehouse)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, updatedWarehouse.Address, result.Address)
		assert.Equal(t, updatedWarehouse.Telephone, result.Telephone)
		assert.Equal(t, updatedWarehouse.WareHouseCode, result.WareHouseCode)
	})

	t.Run("update_non_existent - Returns error if element doesn't exist", func(t *testing.T) {
		// Test: Update non-existent warehouse returns error / Actualizar almacén inexistente retorna error
		// update_non_existent - Si el elemento no existe retorna error
		mockRepo, service := setupWarehouseServiceTest()

		warehouse := models.Warehouse{
			Address:            "Dirección Test",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		mockRepo.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			return models.Warehouse{}, error_message.ErrNotFound
		}

		result, err := service.Update(context.Background(), 999, warehouse)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrNotFound))
		assert.Empty(t, result)
	})

	t.Run("update_with_code_change - Validates uniqueness if code changes", func(t *testing.T) {
		// Test: Code change validates uniqueness successfully / Cambio de código valida unicidad exitosamente
		// update_with_code_change - Si el código cambia, valida unicidad
		mockRepo, service := setupWarehouseServiceTest()

		currentWarehouse := models.Warehouse{
			Id:                 1,
			Address:            "Dirección Original",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		updatedWarehouse := models.Warehouse{
			Address:            "Dirección Actualizada",
			Telephone:          "987654321",
			WareHouseCode:      "WH002", // Different code / Código diferente
			MinimumCapacity:    150,
			MinimumTemperature: 3.0,
			LocalityId:         2,
		}

		mockRepo.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			if id == 1 {
				return currentWarehouse, nil
			}
			return models.Warehouse{}, error_message.ErrNotFound
		}

		mockRepo.ExistsByCodeFunc = func(ctx context.Context, code string) (bool, error) {
			if code == "WH002" {
				return false, nil // Code available / Código disponible
			}
			return true, nil
		}

		mockRepo.UpdateFunc = func(ctx context.Context, id int, warehouse models.Warehouse) (models.Warehouse, error) {
			warehouse.Id = id
			return warehouse, nil
		}

		result, err := service.Update(context.Background(), 1, updatedWarehouse)

		assert.NoError(t, err)
		assert.Equal(t, "WH002", result.WareHouseCode)
	})

	t.Run("update_with_code_conflict - Returns error if new code already exists", func(t *testing.T) {
		// Test: Duplicate new code prevents update / Nuevo código duplicado previene actualización
		// update_with_code_conflict - Si el nuevo código ya existe, retorna error
		mockRepo, service := setupWarehouseServiceTest()

		currentWarehouse := models.Warehouse{
			Id:                 1,
			Address:            "Dirección Original",
			Telephone:          "123456789",
			WareHouseCode:      "WH001",
			MinimumCapacity:    100,
			MinimumTemperature: 5.0,
			LocalityId:         1,
		}

		updatedWarehouse := models.Warehouse{
			Address:            "Dirección Actualizada",
			Telephone:          "987654321",
			WareHouseCode:      "WH002", // Code that already exists / Código que ya existe
			MinimumCapacity:    150,
			MinimumTemperature: 3.0,
			LocalityId:         2,
		}

		mockRepo.GetWarehouseByIDFunc = func(ctx context.Context, id int) (models.Warehouse, error) {
			if id == 1 {
				return currentWarehouse, nil
			}
			return models.Warehouse{}, error_message.ErrNotFound
		}

		mockRepo.ExistsByCodeFunc = func(ctx context.Context, code string) (bool, error) {
			if code == "WH002" {
				return true, nil // Code already exists / Código ya existe
			}
			return false, nil
		}

		result, err := service.Update(context.Background(), 1, updatedWarehouse)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrAlreadyExists))
		assert.Empty(t, result)
	})
}

// TestWarehouseService_Delete - Tests for Delete method / Tests para método Delete
// Cases: successful deletion, non-existent warehouse
// Casos: eliminación exitosa, almacén inexistente
func TestWarehouseService_Delete(t *testing.T) {
	t.Run("delete_ok - Element will not appear in list if deletion is successful", func(t *testing.T) {
		// Test: Valid deletion removes warehouse successfully / Eliminación válida remueve almacén exitosamente
		// delete_ok - Si la eliminación es exitosa el elemento no aparecerá en la lista
		mockRepo, service := setupWarehouseServiceTest()

		mockRepo.DeleteFunc = func(ctx context.Context, id int) error {
			return nil
		}

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
	})

	t.Run("delete_non_existent - Returns null if element to delete doesn't exist", func(t *testing.T) {
		// Test: Delete non-existent warehouse returns error / Eliminar almacén inexistente retorna error
		// delete_non_existent - Si el elemento que se desea eliminar no existe retorna nulo
		mockRepo, service := setupWarehouseServiceTest()

		mockRepo.DeleteFunc = func(ctx context.Context, id int) error {
			return error_message.ErrNotFound
		}

		err := service.Delete(context.Background(), 999)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrNotFound))
	})
}

// TestWarehouseService_ValidateCodeUniqueness - Tests for ValidateCodeUniqueness method
// TestWarehouseService_ValidateCodeUniqueness - Tests para método ValidateCodeUniqueness
// Cases: unique code, existing code, repository error
// Casos: código único, código existente, error de repositorio
func TestWarehouseService_ValidateCodeUniqueness(t *testing.T) {
	t.Run("validate_code_unique - Returns nil if code doesn't exist", func(t *testing.T) {
		// Test: Unique code validation passes / Validación de código único pasa
		// validate_code_unique - Si el código no existe, retorna nil
		mockRepo, service := setupWarehouseServiceTest()

		mockRepo.ExistsByCodeFunc = func(ctx context.Context, code string) (bool, error) {
			return false, nil
		}

		err := service.ValidateCodeUniqueness(context.Background(), "WH001")

		assert.NoError(t, err)
	})

	t.Run("validate_code_exists - Returns error if code already exists", func(t *testing.T) {
		// Test: Existing code validation fails / Validación de código existente falla
		// validate_code_exists - Si el código ya existe, retorna error
		mockRepo, service := setupWarehouseServiceTest()

		mockRepo.ExistsByCodeFunc = func(ctx context.Context, code string) (bool, error) {
			return true, nil
		}

		err := service.ValidateCodeUniqueness(context.Background(), "WH001")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrAlreadyExists))
	})

	t.Run("validate_code_repository_error - Returns internal error if repository fails", func(t *testing.T) {
		// Test: Repository error during validation / Error de repositorio durante validación
		// validate_code_repository_error - Si el repository falla, retorna error interno
		mockRepo, service := setupWarehouseServiceTest()

		mockRepo.ExistsByCodeFunc = func(ctx context.Context, code string) (bool, error) {
			return false, error_message.ErrInternalServerError
		}

		err := service.ValidateCodeUniqueness(context.Background(), "WH001")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
	})
}
