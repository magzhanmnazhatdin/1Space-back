package entities

import "time"

// Booking is the domain entity representing a reservation.
type Booking struct {
	ID         string
	ClubID     string
	UserID     string
	PCNumber   int
	StartTime  time.Time
	EndTime    time.Time
	TotalPrice float64
	Status     string // "active", "cancelled", "completed"
	CreatedAt  time.Time
}
