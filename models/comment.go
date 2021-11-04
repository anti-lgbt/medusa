package models

import (
	"database/sql"
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/models/datatypes"
)

type Comment struct {
	ID        int64         `json:"id" gorm:"primaryKey"`
	UserID    int64         `json:"-" gorm:"bigint;not null;index"`
	MusicID   sql.NullInt64 `json:"-" gorm:"bigint;index"`
	AlbumID   sql.NullInt64 `json:"-" gorm:"bigint;index"`
	Content   string        `json:"content" gorm:"character varying;not null"`
	CreatedAt time.Time     `json:"created_at" gorm:"type:timestamp(0);not null;index"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"type:timestamp(0);not null;index"`
	Musics    []*Music      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Albums    []*Album      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Replys    []*Reply      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
}

func (c *Comment) Delete() {
	config.Database.Delete(&c)
}

func (c *Comment) Reply(user_id int64, content string) *Reply {
	reply := &Reply{
		UserID:    user_id,
		CommentID: c.ID,
		Content:   content,
	}

	config.Database.Create(&reply)

	return reply
}

func (c *Comment) Like(user_id int64) error {
	var l *Like

	if result := config.Database.First(&l, "user_id = ? AND comment_id = ?", user_id, c.ID); result.Error == nil {
		return result.Error
	}

	like := Like{
		UserID: user_id,
		CommentID: datatypes.NullInt64{
			Int64: c.ID,
			Valid: true,
		},
	}

	config.Database.Create(&like)
	return nil
}

func (c *Comment) UnLike(user_id int64) {
	config.Database.Where("user_id = ? AND comment_id = ?", user_id, c.ID).Delete(&Like{})
}

func (c *Comment) LikeCount() int64 {
	var count int64

	config.Database.Model(&Like{}).Where("comment_id = ?", c.ID).Count(&count)

	return count
}
