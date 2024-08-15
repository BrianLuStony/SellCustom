package models

type Category struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	ParentCategory *Category  `json:"parentCategory,omitempty"`
	Products       []*Product `json:"products,omitempty"`
}

type Product struct {
	ID            int                 `json:"id"`
	Name          string              `json:"name"`
	Description   *string             `json:"description,omitempty"`
	Price         float64             `json:"price"`
	StockQuantity int                 `json:"stockQuantity"`
	CategoryID    int                 `json:"-"`
	Category      *Category           `json:"category,omitempty"`
	Images        []*ProductImage     `json:"images,omitempty"`
	Attributes    []*ProductAttribute `json:"attributes,omitempty"`
	Reviews       []*Review           `json:"reviews,omitempty"`
}

type ProductImage struct {
	ID        int    `json:"id"`
	ImageUrl  string `json:"imageUrl"`
	IsPrimary bool   `json:"isPrimary"`
}

type ProductAttribute struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Order struct {
	ID          int          `json:"id"`
	UserID      int          `json:"-"`
	User        *User        `json:"user"`
	TotalAmount float64      `json:"totalAmount"`
	Status      string       `json:"status"`
	Items       []*OrderItem `json:"items"`
	CreatedAt   string       `json:"createdAt"`
}

type OrderItem struct {
	ID          int      `json:"id"`
	ProductID   int      `json:"-"`
	Product     *Product `json:"product"`
	Quantity    int      `json:"quantity"`
	PriceAtTime float64  `json:"priceAtTime"`
}

type Review struct {
	ID        int      `json:"id"`
	ProductID int      `json:"-"`
	Product   *Product `json:"product"`
	UserID    int      `json:"-"`
	User      *User    `json:"user"`
	Rating    int      `json:"rating"`
	Comment   string   `json:"comment"`
	CreatedAt string   `json:"createdAt"`
}
