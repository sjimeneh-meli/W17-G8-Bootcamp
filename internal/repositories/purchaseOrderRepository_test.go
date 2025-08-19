// Package repositories_test - Purchase Order Repository Unit Tests / Tests Unitarios del Repositorio Purchase Order
// Database layer testing for purchase order CRUD operations with MySQL and order number uniqueness
// Testing de capa de datos para operaciones CRUD de órdenes de compra con MySQL y unicidad de número de orden
package repositories_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/repositories"
	"github.com/stretchr/testify/assert"
)

// Test_PurchaseOrder_MySqlPurchaseOrderRepositoryGetAll - Tests for GetAll method / Tests para método GetAll
// Cases: successful retrieval with data, empty result, database query error, row scanning error
// Casos: recuperación exitosa con datos, resultado vacío, error de consulta de base de datos, error de escaneo de filas
func Test_PurchaseOrder_MySqlPurchaseOrderRepositoryGetAll(t *testing.T) {
	t.Run("Successfully returns filled purchase orders map when there is data on db response", func(t *testing.T) {
		// Test: Database with purchase orders returns complete map / Base de datos con órdenes de compra retorna mapa completo
		expectedOrders := map[int]models.PurchaseOrder{
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

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "order_number", "order_date", "tracking_code", "buyer_id", "product_record_id"}).
			AddRow("1", "ORDER-001", time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), "TRACK-001", "1", "1").
			AddRow("2", "ORDER-002", time.Date(2024, 1, 16, 14, 45, 0, 0, time.UTC), "TRACK-002", "2", "2")

		mock.ExpectQuery("select id, order_number, order_date, tracking_code, buyer_id, product_record_id from purchase_orders").WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		dbOrders, err := repository.GetAll(context.Background())

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedOrders, dbOrders, "dbOrders should be a map with two purchase orders")
	})

	t.Run("Successfully returns empty purchase orders map when there isn't data on db response", func(t *testing.T) {
		// Test: Empty database returns empty map / Base de datos vacía retorna mapa vacío
		expectedOrders := map[int]models.PurchaseOrder{}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "order_number", "order_date", "tracking_code", "buyer_id", "product_record_id"})

		mock.ExpectQuery("select id, order_number, order_date, tracking_code, buyer_id, product_record_id from purchase_orders").WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		dbOrders, err := repository.GetAll(context.Background())

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedOrders, dbOrders, "dbOrders should be an empty map of purchase orders")
	})

	t.Run("Fails because of an internal server error querying the database", func(t *testing.T) {
		// Test: Database connection error during query / Error de conexión de base de datos durante consulta
		expectedError := error_message.ErrInternalServerError
		expectedOrders := map[int]models.PurchaseOrder{}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		mock.ExpectQuery("select id, order_number, order_date, tracking_code, buyer_id, product_record_id from purchase_orders").
			WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		dbOrders, err := repository.GetAll(context.Background())

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedOrders, dbOrders, "dbOrders should be an empty map of purchase orders")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})

	t.Run("Fails because of an internal server error scanning database results", func(t *testing.T) {
		// Test: Invalid data types cause row scanning errors / Tipos de datos inválidos causan errores de escaneo
		expectedOrders := map[int]models.PurchaseOrder{}
		expectedError := error_message.ErrInternalServerError
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "order_number", "order_date", "tracking_code", "buyer_id", "product_record_id"}).
			AddRow("1", "ORDER-001", time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), "TRACK-001", "1", "1").
			AddRow("STRING", "ORDER-002", time.Date(2024, 1, 16, 14, 45, 0, 0, time.UTC), "TRACK-002", "2", "2")

		mock.ExpectQuery("select id, order_number, order_date, tracking_code, buyer_id, product_record_id from purchase_orders").WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		dbOrders, err := repository.GetAll(context.Background())

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedOrders, dbOrders, "dbOrders should be an empty map of purchase orders")
		assert.ErrorIs(t, err, expectedError)
	})
}

