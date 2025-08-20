package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProductRecordRepository(t *testing.T) {
	// Reset singleton instance for testing
	productRecordInstance = nil

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	t.Run("should create new instance when none exists", func(t *testing.T) {
		// Reset singleton
		productRecordInstance = nil

		repo := NewProductRecordRepository(db)
		assert.NotNil(t, repo)
		assert.Equal(t, productRecordInstance, repo)
	})

	t.Run("should return existing instance when already created", func(t *testing.T) {
		// First call creates the instance
		repo1 := NewProductRecordRepository(db)
		// Second call should return the same instance
		repo2 := NewProductRecordRepository(db)

		assert.Equal(t, repo1, repo2)
		assert.Same(t, repo1, repo2)
	})
}

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock
}

func TestProductRecordRepository_Create(t *testing.T) {
	testCases := []struct {
		name          string
		productRecord *models.ProductRecord
		setupMock     func(sqlmock.Sqlmock)
		expectedID    int
		expectedError string
	}{
		{
			name: "should create product record successfully",
			productRecord: &models.ProductRecord{
				LastUpdateDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				PurchasePrice:  10.50,
				SalePrice:      15.75,
				ProductID:      123,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES (?, ?, ?, ?)")).
					WithArgs(
						time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
						10.50,
						15.75,
						int64(123),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedID: 1,
		},
		{
			name: "should return error when insert fails",
			productRecord: &models.ProductRecord{
				LastUpdateDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				PurchasePrice:  10.50,
				SalePrice:      15.75,
				ProductID:      123,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES (?, ?, ?, ?)")).
					WithArgs(
						time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
						10.50,
						15.75,
						int64(123),
					).
					WillReturnError(fmt.Errorf("database connection error"))
			},
			expectedError: "error to create product record",
		},
		{
			name: "should return error when getting last insert id fails",
			productRecord: &models.ProductRecord{
				LastUpdateDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				PurchasePrice:  10.50,
				SalePrice:      15.75,
				ProductID:      123,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES (?, ?, ?, ?)")).
					WithArgs(
						time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
						10.50,
						15.75,
						int64(123),
					).
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("last insert id error")))
			},
			expectedError: "error to get last insert id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset singleton for each test
			productRecordInstance = nil

			db, mock := setupMockDB(t)
			defer db.Close()

			tc.setupMock(mock)

			repo := NewProductRecordRepository(db)
			result, err := repo.Create(context.Background(), tc.productRecord)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedID, result.ID)
				assert.Equal(t, tc.productRecord.LastUpdateDate, result.LastUpdateDate)
				assert.Equal(t, tc.productRecord.PurchasePrice, result.PurchasePrice)
				assert.Equal(t, tc.productRecord.SalePrice, result.SalePrice)
				assert.Equal(t, tc.productRecord.ProductID, result.ProductID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProductRecordRepository_GetReportByIdProduct(t *testing.T) {
	testCases := []struct {
		name           string
		productID      int64
		setupMock      func(sqlmock.Sqlmock)
		expectedReport *models.ProductRecordReport
		expectedError  error
	}{
		{
			name:      "should return report for existing product",
			productID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"product_id", "description", "records_count"}).
					AddRow(1, "Test Product", 5)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id as product_id, p.description, count(*) as records_count
	FROM 
		products as p
	LEFT JOIN
		product_records as pr
	ON pr.product_id = p.id
	WHERE p.id = ?
	GROUP BY p.id, p.description`)).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectedReport: &models.ProductRecordReport{
				ProductId:    1,
				Description:  "Test Product",
				RecordsCount: 5,
			},
		},
		{
			name:      "should return error when product not found",
			productID: 999,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id as product_id, p.description, count(*) as records_count
	FROM 
		products as p
	LEFT JOIN
		product_records as pr
	ON pr.product_id = p.id
	WHERE p.id = ?
	GROUP BY p.id, p.description`)).
					WithArgs(int64(999)).
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: error_message.ErrDependencyNotFound,
		},
		{
			name:      "should return error when database query fails",
			productID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id as product_id, p.description, count(*) as records_count
	FROM 
		products as p
	LEFT JOIN
		product_records as pr
	ON pr.product_id = p.id
	WHERE p.id = ?
	GROUP BY p.id, p.description`)).
					WithArgs(int64(1)).
					WillReturnError(fmt.Errorf("database connection error"))
			},
			expectedError: fmt.Errorf("error to scan product record report"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset singleton for each test
			productRecordInstance = nil

			db, mock := setupMockDB(t)
			defer db.Close()

			tc.setupMock(mock)

			repo := NewProductRecordRepository(db)
			result, err := repo.GetReportByIdProduct(context.Background(), tc.productID)

			if tc.expectedError != nil {
				require.Error(t, err)
				if tc.expectedError == error_message.ErrDependencyNotFound {
					assert.Equal(t, tc.expectedError, err)
				} else {
					assert.Contains(t, err.Error(), tc.expectedError.Error())
				}
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedReport.ProductId, result.ProductId)
				assert.Equal(t, tc.expectedReport.Description, result.Description)
				assert.Equal(t, tc.expectedReport.RecordsCount, result.RecordsCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProductRecordRepository_GetReport(t *testing.T) {
	testCases := []struct {
		name            string
		setupMock       func(sqlmock.Sqlmock)
		expectedReports []*models.ProductRecordReport
		expectedError   string
	}{
		{
			name: "should return reports for all products",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"product_id", "description", "records_count"}).
					AddRow(1, "Product 1", 3).
					AddRow(2, "Product 2", 0).
					AddRow(3, "Product 3", 7)

				mock.ExpectPrepare(regexp.QuoteMeta(`SELECT
            p.id as product_id,
            p.description,
            count(pr.product_id) as records_count
        FROM
            products p
        LEFT JOIN
            product_records pr ON pr.product_id = p.id
        GROUP BY
            p.id, p.description`)).
					ExpectQuery().
					WillReturnRows(rows)
			},
			expectedReports: []*models.ProductRecordReport{
				{ProductId: 1, Description: "Product 1", RecordsCount: 3},
				{ProductId: 2, Description: "Product 2", RecordsCount: 0},
				{ProductId: 3, Description: "Product 3", RecordsCount: 7},
			},
		},
		{
			name: "should return empty slice when no products exist",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"product_id", "description", "records_count"})

				mock.ExpectPrepare(regexp.QuoteMeta(`SELECT
            p.id as product_id,
            p.description,
            count(pr.product_id) as records_count
        FROM
            products p
        LEFT JOIN
            product_records pr ON pr.product_id = p.id
        GROUP BY
            p.id, p.description`)).
					ExpectQuery().
					WillReturnRows(rows)
			},
			expectedReports: []*models.ProductRecordReport{},
		},
		{
			name: "should return error when prepare fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta(`SELECT
            p.id as product_id,
            p.description,
            count(pr.product_id) as records_count
        FROM
            products p
        LEFT JOIN
            product_records pr ON pr.product_id = p.id
        GROUP BY
            p.id, p.description`)).
					WillReturnError(fmt.Errorf("prepare statement error"))
			},
			expectedError: "failed to prepare product report query",
		},
		{
			name: "should return error when query execution fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta(`SELECT
            p.id as product_id,
            p.description,
            count(pr.product_id) as records_count
        FROM
            products p
        LEFT JOIN
            product_records pr ON pr.product_id = p.id
        GROUP BY
            p.id, p.description`)).
					ExpectQuery().
					WillReturnError(fmt.Errorf("query execution error"))
			},
			expectedError: "failed to execute product report query",
		},
		{
			name: "should return error when scan fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"product_id", "description", "records_count"}).
					AddRow("invalid", "Product 1", 3) // Invalid type for product_id

				mock.ExpectPrepare(regexp.QuoteMeta(`SELECT
            p.id as product_id,
            p.description,
            count(pr.product_id) as records_count
        FROM
            products p
        LEFT JOIN
            product_records pr ON pr.product_id = p.id
        GROUP BY
            p.id, p.description`)).
					ExpectQuery().
					WillReturnRows(rows)
			},
			expectedError: "failed to scan product record row",
		},
		{
			name: "should return error when rows.Err() returns error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"product_id", "description", "records_count"}).
					AddRow(1, "Product 1", 3).
					RowError(0, fmt.Errorf("row iteration error"))

				mock.ExpectPrepare(regexp.QuoteMeta(`SELECT
            p.id as product_id,
            p.description,
            count(pr.product_id) as records_count
        FROM
            products p
        LEFT JOIN
            product_records pr ON pr.product_id = p.id
        GROUP BY
            p.id, p.description`)).
					ExpectQuery().
					WillReturnRows(rows)
			},
			expectedError: "error iterating product record rows",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset singleton for each test
			productRecordInstance = nil

			db, mock := setupMockDB(t)
			defer db.Close()

			tc.setupMock(mock)

			repo := NewProductRecordRepository(db)
			result, err := repo.GetReport(context.Background())

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tc.expectedReports))

				for i, expectedReport := range tc.expectedReports {
					assert.Equal(t, expectedReport.ProductId, result[i].ProductId)
					assert.Equal(t, expectedReport.Description, result[i].Description)
					assert.Equal(t, expectedReport.RecordsCount, result[i].RecordsCount)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProductRecordRepository_ExistProductRecordByID(t *testing.T) {
	testCases := []struct {
		name           string
		productID      int64
		setupMock      func(sqlmock.Sqlmock)
		expectedExists bool
	}{
		{
			name:      "should return true when product record exists",
			productID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM product_records WHERE id = ? LIMIT 1;")).
					WithArgs(int64(1)).
					WillReturnRows(rows)
			},
			expectedExists: true,
		},
		{
			name:      "should return false when product record does not exist",
			productID: 999,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM product_records WHERE id = ? LIMIT 1;")).
					WithArgs(int64(999)).
					WillReturnError(sql.ErrNoRows)
			},
			expectedExists: false,
		},
		{
			name:      "should return false when database query fails",
			productID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 FROM product_records WHERE id = ? LIMIT 1;")).
					WithArgs(int64(1)).
					WillReturnError(fmt.Errorf("database connection error"))
			},
			expectedExists: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset singleton for each test
			productRecordInstance = nil

			db, mock := setupMockDB(t)
			defer db.Close()

			tc.setupMock(mock)

			repo := NewProductRecordRepository(db)
			result := repo.ExistProductRecordByID(context.Background(), tc.productID)

			assert.Equal(t, tc.expectedExists, result)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Test para verificar que las interfaces están implementadas correctamente
func TestProductRecordRepository_ImplementsInterface(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Reset singleton
	productRecordInstance = nil

	repo := NewProductRecordRepository(db)

	// Verifica que el repositorio implementa la interfaz
	var _ IProductRecordRepository = repo
}

// Tests adicionales para cobertura completa de edge cases
func TestProductRecordRepository_EdgeCases(t *testing.T) {
	t.Run("Create with context cancellation", func(t *testing.T) {
		// Reset singleton
		productRecordInstance = nil

		db, mock := setupMockDB(t)
		defer db.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel context immediately

		productRecord := &models.ProductRecord{
			LastUpdateDate: time.Now(),
			PurchasePrice:  10.50,
			SalePrice:      15.75,
			ProductID:      123,
		}

		// Mock should expect query but context cancellation should prevent execution
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES (?, ?, ?, ?)")).
			WithArgs(sqlmock.AnyArg(), 10.50, 15.75, int64(123)).
			WillReturnError(context.Canceled)

		repo := NewProductRecordRepository(db)
		result, err := repo.Create(ctx, productRecord)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "error to create product record")
	})

	t.Run("GetReportByIdProduct with zero records", func(t *testing.T) {
		// Reset singleton
		productRecordInstance = nil

		db, mock := setupMockDB(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"product_id", "description", "records_count"}).
			AddRow(1, "Product with no records", 0)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id as product_id, p.description, count(*) as records_count
	FROM 
		products as p
	LEFT JOIN
		product_records as pr
	ON pr.product_id = p.id
	WHERE p.id = ?
	GROUP BY p.id, p.description`)).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		repo := NewProductRecordRepository(db)
		result, err := repo.GetReportByIdProduct(context.Background(), 1)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ProductId)
		assert.Equal(t, "Product with no records", result.Description)
		assert.Equal(t, int64(0), result.RecordsCount)
	})
}
