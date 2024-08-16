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

// Resolves a single product by ID
func (r *Resolver) Product(ctx context.Context, args struct{ ID graphql.ID }) (*ProductResolver, error) {
	var p models.Product
	err := db.DB.QueryRow("SELECT id, name, description, price, stock_quantity, category_id FROM products WHERE id = $1", args.ID).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity, &p.CategoryID)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{p}, nil
}

// Resolves a list of products optionally filtered by category or search string
func (r *Resolver) Products(ctx context.Context, args struct{ Category, Search *string }) ([]*ProductResolver, error) {
	var rows *sql.Rows
	var err error

	if args.Category != nil {
		rows, err = db.DB.Query("SELECT id, name, description, price, stock_quantity, category_id FROM products WHERE category_id = $1", *args.Category)
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

// Resolves a list of all categories
func (r *Resolver) Categories(ctx context.Context) ([]*CategoryResolver, error) {
	rows, err := db.DB.Query("SELECT id, name, parent_category_id FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*CategoryResolver
	for rows.Next() {
		var c models.Category
		var parentCategoryID sql.NullInt64
		if err := rows.Scan(&c.ID, &c.Name, &parentCategoryID); err != nil {
			return nil, err
		}

		if parentCategoryID.Valid {
			c.ParentCategory = &models.Category{ID: int(parentCategoryID.Int64)}
		}

		categories = append(categories, &CategoryResolver{c})
	}

	return categories, nil
}

// ProductResolver resolves the Product type
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

// CategoryResolver resolves the Category type
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
	var products []*models.Product
	rows, err := db.DB.Query("SELECT id, name, description, price, stock_quantity FROM products WHERE category_id = $1", r.c.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	var resolvers []*ProductResolver
	for _, p := range products {
		resolvers = append(resolvers, &ProductResolver{*p})
	}
	return resolvers, nil
}

// Implement resolvers for Product.Images, Product.Attributes, Product.Reviews, etc.
