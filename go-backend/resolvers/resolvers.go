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
func (r *OrderItemResolver) Quantity() int {
	return r.oi.Quantity
}

// Resolve PriceAtTime field
func (r *OrderItemResolver) PriceAtTime() float64 {
	return r.oi.PriceAtTime
}

type QueryResolver struct {
	q models.Query
}

func (r *QueryResolver) Product(ctx context.Context, args struct{ ID graphql.ID }) (*ProductResolver, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}
	product, err := r.q.Product(id)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{*product}, nil
}

func (r *QueryResolver) Products(ctx context.Context, args struct {
	Category *graphql.ID
	Search   *string
}) ([]*ProductResolver, error) {
	var categoryID *int
	if args.Category != nil {
		id, err := strconv.Atoi(string(*args.Category))
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %v", err)
		}
		categoryID = &id
	}
	products, err := r.q.Products(categoryID, args.Search)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*ProductResolver, len(products))
	for i, p := range products {
		resolvers[i] = &ProductResolver{*p}
	}
	return resolvers, nil
}

func (r *QueryResolver) Categories(ctx context.Context) ([]*CategoryResolver, error) {
	categories, err := r.q.Categories()
	if err != nil {
		return nil, err
	}
	resolvers := make([]*CategoryResolver, len(categories))
	for i, c := range categories {
		resolvers[i] = &CategoryResolver{*c}
	}
	return resolvers, nil
}

func (r *QueryResolver) Order(ctx context.Context, args struct{ ID graphql.ID }) (*OrderResolver, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}
	order, err := r.q.Order(id)
	if err != nil {
		return nil, err
	}
	return &OrderResolver{*order}, nil
}

func (r *QueryResolver) UserOrders(ctx context.Context, args struct{ UserID graphql.ID }) ([]*OrderResolver, error) {
	userID, err := strconv.Atoi(string(args.UserID))
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}
	orders, err := r.q.UserOrders(userID)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*OrderResolver, len(orders))
	for i, o := range orders {
		resolvers[i] = &OrderResolver{*o}
	}
	return resolvers, nil
}

type MutationResolver struct {
	m models.Mutation
}

func (r *MutationResolver) CreateProduct(ctx context.Context, args struct{ Input models.ProductInput }) (*ProductResolver, error) {
	product, err := r.m.CreateProduct(args.Input)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{*product}, nil
}

func (r *MutationResolver) UpdateProduct(ctx context.Context, args struct {
	ID    graphql.ID
	Input models.ProductInput
}) (*ProductResolver, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}
	product, err := r.m.UpdateProduct(id, args.Input)
	if err != nil {
		return nil, err
	}
	return &ProductResolver{*product}, nil
}

func (r *MutationResolver) DeleteProduct(ctx context.Context, args struct{ ID graphql.ID }) (bool, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return false, fmt.Errorf("invalid ID: %v", err)
	}
	return r.m.DeleteProduct(id)
}

func (r *MutationResolver) CreateOrder(ctx context.Context, args struct{ Input models.OrderInput }) (*OrderResolver, error) {
	order, err := r.m.CreateOrder(args.Input)
	if err != nil {
		return nil, err
	}
	return &OrderResolver{*order}, nil
}

func (r *MutationResolver) CreateReview(ctx context.Context, args struct{ Input models.ReviewInput }) (*ReviewResolver, error) {
	review, err := r.m.CreateReview(args.Input)
	if err != nil {
		return nil, err
	}
	return &ReviewResolver{*review}, nil
}
