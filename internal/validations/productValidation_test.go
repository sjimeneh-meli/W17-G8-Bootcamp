package validations

import (
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/stretchr/testify/assert"
)

// TestGetProductValidation agrupa todos los tests para el constructor GetProductValidation
// TestGetProductValidation groups all tests for the GetProductValidation constructor
func TestGetProductValidation(t *testing.T) {
	// Test case: Successful validation instance creation
	// Caso de prueba: Creación exitosa de instancia de validación
	t.Run("Success_CreatesValidationInstance", func(t *testing.T) {
		// Act - Crear instancia de validación
		// Act - Create validation instance
		validator := GetProductValidation()

		// Assert - Verificar que la instancia se creó correctamente
		// Assert - Verify that the instance was created correctly
		assert.NotNil(t, validator)
		assert.IsType(t, &ProductValidation{}, validator)
	})

	// Test case: Multiple calls return different instances (no singleton)
	// Caso de prueba: Múltiples llamadas retornan instancias diferentes (no singleton)
	t.Run("Success_ReturnsNewInstancesEachTime", func(t *testing.T) {
		// Act - Crear múltiples instancias
		// Act - Create multiple instances
		validator1 := GetProductValidation()
		validator2 := GetProductValidation()

		// Assert - Verificar que son instancias diferentes
		// Assert - Verify they are different instances
		assert.NotNil(t, validator1)
		assert.NotNil(t, validator2)
		// Note: Since ProductValidation is a simple struct with no fields,
		// we can't test if they're different instances easily,
		// but we can verify they're both valid ProductValidation pointers
		assert.IsType(t, &ProductValidation{}, validator1)
		assert.IsType(t, &ProductValidation{}, validator2)
	})
}

