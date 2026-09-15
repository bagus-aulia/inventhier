package product

import (
	"context"
	"database/sql"

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
