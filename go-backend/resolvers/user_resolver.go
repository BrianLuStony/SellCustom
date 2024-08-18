package resolvers

import (
	"context"
	"fmt"
	"go-backend/db"
	"go-backend/models"

	"github.com/graph-gophers/graphql-go"
)

type UserResolver struct {
	u models.User
}

// Resolve ID field
func (r *UserResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprint(r.u.ID))
}

// Resolve Email field
func (r *UserResolver) Email() string {
	return r.u.Email
}

// Resolve FirstName field
func (r *UserResolver) FirstName() *string {
	return &r.u.FirstName
}

// Resolve LastName field
func (r *UserResolver) LastName() *string {
	return &r.u.LastName
}

// Resolve Orders field (optional, if you want to resolve user's orders)
func (r *UserResolver) Orders(ctx context.Context) ([]*OrderResolver, error) {
	rows, err := db.DB.QueryContext(ctx, "SELECT id, user_id, total_amount, status, created_at FROM orders WHERE user_id = $1", r.u.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*OrderResolver
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &OrderResolver{o})
	}

	return orders, nil
}
