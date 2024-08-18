package resolvers

import (
	"context"
	"fmt"
	"go-backend/db"
	"go-backend/models"

	"github.com/graph-gophers/graphql-go"
)

type OrderResolver struct {
	o models.Order
}

// Resolve ID field
func (r *OrderResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprint(r.o.ID))
}

// Resolve User field
func (r *OrderResolver) User(ctx context.Context) (*UserResolver, error) {
	var u models.User
	err := db.DB.QueryRowContext(ctx, "SELECT id, email, first_name, last_name FROM users WHERE id = $1", r.o.UserID).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName)
	if err != nil {
		return nil, err
	}
	return &UserResolver{u}, nil
}

// Resolve TotalAmount field
func (r *OrderResolver) TotalAmount() float64 {
	return r.o.TotalAmount
}

// Resolve Status field
func (r *OrderResolver) Status() string {
	return r.o.Status
}

// Resolve Items field
func (r *OrderResolver) Items(ctx context.Context) ([]*OrderItemResolver, error) {
	rows, err := db.DB.QueryContext(ctx, "SELECT id, product_id, quantity, price_at_time FROM order_items WHERE order_id = $1", r.o.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*OrderItemResolver
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.PriceAtTime); err != nil {
			return nil, err
		}
		items = append(items, &OrderItemResolver{item})
	}

	return items, nil
}

// Resolve CreatedAt field
func (r *OrderResolver) CreatedAt() string {
	return r.o.CreatedAt
}

type OrderItemResolver struct {
	oi models.OrderItem
}

// Resolve ID field
func (r *OrderItemResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprint(r.oi.ID))
}

// Resolve Product field
func (r *OrderItemResolver) Product(ctx context.Context) (*ProductResolver, error) {
	var p models.Product
	err := db.DB.QueryRowContext(ctx, "SELECT id, name, description, price, stock_quantity, category_id FROM products WHERE id = $1", r.oi.ProductID).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity, &p.CategoryID)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{p}, nil
}

// Resolve Quantity field
func (r *OrderItemResolver) Quantity() int32 {
	return int32(r.oi.Quantity)
}

// Resolve PriceAtTime field
func (r *OrderItemResolver) PriceAtTime() float64 {
	return r.oi.PriceAtTime
}
