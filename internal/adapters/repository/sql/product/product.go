package product

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/bagus-aulia/inventhier/internal/core/constants"
	dto "github.com/bagus-aulia/inventhier/internal/core/dto/product"
	"github.com/bagus-aulia/inventhier/internal/core/helpers"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
)

type sqlRepository struct {
	db *sql.DB
}

// NewSQLRepository creates a new relational SQL-based repository using the provided sql.DB connection.
func NewSQLRepository(db *sql.DB) ports.Product {
	return &sqlRepository{
		db: db,
	}
}

func (r *sqlRepository) GetProductBySKU(ctx context.Context, sku string) (*dto.Product, error) {
	logger := helpers.GetZerologWithContext(ctx).
		With().
		Str("repository", "sql.product").
		Str("function", "GetProductBySKU").
		Str("sku", sku).
		Logger()

	query := `SELECT 
			id, 
			sku, 
			name, 
			price, 
			supplier,
			stock, 
			created_at, 
			updated_at
		FROM products 
		WHERE sku = ?`
	row := r.db.QueryRowContext(ctx, query, sku)

	data := &dto.Product{}
	err := row.Scan(
		&data.ID,
		&data.SKU,
		&data.Name,
		&data.Price,
		&data.Supplier,
		&data.Stock,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error().Err(err).Msg("Failed executing operation")
		}

		return nil, err
	}

	return data, nil
}

func (r *sqlRepository) UpdateProductStock(ctx context.Context, sku string, stockIn int, staffUUID string) error {
	logger := helpers.GetZerologWithContext(ctx).
		With().
		Str("repository", "sql.product").
		Str("function", "UpdateProductStock").
		Str("sku", sku).
		Int("stock_in", stockIn).
		Logger()

	// 1. Begin Transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to begin sql transaction")
		return err
	}

	// rollback the transaction if there is an error during the process
	defer tx.Rollback()

	var currentStock int
	// 2.Pessimistic Lock
	queryLock := `SELECT stock 
		FROM products 
		WHERE sku = ? 
		FOR UPDATE`
	err = tx.QueryRowContext(ctx, queryLock, sku).Scan(&currentStock)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Error().Err(err).Msg("Failed executing select operation")
		}

		return err
	}

	// 3. update data
	queryUpdate := `UPDATE products SET 
			stock = stock + ?,
			updated_at = ?,
			updated_by = ?
		WHERE sku = ?`
	_, err = tx.ExecContext(ctx, queryUpdate, stockIn, time.Now(), staffUUID, sku)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update operation")
		return err
	}

	// 4. save the change and lock released
	if err := tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("Failed to commit transaction")
		return err
	}

	return nil
}

func (r *sqlRepository) CreateProduct(ctx context.Context, data dto.Product) error {
	logger := helpers.GetZerologWithContext(ctx).
		With().
		Str("repository", "sql.product").
		Str("function", "CreateProduct").
		Interface("data", data).
		Logger()

	query := `INSERT INTO products 
		(
			sku, 
			name, 
			price, 
			supplier,
			stock, 
			created_at
		) 
		VALUES (?, ?, ?, ?, ?, NOW())`

	var vals []interface{}
	vals = append(
		vals,
		data.SKU,
		data.Name,
		data.Price,
		data.Supplier,
		data.Stock,
	)

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to prepare context")
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, vals...)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to execute context")
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to get rows affected")
		return err
	}
	if affect < 1 {
		logger.Error().
			Msg("No rows affected")
		return constants.ErrInternalServer
	}

	return nil
}
