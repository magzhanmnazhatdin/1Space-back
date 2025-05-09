package entities

// Club is the domain entity representing a computer club.
type Club struct {
	ID           string  `firestore:"id"             json:"id"`
	Name         string  `firestore:"name"           json:"name"`
	Address      string  `firestore:"address"        json:"address"`
	PricePerHour float64 `firestore:"price_per_hour" json:"price_per_hour"`
	AvailablePCs int     `firestore:"available_pcs"  json:"available_pcs"`
}
