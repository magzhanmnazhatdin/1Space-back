package entities

// Computer — доменная сущность компьютера в клубе.
type Computer struct {
	ID          string `firestore:"id"           json:"id"`
	ClubID      string `firestore:"club_id"      json:"club_id"`
	PCNumber    int    `firestore:"pc_number"    json:"pc_number"`
	Description string `firestore:"description"  json:"description"`
	IsAvailable bool   `firestore:"is_available" json:"is_available"`
}
