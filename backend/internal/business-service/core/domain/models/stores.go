package models

import "time"

type Stores struct {
	StoreID     string    `db:"store_id" json:"store_id"`
	OwnerID     string    `db:"user_id" json:"user_id"`
	CoordID     string    `db:"coord_id" json:"coord_id"`
	Name        *string   `db:"name" json:"name"`
	Address     *string   `db:"address" json:"address"`
	Description *string   `db:"description" json:"description,omitempty"`
	Photo       *string   `db:"photo" json:"photo,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
