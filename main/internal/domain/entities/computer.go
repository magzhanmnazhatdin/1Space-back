package entities

// Computer is the domain entity representing a computer in a club.
type Computer struct {
	ID          string
	ClubID      string
	PCNumber    int
	Description string
	IsAvailable bool
}
