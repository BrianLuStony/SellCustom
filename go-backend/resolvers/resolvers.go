package resolvers

import (
	"context"

	"github.com/graphql-go/graphql"
)

type Resolver struct {
}

// Product represents a product in the database
type Product struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"` // Allow null descriptions
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stockQuantity"`
	CategoryID    int       `json:"categoryId"`
	Category      *Category `json:"category,omitempty"` // Optional field to hold category data
	Reviews       []*Review `json:"reviews,omitempty"`  // Optional field to hold reviews
}

// ProductResolver resolves fields for the Product type
type ProductResolver struct {
	p Product
}

// Category represents a category in the database
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CategoryResolver resolves fields for the Category type
type CategoryResolver struct {
	c Category
}

// Review represents a review in the database
type Review struct {
	ID        int    `json:"id"`
	ProductID int    `json:"productId"`
	UserID    int    `json:"userId"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"createdAt"`
}

// ReviewResolver resolves fields for the Review type
type ReviewResolver struct {
	r Review
}

// Query represents the root query type
type Query struct {
	Resolver *Resolver
}

// Mutation represents the root mutation type
type Mutation struct {
	Resolver *Resolver
}

// Product returns a product by ID
func (r *Resolver) Product(ctx context.Context, args struct{ ID graphql.ID }) (*ProductResolver, error) {
	var p Product
	err := db.QueryRow("SELECT id, name, description, price, stock_quantity, category_id FROM products WHERE id = $1", args.ID).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity, &p.CategoryID)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{p}, nil
}

// (Implement other resolvers for Product fields like Category and Reviews)

func (r *ProductResolver) Category(ctx context.Context) (*CategoryResolver, error) {
	// Implement logic to fetch category based on product.CategoryID
	// and return a CategoryResolver instance
	return nil, nil // Replace with actual implementation
}

func (r *ProductResolver) Reviews(ctx context.Context) ([]*ReviewResolver, error) {
	// Implement logic to fetch reviews for the product
	// and return an array of ReviewResolver instances
	return nil, nil // Replace with actual implementation
}

// (Implement resolvers for other types like Category, User, Order, etc.)

// Define other resolver functions and types following similar patterns ...

// var RootQuery = Query{
//   Resolver: &Resolver{db: db}, // Replace with actual database connection
// }

// var RootMutation = Mutation{
//   Resolver: &Resolver{db: db}, // Replace with actual database connection
// }
