package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/SAURABH-CHOUDHARI/memecity/internals/models"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/ws"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidVoteType = errors.New("vote type must be 'up' or 'down'")
	ErrMemeNotFound    = errors.New("meme not found")
)

func VoteOnMeme(conn storage.Repository, userID uuid.UUID, memeID uuid.UUID, voteType string) (string, error) {
	ctx := context.Background()

	// ----- Fetch Meme (try Redis) -----
	var meme models.Meme
	memeCacheKey := "meme:" + memeID.String()
	if val, err := conn.RedisClient.Get(ctx, memeCacheKey).Result(); err == nil {
		_ = json.Unmarshal([]byte(val), &meme)
	} else {
		if err := conn.DB.First(&meme, "id = ?", memeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", ErrMemeNotFound
			}
			return "", err
		}
		data, _ := json.Marshal(meme)
		conn.RedisClient.Set(ctx, memeCacheKey, data, 10*time.Minute)
	}

	// ----- Fetch User (try Redis) -----
	var user models.User
	userCacheKey := "user:" + userID.String()
	if val, err := conn.RedisClient.Get(ctx, userCacheKey).Result(); err == nil {
		_ = json.Unmarshal([]byte(val), &user)
	} else {
		if err := conn.DB.First(&user, "id = ?", userID).Error; err != nil {
			return "", err
		}
		data, _ := json.Marshal(user)
		conn.RedisClient.Set(ctx, userCacheKey, data, 10*time.Minute)
	}

	// ----- Normalize vote type -----
	var isUpvote bool
	switch voteType {
	case "up":
		isUpvote = true
	case "down":
		isUpvote = false
	default:
		return "", ErrInvalidVoteType
	}

	// ----- Check for existing vote -----
	var existing models.Vote
	err := conn.DB.First(&existing, "user_id = ? AND meme_id = ?", userID, memeID).Error

	var action string

	if err == nil {
		// Vote exists
		if existing.Type == isUpvote {
			// Toggle off
			if err := conn.DB.Delete(&existing).Error; err != nil {
				return "", err
			}
			action = "removed"
		} else {
			// Flip
			if err := conn.DB.Model(&existing).Update("type", isUpvote).Error; err != nil {
				return "", err
			}
			action = "flipped"
		}
	} else {
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
		action = "created"
	}

	// ----- Send WebSocket Event -----
	event := ws.VoteEvent{
		Type:       "vote",
		MemeID:     meme.ID,
		MemeTitle:  meme.Title,
		ImageURL:   meme.ImageURL,
		UserID:     user.ID,
		Username:   user.Username,
		ProfilePic: user.ProfilePic,
		VoteType:   voteType,
		Action:     action,
		Timestamp:  time.Now().UTC(),
	}

	payload, _ := json.Marshal(event)
	conn.Hub.Broadcast <- payload

	// ----- Response message -----
	switch action {
	case "created":
		if isUpvote {
			return "Upvoted", nil
		}
		return "Downvoted", nil
	case "flipped":
		if isUpvote {
			return "Flipped to upvote", nil
		}
		return "Flipped to downvote", nil
	case "removed":
		if isUpvote {
			return "Removed upvote", nil
		}
		return "Removed downvote", nil
	default:
		return "Vote action done", nil
	}
}
