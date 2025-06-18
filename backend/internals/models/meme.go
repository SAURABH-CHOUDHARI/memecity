package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Meme struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title     string         `gorm:"not null"`
	ImageURL  string         `gorm:"type:text;not null"`
	Tags      datatypes.JSON `gorm:"type:jsonb"`
	Caption   string
	OwnerID   uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
