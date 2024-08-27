package resolvers

import (
	"context"
	"database/sql"
	"fmt"
	"go-backend/db"
	"go-backend/models"
	"strconv"

	"github.com/graph-gophers/graphql-go"
)

type Resolver struct {
	Query    *QueryResolver
	Mutation *MutationResolver
	ProductR *ProductResolver
	Category *CategoryResolver
	OrderR   *OrderResolver
	Review   *ReviewResolver
	User     *UserResolver
	// Add other resolvers as needed
}

func NewResolver() *Resolver {
	return &Resolver{
		Query:    &QueryResolver{},
		Mutation: &MutationResolver{},
		ProductR: &ProductResolver{},
		Category: &CategoryResolver{},
		OrderR:   &OrderResolver{},
		Review:   &ReviewResolver{},
		User:     &UserResolver{},
		// Initialize other resolvers
	}
}

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

// Resolves a list of orders
func (r *Resolver) Orders(ctx context.Context) ([]*OrderResolver, error) {
	rows, err := db.DB.Query("SELECT id, user_id, total_amount, status, created_at FROM orders")
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

func (r *Resolver) Order(ctx context.Context, args struct{ ID graphql.ID }) (*OrderResolver, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}
	var o models.Order
	err = db.DB.QueryRow("SELECT id, user_id, total_amount, status, created_at FROM orders WHERE id = $1", id).Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &OrderResolver{o}, nil
}

