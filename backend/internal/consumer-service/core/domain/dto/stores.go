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