// Test_PurchaseOrder_Create - Tests for Create method / Tests para método Create
// Cases: successful creation, internal server error executing insert, LastInsertId error
// Casos: creación exitosa, error interno del servidor ejecutando inserción, error de LastInsertId
func Test_PurchaseOrder_Create(t *testing.T) {
	t.Run("Successfully create a new purchase order record on db", func(t *testing.T) {
		// Test: Valid purchase order data creates new order with generated ID / Datos válidos de orden de compra crean nueva orden con ID generado
		expectedOrder := models.PurchaseOrder{
			Id:              17,
			OrderNumber:     "ORDER-017",
			OrderDate:       time.Date(2024, 1, 20, 9, 0, 0, 0, time.UTC),
			TrackingCode:    "TRACK-017",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		orderToInsert := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-017",
			OrderDate:       time.Date(2024, 1, 20, 9, 0, 0, 0, time.UTC),
			TrackingCode:    "TRACK-017",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta("insert into purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id) values (?, ?, ?, ?, ?)")).
			WithArgs(orderToInsert.OrderNumber, orderToInsert.OrderDate, orderToInsert.TrackingCode, orderToInsert.BuyerId, orderToInsert.ProductRecordId).
			WillReturnResult(sqlmock.NewResult(17, 1))

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		orderDb, err := repository.Create(context.Background(), orderToInsert)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedOrder, orderDb)
	})

	t.Run("Fails because of an internal server error executing insert", func(t *testing.T) {
		// Test: Database connection error during purchase order insertion / Error de conexión de base de datos durante inserción de orden de compra
		expectedOrder := models.PurchaseOrder{}
		expectedError := error_message.ErrInternalServerError

		orderToInsert := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-017",
			OrderDate:       time.Date(2024, 1, 20, 9, 0, 0, 0, time.UTC),
			TrackingCode:    "TRACK-017",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta("insert into purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id) values (?, ?, ?, ?, ?)")).
			WithArgs(orderToInsert.OrderNumber, orderToInsert.OrderDate, orderToInsert.TrackingCode, orderToInsert.BuyerId, orderToInsert.ProductRecordId).
			WillReturnError(errors.New("error executing query"))

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		orderDb, err := repository.Create(context.Background(), orderToInsert)

		assert.NotNil(t, err, "err should not be nil")
		assert.ErrorIs(t, err, expectedError)
		assert.Equal(t, expectedOrder, orderDb)
	})

	t.Run("Fails because of an internal server error getting the inserted record Id", func(t *testing.T) {
		// Test: Error retrieving generated ID after successful insert / Error al recuperar ID generado después de inserción exitosa
		expectedOrder := models.PurchaseOrder{}
		expectedError := error_message.ErrInternalServerError

		orderToInsert := models.PurchaseOrder{
			Id:              0,
			OrderNumber:     "ORDER-017",
			OrderDate:       time.Date(2024, 1, 20, 9, 0, 0, 0, time.UTC),
			TrackingCode:    "TRACK-017",
			BuyerId:         1,
			ProductRecordId: 1,
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta("insert into purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id) values (?, ?, ?, ?, ?)")).
			WithArgs(orderToInsert.OrderNumber, orderToInsert.OrderDate, orderToInsert.TrackingCode, orderToInsert.BuyerId, orderToInsert.ProductRecordId).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("error getting LastInsertId")))

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		orderDb, err := repository.Create(context.Background(), orderToInsert)

		assert.NotNil(t, err, "err should not be nil")
		assert.ErrorIs(t, err, expectedError)
		assert.Equal(t, expectedOrder, orderDb)
	})
}

// Test_PurchaseOrder_ExistPurchaseOrderByOrderNumber - Tests for ExistPurchaseOrderByOrderNumber method / Tests para método ExistPurchaseOrderByOrderNumber
// Cases: purchase order exists, purchase order doesn't exist, internal server error
// Casos: orden de compra existe, orden de compra no existe, error interno del servidor
func Test_PurchaseOrder_ExistPurchaseOrderByOrderNumber(t *testing.T) {
	t.Run("Successfully returns true when purchase order exists", func(t *testing.T) {
		// Test: Existing order number returns true / Número de orden existente retorna true
		expectedExists := true
		orderNumber := "ORDER-001"

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		row := mock.NewRows([]string{"1"}).AddRow(1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM purchase_orders WHERE order_number = ? LIMIT 1")).
			WithArgs(orderNumber).
			WillReturnRows(row)

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		exists, err := repository.ExistPurchaseOrderByOrderNumber(context.Background(), orderNumber)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedExists, exists)
	})

	t.Run("Successfully returns false when purchase order doesn't exist", func(t *testing.T) {
		// Test: Non-existent order number returns false / Número de orden inexistente retorna false
		expectedExists := false
		orderNumber := "ORDER-999"

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM purchase_orders WHERE order_number = ? LIMIT 1")).
			WithArgs(orderNumber).
			WillReturnError(sql.ErrNoRows)

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		exists, err := repository.ExistPurchaseOrderByOrderNumber(context.Background(), orderNumber)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedExists, exists)
	})

	t.Run("Fails because of an internal server error querying the database", func(t *testing.T) {
		// Test: Database connection error during existence check / Error de conexión de base de datos durante verificación de existencia
		expectedExists := false
		orderNumber := "ORDER-001"

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM purchase_orders WHERE order_number = ? LIMIT 1")).
			WithArgs(orderNumber).
			WillReturnError(errors.New("database connection error"))

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		exists, err := repository.ExistPurchaseOrderByOrderNumber(context.Background(), orderNumber)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedExists, exists)
	})
}

