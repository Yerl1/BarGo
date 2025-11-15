package dto

import "time"

// Product represents a row in the products table.
// Uses db tags for sqlx / database/sql, json tags for API responses.
type Product struct {
	ProductID   string  `db:"product_id" json:"product_id"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description,omitempty"`
	Photo       *string `db:"photo" json:"photo,omitempty"`
	Price       float64 `db:"price" json:"price"` // exact decimal for money'
	InStock     bool    `db:"in_stock" json:"in_stock"`
}

type ProductInfo struct {
	ProductID   string    `db:"product_id" json:"product_id"`
	StoreID     string    `db:"store_id" json:"store_id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description,omitempty"`
	Photo       *string   `db:"photo" json:"photo,omitempty"`
	Price       float64   `db:"price" json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	InStock     bool      `db:"in_stock" json:"in_stock"`
}

type ProductStoresInfo struct {
	Product Product `json:"product"`
	Stores  []Store `json:"stores"`
}
