package models

import (
	"time"
	"github.com/google/uuid"
)

type Bid struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MemeID    uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	Credits   int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	
	// Relations
	Meme Meme `gorm:"foreignKey:MemeID;references:ID"`
	User User `gorm:"foreignKey:UserID;references:ID"`
}
