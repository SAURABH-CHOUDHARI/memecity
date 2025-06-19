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
	ErrInvalidBid          = errors.New("invalid bid")
	ErrInsufficientCredits = errors.New("not enough credits")
	ErrUnauthorized        = errors.New("cannot bid on your own meme")
	ErrNotFound            = errors.New("meme not found")
)

func PlaceBid(conn storage.Repository, userID uuid.UUID, memeIDStr string, bidAmount int) error {
	if bidAmount <= 0 {
		return ErrInvalidBid
	}

	memeID, err := uuid.Parse(memeIDStr)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// ---------- Meme Lookup ----------
	var meme models.Meme
	memeCacheKey := "meme:" + memeID.String()

	if val, err := conn.RedisClient.Get(ctx, memeCacheKey).Result(); err == nil {
		_ = json.Unmarshal([]byte(val), &meme)
	} else {
		if err := conn.DB.First(&meme, "id = ?", memeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		// Cache meme in Redis
		data, _ := json.Marshal(meme)
		conn.RedisClient.Set(ctx, memeCacheKey, data, 10*time.Minute)
	}

	if !meme.OnSale {
		return ErrInvalidBid
	}
	if meme.OwnerID == userID {
		return ErrUnauthorized
	}

	// ---------- User Lookup ----------
	var user models.User
	userCacheKey := "user:" + userID.String()

	if val, err := conn.RedisClient.Get(ctx, userCacheKey).Result(); err == nil {
		_ = json.Unmarshal([]byte(val), &user)
	} else {
		if err := conn.DB.First(&user, "id = ?", userID).Error; err != nil {
			return err
		}
		// Cache user in Redis
		data, _ := json.Marshal(user)
		conn.RedisClient.Set(ctx, userCacheKey, data, 10*time.Minute)
	}

	if user.Credits < bidAmount {
		return ErrInsufficientCredits
	}

	// ---------- Create Bid ----------
	bid := models.Bid{
		ID:      uuid.New(),
		MemeID:  meme.ID,
		UserID:  user.ID,
		Credits: bidAmount,
	}

	// ---------- WebSocket Event ----------
	event := ws.BidEvent{
		Type:       "bid",
		MemeID:     meme.ID,
		MemeTitle:  meme.Title,
		ImageURL:   meme.ImageURL,
		UserID:     user.ID,
		Username:   user.Username,
		ProfilePic: user.ProfilePic,
		Credits:    bidAmount,
		Timestamp:  time.Now().UTC(),
	}
	payload, _ := json.Marshal(event)
	conn.Hub.Broadcast <- payload

	// ---------- Transaction: Save Bid + Deduct Credits ----------
	return conn.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&bid).Error; err != nil {
			return err
		}
		if err := tx.Model(&user).Update("credits", user.Credits-bidAmount).Error; err != nil {
			return err
		}

		// Fetch updated credits and cache them
	var updatedUser models.User
	if err := tx.First(&updatedUser, "id = ?", user.ID).Error; err == nil {
		data, _ := json.Marshal(updatedUser)
		conn.RedisClient.Set(ctx, userCacheKey, data, 10*time.Minute)
	}

		return nil
	})
}
