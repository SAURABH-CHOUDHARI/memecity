package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SAURABH-CHOUDHARI/memecity/internals/models"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
)

type MemeWithVotes struct {
	models.Meme
	Upvotes   int `json:"upvotes"`
	Downvotes int `json:"downvotes"` // ➕ Added
}

func GetLeaderboardMemes(conn storage.Repository, limit, offset int) ([]MemeWithVotes, error) {
	var memes []MemeWithVotes
	cacheKey := fmt.Sprintf("leaderboard:limit=%d:offset=%d", limit, offset)
	ctx := context.Background()

	// Try Redis cache
	if cached, err := conn.RedisClient.Get(ctx, cacheKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(cached), &memes); err == nil {
			return memes, nil
		}
	}

	// Updated SQL to count both upvotes and downvotes
	err := conn.DB.
		Table("memes").
		Select(`
			memes.*, 
			COUNT(votes.id) FILTER (WHERE votes.type = true) AS upvotes,
			COUNT(votes.id) FILTER (WHERE votes.type = false) AS downvotes
		`).
		Joins("LEFT JOIN votes ON votes.meme_id = memes.id").
		Preload("Owner").
		Group("memes.id").
		Order("upvotes DESC, memes.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&memes).Error

	if err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(memes); err == nil {
		conn.RedisClient.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return memes, nil
}


