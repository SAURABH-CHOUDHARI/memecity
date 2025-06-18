package models

import (
	"github.com/google/uuid"
	"time"
)

type Vote struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MemeID    uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	Type      bool      `gorm:"not null"` // true = upvote (1), false = downvote (0)
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relations
	Meme Meme `gorm:"foreignKey:MemeID;references:ID"`
	User User `gorm:"foreignKey:UserID;references:ID"`
}

// Helper methods for cleaner code
func (v *Vote) IsUpvote() bool {
	return v.Type
}

func (v *Vote) IsDownvote() bool {
	return !v.Type
}

func (v *Vote) SetUpvote() {
	v.Type = true
}

func (v *Vote) SetDownvote() {
	v.Type = false
}