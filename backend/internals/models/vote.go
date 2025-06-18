package models

import (
	"time"
	"github.com/google/uuid"
)

type Vote struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MemeID    uuid.UUID `gorm:"type:uuid;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Type      string    `gorm:"type:varchar(10);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
