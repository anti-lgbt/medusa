package models

import (
	"database/sql"
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/entities"
)

type Reply struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"-" gorm:"type:bigint;not null;index"`
	CommentID int64     `json:"-" gorm:"type:bigint;not null;index"`
	Content   string    `json:"content" gorm:"type:character varying;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp(0);not null;index"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp(0);not null;index"`
	Likes     []*Like   `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	User      *User
	Comment   *Comment
}

func (r *Reply) Delete() {
	config.Database.Delete(&r)
}

func (r *Reply) Liked(user_id int64) bool {
	var l *Like

	result := config.Database.First(&l, "user_id = ? AND reply_id = ?", user_id, r.ID)

	return result.Error == nil
}

func (r *Reply) Like(user_id int64) error {
	var l *Like

	if result := config.Database.First(&l, "user_id = ? AND reply_id = ?", user_id, r.ID); result.Error == nil {
		return result.Error
	}

	like := Like{
		UserID: user_id,
		ReplyID: sql.NullInt64{
			Int64: r.ID,
			Valid: true,
		},
	}

	config.Database.Create(&like)
	return nil
}

func (r *Reply) UnLike(user_id int64) {
	config.Database.Where("user_id = ? AND reply_id = ?", user_id, r.ID).Delete(&Like{})
}

func (r *Reply) LikeCount() int64 {
	var count int64

	config.Database.Model(&Like{}).Where("reply_id = ?", r.ID).Count(&count)

	return count
}

func (r *Reply) ToEntity(user_id sql.NullInt64) *entities.Reply {
	entity := &entities.Reply{
		ID:        r.ID,
		Content:   r.Content,
		LikeCount: r.LikeCount(),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	if user_id.Valid {
		entity.Liked = r.Liked(user_id.Int64)
	}

	return entity
}
