package entities

import "time"

// Booking is the domain entity representing a reservation.
type Booking struct {
	ID         string    `firestore:"id"           json:"id"`
	ClubID     string    `firestore:"club_id"      json:"club_id"`
	UserID     string    `firestore:"user_id"      json:"user_id"`
	PCNumber   int       `firestore:"pc_number"    json:"pc_number"`
	StartTime  time.Time `firestore:"start_time"   json:"start_time"`
	EndTime    time.Time `firestore:"end_time"     json:"end_time"`
	TotalPrice float64   `firestore:"total_price"  json:"total_price"`
	Status     string    `firestore:"status"       json:"status"`
	CreatedAt  time.Time `firestore:"created_at"   json:"created_at"`
}
