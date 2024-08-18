package resolvers

import (
	"context"
	"fmt"
	"go-backend/db"
	"go-backend/models"

	"github.com/graph-gophers/graphql-go"
)

type CategoryResolver struct {
	c models.Category
}

func (r *CategoryResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprint(r.c.ID))
}

func (r *CategoryResolver) Name() string {
	return r.c.Name
}

func (r *CategoryResolver) ParentCategory(ctx context.Context) (*CategoryResolver, error) {
	if r.c.ParentCategory == nil {
		return nil, nil
	}
	return &CategoryResolver{*r.c.ParentCategory}, nil
}

func (r *CategoryResolver) Products(ctx context.Context) ([]*ProductResolver, error) {
	rows, err := db.DB.Query("SELECT id, name, description, price, stock_quantity FROM products WHERE category_id = $1", r.c.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*ProductResolver
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity); err != nil {
			return nil, err
		}
		products = append(products, &ProductResolver{p})
	}

	return products, nil
}
