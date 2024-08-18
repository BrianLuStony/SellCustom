package resolvers

import (
	"context"
	"fmt"
	"go-backend/db"
	"go-backend/models"

	"github.com/graph-gophers/graphql-go"
)

type ReviewResolver struct {
	r models.Review
}

func (r *ReviewResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprint(r.r.ID))
}

func (r *ReviewResolver) Rating() int32 {
	return int32(r.r.Rating)
}

func (r *ReviewResolver) Comment() *string {
	return &r.r.Comment
}
func (r *ReviewResolver) CreatedAt() string {
	return r.r.CreatedAt
}

func (r *ReviewResolver) Product(ctx context.Context) (*ProductResolver, error) {
	var p models.Product
	err := db.DB.QueryRow("SELECT id, name, description, price, stock_quantity, category_id FROM products WHERE id = $1", r.r.ProductID).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity, &p.CategoryID)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{p}, nil
}

func (r *ReviewResolver) User(ctx context.Context) (*UserResolver, error) {
	var u models.User
	err := db.DB.QueryRow("SELECT id, email, first_name, last_name FROM users WHERE id = $1", r.r.UserID).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName)
	if err != nil {
		return nil, err
	}
	return &UserResolver{u}, nil
}
