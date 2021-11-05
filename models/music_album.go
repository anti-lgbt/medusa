package models

import "time"

type MusicAlbum struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	MusicID   int64     `json:"music_id" gorm:"type:bigint;not null;uniqueIndex:idx_music_id_and_album_id;index"`
	AlbumID   int64     `json:"album_id" gorm:"type:bigint;not null;uniqueIndex:idx_strategy_id_and_market_id;index"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp(0);not null;index"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp(0);not null;index"`
}
