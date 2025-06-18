package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username  string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	Credits   int       `gorm:"default:1000"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
