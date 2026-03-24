package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Order represents the database entity
type Order struct {
	ID         string
	CustomerID string
	Amount     float64
	Status     string
}

// OrderRepository handles all database interactions for orders
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
}

type postgresOrderRepository struct {
	db *pgxpool.Pool
}

// NewPostgresOrderRepository creates a new Postgres implementation of OrderRepository
func NewPostgresOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &postgresOrderRepository{db: db}
}

func (r *postgresOrderRepository) CreateOrder(ctx context.Context, order *Order) error {
	query := `
		INSERT INTO orders (id, customer_id, amount, status)
		VALUES ($1, $2, $3, $4)
	`
	
	// Execute the query
	_, err := r.db.Exec(ctx, query, order.ID, order.CustomerID, order.Amount, order.Status)
	if err != nil {
		return err
	}

	return nil
}
