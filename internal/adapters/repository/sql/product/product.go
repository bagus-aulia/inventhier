package product

import (
	"context"
	"database/sql"

	"github.com/bagus-aulia/inventhier/internal/core/domain"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
)

type sqlRepository struct {
	db *sql.DB
}

// NewSQLRepository creates a new relational SQL-based repository using the provided sql.DB connection.
func NewSQLRepository(db *sql.DB) ports.ProductRepository {
	return &sqlRepository{
		db: db,
	}
}

func (r *sqlRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	query := "SELECT id, name, sku, price, stock, created_at, updated_at FROM products WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	var p domain.Product
	err := row.Scan(&p.ID, &p.Name, &p.SKU, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *sqlRepository) Create(ctx context.Context, p *domain.Product) error {
	query := "INSERT INTO products (id, name, sku, price, stock, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.SKU, p.Price, p.Stock, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *sqlRepository) List(ctx context.Context) ([]domain.Product, error) {
	query := "SELECT id, name, sku, price, stock, created_at, updated_at FROM products"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}