// TestProductValidation_ValidateProductRequestStruct agrupa todos los tests para el método ValidateProductRequestStruct
// TestProductValidation_ValidateProductRequestStruct groups all tests for the ValidateProductRequestStruct method
func TestProductValidation_ValidateProductRequestStruct(t *testing.T) {
	// Test case: Valid product request with all required fields
	// Caso de prueba: Request de producto válido con todos los campos requeridos
	t.Run("Success_ValidProductRequest", func(t *testing.T) {
		// Arrange - Configurar validador y request válido
		// Arrange - Set up validator and valid request
		validator := GetProductValidation()
		validRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product Description",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
			SellerID:                       nil, // Optional field
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(validRequest)

		// Assert - Verificar que no hay errores
		// Assert - Verify no errors
		assert.NoError(t, err)
	})

	// Test case: Valid product request with SellerID provided
	// Caso de prueba: Request de producto válido con SellerID proporcionado
	t.Run("Success_ValidProductRequestWithSellerID", func(t *testing.T) {
		// Arrange - Configurar validador y request válido con SellerID
		// Arrange - Set up validator and valid request with SellerID
		validator := GetProductValidation()
		sellerID := int64(123)
		validRequest := requests.ProductRequest{
			ProductCode:                    "PROD002",
			Description:                    "Test Product with Seller",
			Width:                          8.5,
			Height:                         15.3,
			Length:                         12.7,
			NetWeight:                      3.2,
			ExpirationRate:                 0.2,
			RecommendedFreezingTemperature: -20.0,
			FreezingRate:                   0.9,
			ProductTypeID:                  2,
			SellerID:                       &sellerID, // Optional field provided
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(validRequest)

		// Assert - Verificar que no hay errores
		// Assert - Verify no errors
		assert.NoError(t, err)
	})

	// Test case: Invalid request - missing ProductCode
	// Caso de prueba: Request inválido - falta ProductCode
	t.Run("Error_MissingProductCode", func(t *testing.T) {
		// Arrange - Configurar validador y request con ProductCode faltante
		// Arrange - Set up validator and request with missing ProductCode
		validator := GetProductValidation()
		invalidRequest := requests.ProductRequest{
			ProductCode:                    "", // Missing required field
			Description:                    "Test Product Description",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por ProductCode faltante
		// Assert - Verify error for missing ProductCode
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product_code")
	})

	// Test case: Invalid request - missing Description
	// Caso de prueba: Request inválido - falta Description
	t.Run("Error_MissingDescription", func(t *testing.T) {
		// Arrange - Configurar validador y request con Description faltante
		// Arrange - Set up validator and request with missing Description
		validator := GetProductValidation()
		invalidRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "", // Missing required field
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por Description faltante
		// Assert - Verify error for missing Description
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "description")
	})

	// Test case: Invalid request - missing Width (zero value)
	// Caso de prueba: Request inválido - falta Width (valor cero)
	t.Run("Error_MissingWidth", func(t *testing.T) {
		// Arrange - Configurar validador y request con Width en valor cero
		// Arrange - Set up validator and request with zero Width
		validator := GetProductValidation()
		invalidRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product Description",
			Width:                          0, // Zero value for required field
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por Width faltante
		// Assert - Verify error for missing Width
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "width")
	})

	// Test case: Invalid request - missing Height (zero value)
	// Caso de prueba: Request inválido - falta Height (valor cero)
	t.Run("Error_MissingHeight", func(t *testing.T) {
		// Arrange - Configurar validador y request con Height en valor cero
		// Arrange - Set up validator and request with zero Height
		validator := GetProductValidation()
		invalidRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product Description",
			Width:                          10.5,
			Height:                         0, // Zero value for required field
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por Height faltante
		// Assert - Verify error for missing Height
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "height")
	})

	// Test case: Invalid request - missing Length (zero value)
	// Caso de prueba: Request inválido - falta Length (valor cero)
	t.Run("Error_MissingLength", func(t *testing.T) {
		// Arrange - Configurar validador y request con Length en valor cero
		// Arrange - Set up validator and request with zero Length
		validator := GetProductValidation()
		invalidRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product Description",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         0, // Zero value for required field
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por Length faltante
		// Assert - Verify error for missing Length
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "length")
	})

	// Test case: Invalid request - missing NetWeight (zero value)
	// Caso de prueba: Request inválido - falta NetWeight (valor cero)
	t.Run("Error_MissingNetWeight", func(t *testing.T) {
		// Arrange - Configurar validador y request con NetWeight en valor cero
		// Arrange - Set up validator and request with zero NetWeight
		validator := GetProductValidation()
		invalidRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product Description",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      0, // Zero value for required field
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por NetWeight faltante
		// Assert - Verify error for missing NetWeight
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "net_weight")
	})

	// Test case: Invalid request - missing ExpirationRate (zero value)
	// Caso de prueba: Request inválido - falta ExpirationRate (valor cero)
	t.Run("Error_MissingExpirationRate", func(t *testing.T) {
		// Arrange - Configurar validador y request con ExpirationRate en valor cero
		// Arrange - Set up validator and request with zero ExpirationRate
		validator := GetProductValidation()
		invalidRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product Description",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0, // Zero value for required field
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por ExpirationRate faltante
		// Assert - Verify error for missing ExpirationRate
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expiration_rate")
	})

	// Test case: Invalid request - missing RecommendedFreezingTemperature (zero value)
	// Caso de prueba: Request inválido - falta RecommendedFreezingTemperature (valor cero)
	t.Run("Error_MissingRecommendedFreezingTemperature", func(t *testing.T) {
		// Arrange - Configurar validador y request con RecommendedFreezingTemperature en valor cero
		// Arrange - Set up validator and request with zero RecommendedFreezingTemperature
		validator := GetProductValidation()
		invalidRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product Description",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: 0, // Zero value for required field
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por RecommendedFreezingTemperature faltante
		// Assert - Verify error for missing RecommendedFreezingTemperature
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "recommended_freezing_temperature")
	})

	// Test case: Invalid request - missing FreezingRate (zero value)
	// Caso de prueba: Request inválido - falta FreezingRate (valor cero)
	t.Run("Error_MissingFreezingRate", func(t *testing.T) {
		// Arrange - Configurar validador y request con FreezingRate en valor cero
		// Arrange - Set up validator and request with zero FreezingRate
		validator := GetProductValidation()
		invalidRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product Description",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0, // Zero value for required field
			ProductTypeID:                  1,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por FreezingRate faltante
		// Assert - Verify error for missing FreezingRate
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "freezing_rate")
	})

	// Test case: Invalid request - missing ProductTypeID (zero value)
	// Caso de prueba: Request inválido - falta ProductTypeID (valor cero)
	t.Run("Error_MissingProductTypeID", func(t *testing.T) {
		// Arrange - Configurar validador y request con ProductTypeID en valor cero
		// Arrange - Set up validator and request with zero ProductTypeID
		validator := GetProductValidation()
		invalidRequest := requests.ProductRequest{
			ProductCode:                    "PROD001",
			Description:                    "Test Product Description",
			Width:                          10.5,
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  0, // Zero value for required field
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(invalidRequest)

		// Assert - Verificar que hay error por ProductTypeID faltante
		// Assert - Verify error for missing ProductTypeID
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product_type_id")
	})

	// Test case: Invalid request - multiple missing fields
	// Caso de prueba: Request inválido - múltiples campos faltantes
	t.Run("Error_MultipleMissingFields", func(t *testing.T) {
		// Arrange - Configurar validador y request con múltiples campos faltantes
		// Arrange - Set up validator and request with multiple missing fields
		validator := GetProductValidation()
		invalidRequest := requests.ProductRequest{
			ProductCode:                    "", // Missing
			Description:                    "", // Missing
			Width:                          0,  // Missing
			Height:                         20.3,
			Length:                         15.7,
			NetWeight:                      5.2,
			ExpirationRate:                 0.1,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   0.8,
			ProductTypeID:                  1,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(invalidRequest)

		// Assert - Verificar que hay errores por múltiples campos faltantes
		// Assert - Verify errors for multiple missing fields
		assert.Error(t, err)
		errorMsg := err.Error()
		assert.Contains(t, errorMsg, "product_code")
		assert.Contains(t, errorMsg, "description")
		assert.Contains(t, errorMsg, "width")
	})

	// Test case: Valid request - SellerID is optional (confirm it doesn't cause validation error)
	// Caso de prueba: Request válido - SellerID es opcional (confirmar que no causa error de validación)
	t.Run("Success_SellerIDIsOptional", func(t *testing.T) {
		// Arrange - Configurar validador y request sin SellerID
		// Arrange - Set up validator and request without SellerID
		validator := GetProductValidation()
		validRequest := requests.ProductRequest{
			ProductCode:                    "PROD003",
			Description:                    "Test Product without Seller",
			Width:                          12.0,
			Height:                         25.0,
			Length:                         18.0,
			NetWeight:                      7.5,
			ExpirationRate:                 0.15,
			RecommendedFreezingTemperature: -22.0,
			FreezingRate:                   0.95,
			ProductTypeID:                  3,
			// SellerID is intentionally omitted (nil)
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(validRequest)

		// Assert - Verificar que no hay errores (SellerID es opcional)
		// Assert - Verify no errors (SellerID is optional)
		assert.NoError(t, err)
	})

	// Test case: Edge case - negative values for numeric fields (should pass validation)
	// Caso de prueba: Caso límite - valores negativos para campos numéricos (debe pasar validación)
	t.Run("Success_NegativeValuesAllowed", func(t *testing.T) {
		// Arrange - Configurar validador y request con valores negativos
		// Arrange - Set up validator and request with negative values
		validator := GetProductValidation()
		validRequest := requests.ProductRequest{
			ProductCode:                    "PROD004",
			Description:                    "Test Product with Negative Values",
			Width:                          -10.5, // Negative value should be allowed
			Height:                         -20.3, // Negative value should be allowed
			Length:                         -15.7, // Negative value should be allowed
			NetWeight:                      -5.2,  // Negative value should be allowed
			ExpirationRate:                 -0.1,  // Negative value should be allowed
			RecommendedFreezingTemperature: -30.0, // Typical freezing temperature (negative)
			FreezingRate:                   -0.8,  // Negative value should be allowed
			ProductTypeID:                  -1,    // Negative value should be allowed
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(validRequest)

		// Assert - Verificar que no hay errores (los valores negativos son válidos para la validación Required)
		// Assert - Verify no errors (negative values are valid for Required validation)
		assert.NoError(t, err)
	})

	// Test case: Edge case - very small positive values
	// Caso de prueba: Caso límite - valores positivos muy pequeños
	t.Run("Success_VerySmallPositiveValues", func(t *testing.T) {
		// Arrange - Configurar validador y request con valores muy pequeños
		// Arrange - Set up validator and request with very small values
		validator := GetProductValidation()
		validRequest := requests.ProductRequest{
			ProductCode:                    "PROD005",
			Description:                    "Test Product with Small Values",
			Width:                          0.001,
			Height:                         0.001,
			Length:                         0.001,
			NetWeight:                      0.001,
			ExpirationRate:                 0.001,
			RecommendedFreezingTemperature: 0.001,
			FreezingRate:                   0.001,
			ProductTypeID:                  1,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(validRequest)

		// Assert - Verificar que no hay errores
		// Assert - Verify no errors
		assert.NoError(t, err)
	})

	// Test case: Edge case - very large values
	// Caso de prueba: Caso límite - valores muy grandes
	t.Run("Success_VeryLargeValues", func(t *testing.T) {
		// Arrange - Configurar validador y request con valores muy grandes
		// Arrange - Set up validator and request with very large values
		validator := GetProductValidation()
		validRequest := requests.ProductRequest{
			ProductCode:                    "PROD006",
			Description:                    "Test Product with Large Values",
			Width:                          999999.999,
			Height:                         999999.999,
			Length:                         999999.999,
			NetWeight:                      999999.999,
			ExpirationRate:                 999999.999,
			RecommendedFreezingTemperature: 999999.999,
			FreezingRate:                   999999.999,
			ProductTypeID:                  999999,
		}

		// Act - Validar el request
		// Act - Validate the request
		err := validator.ValidateProductRequestStruct(validRequest)

		// Assert - Verificar que no hay errores
		// Assert - Verify no errors
		assert.NoError(t, err)
	})
}
