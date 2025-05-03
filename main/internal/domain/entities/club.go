package entities

// Club is the domain entity representing a computer club.
type Club struct {
	ID           string
	Name         string
	Address      string
	PricePerHour float64
	AvailablePCs int
}
