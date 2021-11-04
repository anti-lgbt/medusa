package models

import "time"

type Popular struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id" gorm:"type:bigint;not null;index"`
	AlbumID   int64     `json:"album_id" gorm:"type:bigint;not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp(0);not null;index"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp(0);not null;index"`
}
