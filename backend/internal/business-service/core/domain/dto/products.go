package dto

// Product represents a row in the products table.
// Uses db tags for sqlx / database/sql, json tags for API responses.
type Product struct {
	ProductID   string  `db:"product_id" json:"product_id"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description,omitempty"`
	Photo       *string `db:"photo" json:"photo,omitempty"`
	Price       float64 `db:"price" json:"price"` // exact decimal for money
}