func (r *Resolver) CreateProduct(ctx context.Context, args struct{ Input models.ProductInput }) (*ProductResolver, error) {
	// Start a transaction

	stockQuantity := int(args.Input.StockQuantity)
	category_id := int(args.Input.CategoryID)

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check if the category exists
	var categoryExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", category_id).Scan(&categoryExists)
	if err != nil {
		return nil, err
	}
	if !categoryExists {
		return nil, fmt.Errorf("category with ID %d does not exist", category_id)
	}

	// Insert the new product
	var productID int32
	err = tx.QueryRow(`
        INSERT INTO products (name, description, price, stock_quantity, category_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, args.Input.Name, args.Input.Description, args.Input.Price, stockQuantity, category_id).Scan(&productID)
	if err != nil {
		return nil, err
	}

	for _, img := range args.Input.Images {
		_, err = tx.Exec(`
            INSERT INTO product_images (product_id, image_url, is_primary)
            VALUES ($1, $2, $3)
        `, productID, img.ImageUrl, img.IsPrimary)
		if err != nil {
			return nil, err
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Fetch the newly created product
	var p models.Product
	err = db.DB.QueryRow(`
        SELECT id, name, description, price, stock_quantity, category_id
        FROM products WHERE id = $1
    `, productID).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity, &p.CategoryID)
	if err != nil {
		return nil, err
	}

	// Initialize the Category field
	p.Category = &models.Category{ID: p.CategoryID}

	return &ProductResolver{p: p}, nil
}

func (r *Resolver) AddProductImage(ctx context.Context, args struct {
	ProductID graphql.ID
	Input     models.ProductImageInput
}) (*ProductImageResolver, error) {
	productID, err := strconv.Atoi(string(args.ProductID))
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %v", err)
	}

	var imageID int32
	err = db.DB.QueryRow(`
        INSERT INTO product_images (product_id, image_url, is_primary)
        VALUES ($1, $2, $3)
        RETURNING id
    `, productID, args.Input.ImageUrl, args.Input.IsPrimary).Scan(&imageID)
	if err != nil {
		return nil, err
	}

	img := models.ProductImage{
		ID:        imageID,
		ImageUrl:  args.Input.ImageUrl,
		IsPrimary: args.Input.IsPrimary,
	}

	return &ProductImageResolver{i: img}, nil
}

func (r *Resolver) UpdateProduct(ctx context.Context, args struct {
	ID    graphql.ID
	Input models.ProductInput
}) (*ProductResolver, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}
	// Start a transaction
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stockQuantity := int(args.Input.StockQuantity)

	// Update the product
	_, err = tx.Exec(`
        UPDATE products
        SET name = $1, description = $2, price = $3, stock_quantity = $4, category_id = $5
        WHERE id = $6
    `, args.Input.Name, args.Input.Description, args.Input.Price, stockQuantity, args.Input.CategoryID, id)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Fetch the updated product
	var p models.Product
	err = db.DB.QueryRow(`
        SELECT id, name, description, price, stock_quantity, category_id
        FROM products WHERE id = $1
    `, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity, &p.CategoryID)
	if err != nil {
		return nil, err
	}

	return &ProductResolver{p}, nil
}

func (r *Resolver) DeleteProduct(ctx context.Context, args struct{ ID graphql.ID }) (bool, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return false, fmt.Errorf("invalid ID: %v", err)
	}

	result, err := db.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *Resolver) CreateOrder(ctx context.Context, args struct{ Input models.OrderInput }) (*OrderResolver, error) {
	// Start a transaction
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert the new order
	var orderID int
	err = tx.QueryRow(`
        INSERT INTO orders (user_id, total_amount, status)
        VALUES ($1, $2, $3)
        RETURNING id
    `, args.Input.UserID, args.Input.TotalAmount, "PENDING").Scan(&orderID)
	if err != nil {
		return nil, err
	}

	// Insert order items
	for _, item := range args.Input.Items {
		_, err = tx.Exec(`
            INSERT INTO order_items (order_id, product_id, quantity, price_at_time)
            VALUES ($1, $2, $3, (SELECT price FROM products WHERE id = $2))
        `, orderID, item.ProductID, item.Quantity)
		if err != nil {
			return nil, err
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Fetch the newly created order
	var o models.Order
	err = db.DB.QueryRow(`
        SELECT id, user_id, total_amount, status, created_at
        FROM orders WHERE id = $1
    `, orderID).Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &OrderResolver{o}, nil
}

func (r *Resolver) CreateReview(ctx context.Context, args struct{ Input models.ReviewInput }) (*ReviewResolver, error) {
	// Insert the new review
	var reviewID int
	err := db.DB.QueryRow(`
        INSERT INTO reviews (product_id, user_id, rating, comment)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, args.Input.ProductID, args.Input.UserID, args.Input.Rating, args.Input.Comment).Scan(&reviewID)
	if err != nil {
		return nil, err
	}

	// Fetch the newly created review
	var rev models.Review
	err = db.DB.QueryRow(`
        SELECT id, product_id, user_id, rating, comment, created_at
        FROM reviews WHERE id = $1
    `, reviewID).Scan(&rev.ID, &rev.ProductID, &rev.UserID, &rev.Rating, &rev.Comment, &rev.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &ReviewResolver{rev}, nil
}

func (r *Resolver) UserOrders(ctx context.Context, args struct{ UserID graphql.ID }) ([]*OrderResolver, error) {
	userID, err := strconv.Atoi(string(args.UserID))
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	rows, err := db.DB.Query("SELECT id, user_id, total_amount, status, created_at FROM orders WHERE user_id = $1", userID)
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

// Resolves a list of all categories
func (r *Resolver) Categories(ctx context.Context) ([]*CategoryResolver, error) {
	rows, err := db.DB.Query("SELECT id, name, parent_id FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*CategoryResolver
	for rows.Next() {
		var c models.Category
		var parentID sql.NullInt32
		if err := rows.Scan(&c.ID, &c.Name, &parentID); err != nil {
			return nil, err
		}

		if parentID.Valid {
			c.ParentCategory = &models.Category{ID: parentID.Int32}
		}

		categories = append(categories, &CategoryResolver{c})
	}

	return categories, nil
}

func (r *Resolver) CreateCategory(ctx context.Context, args struct{ Input models.CategoryInput }) (*CategoryResolver, error) {
	var parentID *int32
	if args.Input.ParentID != nil {
		id, err := strconv.ParseInt(*args.Input.ParentID, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid parent ID: %v", err)
		}
		parentIDInt32 := int32(id)
		parentID = &parentIDInt32
	}

	// Start a transaction
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert the new category
	var categoryID int32
	err = tx.QueryRow(`
        INSERT INTO categories (name, parent_id)
        VALUES ($1, $2)
        RETURNING id
    `, args.Input.Name, parentID).Scan(&categoryID)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Fetch the newly created category
	var c models.Category
	err = db.DB.QueryRow(`
        SELECT id, name, parent_id
        FROM categories WHERE id = $1
    `, categoryID).Scan(&c.ID, &c.Name, &c.ParentCategory)
	if err != nil {
		return nil, err
	}

	return &CategoryResolver{c}, nil
}

// ProductResolver resolves the Product type

// ReviewResolver resolves the Review type

// CategoryResolver resolves the Category type
