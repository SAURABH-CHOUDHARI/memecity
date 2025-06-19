package ws

import (
	"time"

	"github.com/google/uuid"
)

type BidEvent struct {
	Type        string    `json:"type"` // "bid"
	MemeID      uuid.UUID `json:"meme_id"`
	MemeTitle   string    `json:"meme_title"`
	ImageURL    string    `json:"image_url"`
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	ProfilePic  string    `json:"profile_pic"`
	Credits     int       `json:"credits"`
	Timestamp   time.Time `json:"timestamp"`
}

type VoteEvent struct {
	Type       string    `json:"type"`        // should be "vote"
	MemeID     uuid.UUID `json:"meme_id"`
	MemeTitle  string    `json:"meme_title"`
	ImageURL   string    `json:"image_url"`
	UserID     uuid.UUID `json:"user_id"`
	Username   string    `json:"username"`
	ProfilePic string    `json:"profile_pic"`
	VoteType   string    `json:"vote_type"`   // "up" or "down"
	Action     string    `json:"action"`      // "created", "flipped", "removed"
	Timestamp  time.Time `json:"timestamp"`
}