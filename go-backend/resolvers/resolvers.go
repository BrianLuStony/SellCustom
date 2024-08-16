package resolvers

import (
	"context"
	"database/sql"
	"fmt"
	"go-backend/db"
	"go-backend/models"

	"github.com/graph-gophers/graphql-go"
)

type Resolver struct{}

func (r *Resolver) Product(ctx context.Context, args struct{ ID graphql.ID }) (*ProductResolver, error) {
	var p models.Product
	err := db.DB.QueryRow("SELECT id, name, description, price, stock_quantity, category_id FROM products WHERE id = $1", args.ID).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity, &p.CategoryID)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{p}, nil
}

func (r *Resolver) Products(ctx context.Context, args struct{ Category, Search *string }) ([]*ProductResolver, error) {
	// Implement product listing based on category and search
	// Example query below
	var rows *sql.Rows
	var err error
	if args.Category != nil {
		rows, err = db.DB.Query("SELECT id, name, description, price, stock_quantity, category_id FROM products WHERE category_id = $1", args.Category)
	} else if args.Search != nil {
		searchPattern := "%" + *args.Search + "%"
		rows, err = db.DB.Query("SELECT id, name, description, price, stock_quantity, category_id FROM products WHERE name ILIKE $1", searchPattern)
	} else {
		rows, err = db.DB.Query("SELECT id, name, description, price, stock_quantity, category_id FROM products")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*ProductResolver
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, &ProductResolver{p})
	}

	return products, nil
}

func (r *Resolver) Categories(ctx context.Context) ([]*CategoryResolver, error) {
	rows, err := db.DB.Query("SELECT id, name, parent_category_id FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*CategoryResolver
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentCategory.ID); err != nil {
			return nil, err
		}
		categories = append(categories, &CategoryResolver{c})
	}

	return categories, nil
}

// Define other resolvers for Order, Review, etc.

type ProductResolver struct {
	p models.Product
}

func (r *ProductResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprint(r.p.ID))
}

func (r *ProductResolver) Name() string {
	return r.p.Name
}

func (r *ProductResolver) Description() *string {
	return r.p.Description
}

func (r *ProductResolver) Price() float64 {
	return r.p.Price
}

func (r *ProductResolver) StockQuantity() int32 {
	return int32(r.p.StockQuantity)
}

func (r *ProductResolver) Category(ctx context.Context) (*CategoryResolver, error) {
	var c models.Category
	err := db.DB.QueryRow("SELECT id, name, parent_category_id FROM categories WHERE id = $1", r.p.CategoryID).Scan(&c.ID, &c.Name, &c.ParentCategory.ID)
	if err != nil {
		return nil, err
	}
	return &CategoryResolver{c}, nil
}

// Implement resolvers for Product.Images, Product.Attributes, Product.Reviews, etc.

type CategoryResolver struct {
	c models.Category
}
