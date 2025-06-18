package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username   string    `gorm:"not null"`
	Email      string    `gorm:"unique;not null"`
	Credits    int       `gorm:"default:1000"`
	Password   string    `gorm:"not null"`
	ProfilePic string    `gorm:"default:'https://images.unsplash.com/photo-1567963070256-729fb28b079c?q=80&w=576&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
