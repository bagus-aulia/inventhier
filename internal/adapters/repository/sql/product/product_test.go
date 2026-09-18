package product_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bagus-aulia/inventhier/internal/adapters/repository/sql/product"
	dto "github.com/bagus-aulia/inventhier/internal/core/dto/product"
	"github.com/stretchr/testify/assert"
)

var (
	mockProduct = &dto.Product{
		ID:        1,
		SKU:       "SKU-001",
		Name:      "Test Product",
		Price:     100000,
		Supplier:  "Supplier A",
		Stock:     50,
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
)

// TestNewSQLRepository tests the constructor
func TestNewSQLRepository(t *testing.T) {
	t.Run("creates new SQL repository instance", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		assert.NotNil(t, repo)
	})
}

// TestGetProductBySKU tests the GetProductBySKU method
func TestGetProductBySKU(t *testing.T) {
	ctx := context.Background()

	t.Run("returns product successfully when SKU exists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Mock the query
		rows := sqlmock.NewRows([]string{
			"id", "sku", "name", "price", "supplier", "stock", "created_at", "updated_at",
		}).AddRow(
			mockProduct.ID,
			mockProduct.SKU,
			mockProduct.Name,
			mockProduct.Price,
			mockProduct.Supplier,
			mockProduct.Stock,
			mockProduct.CreatedAt,
			mockProduct.UpdatedAt,
		)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, sku, name, price, supplier, stock, created_at, updated_at FROM products WHERE sku = ?")).
			WithArgs("SKU-001").
			WillReturnRows(rows)

		result, err := repo.GetProductBySKU(ctx, "SKU-001")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockProduct.SKU, result.SKU)
		assert.Equal(t, mockProduct.Name, result.Name)
		assert.Equal(t, mockProduct.Price, result.Price)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when SKU does not exist", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Mock the query to return no rows
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, sku, name, price, supplier, stock, created_at, updated_at FROM products WHERE sku = ?")).
			WithArgs("NONEXISTENT").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetProductBySKU(ctx, "NONEXISTENT")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when database query fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		dbErr := errors.New("database connection error")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, sku, name, price, supplier, stock, created_at, updated_at FROM products WHERE sku = ?")).
			WithArgs("SKU-001").
			WillReturnError(dbErr)

		result, err := repo.GetProductBySKU(ctx, "SKU-001")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, dbErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns correct product data with all fields", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		expectedProduct := &dto.Product{
			ID:        2,
			SKU:       "SKU-002",
			Name:      "Premium Product",
			Price:     250000,
			Supplier:  "Supplier B",
			Stock:     100,
			CreatedAt: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC),
		}

		rows := sqlmock.NewRows([]string{
			"id", "sku", "name", "price", "supplier", "stock", "created_at", "updated_at",
		}).AddRow(
			expectedProduct.ID,
			expectedProduct.SKU,
			expectedProduct.Name,
			expectedProduct.Price,
			expectedProduct.Supplier,
			expectedProduct.Stock,
			expectedProduct.CreatedAt,
			expectedProduct.UpdatedAt,
		)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, sku, name, price, supplier, stock, created_at, updated_at FROM products WHERE sku = ?")).
			WithArgs("SKU-002").
			WillReturnRows(rows)

		result, err := repo.GetProductBySKU(ctx, "SKU-002")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedProduct.ID, result.ID)
		assert.Equal(t, expectedProduct.SKU, result.SKU)
		assert.Equal(t, expectedProduct.Stock, result.Stock)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("handles empty SKU string", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, sku, name, price, supplier, stock, created_at, updated_at FROM products WHERE sku = ?")).
			WithArgs("").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetProductBySKU(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestUpdateProductStock tests the UpdateProductStock method
func TestUpdateProductStock(t *testing.T) {
	ctx := context.Background()
	staffUUID := "staff-uuid-123"

	t.Run("successfully updates product stock", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT with FOR UPDATE
		rows := sqlmock.NewRows([]string{"stock"}).AddRow(50)
		mock.ExpectQuery("SELECT stock.*FOR UPDATE").
			WithArgs("SKU-001").
			WillReturnRows(rows)

		// Expect UPDATE with any timestamp and staffUUID
		mock.ExpectExec("UPDATE products SET.*stock.*updated_at.*updated_by").
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect COMMIT
		mock.ExpectCommit()

		err = repo.UpdateProductStock(ctx, "SKU-001", 10, staffUUID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when product SKU does not exist", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT to return no rows
		mock.ExpectQuery("SELECT stock.*FOR UPDATE").
			WithArgs("NONEXISTENT").
			WillReturnError(sql.ErrNoRows)

		// Expect ROLLBACK
		mock.ExpectRollback()

		err = repo.UpdateProductStock(ctx, "NONEXISTENT", 10, staffUUID)

		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when transaction begin fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		beginErr := errors.New("transaction begin failed")
		mock.ExpectBegin().WillReturnError(beginErr)

		err = repo.UpdateProductStock(ctx, "SKU-001", 10, staffUUID)

		assert.Error(t, err)
		assert.Equal(t, beginErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when SELECT FOR UPDATE fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT to fail with different error
		selectErr := errors.New("select lock failed")
		mock.ExpectQuery("SELECT stock.*FOR UPDATE").
			WithArgs("SKU-001").
			WillReturnError(selectErr)

		// Expect ROLLBACK
		mock.ExpectRollback()

		err = repo.UpdateProductStock(ctx, "SKU-001", 10, staffUUID)

		assert.Error(t, err)
		assert.Equal(t, selectErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when UPDATE fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT to succeed
		rows := sqlmock.NewRows([]string{"stock"}).AddRow(50)
		mock.ExpectQuery("SELECT stock.*FOR UPDATE").
			WithArgs("SKU-001").
			WillReturnRows(rows)

		// Expect UPDATE to fail
		updateErr := errors.New("update failed")
		mock.ExpectExec("UPDATE products SET.*").
			WillReturnError(updateErr)

		// Expect ROLLBACK
		mock.ExpectRollback()

		err = repo.UpdateProductStock(ctx, "SKU-001", 10, staffUUID)

		assert.Error(t, err)
		assert.Equal(t, updateErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when commit fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT to succeed
		rows := sqlmock.NewRows([]string{"stock"}).AddRow(50)
		mock.ExpectQuery("SELECT stock.*FOR UPDATE").
			WithArgs("SKU-001").
			WillReturnRows(rows)

		// Expect UPDATE to succeed
		mock.ExpectExec("UPDATE products SET.*").
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect COMMIT to fail
		commitErr := errors.New("commit failed")
		mock.ExpectCommit().WillReturnError(commitErr)

		err = repo.UpdateProductStock(ctx, "SKU-001", 10, staffUUID)

		assert.Error(t, err)
		assert.Equal(t, commitErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("handles negative stock increase (stock out)", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT with FOR UPDATE
		rows := sqlmock.NewRows([]string{"stock"}).AddRow(100)
		mock.ExpectQuery("SELECT stock.*FOR UPDATE").
			WithArgs("SKU-001").
			WillReturnRows(rows)

		// Expect UPDATE with negative value (decrement)
		mock.ExpectExec("UPDATE products SET.*").
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect COMMIT
		mock.ExpectCommit()

		err = repo.UpdateProductStock(ctx, "SKU-001", -20, staffUUID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("handles large stock increase", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT with FOR UPDATE
		rows := sqlmock.NewRows([]string{"stock"}).AddRow(1000)
		mock.ExpectQuery("SELECT stock.*FOR UPDATE").
			WithArgs("SKU-001").
			WillReturnRows(rows)

		// Expect UPDATE with large value
		mock.ExpectExec("UPDATE products SET.*").
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect COMMIT
		mock.ExpectCommit()

		err = repo.UpdateProductStock(ctx, "SKU-001", 50000, staffUUID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("updates both stock and updated_at timestamp", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT with FOR UPDATE
		rows := sqlmock.NewRows([]string{"stock"}).AddRow(50)
		mock.ExpectQuery("SELECT stock.*FOR UPDATE").
			WithArgs("SKU-001").
			WillReturnRows(rows)

		// Expect UPDATE to include updated_at column
		mock.ExpectExec("UPDATE products SET.*updated_at").
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect COMMIT
		mock.ExpectCommit()

		err = repo.UpdateProductStock(ctx, "SKU-001", 10, staffUUID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("tracks staff UUID who performed the update", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)
		differentStaffUUID := "different-staff-uuid"

		// Expect transaction begin
		mock.ExpectBegin()

		// Expect SELECT with FOR UPDATE
		rows := sqlmock.NewRows([]string{"stock"}).AddRow(50)
		mock.ExpectQuery("SELECT stock.*FOR UPDATE").
			WithArgs("SKU-001").
			WillReturnRows(rows)

		// Expect UPDATE with updated_by column
		mock.ExpectExec("UPDATE products SET.*updated_by").
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect COMMIT
		mock.ExpectCommit()

		err = repo.UpdateProductStock(ctx, "SKU-001", 10, differentStaffUUID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCreateProduct tests the CreateProduct method
func TestCreateProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully creates a product", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Expect PREPARE with regex to match any whitespace variations
		mock.ExpectPrepare("INSERT INTO products.*").
			ExpectExec().
			WithArgs(mockProduct.SKU, mockProduct.Name, mockProduct.Price, mockProduct.Supplier, mockProduct.Stock).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateProduct(ctx, *mockProduct)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when prepare fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		prepareErr := errors.New("prepare failed")
		mock.ExpectPrepare("INSERT INTO products.*").
			WillReturnError(prepareErr)

		err = repo.CreateProduct(ctx, *mockProduct)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when execution fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		execErr := errors.New("execution failed")
		mock.ExpectPrepare("INSERT INTO products.*").
			ExpectExec().
			WithArgs(mockProduct.SKU, mockProduct.Name, mockProduct.Price, mockProduct.Supplier, mockProduct.Stock).
			WillReturnError(execErr)

		err = repo.CreateProduct(ctx, *mockProduct)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when no rows affected", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Return result with 0 rows affected
		mock.ExpectPrepare("INSERT INTO products.*").
			ExpectExec().
			WithArgs(mockProduct.SKU, mockProduct.Name, mockProduct.Price, mockProduct.Supplier, mockProduct.Stock).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err = repo.CreateProduct(ctx, *mockProduct)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("creates product with different data", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		newProduct := dto.Product{
			SKU:      "SKU-NEW",
			Name:     "New Product",
			Price:    500000,
			Supplier: "New Supplier",
			Stock:    200,
		}

		// Expect PREPARE and EXEC
		mock.ExpectPrepare("INSERT INTO products.*").
			ExpectExec().
			WithArgs(newProduct.SKU, newProduct.Name, newProduct.Price, newProduct.Supplier, newProduct.Stock).
			WillReturnResult(sqlmock.NewResult(2, 1))

		err = repo.CreateProduct(ctx, newProduct)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when RowsAffected fails", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := product.NewSQLRepository(db)

		// Return a result that fails on RowsAffected
		mock.ExpectPrepare("INSERT INTO products.*").
			ExpectExec().
			WithArgs(mockProduct.SKU, mockProduct.Name, mockProduct.Price, mockProduct.Supplier, mockProduct.Stock).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected failed")))

		err = repo.CreateProduct(ctx, *mockProduct)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Benchmarks
func BenchmarkGetProductBySKU(b *testing.B) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "sku", "name", "price", "supplier", "stock", "created_at", "updated_at",
	}).AddRow(
		mockProduct.ID,
		mockProduct.SKU,
		mockProduct.Name,
		mockProduct.Price,
		mockProduct.Supplier,
		mockProduct.Stock,
		mockProduct.CreatedAt,
		mockProduct.UpdatedAt,
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, sku, name, price, supplier, stock, created_at, updated_at FROM products WHERE sku = ?")).
		WithArgs("SKU-001").
		WillReturnRows(rows)

	repo := product.NewSQLRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		repo.GetProductBySKU(ctx, "SKU-001")
	}
}

func BenchmarkUpdateProductStock(b *testing.B) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	staffUUID := "staff-uuid-123"

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"stock"}).AddRow(50)
	mock.ExpectQuery("SELECT stock.*FOR UPDATE").
		WithArgs("SKU-001").
		WillReturnRows(rows)
	mock.ExpectExec("UPDATE products SET.*").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := product.NewSQLRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		repo.UpdateProductStock(ctx, "SKU-001", 10, staffUUID)
	}
}

func BenchmarkCreateProduct(b *testing.B) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO products.*").
		ExpectExec().
		WithArgs(mockProduct.SKU, mockProduct.Name, mockProduct.Price, mockProduct.Supplier, mockProduct.Stock).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := product.NewSQLRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		repo.CreateProduct(ctx, *mockProduct)
	}
}
