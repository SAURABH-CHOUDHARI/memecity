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

func VoteOnMeme(conn storage.Repository, userID uuid.UUID, memeID uuid.UUID, voteType string) (string, error) {
	// Check if meme exists
	var meme models.Meme
	if err := conn.DB.First(&meme, "id = ?", memeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrMemeNotFound
		}
		return "", err
	}

	// Normalize vote type
	var isUpvote bool
	switch voteType {
	case "up":
		isUpvote = true
	case "down":
		isUpvote = false
	default:
		return "", ErrInvalidVoteType
	}

	// Check for existing vote
	var existing models.Vote
	err := conn.DB.First(&existing, "user_id = ? AND meme_id = ?", userID, memeID).Error

	if err == nil {
		// Vote exists
		if existing.Type == isUpvote {
			// Same vote → toggle off
			if err := conn.DB.Delete(&existing).Error; err != nil {
				return "", err
			}
			if isUpvote {
				return "Removed upvote", nil
			}
			return "Removed downvote", nil
		}

		// Flip vote
		if err := conn.DB.Model(&existing).Update("type", isUpvote).Error; err != nil {
			return "", err
		}
		if isUpvote {
			return "Flipped to upvote", nil
		}
		return "Flipped to downvote", nil
	}

	// New vote
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	newVote := models.Vote{
		ID:     uuid.New(),
		MemeID: memeID,
		UserID: userID,
		Type:   isUpvote,
	}
	if err := conn.DB.Create(&newVote).Error; err != nil {
		return "", err
	}

	if isUpvote {
		return "Upvoted", nil
	}
	return "Downvoted", nil
}

