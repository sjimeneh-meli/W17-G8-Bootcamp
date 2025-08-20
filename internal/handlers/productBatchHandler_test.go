package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/seeders"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostProductBatch(t *testing.T) {

	t.Run("Entity: ProductBatch, Method: POST, Code: 417", func(t *testing.T) {

		expectedResponseCode := 417
		expectedResponseBody := `{
			"message":"json: cannot unmarshal number into Go struct field ProductBatchRequest.batch_number of type string", 
			"status":"Expectation Failed"
		}`

		actualRequest := strings.NewReader(`{
			"batch_number": 123,
			"current_quantity": 1,
			"current_temperature": 3.43,
			"due_date": "2025-08-19",
			"initial_quantity": 1,
			"manufacturing_date": "2025-08-19",
			"manufacturing_hour": 1,
			"minimum_temperature": 3.43,
			"product_id": 1,
			"section_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/productBatches", actualRequest)
		response := httptest.NewRecorder()
		srvMock := tests.GetProductBatchServiceMock()
		sectionSrvMock := tests.GetSectionServiceMock()
		productServiceMock := &tests.ProductServiceMock{}
		vldMock := tests.GetProductBatchValidationMock()

		handler := handlers.GetProductBatchHandler(srvMock, sectionSrvMock, productServiceMock, vldMock)

		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: ProductBatch, Method: POST, Code: 422", func(t *testing.T) {

		expectedResponseCode := 422
		expectedResponseBody := `{
			"message":"batch_number: cannot be blank.", 
			"status":"Unprocessable Entity"
		}`

		actualRequest := strings.NewReader(`{
			"current_quantity": 1,
			"current_temperature": 3.43,
			"due_date": "2025-08-19",
			"initial_quantity": 1,
			"manufacturing_date": "2025-08-19",
			"manufacturing_hour": 1,
			"minimum_temperature": 3.43,
			"product_id": 1,
			"section_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/productBatches", actualRequest)
		response := httptest.NewRecorder()
		srvMock := tests.GetProductBatchServiceMock()
		sectionSrvMock := tests.GetSectionServiceMock()
		productServiceMock := &tests.ProductServiceMock{}
		vldMock := tests.GetProductBatchValidationMock()

		vldMock.On("ValidateProductBatchRequestStruc", seeders.NewProductBatchRequestWithoutNumber).Return(errors.New("batch_number: cannot be blank."))
		handler := handlers.GetProductBatchHandler(srvMock, sectionSrvMock, productServiceMock, vldMock)

		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: ProductBatch, Method: POST, Code: 404 ( Section Service)", func(t *testing.T) {

		expectedResponseCode := 404
		expectedResponseBody := `{
			"message":"section not found", 
			"status":"Not Found"
		}`

		actualRequest := strings.NewReader(`{
			"batch_number": "AAA",
			"current_quantity": 1,
			"current_temperature": 3.43,
			"due_date": "2025-08-19",
			"initial_quantity": 1,
			"manufacturing_date": "2025-08-19",
			"manufacturing_hour": 1,
			"minimum_temperature": 3.43,
			"product_id": 1,
			"section_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/productBatches", actualRequest)
		response := httptest.NewRecorder()
		srvMock := tests.GetProductBatchServiceMock()
		sectionSrvMock := tests.GetSectionServiceMock()
		productServiceMock := &tests.ProductServiceMock{}
		vldMock := tests.GetProductBatchValidationMock()

		vldMock.On("ValidateProductBatchRequestStruc", seeders.NewProductBatchRequest).Return(nil)
		sectionSrvMock.On("ExistWithID", mock.AnythingOfType("*context.timerCtx"), seeders.NewProductBatchRequest.SectionID).Return(false)
		handler := handlers.GetProductBatchHandler(srvMock, sectionSrvMock, productServiceMock, vldMock)

		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: ProductBatch, Method: POST, Code: 404 (Product)", func(t *testing.T) {

		expectedResponseCode := 404
		expectedResponseBody := `{
			"message":"product not found", 
			"status":"Not Found"
		}`

		actualRequest := strings.NewReader(`{
			"batch_number": "AAA",
			"current_quantity": 1,
			"current_temperature": 3.43,
			"due_date": "2025-08-19",
			"initial_quantity": 1,
			"manufacturing_date": "2025-08-19",
			"manufacturing_hour": 1,
			"minimum_temperature": 3.43,
			"product_id": 1,
			"section_id": 5
		}`)

		request := httptest.NewRequest(http.MethodPost, "/productBatches", actualRequest)
		response := httptest.NewRecorder()
		srvMock := tests.GetProductBatchServiceMock()
		sectionSrvMock := tests.GetSectionServiceMock()
		productSrvMock := &tests.ProductServiceMock{}
		vldMock := tests.GetProductBatchValidationMock()

		vldMock.On("ValidateProductBatchRequestStruc", seeders.NewProductBatchRequestTwo).Return(nil)
		sectionSrvMock.On("ExistWithID", mock.AnythingOfType("*context.timerCtx"), seeders.NewProductBatchRequestTwo.SectionID).Return(false)
		productSrvMock.On("ExistById", mock.AnythingOfType("*context.timerCtx"), int64(seeders.NewProductBatchRequestTwo.ProductID)).Return(false, nil)
		handler := handlers.GetProductBatchHandler(srvMock, sectionSrvMock, productSrvMock, vldMock)

		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: ProductBatch, Method: POST, Code: 404 (Service)", func(t *testing.T) {

		expectedResponseCode := 404
		expectedResponseBody := `{
			"message":"section not found", 
			"status":"Not Found"
		}`

		actualRequest := strings.NewReader(`{
			"batch_number": "AAA",
			"current_quantity": 1,
			"current_temperature": 3.43,
			"due_date": "2025-08-19",
			"initial_quantity": 1,
			"manufacturing_date": "2025-08-19",
			"manufacturing_hour": 1,
			"minimum_temperature": 3.43,
			"product_id": 1,
			"section_id": 1
		}`)

		request := httptest.NewRequest(http.MethodPost, "/productBatches", actualRequest)
		response := httptest.NewRecorder()
		srvMock := tests.GetProductBatchServiceMock()
		sectionSrvMock := tests.GetSectionServiceMock()
		productServiceMock := &tests.ProductServiceMock{}
		vldMock := tests.GetProductBatchValidationMock()

		vldMock.On("ValidateProductBatchRequestStruc", seeders.NewProductBatchRequest).Return(nil)
		sectionSrvMock.On("ExistWithID", mock.AnythingOfType("*context.timerCtx"), seeders.NewProductBatchRequest.SectionID).Return(false)
		handler := handlers.GetProductBatchHandler(srvMock, sectionSrvMock, productServiceMock, vldMock)

		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: ProductBatch, Method: POST, Code: 409", func(t *testing.T) {

		expectedResponseCode := 409
		expectedResponseBody := `{
			"message":"already exist a batch with the same number", 
			"status":"Conflict"
		}`

		actualRequest := strings.NewReader(`{
			"batch_number": "AAA",
			"current_quantity": 1,
			"current_temperature": 3.43,
			"due_date": "2025-08-19",
			"initial_quantity": 1,
			"manufacturing_date": "2025-08-19",
			"manufacturing_hour": 1,
			"minimum_temperature": 3.43,
			"product_id": 1,
			"section_id": 5
		}`)

		request := httptest.NewRequest(http.MethodPost, "/productBatches", actualRequest)
		response := httptest.NewRecorder()
		srvMock := tests.GetProductBatchServiceMock()
		sectionSrvMock := tests.GetSectionServiceMock()
		productSrvMock := &tests.ProductServiceMock{}
		vldMock := tests.GetProductBatchValidationMock()

		vldMock.On("ValidateProductBatchRequestStruc", seeders.NewProductBatchRequestTwo).Return(nil)
		sectionSrvMock.On("ExistWithID", mock.AnythingOfType("*context.timerCtx"), seeders.NewProductBatchRequestTwo.SectionID).Return(false)
		productSrvMock.On("ExistById", mock.AnythingOfType("*context.timerCtx"), int64(seeders.NewProductBatchRequestTwo.ProductID)).Return(true, nil)
		srvMock.On("ExistsWithBatchNumber", mock.AnythingOfType("*context.timerCtx"), 0, seeders.NewProductBatchRequestTwo.BatchNumber).Return(true)
		handler := handlers.GetProductBatchHandler(srvMock, sectionSrvMock, productSrvMock, vldMock)

		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: ProductBatch, Method: POST, Code: 417 Batch Number", func(t *testing.T) {

		expectedResponseCode := 417
		expectedResponseBody := `{
			"message":"already exist a batch with the same number", 
			"status":"Expectation Failed"
		}`

		actualRequest := strings.NewReader(`{
			"batch_number": "AAA",
			"current_quantity": 1,
			"current_temperature": 3.43,
			"due_date": "2025-08-19",
			"initial_quantity": 1,
			"manufacturing_date": "2025-08-19",
			"manufacturing_hour": 0,
			"minimum_temperature": 3.43,
			"product_id": 1,
			"section_id": 5
		}`)

		request := httptest.NewRequest(http.MethodPost, "/productBatches", actualRequest)
		response := httptest.NewRecorder()
		srvMock := tests.GetProductBatchServiceMock()
		sectionSrvMock := tests.GetSectionServiceMock()
		productSrvMock := &tests.ProductServiceMock{}
		vldMock := tests.GetProductBatchValidationMock()

		vldMock.On("ValidateProductBatchRequestStruc", seeders.NewProductBatchRequestThree).Return(nil)
		sectionSrvMock.On("ExistWithID", mock.AnythingOfType("*context.timerCtx"), seeders.NewProductBatchRequestThree.SectionID).Return(false)
		productSrvMock.On("ExistById", mock.AnythingOfType("*context.timerCtx"), int64(seeders.NewProductBatchRequestThree.ProductID)).Return(true, nil)
		srvMock.On("ExistsWithBatchNumber", mock.AnythingOfType("*context.timerCtx"), 0, seeders.NewProductBatchRequestThree.BatchNumber).Return(false)
		srvMock.On("Create", mock.AnythingOfType("*context.timerCtx"), seeders.NewProductBatchModelThree).Return(errors.New("already exist a batch with the same number"))
		handler := handlers.GetProductBatchHandler(srvMock, sectionSrvMock, productSrvMock, vldMock)

		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: ProductBatch, Method: POST, Code: 201", func(t *testing.T) {

		expectedResponseCode := 201
		expectedResponseBody := `{
			"data": {
				"batch_number":"AAA", 
				"current_quantity":1, 
				"current_temperature":3.43, 
				"due_date":"2025-08-19 00:00:00 +0000 UTC", 
				"id":0, 
				"initial_quantity":1, 
				"manufacturing_date":"2025-08-19 00:00:00 +0000 UTC", 
				"manufacturing_hour":0, 
				"minimum_temperature":3.43, 
				"product_id":1, 
				"section_id":5
			}
		}`

		actualRequest := strings.NewReader(`{
			"batch_number": "AAA",
			"current_quantity": 1,
			"current_temperature": 3.43,
			"due_date": "2025-08-19",
			"initial_quantity": 1,
			"manufacturing_date": "2025-08-19",
			"manufacturing_hour": 0,
			"minimum_temperature": 3.43,
			"product_id": 1,
			"section_id": 5
		}`)

		request := httptest.NewRequest(http.MethodPost, "/productBatches", actualRequest)
		response := httptest.NewRecorder()
		srvMock := tests.GetProductBatchServiceMock()
		sectionSrvMock := tests.GetSectionServiceMock()
		productSrvMock := &tests.ProductServiceMock{}
		vldMock := tests.GetProductBatchValidationMock()

		vldMock.On("ValidateProductBatchRequestStruc", seeders.NewProductBatchRequestThree).Return(nil)
		sectionSrvMock.On("ExistWithID", mock.AnythingOfType("*context.timerCtx"), seeders.NewProductBatchRequestThree.SectionID).Return(false)
		productSrvMock.On("ExistById", mock.AnythingOfType("*context.timerCtx"), int64(seeders.NewProductBatchRequestThree.ProductID)).Return(true, nil)
		srvMock.On("ExistsWithBatchNumber", mock.AnythingOfType("*context.timerCtx"), 0, seeders.NewProductBatchRequestThree.BatchNumber).Return(false)
		srvMock.On("Create", mock.AnythingOfType("*context.timerCtx"), seeders.NewProductBatchModelThree).Return(nil)
		handler := handlers.GetProductBatchHandler(srvMock, sectionSrvMock, productSrvMock, vldMock)

		handler.Create(response, request)

		assert.Equal(t, expectedResponseCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}

func TestGetProductBatch(t *testing.T) {
	t.Run("Entity: ProductBatch, Method: GET, Code: 417", func(t *testing.T) {
		expectedCode := 417
		id := "100a"
		expectedResponseBody := `{
			"message":"strconv.Atoi: parsing \"100a\": invalid syntax", 
			"status":"Expectation Failed"
		}`

		request, _ := http.NewRequest(http.MethodGet, "/sections/reportProducts", nil)
		response := httptest.NewRecorder()
		query := request.URL.Query()
		query.Add("id", id)
		request.URL.RawQuery = query.Encode()

		srvMock := tests.GetProductBatchServiceMock()
		sectionSrvMock := tests.GetSectionServiceMock()
		productSrvMock := &tests.ProductServiceMock{}
		vldMock := tests.GetProductBatchValidationMock()

		handler := handlers.GetProductBatchHandler(srvMock, sectionSrvMock, productSrvMock, vldMock)

		handler.GetReportProduct(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: ProductBatch, Method: GET, Code: 404", func(t *testing.T) {
		expectedCode := 404
		id := "100"
		expectedResponseBody := `{
			"message":"section not found", 
			"status":"Not Found"
		}`

		request, _ := http.NewRequest(http.MethodGet, "/sections/reportProducts", nil)
		query := request.URL.Query()
		query.Add("id", id)
		request.URL.RawQuery = query.Encode()
		response := httptest.NewRecorder()

		srvMock := tests.GetProductBatchServiceMock()
		sectionSrvMock := tests.GetSectionServiceMock()
		productSrvMock := &tests.ProductServiceMock{}
		vldMock := tests.GetProductBatchValidationMock()

		sectionSrvMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), 100).Return(&models.Section{}, errors.New("section not found"))
		handler := handlers.GetProductBatchHandler(srvMock, sectionSrvMock, productSrvMock, vldMock)

		handler.GetReportProduct(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: ProductBatch, Method: GET, Code: 201", func(t *testing.T) {
		expectedCode := 201
		id := "0"
		expectedResponseBody := `{
			"data": {
				"section_id":0, 
				"section_number":"A-01",
				"products_count": 1
			}
		}`

		request, _ := http.NewRequest(http.MethodGet, "/sections/reportProducts", nil)
		query := request.URL.Query()
		query.Add("id", id)
		request.URL.RawQuery = query.Encode()
		response := httptest.NewRecorder()

		srvMock := tests.GetProductBatchServiceMock()
		sectionSrvMock := tests.GetSectionServiceMock()
		productSrvMock := &tests.ProductServiceMock{}
		vldMock := tests.GetProductBatchValidationMock()

		sectionSrvMock.On("GetByID", mock.AnythingOfType("*context.timerCtx"), 0).Return(seeders.NewSectionModel, nil)
		srvMock.On("GetProductQuantityBySectionId", mock.AnythingOfType("*context.timerCtx"), 0).Return(1, nil)
		handler := handlers.GetProductBatchHandler(srvMock, sectionSrvMock, productSrvMock, vldMock)

		handler.GetReportProduct(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})

	t.Run("Entity: ProductBatch, Method: GET, Code: 417 (No Params)", func(t *testing.T) {
		expectedCode := 417
		expectedResponseBody := `{
			"message":"No data", 
			"status":"Expectation Failed"
		}`

		request := httptest.NewRequest(http.MethodGet, "/sections/reportProducts", nil)
		response := httptest.NewRecorder()

		srvMock := tests.GetProductBatchServiceMock()
		sectionSrvMock := tests.GetSectionServiceMock()
		productSrvMock := &tests.ProductServiceMock{}
		vldMock := tests.GetProductBatchValidationMock()

		sectionSrvMock.On("GetAll", mock.AnythingOfType("*context.timerCtx")).Return([]*models.Section{}, errors.New("No data"))

		handler := handlers.GetProductBatchHandler(srvMock, sectionSrvMock, productSrvMock, vldMock)

		handler.GetReportProduct(response, request)

		assert.Equal(t, expectedCode, response.Code)
		assert.JSONEq(t, expectedResponseBody, response.Body.String())
	})
}
