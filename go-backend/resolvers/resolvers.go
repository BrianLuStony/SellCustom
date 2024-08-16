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

func (r *ProductResolver) Images(ctx context.Context) ([]*ProductImageResolver, error) {
	rows, err := db.DB.Query("SELECT id, url, is_primary FROM images WHERE product_id = $1", r.p.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*ProductImageResolver
	for rows.Next() {
		var img models.ProductImage
		if err := rows.Scan(&img.ID, &img.ImageUrl, &img.IsPrimary); err != nil {
			return nil, err
		}
		images = append(images, &ProductImageResolver{img})
	}

	return images, nil
}

func (r *ProductResolver) Attributes(ctx context.Context) ([]*ProductAttributeResolver, error) {
	rows, err := db.DB.Query("SELECT id, name, value FROM attributes WHERE product_id = $1", r.p.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attributes []*ProductAttributeResolver
	for rows.Next() {
		var attr models.ProductAttribute
		if err := rows.Scan(&attr.ID, &attr.Name, &attr.Value); err != nil {
			return nil, err
		}
		attributes = append(attributes, &ProductAttributeResolver{attr})
	}

	return attributes, nil
}

func (r *ProductResolver) Reviews(ctx context.Context) ([]*ReviewResolver, error) {
	rows, err := db.DB.Query("SELECT id, rating, comment FROM reviews WHERE product_id = $1", r.p.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*ReviewResolver
	for rows.Next() {
		var rev models.Review
		if err := rows.Scan(&rev.ID, &rev.Rating, &rev.Comment); err != nil {
			return nil, err
		}
		reviews = append(reviews, &ReviewResolver{rev})
	}

	return reviews, nil
}

// ProductImageResolver resolves the ProductImage type
type ProductImageResolver struct {
	i models.ProductImage
}

func (r *ProductImageResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprint(r.i.ID))
}

func (r *ProductImageResolver) ImageUrl() string {
	return r.i.ImageUrl
}

func (r *ProductImageResolver) IsPrimary() bool {
	return r.i.IsPrimary
}

// ProductAttributeResolver resolves the ProductAttribute type
type ProductAttributeResolver struct {
	a models.ProductAttribute
}

func (r *ProductAttributeResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprint(r.a.ID))
}

func (r *ProductAttributeResolver) Name() string {
	return r.a.Name
}

func (r *ProductAttributeResolver) Value() string {
	return r.a.Value
}

// ReviewResolver resolves the Review type
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

func (r *ReviewResolver) Product(ctx context.Context) (*ProductResolver, error) {
	var p models.Product
	err := db.DB.QueryRow("SELECT id, name, description, price, stock_quantity, category_id FROM products WHERE id = $1", r.r.ProductID).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity, &p.CategoryID)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{p}, nil
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
