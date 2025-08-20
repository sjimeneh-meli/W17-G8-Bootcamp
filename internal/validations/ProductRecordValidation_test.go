package validations

import (
	"testing"
	"time"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/stretchr/testify/assert"
)

// TestGetProductRecordValidation agrupa todos los tests para el constructor GetProductRecordValidation
// TestGetProductRecordValidation groups all tests for the GetProductRecordValidation constructor
func TestGetProductRecordValidation(t *testing.T) {
	// Test case: Successful validation instance creation
	// Caso de prueba: Creación exitosa de instancia de validación
	t.Run("Success_CreatesValidationInstance", func(t *testing.T) {
		// Act - Crear instancia de validación
		// Act - Create validation instance
		validator := GetProductRecordValidation()

		// Assert - Verificar que la instancia se creó correctamente
		// Assert - Verify that the instance was created correctly
		assert.NotNil(t, validator)
		assert.IsType(t, &ProductRecordValidation{}, validator)
	})

	// Test case: Multiple calls return different instances (no singleton)
	// Caso de prueba: Múltiples llamadas retornan instancias diferentes (no singleton)
	t.Run("Success_ReturnsNewInstancesEachTime", func(t *testing.T) {
		// Act - Crear múltiples instancias
		// Act - Create multiple instances
		validator1 := GetProductRecordValidation()
		validator2 := GetProductRecordValidation()

		// Assert - Verificar que son instancias diferentes
		// Assert - Verify they are different instances
		assert.NotNil(t, validator1)
		assert.NotNil(t, validator2)
		// Note: Since ProductRecordValidation is a simple struct with no fields,
		// we can't test if they're different instances easily,
		// but we can verify they're both valid ProductRecordValidation pointers
		assert.IsType(t, &ProductRecordValidation{}, validator1)
		assert.IsType(t, &ProductRecordValidation{}, validator2)
	})
}

