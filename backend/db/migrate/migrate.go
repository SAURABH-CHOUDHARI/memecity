package migrate

import (
	"log"
	"gorm.io/gorm"
	"github.com/SAURABH-CHOUDHARI/memecity/internals/models"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Meme{},
		&models.Vote{},
		&models.Bid{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("✅ Database migrated successfully")
}
