package models

import (
	"database/sql"
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/models/datatypes"
	"github.com/anti-lgbt/medusa/types"
)

type Music struct {
	ID          int64                `json:"id" gorm:"primaryKey"`
	UserID      int64                `json:"user_id" gorm:"type:bigint;not null;index"`
	Name        string               `json:"name" gorm:"type:character varying;not null;index"`
	Description datatypes.NullString `json:"description" gorm:"type:character varying;not null;index"`
	State       types.MusicState     `json:"state" gorm:"type:character varying(10);not null;index"`
	ViewCount   int64                `json:"view_count" gorm:"type:integer;not null;index;default:0"`
	Path        string               `json:"-" gorm:"type:character varying;not null"`
	Image       datatypes.NullString `json:"-" gorm:"type:character varying"`
	CreatedAt   time.Time            `json:"created_at" gorm:"type:timestamp(0);not null;index"`
	UpdatedAt   time.Time            `json:"updated_at" gorm:"type:timestamp(0);not null;index"`
	MusicAlbums []*MusicAlbum        `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Likes       []*Like              `json:"-" gorm:"foreignKey:MusicID;references:ID;constraint:OnDelete:CASCADE"`
	Comments    []*Comment           `json:"-" gorm:"foreignKey:MusicID;references:ID;constraint:OnDelete:CASCADE"`
}

func (m *Music) Delete() {
	config.Database.Delete(&m)
}

func (m *Music) Comment(user_id int64, content string) *Comment {
	comment := &Comment{
		UserID: user_id,
		MusicID: sql.NullInt64{
			Int64: m.ID,
			Valid: true,
		},
		Content: content,
	}

	config.Database.Create(&comment)

	return comment
}

func (m *Music) Like(user_id int64) error {
	var l *Like

	if result := config.Database.First(&l, "user_id = ? AND music_id = ?", user_id, m.ID); result.Error == nil {
		return result.Error
	}

	like := Like{
		UserID: user_id,
		MusicID: datatypes.NullInt64{
			Int64: m.ID,
			Valid: true,
		},
	}

	config.Database.Create(&like)
	return nil
}

func (m *Music) UnLike(user_id int64) {
	config.Database.Where("user_id = ? AND music_id = ?", user_id, m.ID).Delete(&Like{})
}

func (m *Music) LikeCount() int64 {
	var count int64

	config.Database.Model(&Like{}).Where("music_id = ?", m.ID).Count(&count)

	return count
}
