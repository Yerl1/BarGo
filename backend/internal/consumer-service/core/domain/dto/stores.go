package dto

// Store represents a physical store in the system.
type Store struct {
	StoreID     string  `json:"store_id"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description *string `json:"description,omitempty"`
	Photo       *string `json:"photo,omitempty"`
}

type StoreInfo struct {
	StoreID       string  `db:"store_id" json:"store_id"`
	Name          string  `db:"name" json:"name"`
	Address       string  `db:"address" json:"address"`
	Latitude      float64 `db:"latitude" json:"latitude"`
	Longitude     float64 `db:"longitude" json:"longitude"`
	Description   *string `db:"description" json:"description,omitempty"`
	Photo         *string `db:"photo" json:"photo,omitempty"`
	OwnerID       *string `db:"owner_id" json:"owner_id,omitempty"`
	AverageRating float64 `db:"average_rating" json:"average_rating"`
	ReviewCount   int     `db:"review_count" json:"review_count"`
}
