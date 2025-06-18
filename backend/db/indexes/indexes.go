package indexes

import (
	"gorm.io/gorm"
)

func CreateIndexes(db *gorm.DB) {
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_votes_user_meme ON votes(user_id, meme_id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_votes_meme_id ON votes(meme_id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_bids_meme_id ON bids(meme_id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_memes_owner_id ON memes(owner_id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_memes_onsale ON memes(on_sale)`)
}
