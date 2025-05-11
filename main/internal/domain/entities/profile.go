// internal/domain/entities/profile.go
package entities

// Profile — профиль пользователя, хранится отдельно от Auth.
type Profile struct {
	UserID      string `firestore:"user_id" json:"user_id"`
	ProfileName string `firestore:"profile_name" json:"profile_name"`
	AvatarURL   string `firestore:"avatar_url,omitempty" json:"avatar_url,omitempty"`
	Bio         string `firestore:"bio,omitempty" json:"bio,omitempty"`
}
