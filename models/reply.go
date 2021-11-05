package models

import (
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/models/datatypes"
)

type Reply struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"-" gorm:"type:bigint;not null;index"`
	CommentID int64     `json:"-" gorm:"type:bigint;not null;index"`
	Content   string    `json:"content" gorm:"type:character varying;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp(0);not null;index"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp(0);not null;index"`
	Likes     []*Like   `json:"-" gorm:"foreignKey:CommentID;references:ID;constraint:OnDelete:CASCADE"`
}

func (r *Reply) Delete() {
	config.Database.Delete(&r)
}

func (r *Reply) Like(user_id int64) error {
	var l *Like

	if result := config.Database.First(&l, "user_id = ? AND reply_id = ?", user_id, r.ID); result.Error == nil {
		return result.Error
	}

	like := Like{
		UserID: user_id,
		ReplyID: datatypes.NullInt64{
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
