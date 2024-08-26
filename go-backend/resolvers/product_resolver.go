package resolvers

import (
	"context"
	"database/sql"
	"fmt"
	"go-backend/db"
	"go-backend/models"

	"github.com/graph-gophers/graphql-go"
)

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

	if r.p.CategoryID == 0 {
		return nil, nil
	}

	var c models.Category
	err := db.DB.QueryRow("SELECT id, name FROM categories WHERE id = $1", r.p.CategoryID).Scan(&c.ID, &c.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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
