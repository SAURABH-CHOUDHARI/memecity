package services

import (
	"errors"

	"github.com/SAURABH-CHOUDHARI/memecity/internals/models"
	"github.com/SAURABH-CHOUDHARI/memecity/pkg/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidBid        = errors.New("invalid bid")
	ErrInsufficientCredits = errors.New("not enough credits")
	ErrUnauthorized      = errors.New("cannot bid on your own meme")
	ErrNotFound          = errors.New("meme not found")
)

func PlaceBid(conn storage.Repository, userID uuid.UUID, memeIDStr string, bidAmount int) error {
	if bidAmount <= 0 {
		return ErrInvalidBid
	}

	memeID, err := uuid.Parse(memeIDStr)
	if err != nil {
		return err
	}

	var meme models.Meme
	if err := conn.DB.First(&meme, "id = ?", memeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	if !meme.OnSale {
		return ErrInvalidBid
	}

	if meme.OwnerID == userID {
		return ErrUnauthorized
	}

	var user models.User
	if err := conn.DB.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	if user.Credits < bidAmount {
		return ErrInsufficientCredits
	}

	// Create bid
	bid := models.Bid{
		ID:      uuid.New(),
		MemeID:  memeID,
		UserID:  userID,
		Credits: bidAmount,
	}

	// Transaction: deduct credits and save bid
	return conn.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&bid).Error; err != nil {
			return err
		}
		if err := tx.Model(&user).Update("credits", user.Credits-bidAmount).Error; err != nil {
			return err
		}
		return nil
	})
}
