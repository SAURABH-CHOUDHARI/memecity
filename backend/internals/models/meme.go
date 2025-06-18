package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Meme struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title     string         `gorm:"default:'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRvzTAz0n_p5GyzLkmRP0JZ-r_Cd0dEfUmbjw&s'"`
	ImageURL  string         `gorm:"type:text;not null"`
	Tags      datatypes.JSON `gorm:"type:jsonb"`
	Caption   string         `gorm:"not null"`	
	Price     int            `gorm:"not null;default:0"`
	OnSale    bool           `gorm:"default:false"`
	OwnerID   uuid.UUID      `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	
	// Relations
	Owner User `gorm:"foreignKey:OwnerID;references:ID"`
}