// Test_PurchaseOrder_GetPurchaseOrdersReportByBuyerId - Tests for GetPurchaseOrdersReportByBuyerId method / Tests para método GetPurchaseOrdersReportByBuyerId
// Cases: successful retrieval, buyer not found, internal server error querying, internal server error scanning
// Casos: recuperación exitosa, comprador no encontrado, error interno del servidor consultando, error interno del servidor escaneando
func Test_PurchaseOrder_GetPurchaseOrdersReportByBuyerId(t *testing.T) {
	t.Run("Successfully returns purchase order report for buyer", func(t *testing.T) {
		// Test: Valid buyer ID returns complete report / ID de comprador válido retorna reporte completo
		expectedReport := models.PurchaseOrderReport{
			Id:                 1,
			IdCardNumber:       "CARD-001",
			FirstName:          "Ignacio",
			LastName:           "Garcia",
			PurchaseOrderCount: 3,
		}

		buyerId := 1

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		row := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "purchase_orders_count"}).
			AddRow("1", "CARD-001", "Ignacio", "Garcia", "3")

		mock.ExpectQuery(regexp.QuoteMeta("select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as \"purchase_orders_count\" from productos_frescos.buyers b inner join productos_frescos.purchase_orders po on po.buyer_id = b.id where b.id = ? group by b.id")).
			WithArgs(buyerId).
			WillReturnRows(row)

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		report, err := repository.GetPurchaseOrdersReportByBuyerId(context.Background(), buyerId)

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedReport, report)
	})

	t.Run("Fails because buyer doesn't exist", func(t *testing.T) {
		// Test: Non-existent buyer ID returns not found error / ID de comprador inexistente retorna error not found
		expectedReport := models.PurchaseOrderReport{}
		expectedError := error_message.ErrNotFound
		buyerId := 999

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta("select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as \"purchase_orders_count\" from productos_frescos.buyers b inner join productos_frescos.purchase_orders po on po.buyer_id = b.id where b.id = ? group by b.id")).
			WithArgs(buyerId).
			WillReturnError(sql.ErrNoRows)

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		report, err := repository.GetPurchaseOrdersReportByBuyerId(context.Background(), buyerId)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, expectedError)
		assert.Equal(t, expectedReport, report)
	})

	t.Run("Fails because of an internal server error querying the database", func(t *testing.T) {
		// Test: Database connection error during query / Error de conexión de base de datos durante consulta
		expectedReport := models.PurchaseOrderReport{}
		expectedError := error_message.ErrInternalServerError
		buyerId := 1

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta("select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as \"purchase_orders_count\" from productos_frescos.buyers b inner join productos_frescos.purchase_orders po on po.buyer_id = b.id where b.id = ? group by b.id")).
			WithArgs(buyerId).
			WillReturnError(errors.New("database connection error"))

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		report, err := repository.GetPurchaseOrdersReportByBuyerId(context.Background(), buyerId)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, expectedError)
		assert.Equal(t, expectedReport, report)
	})

	t.Run("Fails because of an internal server error scanning database results", func(t *testing.T) {
		// Test: Invalid data types cause row scanning errors / Tipos de datos inválidos causan errores de escaneo
		expectedReport := models.PurchaseOrderReport{}
		expectedError := error_message.ErrInternalServerError
		buyerId := 1

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		row := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "purchase_orders_count"}).
			AddRow("STRING", "CARD-001", "Ignacio", "Garcia", "3")

		mock.ExpectQuery(regexp.QuoteMeta("select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as \"purchase_orders_count\" from productos_frescos.buyers b inner join productos_frescos.purchase_orders po on po.buyer_id = b.id where b.id = ? group by b.id")).
			WithArgs(buyerId).
			WillReturnRows(row)

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		report, err := repository.GetPurchaseOrdersReportByBuyerId(context.Background(), buyerId)

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.ErrorIs(t, err, expectedError)
		assert.Equal(t, expectedReport, report)
	})
}

