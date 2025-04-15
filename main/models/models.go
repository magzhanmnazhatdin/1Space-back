package models

import "time"

type ComputerClub struct {
	ID           string  `json:"id" firestore:"id"`
	Name         string  `json:"name" firestore:"name"`
	Address      string  `json:"address" firestore:"address"`
	PricePerHour float64 `json:"price_per_hour" firestore:"price_per_hour"`
	AvailablePCs int     `json:"available_pcs" firestore:"available_pcs"`
}

type Booking struct {
	ID         string    `json:"id" firestore:"id"`
	ClubID     string    `json:"club_id" firestore:"club_id"`
	ClubName   string    `json:"club_name,omitempty" firestore:"-"`
	UserID     string    `json:"user_id" firestore:"user_id"`
	PCNumber   int       `json:"pc_number" firestore:"pc_number"`
	StartTime  time.Time `json:"start_time" firestore:"start_time"`
	EndTime    time.Time `json:"end_time" firestore:"end_time"`
	TotalPrice float64   `json:"total_price" firestore:"total_price"`
	Status     string    `json:"status" firestore:"status"` // "active", "cancelled", "completed"
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
}

type Computer struct {
	ID          string `json:"id" firestore:"id"`
	ClubID      string `json:"club_id" firestore:"club_id"`
	PCNumber    int    `json:"pc_number" firestore:"pc_number"`
	Description string `json:"description" firestore:"description"`
	IsAvailable bool   `json:"is_available" firestore:"is_available"`
}
