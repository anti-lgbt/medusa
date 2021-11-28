package models

import "time"

type TrendingMusic struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	MusicID        int64     `gorm:"type:bigint;not null;index"`
	TotalViewCount int64     `gorm:"type:integer;not null;index;default:0"`
	DayViewCount   int64     `gorm:"type:integer;not null;index;default:0"`
	ReleaseAt      time.Time `json:"release_at" gorm:"type:timestamp(0);not null;index"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp(0);not null;index"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:timestamp(0);not null;index"`
	Music          *Music    `json:"-"`
}
