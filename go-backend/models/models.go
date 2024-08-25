package models

type Category struct {
	ID             int32      `json:"id"`
	Name           string     `json:"name"`
	ParentCategory *Category  `json:"parentCategory,omitempty"`
	Products       []*Product `json:"products,omitempty"`
}

type Product struct {
	ID            int32               `json:"id"`
	Name          string              `json:"name"`
	Description   *string             `json:"description,omitempty"`
	Price         float64             `json:"price"`
	StockQuantity int32               `json:"stockQuantity"`
	CategoryID    int32               `json:"-"`
	Category      *Category           `json:"category,omitempty"`
	Images        []*ProductImage     `json:"images,omitempty"`
	Attributes    []*ProductAttribute `json:"attributes,omitempty"`
	Reviews       []*Review           `json:"reviews,omitempty"`
}

type ProductImage struct {
	ID        int32  `json:"id"`
	ImageUrl  string `json:"imageUrl"`
	IsPrimary bool   `json:"isPrimary"`
}

type ProductAttribute struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type User struct {
	ID        int32  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Order struct {
	ID          int32        `json:"id"`
	UserID      int32        `json:"-"`
	User        *User        `json:"user"`
	TotalAmount float64      `json:"totalAmount"`
	Status      string       `json:"status"`
	Items       []*OrderItem `json:"items"`
	CreatedAt   string       `json:"createdAt"`
}

type OrderItem struct {
	ID          int32    `json:"id"`
	ProductID   int32    `json:"-"`
	Product     *Product `json:"product"`
	Quantity    int32    `json:"quantity"`
	PriceAtTime float64  `json:"priceAtTime"`
}

type Review struct {
	ID        int32    `json:"id"`
	ProductID int32    `json:"-"`
	Product   *Product `json:"product"`
	UserID    int32    `json:"-"`
	User      *User    `json:"user"`
	Rating    int32    `json:"rating"`
	Comment   string   `json:"comment"`
	CreatedAt string   `json:"createdAt"`
}

type Query struct {
	Product    func(id int32) (*Product, error)
	Products   func(category *int32, search *string) ([]*Product, error)
	Categories func() ([]*Category, error)
	Order      func(id int32) (*Order, error)
	UserOrders func(userId int32) ([]*Order, error)
}

type Mutation struct {
	CreateProduct  func(input ProductInput) (*Product, error)
	UpdateProduct  func(id int32, input ProductInput) (*Product, error)
	DeleteProduct  func(id int32) (bool, error)
	CreateOrder    func(input OrderInput) (*Order, error)
	CreateReview   func(input ReviewInput) (*Review, error)
	CreateCategory func(name string, parentID *int32) (*Category, error)
}

type ProductInput struct {
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	Price         float64 `json:"price"`
	StockQuantity int32   `json:"stockQuantity"`
	CategoryID    int32   `json:"categoryId"`
}

type OrderInput struct {
	UserID      int32             `json:"userId"`
	Items       []*OrderItemInput `json:"items"`
	TotalAmount float64           `json:"totalAmount"`
}

type CategoryInput struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parentId"`
}

type OrderItemInput struct {
	ProductID int32 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type ReviewInput struct {
	ProductID int32   `json:"productId"`
	UserID    int32   `json:"userId"`
	Rating    int32   `json:"rating"`
	Comment   *string `json:"comment"`
}