// TestProductRecordValidation_ValidateProductRecordRequestStruct agrupa todos los tests para el método ValidateProductRecordRequestStruct
// TestProductRecordValidation_ValidateProductRecordRequestStruct groups all tests for the ValidateProductRecordRequestStruct method
func TestProductRecordValidation_ValidateProductRecordRequestStruct(t *testing.T) {
	// Test case: Valid product record request with all required fields
	// Caso de prueba: Request de registro de producto válido con todos los campos requeridos
	t.Run("Success_ValidProductRecordRequest", func(t *testing.T) {
		// Arrange - Configurar validador y request válido
		// Arrange - Set up validator and valid request
		validator := GetProductRecordValidation()
		validRequest := requests.ProductRecordRequest{
			LastUpdateDate: time.Now(),
			PurchasePrice:  25.50,
			SalePrice:      35.75,
			ProductID:      123,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(validRequest)

		// Assert - Verificar que no hay errores
		// Assert - Verify no errors
		assert.NoError(t, err)
	})

	// Test case: Valid product record request with specific date
	// Caso de prueba: Request de registro de producto válido con fecha específica
	t.Run("Success_ValidProductRecordRequestWithSpecificDate", func(t *testing.T) {
		// Arrange - Configurar validador y request válido con fecha específica
		// Arrange - Set up validator and valid request with specific date
		validator := GetProductRecordValidation()
		specificDate := time.Date(2023, 12, 15, 10, 30, 0, 0, time.UTC)
		validRequest := requests.ProductRecordRequest{
			LastUpdateDate: specificDate,
			PurchasePrice:  100.00,
			SalePrice:      150.00,
			ProductID:      456,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(validRequest)

		// Assert - Verificar que no hay errores
		// Assert - Verify no errors
		assert.NoError(t, err)
	})

	// Test case: Valid product record request with very small prices
	// Caso de prueba: Request de registro de producto válido con precios muy pequeños
	t.Run("Success_ValidProductRecordRequestWithSmallPrices", func(t *testing.T) {
		// Arrange - Configurar validador y request válido con precios pequeños
		// Arrange - Set up validator and valid request with small prices
		validator := GetProductRecordValidation()
		validRequest := requests.ProductRecordRequest{
			LastUpdateDate: time.Now(),
			PurchasePrice:  0.01,
			SalePrice:      0.02,
			ProductID:      789,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(validRequest)

		// Assert - Verificar que no hay errores
		// Assert - Verify no errors
		assert.NoError(t, err)
	})

	// Test case: Valid product record request with large prices
	// Caso de prueba: Request de registro de producto válido con precios grandes
	t.Run("Success_ValidProductRecordRequestWithLargePrices", func(t *testing.T) {
		// Arrange - Configurar validador y request válido con precios grandes
		// Arrange - Set up validator and valid request with large prices
		validator := GetProductRecordValidation()
		validRequest := requests.ProductRecordRequest{
			LastUpdateDate: time.Now(),
			PurchasePrice:  9999.99,
			SalePrice:      19999.99,
			ProductID:      999999,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(validRequest)

		// Assert - Verificar que no hay errores
		// Assert - Verify no errors
		assert.NoError(t, err)
	})

	// Test case: Invalid request - missing LastUpdateDate (zero value)
	// Caso de prueba: Request inválido - falta LastUpdateDate (valor cero)
	t.Run("Error_MissingLastUpdateDate", func(t *testing.T) {
		// Arrange - Configurar validador y request con LastUpdateDate en valor cero
		// Arrange - Set up validator and request with zero LastUpdateDate
		validator := GetProductRecordValidation()
		invalidRequest := requests.ProductRecordRequest{
			LastUpdateDate: time.Time{}, // Zero value for required field
			PurchasePrice:  25.50,
			SalePrice:      35.75,
			ProductID:      123,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por LastUpdateDate faltante
		// Assert - Verify error for missing LastUpdateDate
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "last_update_date")
	})

	// Test case: Invalid request - missing PurchasePrice (zero value)
	// Caso de prueba: Request inválido - falta PurchasePrice (valor cero)
	t.Run("Error_MissingPurchasePrice", func(t *testing.T) {
		// Arrange - Configurar validador y request con PurchasePrice en valor cero
		// Arrange - Set up validator and request with zero PurchasePrice
		validator := GetProductRecordValidation()
		invalidRequest := requests.ProductRecordRequest{
			LastUpdateDate: time.Now(),
			PurchasePrice:  0, // Zero value for required field
			SalePrice:      35.75,
			ProductID:      123,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por PurchasePrice faltante
		// Assert - Verify error for missing PurchasePrice
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "purchase_price")
	})

	// Test case: Invalid request - missing SalePrice (zero value)
	// Caso de prueba: Request inválido - falta SalePrice (valor cero)
	t.Run("Error_MissingSalePrice", func(t *testing.T) {
		// Arrange - Configurar validador y request con SalePrice en valor cero
		// Arrange - Set up validator and request with zero SalePrice
		validator := GetProductRecordValidation()
		invalidRequest := requests.ProductRecordRequest{
			LastUpdateDate: time.Now(),
			PurchasePrice:  25.50,
			SalePrice:      0, // Zero value for required field
			ProductID:      123,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por SalePrice faltante
		// Assert - Verify error for missing SalePrice
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sale_price")
	})

	// Test case: Invalid request - missing ProductID (zero value)
	// Caso de prueba: Request inválido - falta ProductID (valor cero)
	t.Run("Error_MissingProductID", func(t *testing.T) {
		// Arrange - Configurar validador y request con ProductID en valor cero
		// Arrange - Set up validator and request with zero ProductID
		validator := GetProductRecordValidation()
		invalidRequest := requests.ProductRecordRequest{
			LastUpdateDate: time.Now(),
			PurchasePrice:  25.50,
			SalePrice:      35.75,
			ProductID:      0, // Zero value for required field
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por ProductID faltante
		// Assert - Verify error for missing ProductID
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product_id")
	})

	// Test case: Invalid request - multiple missing fields
	// Caso de prueba: Request inválido - múltiples campos faltantes
	t.Run("Error_MultipleMissingFields", func(t *testing.T) {
		// Arrange - Configurar validador y request con múltiples campos faltantes
		// Arrange - Set up validator and request with multiple missing fields
		validator := GetProductRecordValidation()
		invalidRequest := requests.ProductRecordRequest{
			LastUpdateDate: time.Time{}, // Missing
			PurchasePrice:  0,           // Missing
			SalePrice:      0,           // Missing
			ProductID:      123,         // Valid
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(invalidRequest)

		// Assert - Verificar que hay errores por múltiples campos faltantes
		// Assert - Verify errors for multiple missing fields
		assert.Error(t, err)
		errorMsg := err.Error()
		assert.Contains(t, errorMsg, "last_update_date")
		assert.Contains(t, errorMsg, "purchase_price")
		assert.Contains(t, errorMsg, "sale_price")
	})

	// Test case: Invalid request - all fields missing (worst case)
	// Caso de prueba: Request inválido - todos los campos faltantes (peor caso)
	t.Run("Error_AllFieldsMissing", func(t *testing.T) {
		// Arrange - Configurar validador y request con todos los campos faltantes
		// Arrange - Set up validator and request with all fields missing
		validator := GetProductRecordValidation()
		invalidRequest := requests.ProductRecordRequest{
			LastUpdateDate: time.Time{}, // Missing
			PurchasePrice:  0,           // Missing
			SalePrice:      0,           // Missing
			ProductID:      0,           // Missing
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(invalidRequest)

		// Assert - Verificar que hay errores por todos los campos faltantes
		// Assert - Verify errors for all missing fields
		assert.Error(t, err)
		errorMsg := err.Error()
		assert.Contains(t, errorMsg, "last_update_date")
		assert.Contains(t, errorMsg, "purchase_price")
		assert.Contains(t, errorMsg, "sale_price")
		assert.Contains(t, errorMsg, "product_id")
	})

	// Test case: Valid request - negative prices (should pass validation)
	// Caso de prueba: Request válido - precios negativos (debe pasar validación)
	t.Run("Success_NegativePricesAllowed", func(t *testing.T) {
		// Arrange - Configurar validador y request con precios negativos
		// Arrange - Set up validator and request with negative prices
		validator := GetProductRecordValidation()
		validRequest := requests.ProductRecordRequest{
			LastUpdateDate: time.Now(),
			PurchasePrice:  -25.50, // Negative value should be allowed for Required validation
			SalePrice:      -35.75, // Negative value should be allowed for Required validation
			ProductID:      -123,   // Negative value should be allowed for Required validation
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(validRequest)

		// Assert - Verificar que no hay errores (los valores negativos son válidos para la validación Required)
		// Assert - Verify no errors (negative values are valid for Required validation)
		assert.NoError(t, err)
	})

	// Test case: Valid request - edge case with future date
	// Caso de prueba: Request válido - caso límite con fecha futura
	t.Run("Success_FutureDateAllowed", func(t *testing.T) {
		// Arrange - Configurar validador y request con fecha futura
		// Arrange - Set up validator and request with future date
		validator := GetProductRecordValidation()
		futureDate := time.Now().Add(24 * time.Hour) // Tomorrow
		validRequest := requests.ProductRecordRequest{
			LastUpdateDate: futureDate,
			PurchasePrice:  50.00,
			SalePrice:      75.00,
			ProductID:      999,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(validRequest)

		// Assert - Verificar que no hay errores (las fechas futuras son válidas para la validación Required)
		// Assert - Verify no errors (future dates are valid for Required validation)
		assert.NoError(t, err)
	})

	// Test case: Valid request - edge case with very old date
	// Caso de prueba: Request válido - caso límite con fecha muy antigua
	t.Run("Success_VeryOldDateAllowed", func(t *testing.T) {
		// Arrange - Configurar validador y request con fecha muy antigua
		// Arrange - Set up validator and request with very old date
		validator := GetProductRecordValidation()
		oldDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
		validRequest := requests.ProductRecordRequest{
			LastUpdateDate: oldDate,
			PurchasePrice:  10.00,
			SalePrice:      15.00,
			ProductID:      1,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRecordRequestStruct(validRequest)

		// Assert - Verificar que no hay errores (las fechas antiguas son válidas para la validación Required)
		// Assert - Verify no errors (old dates are valid for Required validation)
		assert.NoError(t, err)
	})
}
