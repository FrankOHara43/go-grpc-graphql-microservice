package order

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) PutOrder(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES($1, $2, $3, $4)",
		o.ID,
		o.CreatedAt,
		o.AccountID,
		o.TotalPrice,
	); err != nil {
		tx.Rollback()
		return err
	}

	for _, p := range o.Products {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO order_products(order_id, product_id, quantity) VALUES($1, $2, $3)",
			o.ID,
			p.ID,
			p.Quantity,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
		o.id, 
		o.created_at, 
		o.account_id, 
		o.total_price::money::numeric::float8,
		op.product_id,
		op.quantity
		FROM orders o
		JOIN order_products op ON o.id = op.order_id
		WHERE o.account_id = $1
		ORDER BY o.id`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []Order{}
	lastOrder := &Order{}

	for rows.Next() {
		currentOrder := &Order{}
		orderedProduct := &OrderedProduct{}
		if err = rows.Scan(
			&currentOrder.ID,
			&currentOrder.CreatedAt,
			&currentOrder.AccountID,
			&currentOrder.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}

		if lastOrder.ID != "" && lastOrder.ID != currentOrder.ID {
			newOrder := Order{
				ID:         lastOrder.ID,
				AccountID:  lastOrder.AccountID,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products:   lastOrder.Products,
			}
			orders = append(orders, newOrder)
			lastOrder = &Order{}
		}

		lastOrder.ID = currentOrder.ID
		lastOrder.AccountID = currentOrder.AccountID
		lastOrder.CreatedAt = currentOrder.CreatedAt
		lastOrder.TotalPrice = currentOrder.TotalPrice
		lastOrder.Products = append(lastOrder.Products, OrderedProduct{
			ID:       orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})
	}

	if lastOrder.ID != "" {
		newOrder := Order{
			ID:         lastOrder.ID,
			AccountID:  lastOrder.AccountID,
			CreatedAt:  lastOrder.CreatedAt,
			TotalPrice: lastOrder.TotalPrice,
			Products:   lastOrder.Products,
		}
		orders = append(orders, newOrder)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
