package models

import "time"

type Product struct {
	ProductID   string    `db:"product_id" json:"product_id"`
	StoreID     string    `db:"store_id" json:"store_id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description,omitempty"`
	Photo       *string   `db:"photo" json:"photo,omitempty"`
	Price       float64   `db:"price" json:"price"` // exact decimal for money
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
