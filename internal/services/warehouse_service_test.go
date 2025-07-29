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

func setupWarehouseServiceTest() (*tests.WarehouseRepositoryMock, WarehouseService) {
	mockRepo := &tests.WarehouseRepositoryMock{}
	// Crear una instancia directa del service para evitar el singleton
	service := &WarehouseServiceImpl{warehouseRepository: mockRepo}
	return mockRepo, service
}

func TestWarehouseService_Create(t *testing.T) {
	t.Run("create_ok - Si contiene los campos necesarios se creará", func(t *testing.T) {
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
			warehouse.Id = 1 // Simular ID generado
			return warehouse, nil
		}

		result, err := service.Create(context.Background(), expectedWarehouse)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, expectedWarehouse.Address, result.Address)
		assert.Equal(t, expectedWarehouse.WareHouseCode, result.WareHouseCode)
	})

	t.Run("create_conflict - Si el warehouse_code ya existe no se podrá ser creado", func(t *testing.T) {
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

func TestWarehouseService_GetAll(t *testing.T) {
	t.Run("find_all - Si la lista posee 'n' elementos devolverá una cantidad de los elementos totales", func(t *testing.T) {
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

	t.Run("find_all_empty - Si la lista está vacía devuelve lista vacía", func(t *testing.T) {
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

func TestWarehouseService_GetById(t *testing.T) {
	t.Run("find_by_id_existent - Si el elemento buscado por id existe devolverá la información del elemento solicitado", func(t *testing.T) {
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

	t.Run("find_by_id_non_existent - Si el elemento buscado por id no existe retorna nulo", func(t *testing.T) {
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

func TestWarehouseService_Update(t *testing.T) {
	t.Run("update_existent - Si los campos son actualizados exitosamente devolverá la información del elemento actualizado", func(t *testing.T) {
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
			WareHouseCode:      "WH001", // Mismo código, no necesita validación
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

	t.Run("update_non_existent - Si el elemento no existe retorna error", func(t *testing.T) {
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

	t.Run("update_with_code_change - Si el código cambia, valida unicidad", func(t *testing.T) {
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
			WareHouseCode:      "WH002", // Código diferente
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
				return false, nil // Código disponible
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

	t.Run("update_with_code_conflict - Si el nuevo código ya existe, retorna error", func(t *testing.T) {
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
			WareHouseCode:      "WH002", // Código que ya existe
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
				return true, nil // Código ya existe
			}
			return false, nil
		}

		result, err := service.Update(context.Background(), 1, updatedWarehouse)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrAlreadyExists))
		assert.Empty(t, result)
	})
}

func TestWarehouseService_Delete(t *testing.T) {
	t.Run("delete_ok - Si la eliminación es exitosa el elemento no aparecerá en la lista", func(t *testing.T) {
		mockRepo, service := setupWarehouseServiceTest()

		mockRepo.DeleteFunc = func(ctx context.Context, id int) error {
			return nil
		}

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
	})

	t.Run("delete_non_existent - Si el elemento que se desea eliminar no existe retorna nulo", func(t *testing.T) {
		mockRepo, service := setupWarehouseServiceTest()

		mockRepo.DeleteFunc = func(ctx context.Context, id int) error {
			return error_message.ErrNotFound
		}

		err := service.Delete(context.Background(), 999)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrNotFound))
	})
}

func TestWarehouseService_ValidateCodeUniqueness(t *testing.T) {
	t.Run("validate_code_unique - Si el código no existe, retorna nil", func(t *testing.T) {
		mockRepo, service := setupWarehouseServiceTest()

		mockRepo.ExistsByCodeFunc = func(ctx context.Context, code string) (bool, error) {
			return false, nil
		}

		err := service.ValidateCodeUniqueness(context.Background(), "WH001")

		assert.NoError(t, err)
	})

	t.Run("validate_code_exists - Si el código ya existe, retorna error", func(t *testing.T) {
		mockRepo, service := setupWarehouseServiceTest()

		mockRepo.ExistsByCodeFunc = func(ctx context.Context, code string) (bool, error) {
			return true, nil
		}

		err := service.ValidateCodeUniqueness(context.Background(), "WH001")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrAlreadyExists))
	})

	t.Run("validate_code_repository_error - Si el repository falla, retorna error interno", func(t *testing.T) {
		mockRepo, service := setupWarehouseServiceTest()

		mockRepo.ExistsByCodeFunc = func(ctx context.Context, code string) (bool, error) {
			return false, error_message.ErrInternalServerError
		}

		err := service.ValidateCodeUniqueness(context.Background(), "WH001")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, error_message.ErrInternalServerError))
	})
}
