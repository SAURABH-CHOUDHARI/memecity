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

func VoteOnMeme(conn storage.Repository, userID uuid.UUID, memeID uuid.UUID, voteType string) (string, models.Meme, error) {
    ctx := context.Background()

    var meme models.Meme
    memeCacheKey := "meme:" + memeID.String()

    if val, err := conn.RedisClient.Get(ctx, memeCacheKey).Result(); err == nil {
        _ = json.Unmarshal([]byte(val), &meme)
    } else {
        if err := conn.DB.First(&meme, "id = ?", memeID).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return "", meme, ErrMemeNotFound
            }
            return "", meme, err
        }
    }

    var user models.User
    userCacheKey := "user:" + userID.String()
    if val, err := conn.RedisClient.Get(ctx, userCacheKey).Result(); err == nil {
        _ = json.Unmarshal([]byte(val), &user)
    } else {
        if err := conn.DB.First(&user, "id = ?", userID).Error; err != nil {
            return "", meme, err
        }
    }

    var isUpvote bool
    switch voteType {
    case "up":
        isUpvote = true
    case "down":
        isUpvote = false
    default:
        return "", meme, ErrInvalidVoteType
    }

    var existing models.Vote
    err := conn.DB.First(&existing, "user_id = ? AND meme_id = ?", userID, memeID).Error

    var action string
    var message string

    if err == nil {
        if existing.Type == isUpvote {
            if err := conn.DB.Delete(&existing).Error; err != nil {
                return "", meme, err
            }
            action = "removed"
            message = "Vote removed successfully"
        } else {
            if err := conn.DB.Model(&existing).Update("type", isUpvote).Error; err != nil {
                return "", meme, err
            }
            action = "flipped"
            message = "Vote flipped successfully"
        }
    } else {
        if !errors.Is(err, gorm.ErrRecordNotFound) {
            return "", meme, err
        }
        newVote := models.Vote{
            ID:     uuid.New(),
            MemeID: memeID,
            UserID: userID,
            Type:   isUpvote,
        }
        if err := conn.DB.Create(&newVote).Error; err != nil {
            return "", meme, err
        }
        action = "created"
        message = "Vote created successfully"
    }

    // Count upvotes/downvotes fresh from DB
    var upvotes, downvotes int64
    conn.DB.Model(&models.Vote{}).Where("meme_id = ? AND type = ?", memeID, true).Count(&upvotes)
    conn.DB.Model(&models.Vote{}).Where("meme_id = ? AND type = ?", memeID, false).Count(&downvotes)

    // Update meme object with current counts
    meme.Upvotes = int(upvotes)
    meme.Downvotes = int(downvotes)

    // Re-cache updated meme
    data, _ := json.Marshal(meme)
    conn.RedisClient.Set(ctx, memeCacheKey, data, 10*time.Minute)

    // WebSocket broadcast
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

    return message, meme, nil
}