// Test_PurchaseOrder_GetAllPurchaseOrdersReports - Tests for GetAllPurchaseOrdersReports method / Tests para método GetAllPurchaseOrdersReports
// Cases: successful retrieval with data, empty result, database query error, row scanning error
// Casos: recuperación exitosa con datos, resultado vacío, error de consulta de base de datos, error de escaneo de filas
func Test_PurchaseOrder_GetAllPurchaseOrdersReports(t *testing.T) {
	t.Run("Successfully returns filled purchase order reports when there is data on db response", func(t *testing.T) {
		// Test: Database with reports returns complete list / Base de datos con reportes retorna lista completa
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

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "purchase_orders_count"}).
			AddRow("1", "CARD-001", "Ignacio", "Garcia", "3").
			AddRow("2", "CARD-002", "Jesus", "Ortega", "2")

		mock.ExpectQuery(regexp.QuoteMeta("select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as \"purchase_orders_count\" from productos_frescos.buyers b inner join productos_frescos.purchase_orders po on po.buyer_id = b.id group by b.id order by b.id")).
			WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		reports, err := repository.GetAllPurchaseOrdersReports(context.Background())

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedReports, reports, "reports should be a slice with two purchase order reports")
	})

	t.Run("Successfully returns empty purchase order reports when there isn't data on db response", func(t *testing.T) {
		// Test: Empty database returns empty slice / Base de datos vacía retorna slice vacío
		expectedReports := []models.PurchaseOrderReport{}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "purchase_orders_count"})

		mock.ExpectQuery(regexp.QuoteMeta("select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as \"purchase_orders_count\" from productos_frescos.buyers b inner join productos_frescos.purchase_orders po on po.buyer_id = b.id group by b.id order by b.id")).
			WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		reports, err := repository.GetAllPurchaseOrdersReports(context.Background())

		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedReports, reports, "reports should be an empty slice of purchase order reports")
	})

	t.Run("Fails because of an internal server error querying the database", func(t *testing.T) {
		// Test: Database connection error during query / Error de conexión de base de datos durante consulta
		expectedError := error_message.ErrInternalServerError
		expectedReports := []models.PurchaseOrderReport{}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta("select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as \"purchase_orders_count\" from productos_frescos.buyers b inner join productos_frescos.purchase_orders po on po.buyer_id = b.id group by b.id order by b.id")).
			WillReturnError(error_message.ErrInternalServerError)

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		reports, err := repository.GetAllPurchaseOrdersReports(context.Background())

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedReports, reports, "reports should be an empty slice of purchase order reports")
		assert.ErrorIs(t, err, expectedError, "err should be of type ErrInternalServerError")
	})

	t.Run("Fails because of an internal server error scanning database results", func(t *testing.T) {
		// Test: Invalid data types cause row scanning errors / Tipos de datos inválidos causan errores de escaneo
		expectedReports := []models.PurchaseOrderReport{}
		expectedError := error_message.ErrInternalServerError

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		rows := mock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "purchase_orders_count"}).
			AddRow("1", "CARD-001", "Ignacio", "Garcia", "3").
			AddRow("STRING", "CARD-002", "Jesus", "Ortega", "2")

		mock.ExpectQuery(regexp.QuoteMeta("select b.id, b.id_card_number, b.first_name, b.last_name, count(po.id) as \"purchase_orders_count\" from productos_frescos.buyers b inner join productos_frescos.purchase_orders po on po.buyer_id = b.id group by b.id order by b.id")).
			WillReturnRows(rows).RowsWillBeClosed()

		repository := repositories.MySqlPurchaseOrderRepository{
			Db: db,
		}

		reports, err := repository.GetAllPurchaseOrdersReports(context.Background())

		assert.NotNil(t, err, "err shouldn't be nil")
		assert.Equal(t, expectedReports, reports, "reports should be an empty slice of purchase order reports")
		assert.ErrorIs(t, err, expectedError)
	})
}

// Test_PurchaseOrder_GetNewPurchaseOrderMySQLRepository - Tests for repository constructor / Tests para constructor del repositorio
// Cases: successful creation, singleton pattern validation
// Casos: creación exitosa, validación de patrón singleton
func Test_PurchaseOrder_GetNewPurchaseOrderMySQLRepository(t *testing.T) {
	t.Run("Successfully returns a new purchase order repository", func(t *testing.T) {
		// Test: Constructor creates valid repository instance / Constructor crea instancia válida del repositorio
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		repository := repositories.GetNewPurchaseOrderMySQLRepository(db)

		assert.NotNil(t, repository, "repository shouldn't be nil")
	})

	t.Run("Multiple calls to GetNewPurchaseOrderMySQLRepository should return the same instance", func(t *testing.T) {
		// Test: Singleton pattern ensures same instance is returned / Patrón singleton asegura que se retorna la misma instancia
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to open sqlmock database: %v", err)
		}
		defer db.Close()

		repository := repositories.GetNewPurchaseOrderMySQLRepository(db)

		assert.NotNil(t, repository, "repository shouldn't be nil")

		repository2 := repositories.GetNewPurchaseOrderMySQLRepository(nil)

		assert.Equal(t, repository, repository2, "repository and repository2 should be the same")
	})
}
