package dto

type Store struct {
	StoreID     string  `db:"store_id" json:"store_id"`
	Name        *string `db:"name" json:"name"`
	Address     *string `db:"address" json:"address"`
	Photo       *string `db:"photo" json:"photo,omitempty"`
	Description *string `db:"description" json:"description,omitempty"`
	Latitude    float64 `db:"latitude" json:"latitude"`
	Longitude   float64 `db:"longitude" json:"longitude"`
}

type Coord struct {
	CoordID    string  `json:"coord_id"`
	EntityID   string  `json:"entity_id"`
	EntityType string  `json:"entity_type"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Address    string  `json:"address,omitempty"`
}
