package services

import (
	"errors"

	"github.com/SAURABH-CHOUDHARI/memecity/internals/models"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidVoteType = errors.New("vote type must be 'up' or 'down'")
	ErrMemeNotFound    = errors.New("meme not found")
)

func VoteOnMeme(conn storage.Repository, userID uuid.UUID, memeID uuid.UUID, voteType string) error {
	var meme models.Meme
	if err := conn.DB.First(&meme, "id = ?", memeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMemeNotFound
		}
		return err
	}

	isUpvote := false
	switch voteType {
	case "up":
		isUpvote = true
	case "down":
		isUpvote = false
	default:
		return ErrInvalidVoteType
	}

	var existing models.Vote
	err := conn.DB.First(&existing, "user_id = ? AND meme_id = ?", userID, memeID).Error
	if err == nil {
		// Already voted
		if existing.Type == isUpvote {
			// Same vote → toggle off (delete)
			return conn.DB.Delete(&existing).Error
		}
		// Flip vote
		return conn.DB.Model(&existing).Update("type", isUpvote).Error
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// New vote
	vote := models.Vote{
		ID:     uuid.New(),
		MemeID: memeID,
		UserID: userID,
		Type:   isUpvote,
	}
	return conn.DB.Create(&vote).Error
}
