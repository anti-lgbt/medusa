package models

import (
	"database/sql"
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/entities"
)

type Comment struct {
	ID        int64         `json:"id" gorm:"primaryKey"`
	UserID    int64         `json:"-" gorm:"bigint;not null;index"`
	MusicID   sql.NullInt64 `json:"-" gorm:"bigint;index"`
	AlbumID   sql.NullInt64 `json:"-" gorm:"bigint;index"`
	Content   string        `json:"content" gorm:"character varying;not null"`
	CreatedAt time.Time     `json:"created_at" gorm:"type:timestamp(0);not null;index"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"type:timestamp(0);not null;index"`
	Likes     []*Like       `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	User      *User         `json:"-"`
	Music     *Music        `json:"-"`
	Album     *Album        `json:"-"`
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
		CommentID: sql.NullInt64{
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

func (c *Comment) ToEntity() *entities.Comment {
	var replies []*Reply

	config.Database.Find(&replies, "comment_id = ?", c.ID)

	reply_entities := make([]*entities.Reply, 0)

	for _, reply := range replies {
		reply_entities = append(reply_entities, reply.ToEntity())
	}

	return &entities.Comment{
		ID:        c.ID,
		Content:   c.Content,
		LikeCount: c.LikeCount(),
		Replies:   reply_entities,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
